package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	desc "github.com/vizurth/auth/pkg/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"time"
)

type Server struct {
	desc.UnimplementedUserServer
}

func NewServer() *Server {
	return &Server{}
}

type User struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

var value map[int64]*User = make(map[int64]*User)

func (s *Server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	rand.Seed(time.Now().UnixNano())
	temp := &User{
		Id:       rand.Int63(),
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Role:     int(req.GetRole()),
	}

	if _, ok := value[temp.Id]; ok {
		return nil, errors.New("User already exists")
	} else {
		value[temp.Id] = temp
	}

	return &desc.CreateResponse{
		Id: temp.Id,
	}, nil

}

func (s *Server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, ok := value[req.GetId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	return &desc.GetResponse{
		Id:        req.GetId(),
		Name:      user.Name,
		Email:     user.Email,
		Role:      desc.Role(user.Role),
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

func (s *Server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	fmt.Println(req.GetId())
	user, ok := value[req.GetId()]
	if !ok {
		return nil, errors.New("User not found")
	}
	if req.Name != nil {
		user.Name = req.Name.GetValue()
	}
	if req.Email != nil {
		user.Email = req.Email.GetValue()
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	_, ok := value[req.GetId()]
	if !ok {
		return nil, errors.New("User not found")
	}
	delete(value, req.GetId())
	return &emptypb.Empty{}, nil
}
