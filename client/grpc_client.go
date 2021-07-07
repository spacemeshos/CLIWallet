package client

import (
	"context"
	"time"

	"github.com/fullstorydev/grpcurl"
	"google.golang.org/grpc/credentials"

	apitypes "github.com/spacemeshos/api/release/go/spacemesh/v1"
	"google.golang.org/grpc"
)

const DefaultGRPCServer = "localhost:9092"
const DefaultSecureConnection = false

type gRPCClient struct {
	connection               *grpc.ClientConn
	server                   string
	secureConnection         bool
	nodeServiceClient        apitypes.NodeServiceClient
	debugServiceClient       apitypes.DebugServiceClient
	meshServiceClient        apitypes.MeshServiceClient
	globalStateServiceClient apitypes.GlobalStateServiceClient
	transactionServiceClient apitypes.TransactionServiceClient
	smesherServiceClient     apitypes.SmesherServiceClient
}

func newGRPCClient(server string, secureConnection bool) *gRPCClient {
	return &gRPCClient{
		nil,
		server,
		secureConnection,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	}
}

func (c *gRPCClient) Connect() error {
	if c.connection != nil {
		_ = c.connection.Close()
	}

	var conn *grpc.ClientConn
	var err error
	if !c.secureConnection {
		// simple grpc dial
		conn, err = grpc.Dial(c.server, grpc.WithInsecure())
	} else {
		// secure connection without client cert or server cert validation
		conn, err = c.dial(c.server)
	}

	if err != nil {
		return err
	}
	c.connection = conn
	return nil
}

// dial is a secure tls dialing helper
func (c *gRPCClient) dial(address string) (*grpc.ClientConn, error) {
	dialTime := 60 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), dialTime)
	defer cancel()

	var creds credentials.TransportCredentials
	var err error
	creds, err = grpcurl.ClientTransportCredentials(false, "", "", "")
	if err != nil {
		return nil, err
	}

	// todo: set release version in user agent
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithUserAgent("sm-cli-wallet/dev-build"))

	cc, err := grpcurl.BlockingDial(ctx, "tcp", address, creds, opts...)
	if err != nil {
		return nil, err
	}

	return cc, nil
}

func (c *gRPCClient) Close() error {
	if c.connection != nil {
		return c.connection.Close()
	}
	return nil
}

// ServerInfo returns basic api server and connection info
func (c *gRPCClient) ServerInfo() string {
	s := c.server + " (GRPC API 1.1)"
	if c.secureConnection {
		s += ". Secure Connection."
	} else {
		s += ". >> Insecure Connection. Use only with a local trusted server <<"
	}

	return s
}

//// services clients

func (c *gRPCClient) getNodeServiceClient() apitypes.NodeServiceClient {
	if c.nodeServiceClient == nil {
		c.nodeServiceClient = apitypes.NewNodeServiceClient(c.connection)
	}
	return c.nodeServiceClient
}

func (c *gRPCClient) getDebugServiceClient() apitypes.DebugServiceClient {
	if c.debugServiceClient == nil {
		c.debugServiceClient = apitypes.NewDebugServiceClient(c.connection)
	}
	return c.debugServiceClient
}

func (c *gRPCClient) getMeshServiceClient() apitypes.MeshServiceClient {
	if c.meshServiceClient == nil {
		c.meshServiceClient = apitypes.NewMeshServiceClient(c.connection)
	}
	return c.meshServiceClient
}

func (c *gRPCClient) getGlobalStateServiceClient() apitypes.GlobalStateServiceClient {
	if c.globalStateServiceClient == nil {
		c.globalStateServiceClient = apitypes.NewGlobalStateServiceClient(c.connection)
	}
	return c.globalStateServiceClient

}

func (c *gRPCClient) getTransactionServiceClient() apitypes.TransactionServiceClient {
	if c.transactionServiceClient == nil {
		c.transactionServiceClient = apitypes.NewTransactionServiceClient(c.connection)
	}
	return c.transactionServiceClient
}

func (c *gRPCClient) getSmesherServiceClient() apitypes.SmesherServiceClient {
	if c.smesherServiceClient == nil {
		c.smesherServiceClient = apitypes.NewSmesherServiceClient(c.connection)
	}
	return c.smesherServiceClient
}
