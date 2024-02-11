package groupcache

import (
	"context"
	"fmt"

	"log"
	"time"
)

func ExampleUsage() {
	/*
			// Make sure each node running the groupcache has the env vars properly set:
			// GROUP_PORT, NODE_NAME and NODE_IP
			// kindly refer to the readme notes.


			// Create an instance of the discovery server.
			For instance let us say we are using kubernetes
			provider := kubernetes.New()

			// Create the discovery options
			For kubernetes we only need the namespace and the application name
			namespace := "default"
			application := "users"

			options := discovery.Config{
		    kubernetes.ApplicationName: application,
		    kubernetes.Namespace:       namespace,
			}

			// Create an instance of the service discovery
			serviceDiscovery := discovery.NewServiceDiscovery(provider, options)

			// Create an instance of the cluster
			ctx := context.Background()
			cluster := cluster.New(ctx, serviceDiscovery)

			// Start the cluster
			err := cluster.Start(ctx)

			// Stop the cluster
			defer cluster.Stop(ctx)
	*/

	// Create a new group cache with a max cache size of 3MB
	group := NewGroup("users", 3000000, GetterFunc(
		func(ctx context.Context, id string, dest Sink) error {

			// In a real scenario we might fetch the value from a database.
			/*if user, err := fetchUserFromMongo(ctx, id); err != nil {
				return err
			}*/

			user := User{
				Id:      "12345",
				Name:    "John Doe",
				Age:     40,
				IsSuper: true,
			}

			// Set the user in the groupcache to expire after 5 minutes
			if err := dest.SetProto(&user, time.Now().Add(time.Minute*5)); err != nil {
				return err
			}
			return nil
		},
	))

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := group.Get(ctx, "12345", ProtoSink(&user)); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("-- User --\n")
	fmt.Printf("Id: %s\n", user.Id)
	fmt.Printf("Name: %s\n", user.Name)
	fmt.Printf("Age: %d\n", user.Age)
	fmt.Printf("IsSuper: %t\n", user.IsSuper)

	/*
		// Remove the key from the groupcache
		if err := group.Remove(ctx, "12345"); err != nil {
			fmt.Printf("Remove Err: %s\n", err)
			log.Fatal(err)
		}
	*/

	// Output: -- User --
	// Id: 12345
	// Name: John Doe
	// Age: 40
	// IsSuper: true
}
