package client

import (
	"errors"
	pb "github.com/spacemeshos/api/release/go/spacemesh/v1"
	"golang.org/x/net/context"
)

func (c *GRPCClient) Sanity() error {
	service := c.nodeServiceClient()

	const msg = "hello spacemesh"

	resp, err := service.Echo(context.Background(), &pb.EchoRequest{
		Msg: &pb.SimpleString{Value: msg}})

	if err != nil {
		return err
	}

	if resp.Msg.Value != msg {
		return errors.New("unexpected node service echo response")
	}

	return nil
}

func (c *GRPCClient) NodeInfo() (*NodeInfo, error) {
	return nil, nil
}