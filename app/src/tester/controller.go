package main

import "os"
import "app/config"
import "github.com/fsouza/go-dockerclient"

const DockerStopTimeout = 3

type Controller struct {
	docker    *docker.Client
	container string
	image     string
}

func NewController(container, image string) (*Controller, error) {
	endpoint := os.Getenv("DOCKER_HOST")
	if endpoint == "" {
		endpoint = config.Get().DockerEndpoint
	}

	var client *docker.Client
	var e error

	certpath := os.Getenv("DOCKER_CERT_PATH")
	if certpath == "" {
		client, e = docker.NewClient(endpoint)
		if e != nil {
			return nil, e
		}

	} else {
		client, e = docker.NewTLSClient(endpoint,
			certpath+"/cert.pem",
			certpath+"/key.pem",
			certpath+"/ca.pem")

		if e != nil {
			return nil, e
		}
	}

	return &Controller{client, container, image}, nil
}

func (c *Controller) Start() error {
	return c.docker.StartContainer(c.container, nil)
}

func (c *Controller) Started() (bool, error) {
	container, e := c.docker.InspectContainer(c.container)
	if e != nil {
		return false, e
	}

	return container.State.Running, nil
}

func (c *Controller) Stop() error {
	return c.docker.StopContainer(c.container, DockerStopTimeout)
}

func (c *Controller) Create() error {
	_, e := c.docker.CreateContainer(docker.CreateContainerOptions{
		Name:   c.container,
		Config: &docker.Config{Image: c.image},
	})
	return e
}

func (c *Controller) Destroy() error {
	return c.docker.RemoveContainer(docker.RemoveContainerOptions{
		ID:            c.container,
		RemoveVolumes: true,
		Force:         true,
	})
}
