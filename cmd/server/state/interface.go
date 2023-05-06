package state

import (
	"database/sql"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type State interface {
	MQTT() mqtt.Client
	DB() *sql.DB
}
