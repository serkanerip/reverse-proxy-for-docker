package main

var ingressess = map[string]Ingress{}

func (d *DockerClient) GenareteSitesFromContainers() {
	ingressess = make(map[string]Ingress)
	containers, _ := d.GetContainers()

	for _, container := range containers {
		if match, exists := container.Labels["erip-proxy.match"]; exists {
			backends := ingressess[match].Backends
			backends = append(backends, Backend{
				Ip:   container.NetworkSettings.Networks["go-docker_default"].IPAddress,
				Port: container.Labels["erip-proxy.port"],
			})
			ingressess[match] = Ingress{
				Path:     match,
				Host:     container.Labels["erip-proxy.host"],
				Backends: backends,
			}
		}
	}
}
