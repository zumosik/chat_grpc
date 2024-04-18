package main

import (
	"context"
	"log"
	"os"

	"github.com/zumosik/grpc_chat_protos/go/notifications"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:7778", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := notifications.NewNotificationServiceClient(conn)

	c, err := client.SendNotification(context.Background(),
		&notifications.NotificationRequest{
			Notification: &notifications.NotificationRequest_ConfirmEmail_{
				ConfirmEmail: &notifications.NotificationRequest_ConfirmEmail{
					Email:  os.Getenv("TEST_EMAIL"),
					UserId: "123",
				},
			},
		},
	)
	if err != nil {
		panic(err)
	}

	log.Println(c.GetStatus())
}
