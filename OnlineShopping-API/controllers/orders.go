package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/melardev/GoGonicEcommerceApi/dtos"
	"github.com/melardev/GoGonicEcommerceApi/middlewares"
	"github.com/melardev/GoGonicEcommerceApi/models"
	"github.com/melardev/GoGonicEcommerceApi/services"
	"net/http"
	"strconv"
)

func RegisterOrderRoutes(router *gin.RouterGroup) {
	router.GET("/list", ListOrders)
}

func ListOrders(c *gin.Context) {
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
	userId := c.MustGet("currentUserId").(uint)

	orders, totalCommentCount, err := services.FetchOrdersPage(userId, page, pageSize)

	c.JSON(http.StatusOK, dtos.CreateOrderPagedResponse(c.Request, orders, page, pageSize, totalCommentCount, false, false))
}