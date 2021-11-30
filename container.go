package main

import (
	"github.com/docker/docker/api/types"
)

func (d *DockerClient) GetContainers() (containers []types.Container, err error) {
	containers, err = d.client.ContainerList(ctx, types.ContainerListOptions{})
	return
}
