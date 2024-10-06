package http

import (
	"grpc-hexa/internal/core/common/router"
	"grpc-hexa/internal/core/entity/error_code"
	"grpc-hexa/internal/core/model/request"
	"grpc-hexa/internal/core/model/response"
	"grpc-hexa/internal/core/port/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	invalidRequestResponse = &response.Response{
		ErrorCode:    error_code.InvalidRequest,
		ErrorMessage: error_code.InvalidRequestErrMsg,
		Status:       false,
	}
)

type UserController struct {
	gin         *gin.Engine
	userService service.UserService
}

func NewUserController(
	gin *gin.Engine,
	userService service.UserService,
) UserController {
	return UserController{
		gin:         gin,
		userService: userService,
	}
}

func (u UserController) InitRouter() {
	api := u.gin.Group("/api/v1")
	router.Post(api, "/signup", u.signUp)
}

func (u UserController) signUp(c *gin.Context) {
	req, err := u.parseRequest(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, &invalidRequestResponse)
		return
	}

	resp := u.userService.SignUp(req)
	c.JSON(http.StatusOK, resp)
}

func (u UserController) parseRequest(ctx *gin.Context) (*request.SignUpRequest, error) {
	var req request.SignUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	return &req, nil
}
