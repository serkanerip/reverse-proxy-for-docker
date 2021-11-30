package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
)

func (d *DockerClient) WatchChanges() chan bool {
	eventStream, errStream := d.client.Events(ctx, types.EventsOptions{})

	changed := make(chan bool)
	go func() {
		for {
			select {
			case event := <-eventStream:
				fmt.Printf("[%s] Type: %s ID: %s\n", strings.ToUpper(event.Action), event.Type, event.ID)
				if event.Type != events.ContainerEventType {
					continue
				}
				changed <- true
				// switch event.Action {
				// case "create":
				// 	container, _ := d.client.ContainerInspect(ctx, event.ID)
				// 	if match, exists := container.Config.Labels["erip-proxy.match"]; exists {
				// 		newSite <- Site{
				// 			Path:        match,
				// 			ServiceName: container.Config.Labels["com.docker.compose.service"],
				// 			ServicePort: container.Config.Labels["erip-proxy.port"],
				// 		}
				// 	}
				// }
			case err := <-errStream:
				fmt.Printf("an error has occured err is: %v\n", err)
			}
			time.Sleep(time.Millisecond * 300) // backoff
		}
	}()

	return changed
}
