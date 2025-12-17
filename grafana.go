package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

type GrafanaOrgUser struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type OrgUsers struct {
	Users       []GrafanaOrgUser
	OrgId       string
	lastChecked time.Time
}
// todo: do not call grafana api on each incoming request
type UsersCache struct {
	users    map[string]OrgUsers
	usersAPI string
	auth     string
	client   *http.Client
}

func NewUserCache() *UsersCache {
	auth := os.Getenv("GRAFANA_AUTH")
	if auth == ""  {
		log.Fatalln("please provide GRAFANA_AUTH env")
	}

	grafanaURL, err := url.Parse(os.Getenv("GRAFANA_BASE_URL"))
	if err != nil {
		log.Fatalf("Error parsing GRAFANA_BASE_URL: %s", err)
	}
	grafanaURL.Path = path.Join(grafanaURL.Path, "/api/org/users")

	return &UsersCache{
		users: make(map[string]OrgUsers),
		usersAPI: grafanaURL.String(),
		auth: auth,
		client:  &http.Client{},
	}
}

func (cache UsersCache) IsViewer(email string, orgID string) bool {
	req, err := http.NewRequest("GET", cache.usersAPI, nil)
	if err != nil {
		log.Fatalf("Error creating request to Grafana API: %v\n", err)
	}
	req.Header = http.Header{
		"Authorization":    {cache.auth},
		"X-Grafana-Org-Id": {orgID},
	}
	res, err := cache.client.Do(req)
	if err != nil {
		log.Printf("Error executing request to Grafana API: %v\n", err)
		return true
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("HTTP%d: %s", res.StatusCode, string(bodyBytes))
	}

	var grafanaOrgUsers []GrafanaOrgUser
	err = json.Unmarshal(bodyBytes, &grafanaOrgUsers)
	if err != nil {
		log.Println(string(bodyBytes))
		log.Fatalf("Error decoding response from Grafana API: %v\n", err)
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
