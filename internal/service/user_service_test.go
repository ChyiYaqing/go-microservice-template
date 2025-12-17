package service

import (
	"context"
	"testing"

	v1 "github.com/ChyiYaqing/go-microservice-template/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateUser(t *testing.T) {
	svc := NewUserService()
	ctx := context.Background()

	tests := []struct {
		name    string
		req     *v1.CreateUserRequest
		wantErr bool
	}{
		{
			name: "valid user",
			req: &v1.CreateUserRequest{
				User: &v1.User{
					Email:       "test@example.com",
					DisplayName: "Test User",
				},
			},
			wantErr: false,
		},
		{
			name: "missing email",
			req: &v1.CreateUserRequest{
				User: &v1.User{
					DisplayName: "Test User",
				},
			},
			wantErr: true,
		},
		{
			name:    "nil user",
			req:     &v1.CreateUserRequest{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.CreateUser(ctx, tt.req)
			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateUser() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("CreateUser() unexpected error: %v", err)
				}
				if user == nil {
					t.Errorf("CreateUser() returned nil user")
				}
				if user != nil && user.Name == "" {
					t.Errorf("CreateUser() returned user with empty name")
				}
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	svc := NewUserService()
	ctx := context.Background()

	// Create a user first
	createReq := &v1.CreateUserRequest{
		User: &v1.User{
			Email:       "test@example.com",
			DisplayName: "Test User",
		},
	}
	createdUser, err := svc.CreateUser(ctx, createReq)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	tests := []struct {
		name    string
		req     *v1.GetUserRequest
		wantErr bool
	}{
		{
			name: "existing user",
			req: &v1.GetUserRequest{
				Name: createdUser.Name,
			},
			wantErr: false,
		},
		{
			name: "non-existing user",
			req: &v1.GetUserRequest{
				Name: "users/999",
			},
			wantErr: true,
		},
		{
			name:    "empty name",
			req:     &v1.GetUserRequest{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.GetUser(ctx, tt.req)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetUser() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("GetUser() unexpected error: %v", err)
				}
				if user == nil {
					t.Errorf("GetUser() returned nil user")
				}
			}
		})
	}
}

func TestListUsers(t *testing.T) {
	svc := NewUserService()
	ctx := context.Background()

	// Create some users
	for i := 0; i < 5; i++ {
		_, err := svc.CreateUser(ctx, &v1.CreateUserRequest{
			User: &v1.User{
				Email:       "test@example.com",
				DisplayName: "Test User",
			},
		})
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
	}

	tests := []struct {
		name     string
		req      *v1.ListUsersRequest
		wantErr  bool
		minUsers int
	}{
		{
			name:     "list all users",
			req:      &v1.ListUsersRequest{},
			wantErr:  false,
			minUsers: 5,
		},
		{
			name: "list with page size",
			req: &v1.ListUsersRequest{
				PageSize: 2,
			},
			wantErr:  false,
			minUsers: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.ListUsers(ctx, tt.req)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ListUsers() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ListUsers() unexpected error: %v", err)
				}
				if resp == nil {
					t.Errorf("ListUsers() returned nil response")
				}
				if resp != nil && len(resp.Users) < tt.minUsers {
					t.Errorf("ListUsers() returned %d users, want at least %d", len(resp.Users), tt.minUsers)
				}
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	svc := NewUserService()
	ctx := context.Background()

	// Create a user first
	createReq := &v1.CreateUserRequest{
		User: &v1.User{
			Email:       "test@example.com",
			DisplayName: "Test User",
		},
	}
	createdUser, err := svc.CreateUser(ctx, createReq)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	tests := []struct {
		name    string
		req     *v1.UpdateUserRequest
		wantErr bool
	}{
		{
			name: "valid update",
			req: &v1.UpdateUserRequest{
				User: &v1.User{
					Name:        createdUser.Name,
					Email:       "updated@example.com",
					DisplayName: "Updated User",
				},
			},
			wantErr: false,
		},
		{
			name: "non-existing user",
			req: &v1.UpdateUserRequest{
				User: &v1.User{
					Name:  "users/999",
					Email: "test@example.com",
				},
			},
			wantErr: true,
		},
		{
			name:    "nil user",
			req:     &v1.UpdateUserRequest{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.UpdateUser(ctx, tt.req)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateUser() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("UpdateUser() unexpected error: %v", err)
				}
				if user == nil {
					t.Errorf("UpdateUser() returned nil user")
				}
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	svc := NewUserService()
	ctx := context.Background()

	// Create a user first
	createReq := &v1.CreateUserRequest{
		User: &v1.User{
			Email:       "test@example.com",
			DisplayName: "Test User",
		},
	}
	createdUser, err := svc.CreateUser(ctx, createReq)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	tests := []struct {
		name    string
		req     *v1.DeleteUserRequest
		wantErr bool
	}{
		{
			name: "existing user",
			req: &v1.DeleteUserRequest{
				Name: createdUser.Name,
			},
			wantErr: false,
		},
		{
			name: "non-existing user",
			req: &v1.DeleteUserRequest{
				Name: "users/999",
			},
			wantErr: true,
		},
		{
			name:    "empty name",
			req:     &v1.DeleteUserRequest{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.DeleteUser(ctx, tt.req)
			if tt.wantErr {
				if err == nil {
					t.Errorf("DeleteUser() expected error, got nil")
				}
				// Check that we get NotFound error for non-existing user
				if tt.name == "non-existing user" {
					if st, ok := status.FromError(err); ok {
						if st.Code() != codes.NotFound {
							t.Errorf("DeleteUser() expected NotFound error, got %v", st.Code())
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("DeleteUser() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestBatchGetUsers(t *testing.T) {
	svc := NewUserService()
	ctx := context.Background()

	// Create some users
	var userNames []string
	for i := 0; i < 3; i++ {
		user, err := svc.CreateUser(ctx, &v1.CreateUserRequest{
			User: &v1.User{
				Email:       "test@example.com",
				DisplayName: "Test User",
			},
		})
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		userNames = append(userNames, user.Name)
	}

	tests := []struct {
		name     string
		req      *v1.BatchGetUsersRequest
		wantErr  bool
		wantSize int
	}{
		{
			name: "existing users",
			req: &v1.BatchGetUsersRequest{
				Names: userNames,
			},
			wantErr:  false,
			wantSize: 3,
		},
		{
			name:    "empty names",
			req:     &v1.BatchGetUsersRequest{},
			wantErr: true,
		},
		{
			name: "mixed existing and non-existing",
			req: &v1.BatchGetUsersRequest{
				Names: append(userNames, "users/999"),
			},
			wantErr:  false,
			wantSize: 3, // Should only return existing users
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.BatchGetUsers(ctx, tt.req)
			if tt.wantErr {
				if err == nil {
					t.Errorf("BatchGetUsers() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("BatchGetUsers() unexpected error: %v", err)
				}
				if resp == nil {
					t.Errorf("BatchGetUsers() returned nil response")
				}
				if resp != nil && len(resp.Users) != tt.wantSize {
					t.Errorf("BatchGetUsers() returned %d users, want %d", len(resp.Users), tt.wantSize)
				}
			}
		})
	}
}
