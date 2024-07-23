package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber"
	"github.com/golang-jwt/jwt/v4"

	"github.com/sandipbera35/jwt_authservice/database"
	"github.com/sandipbera35/jwt_authservice/models"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	data := new(struct {
		EmailID  string `json:"email_id"`
		Password string `json:"password"`
	})

	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	var user models.User
	if err := database.Connect.Where("email_id = ?", data.EmailID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(data.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}
	ExpiresAt := jwt.NewNumericDate(time.Now().Add(time.Hour * 1))
	IssuedAt := jwt.NewNumericDate(time.Now())
	token, err := GenerateJWT(user, ExpiresAt, IssuedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not login",
		})
	}

	return c.JSON(fiber.Map{
		"user_id":    user.ID,
		"expires_at": ExpiresAt,
		"issued_at":  IssuedAt,
		"token":      token,
	})
}

func GetProfile(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	claims, err := VerifyJWT(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	return c.Status(fiber.StatusOK).JSON(claims)
}

func GetUserFromToken(token string) (models.User, error) {
	claims, err := VerifyJWT(token)
	if err != nil {
		return models.User{}, err
	}
	var user models.User
	if err := database.Connect.Where("id = ?", claims.UserId).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil

}

var jwtSecret = []byte("supersecretkey")

type CustomClaims struct {
	UserId    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
	UserName  string `json:"user_name"`
	MobileNo  string `json:"mobile_no"`
	EmailID   string `json:"email_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userDetails models.User, ExpiresAt *jwt.NumericDate, IssuedAt *jwt.NumericDate) (string, error) {
	claims := CustomClaims{
		UserId:    userDetails.ID.String(),
		FirstName: userDetails.FirstName,
		LastName:  userDetails.LastName,
		Gender:    userDetails.Gender,
		BirthDate: userDetails.BirthDate.String(),
		UserName:  userDetails.UserName,
		MobileNo:  userDetails.MobileNo,
		EmailID:   userDetails.EmailID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: ExpiresAt,
			IssuedAt:  IssuedAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func VerifyJWT(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
