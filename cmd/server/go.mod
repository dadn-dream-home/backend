module github.com/dadn-dream-home/x/server

go 1.20

replace github.com/dadn-dream-home/x/protobuf => ../../protobuf

require (
	github.com/dadn-dream-home/x/protobuf v1.0.0
	github.com/eclipse/paho.mqtt.golang v1.4.2
	github.com/golang-migrate/migrate/v4 v4.15.2
	github.com/google/uuid v1.3.0
	github.com/mattn/go-sqlite3 v1.14.16
	google.golang.org/grpc v1.54.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
