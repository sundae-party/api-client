package main

import (
	"context"
	"log"

	sundaeapi "github.com/sundae-party/sundae-grpc-client/sundae-api"
)

func main() {
	ctx := context.Background()
	err, conn, ready := sundaeapi.Connect(ctx, "192.168.1.61", 8443, "ssl/integration01.pem", "ssl/integration01.key", "ssl/ca.pem")
	if err != nil {
		log.Fatalln(err)
	}

	event := sundaeapi.InitIntegrationHandler(ctx, conn, ready)

	for {
		select {
		case integrationEvent := <-event:
			log.Printf("%s\n", integrationEvent.Name)

		}
	}
}
