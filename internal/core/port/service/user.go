package service

import (
	"grpc-hexa/internal/core/model/request"
	"grpc-hexa/internal/core/model/response"
)

type UserService interface {
	SignUp(request *request.SignUpRequest) *response.Response
}
