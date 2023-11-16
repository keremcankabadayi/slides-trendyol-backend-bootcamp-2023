package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

type Application struct {
	Couchbase CouchbaseConfig
}

type CouchbaseConfig struct {
	Hosts         string
	Username      string
	Password      string
	ConnectionURI string
	BucketName    string
	Timeout       time.Duration
}

func New() (*Application, error) {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "stage"
	}

	cfg := Application{}

	v := viper.New()
	v.SetConfigFile(fmt.Sprintf("resources/config.yml"))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	sub := v.Sub(env)

	if err := envSubst(sub); err != nil {
		return nil, err
	}

	if err := sub.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func envSubst(sub *viper.Viper) error {
	for _, k := range sub.AllKeys() {
		value := sub.GetString(k)

		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			key := strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")
			if key == "" {
				return errors.New("dynamic key not found")
			}
			sub.Set(k, os.Getenv(key))
		}
	}

	return nil
}
