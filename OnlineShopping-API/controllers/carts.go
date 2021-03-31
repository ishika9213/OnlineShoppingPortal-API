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

func RegisterCartRoutes(router *gin.RouterGroup) {
	router.POST("/add", CreateCart)
	router.GET("/:cartId/complete", CreateOrder)
	router.GET("/list", ShowOrder)
}

func CreateOrder(c *gin.Context) {
	var orderRequest dtos.CreateOrderRequestDto
	if err := c.ShouldBind(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	userObj, userLoggedIn := c.Get("currentUser")
	var user models.User
	if userLoggedIn {
		user = (userObj).(models.User)
	}

	order := models.Order{
		TrackingNumber: randomString(16),
		OrderStatus:    0,
		Address:        address,
		AddressId:      address.ID,
	}

	if userLoggedIn {
		order.UserId = user.ID
		order.User = user
	}

	var productIds = make([]uint, len(orderRequest.CartItems))
	for i := 0; i < len(orderRequest.CartItems); i++ {
		productIds[i] = orderRequest.CartItems[i].Id
	}

	products, err := services.FetchProductsIdNameAndPrice(productIds)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("db_error", err))
		return
	}

	if len(products) != len(orderRequest.CartItems) {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateErrorDtoWithMessage("make sure all products are still available"))
		return
	}
	orderItems := make([]models.OrderItem, len(products))

	for i := 0; i < len(products); i++ {
		// I am assuming product ids returned are in the same order as the cart_items, TODO: implement a more robust code to ensure
		orderItems[i] = models.OrderItem{
			ProductId:   products[i].ID,
			ProductName: products[i].Name,
			Slug:        products[i].Slug,
			Quantity:    orderRequest.CartItems[i].Quantity,
		}
	}

	order.OrderItems = orderItems
	err = services.CreateOne(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, dtos.CreateOrderCreatedDto(&order))

}

func ShowOrder(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Param("id"))
	user := c.MustGet("currentUser").(models.User)
	order, err := services.FetchOrderDetails(uint(orderId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.CreateDetailedErrorDto("db_error", err))
		return
	}

	if order.UserId == user.ID || user.IsAdmin() {
		c.JSON(http.StatusOK, dtos.CreateOrderDetailsDto(&order))
	} else {
		c.JSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("Permission denied, you can not view this order"))
		return
	}
}
