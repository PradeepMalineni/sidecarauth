package broadcast

import (
	"log"
	"net"
	"sidecarauth/broadcaster/pb/broadcast"

	oauth_token "sidecarauth/broadcaster/pb/oauth_token"
	logger "sidecarauth/utility"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	broadcast.UnimplementedBroadcastServiceServer
}

func (s *server) BroadcastToken(req *oauth_token.OAuthTokenResponse, stream broadcast.BroadcastService_BroadcastTokenServer) error {
	log.Printf("Received token: %s", req.AccessToken)
	for {
		if err := stream.Send(&broadcast.BroadcastResponse{Message: "Token received"}); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
}

func startServer(ip string) *grpc.Server {
	logger.Log("Starting the GRPC Server")
	lis, err := net.Listen("tcp", ip)
	if err != nil {
		logger.LogF("Starting the GRPC Server", err)

		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	broadcast.RegisterBroadcastServiceServer(s, &server{})
	log.Printf("Server listening on %s", ip)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	return s
}
