module github.com/zdnscloud/zke

go 1.13

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.18+incompatible
	github.com/coreos/go-semver v0.3.0
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v1.4.2-0.20180612054059-a9fbbdc8dd87
	github.com/docker/go-connections v0.4.0
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.12.2 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.8.1
	github.com/pkg/sftp v1.10.1
	github.com/prometheus/client_golang v1.4.0 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200122045848-3419fae592fc // indirect
	github.com/urfave/cli v1.20.0
	github.com/zdnscloud/cement v0.0.0-20191114151602-e70d68aebdad
	github.com/zdnscloud/gok8s v0.0.0-20200212071629-b06587f54ee6
	go.uber.org/zap v1.13.0 // indirect
	golang.org/x/crypto v0.0.0-20191119213627-4f8c1d86b1ba
	golang.org/x/net v0.0.0-20191119073136-fc4aabc6c914 // indirect
	google.golang.org/grpc v1.27.0 // indirect
	gopkg.in/yaml.v2 v2.2.7
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
)

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
