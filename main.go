package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool)

	var m sync.Mutex

	dockerClient := NewDockerClient()

	changed := dockerClient.WatchChanges()

	httpServer := NewHttpServer()

	httpServer.r.Get("/admin/ingressses", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		b, _ := json.Marshal(&ingressess)
		rw.Write(b)
	})

	go func() {
		dockerClient.GenareteSitesFromContainers()
		for _, ingress := range ingressess {
			httpServer.RegisterIngress(ingress)
		}
		httpServer.StartServer()
	}()

	go func() {
		for {
			<-changed
			fmt.Println("New changes detected")
			m.Lock()
			dockerClient.GenareteSitesFromContainers()
			for _, ingress := range ingressess {
				httpServer.RegisterIngress(ingress)
			}
			m.Unlock()
		}
	}()

	go func() {
		sig := <-sigs
		fmt.Printf("new system signal: %v", sig)
		done <- true
	}()

	<-done
}
