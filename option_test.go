package groupcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tochemey/groupcache/v2/log"
)

func TestOptions(t *testing.T) {
	testCases := []struct {
		name     string
		option   Option
		expected Node
	}{
		{
			name:     "WithReplicaCount",
			option:   WithReplicaCount(2),
			expected: Node{replicaCount: 2},
		},
		{
			name:     "WithLogger",
			option:   WithLogger(log.DefaultLogger),
			expected: Node{logger: log.DefaultLogger},
		},

		{
			name:     "WithShutdownTimeout",
			option:   WithShutdownTimeout(2 * time.Minute),
			expected: Node{shutdownTimeout: 2 * time.Minute},
		},
		{
			name:     "WithBasePath",
			option:   WithBasePath("/path"),
			expected: Node{basePath: "/path"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cl Node
			tc.option.Apply(&cl)
			assert.Equal(t, tc.expected, cl)
		})
	}
}
