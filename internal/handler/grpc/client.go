package grpc

import (
	"context"
	"fmt"

	"github.com/Woodfyn/auditLog/pkg/core"
	audit "github.com/Woodfyn/auditLog/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	conn        *grpc.ClientConn
	auditClient audit.AuditClient
}

func NewClient(port string) (*Client, error) {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:        conn,
		auditClient: audit.NewAuditClient(conn),
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) SendLogRequest(ctx context.Context, req core.LogItem) error {
	action, err := core.ToPbAction(req.Action)
	if err != nil {
		return err
	}

	entity, err := core.ToPbEntity(req.Entity)
	if err != nil {
		return err
	}

	_, err = c.auditClient.Log(ctx, &audit.LogRequest{
		Action:    action,
		Entity:    entity,
		EntityId:  req.EntityID,
		Timestamp: timestamppb.New(req.Timestamp),
	})

	return err
}
