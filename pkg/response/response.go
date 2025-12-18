package response

import (
	apiv1 "github.com/ChyiYaqing/go-microservice-template/api/proto/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// Error codes
const (
	CodeSuccess            = 0
	CodeInvalidArgument    = 400
	CodeNotFound           = 404
	CodeInternalError      = 500
	CodeAlreadyExists      = 409
	CodePermissionDenied   = 403
	CodeUnauthenticated    = 401
	CodeResourceExhausted  = 429
	CodeUnimplemented      = 501
)

// Error messages
const (
	MsgSuccess           = "success"
	MsgInvalidArgument   = "invalid argument"
	MsgNotFound          = "resource not found"
	MsgInternalError     = "internal server error"
	MsgAlreadyExists     = "resource already exists"
	MsgPermissionDenied  = "permission denied"
	MsgUnauthenticated   = "unauthenticated"
	MsgResourceExhausted = "resource exhausted"
	MsgUnimplemented     = "unimplemented"
)

// Success creates a successful response with data
func Success(data interface{}) (*apiv1.CommonResponse, error) {
	structData, err := structpb.NewStruct(map[string]interface{}{
		"result": data,
	})
	if err != nil {
		return nil, err
	}

	return &apiv1.CommonResponse{
		ErrorCode: CodeSuccess,
		ErrorMsg:  MsgSuccess,
		Data:      structData,
	}, nil
}

// Error creates an error response
func Error(code int32, message string) *apiv1.CommonResponse {
	return &apiv1.CommonResponse{
		ErrorCode: code,
		ErrorMsg:  message,
		Data:      nil,
	}
}

// InvalidArgument creates an invalid argument error response
func InvalidArgument(message string) *apiv1.CommonResponse {
	if message == "" {
		message = MsgInvalidArgument
	}
	return Error(CodeInvalidArgument, message)
}

// NotFound creates a not found error response
func NotFound(message string) *apiv1.CommonResponse {
	if message == "" {
		message = MsgNotFound
	}
	return Error(CodeNotFound, message)
}

// InternalError creates an internal error response
func InternalError(message string) *apiv1.CommonResponse {
	if message == "" {
		message = MsgInternalError
	}
	return Error(CodeInternalError, message)
}

// AlreadyExists creates an already exists error response
func AlreadyExists(message string) *apiv1.CommonResponse {
	if message == "" {
		message = MsgAlreadyExists
	}
	return Error(CodeAlreadyExists, message)
}

// SuccessEmpty creates a successful response with empty data
func SuccessEmpty() *apiv1.CommonResponse {
	return &apiv1.CommonResponse{
		ErrorCode: CodeSuccess,
		ErrorMsg:  MsgSuccess,
		Data:      nil,
	}
}
