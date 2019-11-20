package version

// Version is updated automatically as part of the build process, and is the
// ground source of truth for the current process's build version.
//
// DO NOT EDIT
var Version = "stable-2.6.0"

// ProxyInitVersion is the pinned version of the proxy-init, from
// https://github.com/linkerd/linkerd2-proxy-init
// This has to be kept in sync with the constraint version for
// github.com/linkerd/linkerd2-proxy-init in /go.mod
var ProxyInitVersion = "v1.2.0"
