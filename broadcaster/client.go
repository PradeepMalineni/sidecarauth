package broadcast

import (
	"context"
	"log"
	"sidecarauth/broadcaster/pb/broadcast"
	oauth_token "sidecarauth/broadcaster/pb/oauth_token"
	"time"

	"google.golang.org/grpc"
)

func startClient(ip string, serverIP string) {
	for {
		conn, err := grpc.Dial(serverIP, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Printf("Did not connect: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		defer conn.Close()
		c := broadcast.NewBroadcastServiceClient(conn)

		// Example: Call to BroadcastToken
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			req := &oauth_token.OAuthTokenResponse{
				TokenType:   "Bearer",
				AccessToken: "example-access-token",
				IssuedAt:    time.Now().Unix(),
				ExpiresIn:   3600,
				Scope:       "read write",
			}
			stream, err := c.BroadcastToken(ctx, req)
			if err != nil {
				log.Fatalf("could not broadcast token: %v", err)
			}
			for {
				res, err := stream.Recv()
				if err != nil {
					log.Fatalf("Failed to receive broadcast response: %v", err)
				}
				log.Printf("Broadcast response: %v", res.Message)
			}
		}()
		time.Sleep(5 * time.Second)
	}
}
