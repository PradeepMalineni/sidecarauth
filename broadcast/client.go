package broadcast

import (
	"context"
	"log"
	"time"

	pb "path/to/proto"

	"google.golang.org/grpc"
)

func StartClient(ip string, serverIP string) {
	for {
		conn, err := grpc.Dial(serverIP, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Printf("Did not connect: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		defer conn.Close()
		c := pb.NewBroadcastServiceClient(conn)
		for {
			message := "Hello from " + ip
			stream, err := c.BroadcastMessage(context.Background(), &pb.BroadcastRequest{Message: message})
			if err != nil {
				log.Printf("Error on broadcast message: %v", err)
				break
			}
			for {
				res, err := stream.Recv()
				if err != nil {
					log.Printf("Error on receive: %v", err)
					break
				}
				log.Printf("Received from server: %s", res.Message)
			}
		}
		time.Sleep(5 * time.Second)
	}
}
