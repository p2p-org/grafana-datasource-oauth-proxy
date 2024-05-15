package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type ForbidViewerProxy struct {
	proxy httputil.ReverseProxy
}

func (m ForbidViewerProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-ID-Token")
	if token == "" {
		http.Error(w, "403 - JWT Token not found", http.StatusForbidden)
	} else {
		email, err := GetEmailFromGoogleJWT(token)
		if err != nil {
			log.Printf("Error getting email from Google JWT: %s\n", err)
			http.Error(w, "403 - Wrong JWT Token", http.StatusForbidden)
		} else {
			orgID := r.Header.Get("X-Grafana-Org-ID")
			if isViewer(email, orgID) {
				log.Printf("Viewer %s in orgId %s not allowed to use datasource\n", email, orgID)
				http.Error(w, "403 - Viewer not allowed to use datasource", http.StatusForbidden)
			} else {
				m.proxy.ServeHTTP(w, r)
			}
		}
	}
}

func NewForbidViewerProxy() ForbidViewerProxy {
	target, err := url.Parse(os.Getenv("PROXY_ORIGIN_SERVER"))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("forwarding to -> %s\n", target)

	proxy := httputil.NewSingleHostReverseProxy(target)

	d := proxy.Director
	proxy.Director = func(r *http.Request) {
		d(r) // call default director

		r.Host = target.Host // set Host header as expected by target
	}
	return ForbidViewerProxy{*proxy}
}

func main() {
	log.Fatal(http.ListenAndServe(":8989", NewForbidViewerProxy()))
}
