module github.com/airbloc/flux-slack-alert

go 1.13

require (
	github.com/airbloc/logger v1.1.3
	github.com/docker/distribution v0.0.0-20180611183926-749f6afb4572 // indirect
	github.com/fluxcd/flux v1.15.0
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/pkg/errors v0.8.1
	github.com/slack-go/slack v0.6.5
	github.com/ugorji/go v1.1.7 // indirect
	golang.org/x/sys v0.0.0-20191010194322-b09406accb47 // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace github.com/docker/distribution => github.com/2opremio/distribution v0.0.0-20190419185413-6c9727e5e5de
