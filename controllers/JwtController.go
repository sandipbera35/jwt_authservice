package controllers

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sandipbera35/jwt_authservice/database"
	"github.com/sandipbera35/jwt_authservice/models"
)

var jwtSecret = []byte(os.Getenv("JWTSECRET"))

type CustomClaims struct {
	UserId       string       `json:"user_id"`
	FirstName    string       `json:"first_name"`
	LastName     string       `json:"last_name"`
	Gender       string       `json:"gender"`
	BirthDate    string       `json:"birth_date"`
	UserName     string       `json:"user_name"`
	MobileNo     string       `json:"mobile_no"`
	EmailID      string       `json:"email_id"`
	ProfileImage *models.File `json:"profile_image"`
	DisplayIamge *models.File `json:"display_image"`
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
