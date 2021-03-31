package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/melardev/GoGonicEcommerceApi/dtos"
	"github.com/melardev/GoGonicEcommerceApi/services"

	"github.com/melardev/GoGonicEcommerceApi/models"

	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func RegisterUserRoutes(router *gin.RouterGroup) {
	router.POST("/create", UsersRegistration)
	router.POST("/list", UsersList)
	router.POST("/login", UsersLogin)
}

func UsersRegistration(c *gin.Context) {

	var json dtos.RegisterRequestDto
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
	if err := services.CreateOne(&models.User{
		Username:  json.Username,
		Password:  string(password),
		FirstName: json.FirstName,
		LastName:  json.LastName,
		Email:     json.Email,
	}); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("database", err))
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success":       true,
		"full_messages": []string{"User created successfully"}})
}

func UsersLogin(c *gin.Context) {

	var json dtos.LoginRequestDto
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	user, err := services.FindOneUser(&models.User{Username: json.Username})

	if err != nil {
		c.JSON(http.StatusForbidden, dtos.CreateDetailedErrorDto("login_error", err))
		return
	}

	if user.IsValidPassword(json.Password) != nil {
		c.JSON(http.StatusForbidden, dtos.CreateDetailedErrorDto("login", errors.New("invalid credentials")))
		return
	}

	c.JSON(http.StatusOK, dtos.CreateLoginSuccessful(&user))

}

func UsersList(c *gin.Context) {

	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 5
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	productModels, modelCount, commentsCount, err := services.FetchUsersPage(page, pageSize)
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("products", errors.New("Invalid param")))
		return
	}

	c.JSON(http.StatusOK, dtos.CreatedUserPagedResponse(c.Request, productModels, page, pageSize, modelCount, commentsCount))
}
