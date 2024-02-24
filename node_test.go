package groupcache

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/stretchr/testify/require"
	"github.com/tochemey/groupcache/v2/discovery"
	"github.com/tochemey/groupcache/v2/discovery/nats"
	"github.com/tochemey/groupcache/v2/example"
	"github.com/travisjeffery/go-dynaport"
)

func TestNode(t *testing.T) {
	t.Run("With Single Node", func(t *testing.T) {
		t.Skip()
		server := startNatsServer(t)
		serverAddr := server.Addr().String()
		// start node1
		node := startNode(t, "node", serverAddr)
		require.NotNil(t, node)

		group := NewGroup("users", 3000000, GetterFunc(
			func(ctx context.Context, id string, dest Sink) error {
				user := &example.User{
					Id:      id,
					Name:    "test",
					Age:     20,
					IsSuper: false,
				}
				// Set the user in the groupcache to expire after 5 minutes
				return dest.SetProto(user, time.Now().Add(time.Minute*5))
			},
		))

		user := new(example.User)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
		defer cancel()

		if err := group.Get(ctx, "12345", ProtoSink(user)); err != nil {
			t.Fail()
		}

		t.Cleanup(func() {
			require.NoError(t, node.Stop(context.TODO()))
			server.Shutdown()
			httpPoolMade = false
		})
	})
}

func startNode(t *testing.T, nodeName, serverAddr string) *Node {
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
	config := discovery.Options{
		nats.ApplicationName: applicationName,
		nats.Server:          serverAddr,
		nats.Subject:         natsSubject,
	}

	// create the startClusterNode
	node, err := NewNode(ctx, discovery.New(provider, config))
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
