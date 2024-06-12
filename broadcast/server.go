package broadcast

import (
	"log"
	"net"
	"time"

	pb "path/to/proto"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedBroadcastServiceServer
}

func (s *Server) BroadcastMessage(req *pb.BroadcastRequest, stream pb.BroadcastService_BroadcastMessageServer) error {
	log.Printf("Received message: %s", req.Message)
	for {
		if err := stream.Send(&pb.BroadcastResponse{Message: req.Message}); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
}

func StartServer(ip string) *grpc.Server {
	lis, err := net.Listen("tcp", ip)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBroadcastServiceServer(s, &Server{})
	log.Printf("Server listening on %s", ip)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	return s
}
