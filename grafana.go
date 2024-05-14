package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

type GrafanaOrgUser struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func isViewer(email string, orgID string) bool {
	auth := os.Getenv("GRAFANA_AUTH")
	client := &http.Client{}
	grafanaURL, err := url.Parse(os.Getenv("GRAFANA_BASE_URL"))
	if err != nil {
		log.Printf("Error parsing GRAFANA_BASE_URL: %s\n", err)
		return true
	}
	grafanaURL.Path = path.Join(grafanaURL.Path, "/api/org/users")
	req, err := http.NewRequest("GET", grafanaURL.String(), nil)
	if err != nil {
		log.Printf("Error creating request to Grafana API: %v\n", err)
		return true
	}
	req.Header = http.Header{
		"Authorization":    {auth},
		"X-Grafana-Org-Id": {orgID},
	}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error executing request to Grafana API: %v\n", err)
		return true
	}
	var grafanaOrgUsers []GrafanaOrgUser
	err = json.NewDecoder(res.Body).Decode(&grafanaOrgUsers)
	if err != nil {
		log.Printf("Error decoding response from Grafana API: %v\n", err)
		return true
	}
	for _, grafanaOrgUser := range grafanaOrgUsers {
		if grafanaOrgUser.Email == email {
			if grafanaOrgUser.Role == "Viewer" {
				return true
			} else {
				return false
			}
		}
	}
	return true
}
