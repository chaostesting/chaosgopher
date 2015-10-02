package config

import "os"

type Config struct {
	APIEndpoint    string
	ETCDEndpoint   string
	DockerEndpoint string
}

var config *Config

var defaults = &Config{
	APIEndpoint:    "http://172.17.42.1:8080",
	ETCDEndpoint:   "http://172.17.42.1:2379",
	DockerEndpoint: "unix:///var/run/docker.sock",
}

func Get() *Config {
	if config != nil {
		return config
	}

	config = &Config{}
	*config = *defaults

	apiEndpoint := os.Getenv("API_ENDPOINT")
	if apiEndpoint != "" {
		config.APIEndpoint = apiEndpoint
	}

	etcdEndpoint := os.Getenv("ETCD_ENDPOINT")
	if etcdEndpoint != "" {
		config.ETCDEndpoint = etcdEndpoint
	}

	dockerHost := os.Getenv("DOCKER_HOST")
	if dockerHost != "" {
		config.DockerEndpoint = dockerHost
	}

	return config
}
