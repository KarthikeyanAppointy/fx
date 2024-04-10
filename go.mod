module go.uber.org/fx

go 1.20

require (
	github.com/stretchr/testify v1.8.1
	go.uber.org/dig v1.17.1
	go.uber.org/goleak v1.3.0
	go.uber.org/multierr v1.10.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	go.saastack.io/log-error => /Users/karthikarthi/Desktop/saastack-Backend-log-error
	go.uber.org/dig => github.com/paullen/dig v1.7.1-0.20200612122854-3efded4be010
	go.uber.org/fx => github.com/appointy/fx v1.9.1-0.20190624110333-490d04d33ef6
)
