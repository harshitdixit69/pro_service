package routes

import (
	"example/middleware"
	mongoControllers "example/mongo_controller"
	mysqlControllers "example/mysql_controller"

	"github.com/gin-gonic/gin"
)

func MongodbUserRoute(router *gin.Engine) {
	router.POST("/users", mongoControllers.CreateUser())
	router.GET("/users/:userId", mongoControllers.GetAUser())
	router.PUT("/users/:userId", mongoControllers.EditAUser())
	router.DELETE("/users/:userId", mongoControllers.DeleteAUser())
	router.GET("/userss", mongoControllers.GetAllUsers())
}
func MysqlUserRoute(router *gin.Engine) {
	v1 := router.Group("v1")
	{
		userGroup := v1.Group("user")
		userGroup.POST("/signup", mysqlControllers.CreateUser())
		userGroup.POST("/login", mysqlControllers.LoginUser())
		userGroup.POST("/get_otp_forgot_pass", mysqlControllers.GetOTP())
		userGroup.POST("/submit_otp_forgot_pass", mysqlControllers.SubmitOTP())
		userGroup.PUT("/forgot_password", mysqlControllers.ForgotPassword())

		serviceGroup := v1.Group("service")

		serviceGroup.POST("/dashboard", middleware.RequireAuth, mysqlControllers.GetUserDashboard())
		serviceGroup.POST("/serviceListBytype", middleware.RequireAuth, mysqlControllers.GetServiceListBytype())
		serviceGroup.POST("/serviceList", middleware.RequireAuth, mysqlControllers.GetServiceList())
		serviceGroup.POST("/availability", middleware.RequireAuth, mysqlControllers.GetAvailability())
		//to do : payment method
		serviceGroup.POST("/appointment", middleware.RequireAuth, mysqlControllers.GetAppointment())

		serviceProviderGroup := v1.Group("service_provider")

		serviceProviderGroup.POST("/dashboard", middleware.RequireAuth, mysqlControllers.GetServiceProviderDashboard())
		// controller.GetAccountEntries()
		// controller.CreateData()
		router.POST("/account", mysqlControllers.CreateAccount())
		router.GET("/account/:accountId", mysqlControllers.GetAccount())
		router.GET("/accounts", mysqlControllers.GetListAccount())
	}
}
