FROM golang:1.22.3-bookworm as build

WORKDIR /go/src/grafana-datasource-oauth-proxy
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -v -o /go/bin/grafana-datasource-oauth-proxy ./...

FROM gcr.io/distroless/static-debian12:latest

COPY --from=build /go/bin/grafana-datasource-oauth-proxy /

CMD ["/grafana-datasource-oauth-proxy"]
