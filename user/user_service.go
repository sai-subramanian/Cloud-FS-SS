package user

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sai-subramanian/21BCE0040_Backend.git/configl"
	"github.com/sai-subramanian/21BCE0040_Backend.git/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context){
	//get the email and password off body

	costStr  := os.Getenv("hashCost")
	cost, err := strconv.Atoi(costStr)
	if err!= nil{
		c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to parse hash cost",
        })
        return
    }
	var body struct {
		Email string
		Password string
		Name string
		PhoneNumber string

	}

	if c.Bind(&body) != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	} 
	//hash the password
	hash, err :=	bcrypt.GenerateFromPassword([]byte(body.Password),cost)

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}


	//create the user
	user := models.User{Email: body.Email, Password: string(hash)}
	//tem during next push - to be converted to obj.assing() for scalability
	user.Name = body.Name
    user.PhoneNumber = body.PhoneNumber
    
	result := configl.DB.Create(&user)

	if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})

		return
	}


	//response
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
		"user":user,
	})
}

func GetAllUsers(c *gin.Context){
	var users []models.User
	configl.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func Login(c *gin.Context){
	//get the email and pass off the body
	var body struct {
		Email string
		Password string
	}

	if c.Bind(&body) != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	} 
	//look up the requested user
	var user models.User
	configl.DB.First(&user, "email = ?", body.Email)
	fmt.Println("Queried user:", user)

	if user.ID == 0{
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid email ",
		})
		return
	}

	//compare sent password with saved in password hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil{
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": " password ",
		})
		return
	}
	//generate a jwt token


	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),

	})
	
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("srirams-secret1912232480823ifbefw08fhu"))

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	//send it back
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})

}

func Validate(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"message": "Validated!!",
	})
	fmt.Print("Logged in")
}