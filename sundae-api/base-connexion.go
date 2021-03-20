package sundaeapi

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"

	"github.com/sundae-party/api-server/pkg/apis/core/types"
	ssl_utils "github.com/sundae-party/pki/utils"
)

func Connect(ctx context.Context, addr string, port int16, certPath string, keyPath string, caCertPath string) (err error, conn *grpc.ClientConn, connReady chan bool) {

	var uri string

	// Build base uri for unix or TCP socket
	if strings.HasPrefix(addr, "unix:////") {
		uri = addr
	} else {
		uri = fmt.Sprintf("%s:%d", addr, port)
	}

	// Load gRPC client conf for mTLS auth and CA trust
	cliTlsConf, err := ssl_utils.LoadKeyPair(certPath, keyPath, caCertPath)
	if err != nil {
		return err, nil, nil
	}

	conn, err = grpc.DialContext(ctx, uri, grpc.WithTransportCredentials(cliTlsConf))
	if err != nil {
		return err, nil, nil
	}

	// Check gRPC connection state
	// If state is not READY send false to ConnReady chan
	// If state is READY send true to ConnReady chan
	// That permit to stop watch event loop when gRPC connexion is close.
	// If gRPC is not READY waiting loop is started.
	connReady = make(chan bool)
	go func(c chan bool) {
		for {
			time.Sleep(1 * time.Second)
			if conn.GetState().String() == "READY" {
				// TODO: debug mode only
				log.Println(conn.GetState().String())
				c <- true
			}
			c <- false
		}
	}(connReady)

	return nil, conn, connReady
}

func InitIntegrationHandler(ctx context.Context, conn *grpc.ClientConn, connReady chan bool) chan *types.Integration {
	ih := types.NewIntegrationHandlerClient(conn)
	return watchIntegrationEvent(ctx, ih, connReady)
}
