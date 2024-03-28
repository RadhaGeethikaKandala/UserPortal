package router

import (
	"userportal/internal/app/database"
	"userportal/internal/app/handler"
	"userportal/internal/app/repository"
	"userportal/internal/app/service"

	"github.com/gin-gonic/gin"
)

func Router(engine *gin.Engine) {

	repository := repository.NewUserRepository(database.CreateDatabaseConn())
	service := service.NewUserService(repository)
	handler := handler.NewUserHandler(service)

	movieRentalApiGroup := engine.Group("/api/users/")
	{
		movieRentalApiGroup.GET("/", handler.GetAllUsers)
		movieRentalApiGroup.GET("/:email", handler.GetUserByEmail)
		movieRentalApiGroup.POST("/add", handler.CreateUser)
		movieRentalApiGroup.POST("/add-many", handler.CreateUsers)
		movieRentalApiGroup.PUT("/update", handler.UpdateUser)
		movieRentalApiGroup.DELETE("/:email", handler.DeleteUserByEmail)

	}
}
