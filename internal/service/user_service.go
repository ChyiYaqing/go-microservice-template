package service

import (
	"context"
	"fmt"
	"sync"

	apiv1 "github.com/ChyiYaqing/go-microservice-template/api/proto/v1"
	"github.com/ChyiYaqing/go-microservice-template/pkg/response"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserService implements the UserServiceServer interface
type UserService struct {
	apiv1.UnimplementedUserServiceServer
	users map[string]*apiv1.User
	mu    sync.RWMutex
	nextID int
}

// NewUserService creates a new UserService
func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*apiv1.User),
		nextID: 1,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *apiv1.CreateUserRequest) (*apiv1.CommonResponse, error) {
	if req.GetUser() == nil {
		return response.InvalidArgument("user is required"), nil
	}

	if req.GetUser().GetEmail() == "" {
		return response.InvalidArgument("email is required"), nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate resource name
	userID := fmt.Sprintf("%d", s.nextID)
	s.nextID++

	now := timestamppb.Now()
	user := &apiv1.User{
		Name:        fmt.Sprintf("users/%s", userID),
		Email:       req.GetUser().GetEmail(),
		DisplayName: req.GetUser().GetDisplayName(),
		PhoneNumber: req.GetUser().GetPhoneNumber(),
		CreateTime:  now,
		UpdateTime:  now,
		IsActive:    true,
	}

	s.users[user.Name] = user
	return response.Success(user)
}

// GetUser retrieves a user by resource name
func (s *UserService) GetUser(ctx context.Context, req *apiv1.GetUserRequest) (*apiv1.CommonResponse, error) {
	if req.GetName() == "" {
		return response.InvalidArgument("name is required"), nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[req.GetName()]
	if !exists {
		return response.NotFound(fmt.Sprintf("user %s not found", req.GetName())), nil
	}

	return response.Success(user)
}

// ListUsers lists users with pagination
func (s *UserService) ListUsers(ctx context.Context, req *apiv1.ListUsersRequest) (*apiv1.CommonResponse, error) {
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
	var allUsers []*apiv1.User
	for _, user := range s.users {
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

	return response.Success(map[string]interface{}{
		"users":           users,
		"next_page_token": nextPageToken,
		"total_size":      len(allUsers),
	})
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, req *apiv1.UpdateUserRequest) (*apiv1.CommonResponse, error) {
	if req.GetUser() == nil {
		return response.InvalidArgument("user is required"), nil
	}

	if req.GetUser().GetName() == "" {
		return response.InvalidArgument("user.name is required"), nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[req.GetUser().GetName()]
	if !exists {
		return response.NotFound(fmt.Sprintf("user %s not found", req.GetUser().GetName())), nil
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
	return response.Success(user)
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, req *apiv1.DeleteUserRequest) (*apiv1.CommonResponse, error) {
	if req.GetName() == "" {
		return response.InvalidArgument("name is required"), nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[req.GetName()]; !exists {
		return response.NotFound(fmt.Sprintf("user %s not found", req.GetName())), nil
	}

	delete(s.users, req.GetName())
	return response.SuccessEmpty(), nil
}

// BatchGetUsers retrieves multiple users
func (s *UserService) BatchGetUsers(ctx context.Context, req *apiv1.BatchGetUsersRequest) (*apiv1.CommonResponse, error) {
	if len(req.GetNames()) == 0 {
		return response.InvalidArgument("names is required"), nil
	}

	if len(req.GetNames()) > 1000 {
		return response.InvalidArgument("cannot retrieve more than 1000 users at once"), nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var users []*apiv1.User
	for _, name := range req.GetNames() {
		if user, exists := s.users[name]; exists {
			users = append(users, user)
		}
	}

	return response.Success(map[string]interface{}{
		"users": users,
	})
}

// updateUserWithMask updates user fields based on field mask
func updateUserWithMask(dst, src *apiv1.User, mask *fieldmaskpb.FieldMask) {
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
