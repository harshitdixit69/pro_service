package controller

import (
	"context"
	"example/db"
	"example/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetServiceProviderDashboard() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// user, _ := c.Get("user")

		services, err := db.Query.GetListService(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		serviceProviderType := map[string][]models.ServiceProvider{}
		for _, v := range services {
			servicesType, err := db.Query.GetListServiceType(ctx, v.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			serviceProviderType[v.ServiceType] = servicesType
		}
		var results result
		results.services = services
		results.serviceProviderType = serviceProviderType
		c.JSON(http.StatusCreated, models.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"service": results.services, "serviceProviderType": results.serviceProviderType}})
	}
}

// func UpdateServiceProvider() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 		defer cancel()
// 		// user, _ := c.Get("user")

// 		services, err := db.Query.GetListService(ctx)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
// 			return
// 		}
// 		serviceProviderType := map[string][]models.ServiceProvider{}
// 		for _, v := range services {
// 			servicesType, err := db.Query.GetListServiceType(ctx, v.ID)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
// 				return
// 			}
// 			serviceProviderType[v.ServiceType] = servicesType
// 		}
// 		var results result
// 		results.services = services
// 		results.serviceProviderType = serviceProviderType
// 		c.JSON(http.StatusCreated, models.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"service": results.services, "serviceProviderType": results.serviceProviderType}})
// 	}
// }
