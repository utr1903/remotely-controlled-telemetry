package main

import "github.com/utr1903/remotely-controlled-telemetry/apps/client/controller"

func main() {

	ctrl := controller.New("ws://localhost:8081/ws")
	ctrl.Run()
}
