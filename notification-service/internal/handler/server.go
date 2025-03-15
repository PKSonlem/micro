package handler

import (
	"context"
	"errors"

	pb "github.com/timurzdev/mentorship-test-task/notification/pkg/proto"
)

type server struct {
	pb.UnimplementedNotificationServiceServer
}

func (s *server) Send(ctx context.Context, notification *pb.Notification) (*pb.NotificationResponse, error) {
	return nil, errors.New("not implemented")
}
