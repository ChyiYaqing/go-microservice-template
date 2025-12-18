package service

import (
	"context"
	"testing"

	apiv1 "github.com/ChyiYaqing/go-microservice-template/api/proto/v1"
	"github.com/ChyiYaqing/go-microservice-template/pkg/response"
)

func TestCreateUser(t *testing.T) {
	svc := NewUserService()
	ctx := context.Background()

	tests := []struct {
		name          string
		req           *apiv1.CreateUserRequest
		wantErrorCode int32
	}{
		{
			name: "valid user",
			req: &apiv1.CreateUserRequest{
				User: &apiv1.User{
					Email:       "test@example.com",
					DisplayName: "Test User",
				},
			},
			wantErrorCode: response.CodeSuccess,
		},
		{
			name: "missing email",
			req: &apiv1.CreateUserRequest{
				User: &apiv1.User{
					DisplayName: "Test User",
				},
			},
			wantErrorCode: response.CodeInvalidArgument,
		},
		{
			name:          "nil user",
			req:           &apiv1.CreateUserRequest{},
			wantErrorCode: response.CodeInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.CreateUser(ctx, tt.req)
			if err != nil {
				t.Errorf("CreateUser() unexpected error: %v", err)
			}
			if resp == nil {
				t.Errorf("CreateUser() returned nil response")
				return
			}
			if resp.ErrorCode != tt.wantErrorCode {
				t.Errorf("CreateUser() error_code = %d, want %d", resp.ErrorCode, tt.wantErrorCode)
			}
			if tt.wantErrorCode == response.CodeSuccess && resp.Data == nil {
				t.Errorf("CreateUser() success response should have data")
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	svc := NewUserService()
	ctx := context.Background()

	// Create a user first
	createResp, _ := svc.CreateUser(ctx, &apiv1.CreateUserRequest{
		User: &apiv1.User{
			Email:       "test@example.com",
			DisplayName: "Test User",
		},
	})

	var userName string
	if createResp != nil && createResp.Data != nil {
		if result, ok := createResp.Data.Fields["result"]; ok {
			if userStruct, ok := result.GetStructValue().Fields["name"]; ok {
				userName = userStruct.GetStringValue()
			}
		}
	}

	tests := []struct {
		name          string
		req           *apiv1.GetUserRequest
		wantErrorCode int32
	}{
		{
			name: "existing user",
			req: &apiv1.GetUserRequest{
				Name: userName,
			},
			wantErrorCode: response.CodeSuccess,
		},
		{
			name: "non-existing user",
			req: &apiv1.GetUserRequest{
				Name: "users/999",
			},
			wantErrorCode: response.CodeNotFound,
		},
		{
			name:          "empty name",
			req:           &apiv1.GetUserRequest{},
			wantErrorCode: response.CodeInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.GetUser(ctx, tt.req)
			if err != nil {
				t.Errorf("GetUser() unexpected error: %v", err)
			}
			if resp == nil {
				t.Errorf("GetUser() returned nil response")
				return
			}
			if resp.ErrorCode != tt.wantErrorCode {
				t.Errorf("GetUser() error_code = %d, want %d", resp.ErrorCode, tt.wantErrorCode)
			}
		})
	}
}

func TestListUsers(t *testing.T) {
	svc := NewUserService()
	ctx := context.Background()

	// Create some users
	for i := 0; i < 5; i++ {
		svc.CreateUser(ctx, &apiv1.CreateUserRequest{
			User: &apiv1.User{
				Email:       "test@example.com",
				DisplayName: "Test User",
			},
		})
	}

	tests := []struct {
		name          string
		req           *apiv1.ListUsersRequest
		wantErrorCode int32
		minUsers      int
	}{
		{
			name:          "list all users",
			req:           &apiv1.ListUsersRequest{},
			wantErrorCode: response.CodeSuccess,
			minUsers:      5,
		},
		{
			name: "list with page size",
			req: &apiv1.ListUsersRequest{
				PageSize: 2,
			},
			wantErrorCode: response.CodeSuccess,
			minUsers:      2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.ListUsers(ctx, tt.req)
			if err != nil {
				t.Errorf("ListUsers() unexpected error: %v", err)
			}
			if resp == nil {
				t.Errorf("ListUsers() returned nil response")
				return
			}
			if resp.ErrorCode != tt.wantErrorCode {
				t.Errorf("ListUsers() error_code = %d, want %d", resp.ErrorCode, tt.wantErrorCode)
			}
			if tt.wantErrorCode == response.CodeSuccess && resp.Data == nil {
				t.Errorf("ListUsers() success response should have data")
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	svc := NewUserService()
	ctx := context.Background()

	// Create a user first
	createResp, _ := svc.CreateUser(ctx, &apiv1.CreateUserRequest{
		User: &apiv1.User{
			Email:       "test@example.com",
			DisplayName: "Test User",
		},
	})

	var userName string
	if createResp != nil && createResp.Data != nil {
		if result, ok := createResp.Data.Fields["result"]; ok {
			if userStruct, ok := result.GetStructValue().Fields["name"]; ok {
				userName = userStruct.GetStringValue()
			}
		}
	}

	tests := []struct {
		name          string
		req           *apiv1.UpdateUserRequest
		wantErrorCode int32
	}{
		{
			name: "valid update",
			req: &apiv1.UpdateUserRequest{
				User: &apiv1.User{
					Name:        userName,
					Email:       "updated@example.com",
					DisplayName: "Updated User",
				},
			},
			wantErrorCode: response.CodeSuccess,
		},
		{
			name: "non-existing user",
			req: &apiv1.UpdateUserRequest{
				User: &apiv1.User{
					Name:  "users/999",
					Email: "test@example.com",
				},
			},
			wantErrorCode: response.CodeNotFound,
		},
		{
			name:          "nil user",
			req:           &apiv1.UpdateUserRequest{},
			wantErrorCode: response.CodeInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.UpdateUser(ctx, tt.req)
			if err != nil {
				t.Errorf("UpdateUser() unexpected error: %v", err)
			}
			if resp == nil {
				t.Errorf("UpdateUser() returned nil response")
				return
			}
			if resp.ErrorCode != tt.wantErrorCode {
				t.Errorf("UpdateUser() error_code = %d, want %d", resp.ErrorCode, tt.wantErrorCode)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	svc := NewUserService()
	ctx := context.Background()

	// Create a user first
	createResp, _ := svc.CreateUser(ctx, &apiv1.CreateUserRequest{
		User: &apiv1.User{
			Email:       "test@example.com",
			DisplayName: "Test User",
		},
	})

	var userName string
	if createResp != nil && createResp.Data != nil {
		if result, ok := createResp.Data.Fields["result"]; ok {
			if userStruct, ok := result.GetStructValue().Fields["name"]; ok {
				userName = userStruct.GetStringValue()
			}
		}
	}

	tests := []struct {
		name          string
		req           *apiv1.DeleteUserRequest
		wantErrorCode int32
	}{
		{
			name: "existing user",
			req: &apiv1.DeleteUserRequest{
				Name: userName,
			},
			wantErrorCode: response.CodeSuccess,
		},
		{
			name: "non-existing user",
			req: &apiv1.DeleteUserRequest{
				Name: "users/999",
			},
			wantErrorCode: response.CodeNotFound,
		},
		{
			name:          "empty name",
			req:           &apiv1.DeleteUserRequest{},
			wantErrorCode: response.CodeInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.DeleteUser(ctx, tt.req)
			if err != nil {
				t.Errorf("DeleteUser() unexpected error: %v", err)
			}
			if resp == nil {
				t.Errorf("DeleteUser() returned nil response")
				return
			}
			if resp.ErrorCode != tt.wantErrorCode {
				t.Errorf("DeleteUser() error_code = %d, want %d, msg = %s",
					resp.ErrorCode, tt.wantErrorCode, resp.ErrorMsg)
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
		createResp, _ := svc.CreateUser(ctx, &apiv1.CreateUserRequest{
			User: &apiv1.User{
				Email:       "test@example.com",
				DisplayName: "Test User",
			},
		})

		if createResp != nil && createResp.Data != nil {
			if result, ok := createResp.Data.Fields["result"]; ok {
				if userStruct, ok := result.GetStructValue().Fields["name"]; ok {
					userNames = append(userNames, userStruct.GetStringValue())
				}
			}
		}
	}

	tests := []struct {
		name          string
		req           *apiv1.BatchGetUsersRequest
		wantErrorCode int32
	}{
		{
			name: "existing users",
			req: &apiv1.BatchGetUsersRequest{
				Names: userNames,
			},
			wantErrorCode: response.CodeSuccess,
		},
		{
			name:          "empty names",
			req:           &apiv1.BatchGetUsersRequest{},
			wantErrorCode: response.CodeInvalidArgument,
		},
		{
			name: "mixed existing and non-existing",
			req: &apiv1.BatchGetUsersRequest{
				Names: append(userNames, "users/999"),
			},
			wantErrorCode: response.CodeSuccess,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.BatchGetUsers(ctx, tt.req)
			if err != nil {
				t.Errorf("BatchGetUsers() unexpected error: %v", err)
			}
			if resp == nil {
				t.Errorf("BatchGetUsers() returned nil response")
				return
			}
			if resp.ErrorCode != tt.wantErrorCode {
				t.Errorf("BatchGetUsers() error_code = %d, want %d", resp.ErrorCode, tt.wantErrorCode)
			}
		})
	}
}
