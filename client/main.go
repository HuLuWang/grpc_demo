package main

import (
	"context"
	"google.golang.org/grpc"
	pb "grpc_demo/proto"
	"io"
	"log"
)

func main() {
	conn, _ := grpc.Dial(":"+"4001", grpc.WithInsecure())
	defer conn.Close()
	client := pb.NewGreeterClient(conn)
	// unary rpc
	err := SayHello(client)
	if err != nil {
		log.Fatalf("SayHello err: %v", err)
	}
	// server-side streaming rpc
	err = SayList(client, &pb.HelloRequest{
		Name: "server side streaming rpc",
	})
	if err != nil {
		log.Fatalf("SayList err: %v", err)
	}
	// client-side streaming rpc
	err = SayRecord(client, &pb.HelloRequest{
		Name: "client side streaming rpc",
	})
	if err != nil {
		log.Fatalf("SayRecord err: %v", err)
	}
	// bidirectional streaming rpc
	err = SayRoute(client, &pb.HelloRequest{
		Name: "bidirectional streaming rpc",
	})
	if err != nil {
		log.Fatalf("SayRoute err: %v", err)
	}
}

func SayHello(client pb.GreeterClient) error {
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{
		Name: "unary grpc",
	})
	if err != nil {
		return err
	}
	log.Printf("client.SayHello resp: %s", resp.Message)
	return nil
}

// server side streaming rpc
// stream.Recv 是阻塞等待的
func SayList(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, err := client.SayList(context.Background(), r)
	if err != nil {
		return err
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil
		}
		log.Printf("resp: %v", resp)
	}
	return nil
}

// client side streaming rpc
func SayRecord(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, err := client.SayRecord(context.Background())
	if err != nil {
		return err
	}
	for n := 0; n < 6; n++ {
		err := stream.Send(r)
		if err != nil {
			return err
		}
	}
	// 结束发送并接受
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("resp err: %v", resp)
	return nil
}

// bidirectional streaming rpc
func SayRoute(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, err := client.SayRoute(context.Background())
	if err != nil {
		return err
	}
	for n := 0; n <= 6; n++ {
		err = stream.Send(r)
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp err: %v", resp)
	}
	return nil
}
