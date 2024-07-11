module github.com/glethuillier/mps/client

go 1.22.4

replace github.com/glethuillier/mps/lib => ../lib

require github.com/glethuillier/mps/lib v0.0.0-00010101000000-000000000000

require (
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.2
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/stretchr/testify v1.8.1
	go.uber.org/zap v1.27.0
	google.golang.org/protobuf v1.34.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
