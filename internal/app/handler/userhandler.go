package handler

import (
	"fmt"
	"userportal/internal/app/dto"
	"userportal/internal/app/service"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	CreateUser(ctx *gin.Context)
	CreateUsers(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
	GetUserByEmail(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUserByEmail(ctx *gin.Context)
}

type userhandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) UserHandler {
	return &userhandler{
		service: service,
	}
}

func (uh *userhandler) CreateUser(ctx *gin.Context) {
	var user dto.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		fmt.Println("error in converting request")
		buildErrorResponse(ctx, err.Error())
		return
	}
	err = uh.service.CreateUser(user)
	if err != nil {
		buildErrorResponse(ctx, err.Error())
		return
	}

	ctx.JSON(200, dto.ApiResponse{
		Status:  "success",
		Message: "User added successfully",
		Code:    "200",
	})
}

func (uh *userhandler) CreateUsers(ctx *gin.Context) {
	var users []dto.User
	err := ctx.ShouldBindJSON(&users)
	if err != nil {
		buildErrorResponse(ctx, err.Error())
		return
	}
	err = uh.service.CreateUsers(users)
	if err != nil {
		buildErrorResponse(ctx, err.Error())
		return
	}

	ctx.JSON(200, dto.ApiResponse{
		Status:  "success",
		Message: "Users added successfully",
		Code:    "200",
	})
}

func (uh *userhandler) GetAllUsers(ctx *gin.Context) {
	users, err := uh.service.GetAllUsers()
	if err != nil {
		buildErrorResponse(ctx, err.Error())
		return
	}
	ctx.JSON(200, users)
}

func (uh *userhandler) GetUserByEmail(ctx *gin.Context) {

	user, err := uh.service.GetUserByEmail(ctx.Param("email"))
	if err != nil {
		buildErrorResponse(ctx, err.Error())
		return
	}
	ctx.JSON(200, user)
}

func (uh *userhandler) UpdateUser(ctx *gin.Context) {
	var user dto.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		fmt.Println("error in converting request")
		buildErrorResponse(ctx, err.Error())
		return
	}
	err = uh.service.UpdateUser(user)
	if err != nil {
		buildErrorResponse(ctx, err.Error())
		return
	}

	ctx.JSON(200, dto.ApiResponse{
		Status:  "success",
		Message: "User updated successfully",
		Code:    "200",
	})
}

func (uh *userhandler) DeleteUserByEmail(ctx *gin.Context) {

	err := uh.service.DeleteUserByEmail(ctx.Param("email"))
	if err != nil {
		buildErrorResponse(ctx, err.Error())
		return
	}
	ctx.JSON(200, dto.ApiResponse{
		Status:  "success",
		Message: "User has been deleted succesfully",
		Code:    "200",
	})
}

func buildErrorResponse(ctx *gin.Context, message string) {
	ctx.JSON(400, dto.ApiResponse{
		Status:  "failure",
		Message: message,
		Code:    "400",
	})
}
