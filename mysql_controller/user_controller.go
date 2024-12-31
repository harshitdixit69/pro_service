package controller

import (
	"context"
	"database/sql"
	"example/db"
	sqlc "example/db/sqlc"
	sqlc2 "example/db/sqlc2"
	tx "example/db/sqlc2"
	"example/models"
	"example/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"

	"github.com/gin-gonic/gin"
)

var validate = validator.New()

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		var user models.User
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, models.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}
		maskedAadhar, err := utils.MaskAadhaar(user.AadhaarNo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		newUser := sqlc.CreateusersParams{
			ServiceId:         sql.NullInt64{Int64: user.ServiceId, Valid: true},
			Username:          sql.NullString{String: user.Username, Valid: true},
			FirstName:         sql.NullString{String: user.FirstName, Valid: true},
			LastName:          sql.NullString{String: user.LastName, Valid: true},
			Password:          sql.NullString{String: user.Password, Valid: true},
			EncryptedPasswrod: sql.NullString{String: utils.HashPasswordWithSalt(user.Password), Valid: true},
			Address:           sql.NullString{String: user.Address, Valid: true},
			PhoneNo:           sql.NullString{String: user.PhoneNo, Valid: true},
			AadhaarNo:         sql.NullString{String: maskedAadhar, Valid: true},
			IsActive:          sql.NullInt32{Int32: 0, Valid: true},
			Status:            sql.NullInt32{Int32: 1, Valid: true},
		}
		lastId, err := db.Query.Createusers(ctx, newUser, int(user.IsUser))
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		userCreated, err := db.Query.Getusers(ctx, sqlc.GetusersParams{ID: int32(lastId), IsUser: int64(user.IsUser)})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		res := utils.SendSms(user.PhoneNo, "")
		phoneNo, _ := utils.MaskAadhaar(userCreated.PhoneNo.String)
		if res["status"] == "pending" {
			c.JSON(http.StatusOK, models.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": fmt.Sprintf("Code send at Phone Number  %v ", phoneNo)}})
			return
		} else {
			c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
	}
}
func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		var user models.LoginUser
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, models.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		userGet, err := db.Query.Getusers(ctx, sqlc.GetusersParams{PhoneNo: sql.NullString{String: user.PhoneNo, Valid: true}, IsUser: int64(user.IsUser)})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		if userGet.EncryptedPasswrod.String == utils.HashPasswordWithSalt(user.Password) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": models.User{ID: userGet.ID, IsUser: user.IsUser},
				"exp": time.Now().Add(time.Hour * 2).Unix(),
			})

			// Sign and get the complete encoded token as a string using the secret
			tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to create token",
				})
				return
			}

			// Send it back
			c.SetSameSite(http.SameSiteLaxMode)
			c.SetCookie("Authorization", tokenString, 3600*2, "", "", false, true)
			c.JSON(http.StatusCreated, models.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": fmt.Sprintf("User %v Logged In", userGet.Username.String), "Authorization": tokenString}})
		} else {
			c.JSON(http.StatusCreated, models.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": "Phone No Or Password Incorrect"}})

		}
	}
}

func GetOTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		var accountInfo models.AccountInfo
		defer cancel()
		// username := c.Param("username")
		// username = utils.Decoded(username)
		//validate the request body
		if err := c.BindJSON(&accountInfo); err != nil {
			c.JSON(http.StatusBadRequest, models.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		user, err := db.Query.Getusers(ctx, sqlc.GetusersParams{PhoneNo: sql.NullString{String: accountInfo.Email_Mobile_No, Valid: true}, IsUser: int64(accountInfo.IsUser)})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		if user.ID != 0 {
			res := utils.SendSms(accountInfo.Email_Mobile_No, "")
			if res["status"] == "pending" {
				c.JSON(http.StatusOK, models.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": fmt.Sprintf("Code send at Phone Number  %v ", accountInfo.Email_Mobile_No)}})
				return
			} else {
				c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}
	}
}
func SubmitOTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		var accountInfo models.AccountInfo
		defer cancel()
		// username := c.Param("username")
		// username = utils.Decoded(username)
		//validate the request body
		if err := c.BindJSON(&accountInfo); err != nil {
			c.JSON(http.StatusBadRequest, models.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		user, err := db.Query.Getusers(ctx, sqlc.GetusersParams{PhoneNo: sql.NullString{String: accountInfo.Email_Mobile_No, Valid: true}, IsUser: accountInfo.IsUser})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		if user.ID != 0 {
			res := utils.SendSms(accountInfo.Email_Mobile_No, accountInfo.Code)
			if res["status"] == "approved" {
				updatedUser := sqlc.UpdateusersParams{
					Username:          user.Username,
					FirstName:         user.FirstName,
					LastName:          user.LastName,
					PrePassword:       user.Password,
					PhoneNo:           user.PhoneNo,
					Password:          user.Password,
					EncryptedPasswrod: sql.NullString{String: utils.HashPasswordWithSalt(user.Password.String), Valid: true},
					Username_2:        sql.NullString{String: user.Username.String, Valid: true},
				}
				_, err = db.Query.Updateusers(ctx, updatedUser)
				if err != nil {
					c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
					return
				}
				c.JSON(http.StatusOK, models.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": fmt.Sprintf("Code Verified ")}})
				return
			} else {
				c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "error"}})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
	}
}
func ForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		var forgotPass models.ForgotPass
		defer cancel()
		//validate the request body
		if err := c.BindJSON(&forgotPass); err != nil {
			c.JSON(http.StatusBadRequest, models.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		user, err := db.Query.Getusers(ctx, sqlc.GetusersParams{PhoneNo: sql.NullString{String: forgotPass.Email_Mobile_No, Valid: true}, IsUser: forgotPass.IsUser})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		//use the validator library to validate required fields
		if validationErr := validate.Struct(&forgotPass); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}
		updatedUser := sqlc.UpdateusersParams{
			Username:          user.Username,
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			PrePassword:       user.Password,
			PhoneNo:           user.PhoneNo,
			Password:          sql.NullString{String: forgotPass.Password, Valid: true},
			EncryptedPasswrod: sql.NullString{String: utils.HashPasswordWithSalt(forgotPass.Password), Valid: true},
			Username_2:        sql.NullString{String: user.Username.String, Valid: true},
			IsUser:            forgotPass.IsUser,
		}
		_, err = db.Query.Updateusers(ctx, updatedUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, models.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": fmt.Sprintf("User %v Password Updated", user.Username.String)}})
	}
}
func CreateAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		arg := sqlc2.CreateAccountParams{
			Owner:   "as",
			Balance: 123,
		}
		err := db.Query2.CreateAccount(context.Background(), arg)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}
func GetAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		act, err := db.Query2.GetAccount(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("", act)
	}
}
func GetAccountEntries() {
	act, err := db.Query2.ScoreAndTests(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("", act)
}
func GetListAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		act, err := db.Query2.GetListAccount(context.Background(), sqlc2.GetListAccountParams{Limit: 1, Offset: 1})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("", act)
	}
}
func DeleteAccountById() {
	err := db.Query2.DeleteAccountById(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
func CreateData() {
	store := tx.NewStore(db.DB)
	store.TransferTx(context.Background(), sqlc2.TransferTxParams{
		FromAccountID: 19,
		ToAccountID:   19,
		Amount:        10,
	})
}
