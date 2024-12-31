package controller

import (
	"context"
	"example/db"
	"example/models"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	sqlc "example/db/sqlc"

	"github.com/gin-gonic/gin"
)

type result struct {
	services            []models.Service
	serviceProviderType map[string][]models.ServiceProvider
}

func GetAppointment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		user, _ := c.Get("user")
		users := user.(sqlc.User)
		fmt.Println("", users.AadhaarNo)

		defer cancel()
		appointments, err := db.Query.GetAppointment(ctx, users)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusCreated, models.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"appointments": appointments}})
	}
}
func GetAvailability() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
		var dateTime models.DateTime
		defer cancel()
		if err := c.BindJSON(&dateTime); err != nil {
			c.JSON(http.StatusBadRequest, models.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		user, _ := c.Get("user")
		users := user.(sqlc.User)
		typeAvailability, err := db.Query.GetTypeAvailability(ctx, dateTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		ss := []models.ServiceProvider{}
		serviceProviderType := map[string][]models.ServiceProvider{}
		for _, k := range typeAvailability {
			servicesType, err := db.Query.GetListServiceTypeByDate(ctx, k)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			for _, v := range servicesType {
				v.Date = k.Date
				v.Time = k.Time
				val := haversine(users.Lat.Float64, users.Lon.Float64, v.Lat, v.Lon)
				val = strings.TrimSpace(val)
				km, _ := strconv.ParseFloat(val, 64)
				if km <= 10 {
					ss = append(ss, v)
					serviceProviderType[dateTime.ServiceType] = ss

				}
			}
		}
		c.JSON(http.StatusCreated, models.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"serviceProviderType": serviceProviderType}})
	}
}

//	func GetAvailability() gin.HandlerFunc {
//		return func(c *gin.Context) {
//			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//			var availability models.Availability
//			defer cancel()
//			if err := c.BindJSON(&availability); err != nil {
//				c.JSON(http.StatusBadRequest, models.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
//				return
//			}
//			typeAvailability, err := db.Query.GetTypeAvailability(ctx, availability)
//			if err != nil {
//				c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
//				return
//			}
//			c.JSON(http.StatusCreated, models.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"typeAvailability": typeAvailability}})
//		}
//	}
func GetServiceList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		services, err := db.Query.GetListService(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusCreated, models.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"services": services}})
	}
}
func GetServiceListBytype() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var id models.Service
		user, _ := c.Get("user")
		users := user.(sqlc.User)
		defer cancel()
		if err := c.BindJSON(&id); err != nil {
			c.JSON(http.StatusBadRequest, models.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		servicesType, err := db.Query.GetListServiceType(ctx, id.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		ss := []models.ServiceProvider{}
		serviceProviderType := map[string][]models.ServiceProvider{}
		for _, v := range servicesType {
			val := haversine(users.Lat.Float64, users.Lon.Float64, v.Lat, v.Lon)
			val = strings.TrimSpace(val)
			km, _ := strconv.ParseFloat(val, 64)
			if km <= 10 {
				ss = append(ss, v)
				serviceProviderType[id.ServiceType] = ss

			}
		}

		c.JSON(http.StatusCreated, models.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"serviceProviderType": serviceProviderType}})

	}
}
func GetUserDashboard() gin.HandlerFunc {
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
func haversine(lat1, lon1, lat2, lon2 float64) string {
	const earthRadius = 6371 // Earth radius in kilometers

	// Convert degrees to radians
	toRadians := func(degree float64) float64 {
		return degree * math.Pi / 180
	}

	lat1Rad, lon1Rad := toRadians(lat1), toRadians(lon1)
	lat2Rad, lon2Rad := toRadians(lat2), toRadians(lon2)

	// Haversine formula
	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Distance in kilometers
	distance := earthRadius * c
	return fmt.Sprintf(" %.2f", distance)
}
