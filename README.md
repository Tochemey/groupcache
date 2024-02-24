# groupcache

[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/Tochemey/groupcache/on-pull-request.yaml?branch=main)](https://github.com/Tochemey/groupcache/actions/workflows/on-pull-request.yaml)

groupcache is a caching and cache-filling library, intended as a replacement for memcached in many cases.

For API docs and examples, see http://godoc.org/github.com/tochemey/groupcache/v2

## Table of Content

- [Overview](#overview)
- [Modifications](#modifications-from-original-library)
- [Comparison to Memcached](#comparing-groupcache-to-memcached)
- [Loading Process](#loading-process)
- [Example](#example)
- [Clustering](#clustering)
    - [Discovery Providers](#built-in-discovery-providers)
        - [Kubernetes](#kubernetes-discovery-provider-setup)
        - [mDNS](#mdns-discovery-provider-setup)
        - [NATS](#nats-discovery-provider-setup)

## Overview

A modified version of [group cache](https://github.com/mailgun/groupcache) with support:
- upgrade the protobuf API.
- service discovery
- simple logger interface with a default logger based upon uber zap library
- reorganise the proto definitions into a folder called: [protos](./protos). See how to generate the pbs [here](#protobuf-stub-generation)
   
### Modifications from original library

In addition to these [modifications](https://github.com/mailgun/groupcache?tab=readme-ov-file#modifications-from-original-library),
* Logger interface to help add custom logging framework
* Service Discovery to help discover other group cache automatically. At the moment the following providers are implemented:
  - the [kubernetes](https://kubernetes.io/docs/home/) [api integration](./discovery/kubernetes) is fully functional
  - the [mDNS](https://datatracker.ietf.org/doc/html/rfc6762) and [DNS-SD](https://tools.ietf.org/html/rfc6763)
  - the [NATS](https://nats.io/) [integration](./discovery/nats) is fully functional
* `HTTPPoolOpts` is no longer exposed to the caller. It has been renamed to `httpPoolOpts`. The same applies to `HTTPPool`. This allows to start groupcache with a service discovery provider.
* Upgrade the protocol buffer API

## Comparing Groupcache to memcached

### **Like memcached**, groupcache:

 * shards by key to select which peer is responsible for that key

### **Unlike memcached**, groupcache:

 * does not require running a separate set of servers, thus massively
   reducing deployment/configuration pain.  groupcache is a client
   library as well as a server.  It connects to its own peers.

 * comes with a cache filling mechanism.  Whereas memcached just says
   "Sorry, cache miss", often resulting in a thundering herd of
   database (or whatever) loads from an unbounded number of clients
   (which has resulted in several fun outages), groupcache coordinates
   cache fills such that only one load in one process of an entire
   replicated set of processes populates the cache, then multiplexes
   the loaded value to all callers.

 * does not support versioned values.  If key "foo" is value "bar",
   key "foo" must always be "bar".

## Loading process

In a nutshell, a groupcache lookup of **Get("foo")** looks like:

(On machine #5 of a set of N machines running the same code)

 1. Is the value of "foo" in local memory because it's super hot?  If so, use it.

 2. Is the value of "foo" in local memory because peer #5 (the current
    peer) is the owner of it?  If so, use it.

 3. Amongst all the peers in my set of N, am I the owner of the key
    "foo"?  (e.g. does it consistent hash to 5?)  If so, load it.  If
    other callers come in, via the same process or via RPC requests
    from peers, they block waiting for the load to finish and get the
    same answer.  If not, RPC to the peer that's the owner and get
    the answer.  If the RPC fails, just load it locally (still with
    local dup suppression).

## Example

```go
import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/tochemey/groupcache/v2"
    "github.com/tochemey/groupcache/v2/discovery"
    "github.com/tochemey/groupcache/v2/discovery/kubernetes"
)

func ExampleUsage() {

    // NOTE: It is important each node running the groupcache has the env vars properly set:
    // GROUP_PORT, NODE_NAME and NODE_IP
    // That the service discovery can properly identify the running instance

    // Create an instance of the discovery service.
    // For instance let us use kubernetes
    provider := kubernetes.New()

    // Create the discovery options
    // For kubernetes we only need the namespace and the application name
    application := "users"
    namespace := "default"
	 
    options := discovery.Config{
        kubernetes.ApplicationName: application,
        kubernetes.Namespace:       namespace,
    }

    // Create an instance of the service discovery
    serviceDiscovery := discovery.NewServiceDiscovery(provider, options)

    // Create an instance of the cluster
    ctx := context.Background()
    node := groupcache.NewNode(ctx, serviceDiscovery)
    
    // Start the cluster node
    err := node.Start(ctx)
    
    // Stop the cluster node
    defer node.Stop(ctx)
	 
    // Create a new group cache with a max cache size of 3MB
    group := groupcache.NewGroup("users", 3000000, groupcache.GetterFunc(
        func(ctx context.Context, id string, dest groupcache.Sink) error {

            // Returns a protobuf struct `User`
            user, err := fetchUserFromMongo(ctx, id)
            if err != nil {
                return err
            }

            // Set the user in the groupcache to expire after 5 minutes
            return dest.SetProto(&user, time.Now().Add(time.Minute*5))
        },
    ))

    user := new(User)

    ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
    defer cancel()

    if err := group.Get(ctx, "12345", groupcache.ProtoSink(user)); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("-- User --\n")
    fmt.Printf("Id: %s\n", user.Id)
    fmt.Printf("Name: %s\n", user.Name)
    fmt.Printf("Age: %d\n", user.Age)
    fmt.Printf("IsSuper: %t\n", user.IsSuper)

    // Remove the key from the groupcache
    if err := group.Remove(ctx, "12345"); err != nil {
        log.Fatal(err)
    }
}

```

## Clustering

The cluster engine depends upon the [discovery](./discovery/discovery.go) mechanism to find other nodes in the cluster.

At the moment the following providers are implemented:

- the [kubernetes](https://kubernetes.io/docs/home/) [api integration](./discovery/kubernetes) is fully functional
- the [mDNS](https://datatracker.ietf.org/doc/html/rfc6762) and [DNS-SD](https://tools.ietf.org/html/rfc6763)
- the [NATS](https://nats.io/) [integration](./discovery/nats) is fully functional

Note: One can add additional discovery providers using the following [interface](./discovery/provider.go)

In addition, one needs to set the following environment variables irrespective of the discovery provider to help
identify the host node on which the cluster service is running:

- `NODE_NAME`: the node name. For instance in kubernetes one can just get it from the `metadata.name`
- `NODE_IP`: the node host address. For instance in kubernetes one can just get it from the `status.podIP`
- `GROUP_PORT`: the port used by the discovery provider to communicate.

_Note: Depending upon the discovery provider implementation, the `NODE_NAME` and `NODE_IP` can be the same._

### Built-in Discovery Providers

#### Kubernetes Discovery Provider Setup

To get the kubernetes discovery working as expected, the following pod labels need to be set:

- `app.kubernetes.io/part-of`: set this label with the actor system name
- `app.kubernetes.io/component`: set this label with the application name
- `app.kubernetes.io/name`: set this label with the application name

In addition, each node _is required to have the following port open_ with the following ports name for the cluster
engine to work as expected:

- `group-port`: help the gossip protocol engine. This is actually the kubernetes discovery port

##### Get Started

```go
const (
    namespace = "default"
    applicationName = "accounts"
)
// instantiate the k8 discovery provider
disco := kubernetes.NewDiscovery()
// define the discovery options
discoOptions := discovery.Config{
    kubernetes.ApplicationName: applicationName,
    kubernetes.Namespace:       namespace,
}
// define the service discovery
serviceDiscovery := discovery.NewServiceDiscovery(disco, discoOptions)
// start the cluster
```

##### Role Based Access

You’ll also have to grant the Service Account that your pods run under access to list pods. The following configuration
can be used as a starting point.
It creates a Role, pod-reader, which grants access to query pod information. It then binds the default Service Account
to the Role by creating a RoleBinding.
Adjust as necessary:

```
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: pod-reader
rules:
  - apiGroups: [""] # "" indicates the core API group
    resources: ["pods"]
    verbs: ["get", "watch", "list"]
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: read-pods
subjects:
  # Uses the default service account. Consider creating a new one.
  - kind: ServiceAccount
    name: default
roleRef:
  kind: Role
  name: pod-reader
  apiGroup: rbac.authorization.k8s.io
```

#### mDNS Discovery Provider Setup

- `Service Name`: the service name
- `Domain`: The mDNS discovery domain
- `Port`: The mDNS discovery port
- `IPv6`: States whether to lookup for IPv6 addresses.

#### NATS Discovery Provider Setup

To use the NATS discovery provider one needs to provide the following:

- `NATS Server Address`: the NATS Server address
- `NATS Subject`: the NATS subject to use
- `Application Name`: the application name

```go
const (
    natsServerAddr = "nats://localhost:4248"
    natsSubject = "groupcache-gossip"
    applicationName = "accounts"
)
// instantiate the NATS discovery provider
disco := nats.NewDiscovery()
// define the discovery options
discoOptions := discovery.Config{
    ApplicationName: applicationName,
    NatsServer:      natsServer,
    NatsSubject:     natsSubject,
}
// define the service discovery
serviceDiscovery := discovery.NewServiceDiscovery(disco, discoOptions)
// start the cluster
```

## Contribution

Contributions are welcome!
The project adheres to [Semantic Versioning](https://semver.org)
and [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).
This repo uses [Earthly](https://earthly.dev/get-earthly).

To contribute please:

- Fork the repository
- Create a feature branch
- Submit a [pull request](https://help.github.com/articles/using-pull-requests)

### Protobuf stub Generation

To generate any proto definition add the proto definition into the [`protos`](./protos) folder and run the following command:

```bash
earthly +protogen
```
