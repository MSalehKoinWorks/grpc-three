package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/MSalehKoinWorks/grpc-three/proto"
)

type mockDepositServer struct {
	pb.UnimplementedDepositServiceServer
}

func (*mockDepositServer) Deposit(ctx context.Context, req *pb.DepositRequest) (*pb.DepositResponse, error) {
	if req.GetAmount() < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "cannot deposit %v", req.GetAmount())
	}

	return &pb.DepositResponse{Ok: true}, nil
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pb.RegisterDepositServiceServer(server, &mockDepositServer{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestDepositClient_Deposit(t *testing.T) {
	tests := []struct {
		name   string
		amount float32
		res    bool
		err    error
	}{
		{
			"invalid request with negative amount",
			-1.11,
			false,
			fmt.Errorf("grpc: InvalidArgument, cannot deposit %v", -1.11),
		},
		{
			"valid request with non negative amount",
			0.00,
			true,
			nil,
		},
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := NewDepositClient(conn, time.Second).Deposit(context.Background(), tt.amount)

			if response != tt.res {
				t.Error("error: expected", tt.res, "received", response)
			}

			if err != nil && errors.Is(err, tt.err) {
				t.Error("error: expected", tt.err, "received", err)
			}
		})
	}
}
