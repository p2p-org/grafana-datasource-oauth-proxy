# grafana-datasource-oauth-proxy

This project aims to restrict access to grafana datasource for users with "Viewer" role.

Use case: you have a critical datasource (f.e. logs) in same grafana organisation with common datasources
(prometheus, graphite, etc..). You would like to allow sharing dashboards to unprivileged or even
anonymous users, but prevent critical data leaking

Alternative: [Data source permissions](https://grafana.com/docs/grafana/latest/administration/data-source-management/#data-source-permissions)
available in Grafana Enterprise and Grafana Cloud.

Requirements:
- Grafana OSS server
- Google OAuth2 Provider configured
- Grafana Server Administrator credentials

Environment variables:
- `GRAFANA_BASE_URL` - grafana oss server base url (like https://grafana.example.com or https://infra.example.com/grafana)
- `GRAFANA_AUTH` - server administrator basic auth string in format `Basic $(base64($user:$password)`
- `PROXY_ORIGIN_SERVER` - grafana datasource url (like http://127.0.0.1:7682, http://loki.logging.svc.cluster.local, https://ds.example.com or any other http(s) resource)
