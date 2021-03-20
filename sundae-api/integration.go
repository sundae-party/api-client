package sundaeapi

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/sundae-party/api-server/pkg/apis/core/types"
)

func watchIntegrationEvent(ctx context.Context, ih types.IntegrationHandlerClient, connReady chan bool) (integrationEvent chan *types.Integration) {
	go func() {
		for {
			select {
			// If gRPC connection is ok, open event stream
			case ok := <-connReady:
				if ok {
					// TODO: Create new stream
					int := &types.Integration{
						Name: "test",
					}
					log.Println("Subscribe to the stream integration event")
					stream, err := ih.SubscribeEvents(ctx, int)
					if err != nil {
						log.Printf("Integration stream event error : %s\n", err)
						break
					}
					log.Println("Connected to the apiserver")

					// Watch for event from stream, if error, break this loop to create new stream.
					for {
						res, err := stream.Recv()
						if err == io.EOF {
							// End of stream
							// This Recv method returns (nil, io.EOF) once the server-to-client stream has been completely read through.
							log.Println(err)
							break
						}
						if err != nil {
							// Detect connexion close (server stop, network error)
							// stop event loop and wait for gRPC connexion ok
							// EX err when server stop:
							// -- cannot receive response: rpc error: code = Unavailable desc = transport is closing --
							log.Printf("cannot receive response: %s\n", err)
							break
						}
						integrationEvent <- res
					}
				}
			default:
				// Waitting for gRPC connexion ok
				time.Sleep(time.Second * 1)
			}
		}
	}()
	return integrationEvent
}
