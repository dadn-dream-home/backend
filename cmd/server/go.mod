module github.com/dadn-dream-home/x/server

go 1.20

replace github.com/dadn-dream-home/x/protobuf => ../../protobuf

require (
	github.com/dadn-dream-home/x/protobuf v1.0.0
	github.com/eclipse/paho.mqtt.golang v1.4.2
	github.com/golang-migrate/migrate/v4 v4.15.2
	github.com/google/uuid v1.3.0
	github.com/hashicorp/hcl/v2 v2.16.2
	github.com/mattn/go-sqlite3 v1.14.16
	go.uber.org/zap v1.24.0
	golang.org/x/sync v0.2.0
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1
	google.golang.org/grpc v1.54.0
	google.golang.org/protobuf v1.30.0
	github.com/robfig/cron/v3 v3.0.1
)

require (
	github.com/agext/levenshtein v1.2.1 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/mitchellh/go-wordwrap v0.0.0-20150314170334-ad45545899c7 // indirect
	github.com/stretchr/testify v1.8.2 // indirect
	github.com/zclconf/go-cty v1.12.1 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
)
