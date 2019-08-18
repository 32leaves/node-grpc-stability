//go:generate sh -c "mkdir hello || echo"
//go:generate protoc -I ../protocol --go_out=plugins=grpc:./hello hello.proto

package main

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/32leaves/node-grpc-stability/hello"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{
	StopChan chan struct{}
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(req *pb.HelloRequest, srv pb.Greeter_SayHelloServer) error {
	t := time.NewTicker(1 * time.Second)
	var i int
	for srv.Context().Err() == nil {
		msg := fmt.Sprintf("Hello %s at %d", req.Name, i)
		err := srv.Send(&pb.HelloReply{Message: msg})
		if err != nil {
			return err
		}
		fmt.Println(msg)
		i++

		if i == 5 {
			s.StopChan <- struct{}{}
		}

		<-t.C
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var srv server
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &srv)

	stop := make(chan struct{})
	srv.StopChan = stop
	go func() {
		<-stop
		s.Stop()
	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
