package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"time"

	sqlc "jwt-go/db/sqlc"
	"jwt-go/initializers"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	//get email and password from req body
	var body sqlc.QueryAddUserParams

	err := c.ShouldBind(&body)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "we are not able to bind request body",
			"error":   err})
		return
	}

	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.UserPassword.String), 10)

	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{
			"message": "error in bcrypt",
			"error":   err})
		return
	}
	user := sqlc.QueryAddUserParams{Email: body.Email, UserPassword: sql.NullString{string(hash), true}}
	user.Created = sql.NullTime{time.Now(), true}
	user.Updated = sql.NullTime{time.Now(), true}
	db := sqlc.New(initializers.DB)
	err = db.QueryAddUser(context.Background(), user)
	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{
			"message": "error in creating user",
			"error":   err})
		return
	}

	//create the user

	//Respond
	c.IndentedJSON(http.StatusCreated, gin.H{
		"messsage": "user sucessfully created"})

}

func Login(c *gin.Context) {
	//get the email and password
	var body struct {
		Email    string
		Password string
	}
	err := c.ShouldBind(&body)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "we are not able to bind request body",
			"error":   err})
		return
	}

	//look up requested user
	db := sqlc.New(initializers.DB)
	result, err := db.QueryCheckUserByEmail(context.Background(), body.Email)
	if err != nil || result.ID == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "invalid email or password",
			"error":   err})
		return
	}

	//compare sent pass with saved user hash pass

	err = bcrypt.CompareHashAndPassword([]byte(result.UserPassword.String), []byte(body.Password))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "invalid email or password",
			"error":   err})
		return
	}

	//generate a jwt token

	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": result.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "failed to create token",
			"error":   err})
		return

	}
	// send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization",tokenString,3600*24*30,"","",false,true)
	c.IndentedJSON(http.StatusOK, gin.H{})
}


func Validate(c *gin.Context){

user,_:=c.Get("user")
	c.JSON(http.StatusOK,gin.H{"message":user})
}

func Logout(c *gin.Context){
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization","",-1,"","",false,true)

	c.JSON(http.StatusOK,gin.H{"message":"Log out sucessfully"})
}
