## Developing

### Prerequisite

Install [mosquitto], [Go] and  [SQLite3].

[mosquitto]: https://mosquitto.org/
[Go]: https://golang.org
[SQLite3]: https://sqlite.org

### Running

1. Run `mosquitto` in the background

    ```bash
    $ mosquitto -v
    ```

2. Generate the protobuf files, if changed `backend.proto`:

    ```bash
    # or ./scripts/gen.bat
    $ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protobuf/backend.proto
    ```

3. Run backend:

    ```bash
    $ go run cmd/server
    ```

4. Use `mosquitto_pub` to send data to MQTT server:

    ```bash
    $ mosquitto_pub -t 'temp1' -r -m 43
    ```

5. Use a gRPC client like Postman (with server reflection) to interact with server.