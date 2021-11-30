package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
)

type Backend struct {
	Ip   string
	Port string
}

type Ingress struct {
	Path                 string
	Host                 string
	Backends             []Backend
	LastUsedBackendIndex uint
}

type HttpServer struct {
	r *chi.Mux
}

func NewHttpServer() *HttpServer {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	return &HttpServer{
		r: r,
	}
}

func (h *HttpServer) StartServer() {
	http.ListenAndServe(":80", h.r)
}

func (h *HttpServer) RegisterIngress(i Ingress) {
	h.r.Get(i.Path, func(rw http.ResponseWriter, r *http.Request) {
		backend := i.Backends[i.LastUsedBackendIndex%uint(len(i.Backends))]
		remoteURL := fmt.Sprintf("http://%s:%s", backend.Ip, backend.Port)

		remote, err := url.Parse(remoteURL)
		if err != nil {
			fmt.Printf("cannot parse url: %s  err is: %v", remoteURL, err)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Transport = &http.Transport{
			ResponseHeaderTimeout: 5 * time.Second,
		}
		r.Header.Add("proxy", "erip-proxy")
		r.Host = i.Host
		proxy.ServeHTTP(rw, r)
		i.LastUsedBackendIndex += 1
	})

	backendsStr := ""
	for _, backend := range i.Backends {
		backendsStr += fmt.Sprintf("%s:%s, ", backend.Ip, backend.Port)
	}
	fmt.Printf("Proxy requests to '%s' -> %s\n", i.Path, backendsStr)
}
