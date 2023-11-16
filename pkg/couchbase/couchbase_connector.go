package couchbase

import (
	"be-bootcamp-2023/pkg/config"
	"github.com/couchbase/gocb/v2"
	"strings"
	"time"
)

type Cluster = gocb.Cluster
type Bucket = gocb.Bucket
type Collection = gocb.Collection

func ConnectCluster(host, user, password string) (*Cluster, error) {
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: user,
			Password: password,
		},
		TimeoutsConfig: gocb.TimeoutsConfig{
			ViewTimeout:       500 * time.Millisecond,
			QueryTimeout:      500 * time.Millisecond,
			AnalyticsTimeout:  500 * time.Millisecond,
			SearchTimeout:     500 * time.Millisecond,
			ManagementTimeout: 500 * time.Millisecond,
			KVTimeout:         500 * time.Millisecond,
			KVDurableTimeout:  500 * time.Millisecond,
		},
		InternalConfig: gocb.InternalConfig{
			ConnectionBufferSize: 1024 * 1024,
		},
		CircuitBreakerConfig: gocb.CircuitBreakerConfig{Disabled: true},
	}
	h := strings.ReplaceAll(host, " ", "")
	cluster, err := gocb.Connect(h, opts)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func ConnectClusterWithConfig(config *config.CouchbaseConfig) (*Cluster, error) {
	return ConnectCluster(config.Hosts, config.Username, config.Password)
}
