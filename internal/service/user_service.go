package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	v1 "github.com/ChyiYaqing/go-microservice-template/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserService implements the UserServiceServer interface
type UserService struct {
	v1.UnimplementedUserServiceServer
	users map[string]*v1.User
	mu    sync.RWMutex
	nextID int
}

// NewUserService creates a new UserService
func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*v1.User),
		nextID: 1,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.User, error) {
	if req.GetUser() == nil {
		return nil, status.Error(codes.InvalidArgument, "user is required")
	}

	if req.GetUser().GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate resource name
	userID := fmt.Sprintf("%d", s.nextID)
	s.nextID++

	now := timestamppb.Now()
	user := &v1.User{
		Name:        fmt.Sprintf("users/%s", userID),
		Email:       req.GetUser().GetEmail(),
		DisplayName: req.GetUser().GetDisplayName(),
		PhoneNumber: req.GetUser().GetPhoneNumber(),
		CreateTime:  now,
		UpdateTime:  now,
		IsActive:    true,
	}

	s.users[user.Name] = user
	return user, nil
}

// GetUser retrieves a user by resource name
func (s *UserService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.User, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[req.GetName()]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user %s not found", req.GetName())
	}

	return user, nil
}

// ListUsers lists users with pagination
func (s *UserService) ListUsers(ctx context.Context, req *v1.ListUsersRequest) (*v1.ListUsersResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 1000 {
		pageSize = 1000
	}

	// Convert map to slice
	var allUsers []*v1.User
	for _, user := range s.users {
		allUsers = allUsers
		allUsers = append(allUsers, user)
	}

	// Simple pagination (in production, use a more robust approach)
	start := 0
	if req.GetPageToken() != "" {
		// Parse page token (simplified)
		fmt.Sscanf(req.GetPageToken(), "%d", &start)
	}

	end := start + int(pageSize)
	if end > len(allUsers) {
		end = len(allUsers)
	}

	users := allUsers[start:end]

	var nextPageToken string
	if end < len(allUsers) {
		nextPageToken = fmt.Sprintf("%d", end)
	}

	return &v1.ListUsersResponse{
		Users:         users,
		NextPageToken: nextPageToken,
		TotalSize:     int32(len(allUsers)),
	}, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.User, error) {
	if req.GetUser() == nil {
		return nil, status.Error(codes.InvalidArgument, "user is required")
	}

	if req.GetUser().GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "user.name is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[req.GetUser().GetName()]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user %s not found", req.GetUser().GetName())
	}

	// Apply field mask if provided
	if req.GetUpdateMask() != nil {
		updateUserWithMask(user, req.GetUser(), req.GetUpdateMask())
	} else {
		// Update all fields if no mask provided
		if req.GetUser().GetEmail() != "" {
			user.Email = req.GetUser().GetEmail()
		}
		if req.GetUser().GetDisplayName() != "" {
			user.DisplayName = req.GetUser().GetDisplayName()
		}
		if req.GetUser().GetPhoneNumber() != "" {
			user.PhoneNumber = req.GetUser().GetPhoneNumber()
		}
		user.IsActive = req.GetUser().GetIsActive()
	}

	user.UpdateTime = timestamppb.Now()
	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*emptypb.Empty, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[req.GetName()]; !exists {
		return nil, status.Errorf(codes.NotFound, "user %s not found", req.GetName())
	}

	delete(s.users, req.GetName())
	return &emptypb.Empty{}, nil
}

// BatchGetUsers retrieves multiple users
func (s *UserService) BatchGetUsers(ctx context.Context, req *v1.BatchGetUsersRequest) (*v1.BatchGetUsersResponse, error) {
	if len(req.GetNames()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "names is required")
	}

	if len(req.GetNames()) > 1000 {
		return nil, status.Error(codes.InvalidArgument, "cannot retrieve more than 1000 users at once")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var users []*v1.User
	for _, name := range req.GetNames() {
		if user, exists := s.users[name]; exists {
			users = append(users, user)
		}
	}

	return &v1.BatchGetUsersResponse{
		Users: users,
	}, nil
}

// updateUserWithMask updates user fields based on field mask
func updateUserWithMask(dst, src *v1.User, mask *fieldmaskpb.FieldMask) {
	for _, path := range mask.GetPaths() {
		switch path {
		case "email":
			dst.Email = src.Email
		case "display_name":
			dst.DisplayName = src.DisplayName
		case "phone_number":
			dst.PhoneNumber = src.PhoneNumber
		case "is_active":
			dst.IsActive = src.IsActive
		}
	}
}
