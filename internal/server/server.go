package server

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
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
	db *pgxpool.Pool
}

func NewServer(db *pgxpool.Pool) *Server {
	return &Server{
		UnimplementedUserServer: desc.UnimplementedUserServer{},
		db:                      db,
	}
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

	var count int

	countUser := sq.Select("count(*)").PlaceholderFormat(sq.Dollar).From("users").Where(sq.Eq{"email": req.Email})

	query, args, err := countUser.ToSql()
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRow(ctx, query, args...).Scan(&count)

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password", "role").
		Values(req.GetName(), req.GetEmail(), req.GetPassword(), desc.Role_name[int32(req.GetRole())]).
		Suffix("RETURNING id")

	query, args, err = builderInsert.ToSql()

	if err != nil {
		return nil, err
	}

	var returnId int64

	err = s.db.QueryRow(ctx, query, args...).Scan(&returnId)

	return &desc.CreateResponse{
		Id: returnId,
	}, nil

}

func (s *Server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	returnUser := &desc.GetResponse{}

	var (
		id        int64
		name      string
		email     string
		role      string
		createdAt time.Time
		updatedAt time.Time
	)

	builderGet := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"id": req.GetId()}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builderGet.ToSql()
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	// Заполняем returnUser после успешного сканирования
	returnUser.Id = id
	returnUser.Name = name
	returnUser.Email = email
	returnUser.Role = desc.Role(desc.Role_value[role])
	// Если CreatedAt и UpdatedAt — строки, можно форматировать time.Time в строку, например так:
	returnUser.CreatedAt = timestamppb.New(createdAt)
	returnUser.UpdatedAt = timestamppb.New(updatedAt)

	return returnUser, nil
}

func (s *Server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	builderUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	// Добавляем поля динамически
	if req.Name != nil {
		builderUpdate = builderUpdate.Set("name", req.Name.GetValue())
	}

	if req.Email != nil {
		builderUpdate = builderUpdate.Set("email", req.Email.GetValue())
	}

	query, args, err := builderUpdate.ToSql()

	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	builderDelete := sq.Delete("users").Where(sq.Eq{"id": req.GetId()}).PlaceholderFormat(sq.Dollar)

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
