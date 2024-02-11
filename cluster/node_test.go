/*
 * Copyright 2024 Tochemey
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cluster

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/tochemey/groupcache/v2"

	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/stretchr/testify/require"
	"github.com/tochemey/groupcache/v2/discovery"
	"github.com/tochemey/groupcache/v2/discovery/nats"
	"github.com/travisjeffery/go-dynaport"
)

func TestCluster(t *testing.T) {
	t.Run("With Single Node", func(t *testing.T) {
		server := startNatsServer(t)
		serverAddr := server.Addr().String()
		// start node1
		node := startClusterNode(t, "node", serverAddr)
		require.NotNil(t, node)

		group := groupcache.NewGroup("users", 3000000, groupcache.GetterFunc(
			func(ctx context.Context, id string, dest groupcache.Sink) error {
				user := &groupcache.User{
					Id:      id,
					Name:    "test",
					Age:     20,
					IsSuper: false,
				}
				// Set the user in the groupcache to expire after 5 minutes
				return dest.SetProto(user, time.Now().Add(time.Minute*5))
			},
		))

		user := new(groupcache.User)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
		defer cancel()

		if err := group.Get(ctx, "12345", groupcache.ProtoSink(user)); err != nil {
			t.Fail()
		}

		t.Cleanup(func() {
			require.NoError(t, node.Stop(context.TODO()))
			server.Shutdown()
		})
	})
}

func startNatsServer(t *testing.T) *natsserver.Server {
	t.Helper()
	serv, err := natsserver.NewServer(&natsserver.Options{
		Host: "127.0.0.1",
		Port: -1,
	})

	require.NoError(t, err)

	ready := make(chan bool)
	go func() {
		ready <- true
		serv.Start()
	}()
	<-ready

	if !serv.ReadyForConnections(2 * time.Second) {
		t.Fatalf("nats-io server failed to start")
	}

	return serv
}

func startClusterNode(t *testing.T, nodeName, serverAddr string) *Node {
	// create a context
	ctx := context.TODO()

	// generate the ports for the single startClusterNode
	ports := dynaport.Get(3)
	port := ports[0]

	// create a Cluster startClusterNode
	host := "127.0.0.1"
	// set the environments
	require.NoError(t, os.Setenv("GROUP_PORT", strconv.Itoa(port)))
	require.NoError(t, os.Setenv("NODE_NAME", nodeName))
	require.NoError(t, os.Setenv("NODE_IP", host))

	// create the various config option
	applicationName := "accounts"
	natsSubject := "some-subject"
	// create the instance of provider
	provider := nats.NewDiscovery()

	// create the config
	config := discovery.Config{
		nats.ApplicationName: applicationName,
		nats.Server:          serverAddr,
		nats.Subject:         natsSubject,
	}

	// create the startClusterNode
	node, err := New(ctx, discovery.NewServiceDiscovery(provider, config))
	require.NoError(t, err)
	require.NotNil(t, node)

	// start the node
	require.NoError(t, node.Start(ctx))

	// clear the env var
	require.NoError(t, os.Unsetenv("GROUP_PORT"))
	require.NoError(t, os.Unsetenv("NODE_NAME"))
	require.NoError(t, os.Unsetenv("NODE_IP"))

	// return the cluster startClusterNode
	return node
}
