package service

import (
	"grpc-hexa/internal/core/common/utils"
	"grpc-hexa/internal/core/dto"
	"grpc-hexa/internal/core/entity/error_code"
	"grpc-hexa/internal/core/model/request"
	"grpc-hexa/internal/core/model/response"
	"grpc-hexa/internal/core/port/repository"
	"grpc-hexa/internal/core/port/service"
	"time"
)

const (
	invalidUserNameErrMsg = "invalid username"
	invalidPasswordErrMsg = "invalid password"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) service.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (u userService) SignUp(request *request.SignUpRequest) *response.Response {
	// validate request
	if len(request.Username) == 0 {
		return u.createFailedResponse(error_code.InvalidRequest, invalidUserNameErrMsg)
	}

	if len(request.Password) == 0 {
		return u.createFailedResponse(error_code.InvalidRequest, invalidPasswordErrMsg)
	}

	now := time.Now()
	userDTO := dto.UserDTO{
		UserName:    request.Username,
		Password:    request.Password,
		DisplayName: u.getRandomDisplayName(request.Username),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// save a new user
	err := u.userRepo.Insert(userDTO)
	if err != nil {
		if err == repository.DuplicateUser {
			return u.createFailedResponse(error_code.DuplicateUser, err.Error())
		}
		return u.createFailedResponse(error_code.InternalError, error_code.InternalErrMsg)
	}

	// create data response
	signUpData := response.SignUpDataResponse{
		DisplayName: userDTO.DisplayName,
	}
	return u.createSuccessResponse(signUpData)
}

func (u userService) getRandomDisplayName(username string) string {
	return username + utils.GetUUID()
}

func (u userService) createFailedResponse(
	code error_code.ErrorCode, message string,
) *response.Response {
	return &response.Response{
		Status:       false,
		ErrorCode:    code,
		ErrorMessage: message,
	}
}

func (u userService) createSuccessResponse(data response.SignUpDataResponse) *response.Response {
	return &response.Response{
		Data:         data,
		Status:       true,
		ErrorCode:    error_code.Success,
		ErrorMessage: error_code.SuccessErrMsg,
	}
}
