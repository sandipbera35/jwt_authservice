package database

import (
	"github.com/google/uuid"
	"github.com/sandipbera35/jwt_authservice/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateAdmin() {

	findQ := Connect.Where("email_id = ?", "admin@admin.admin").Find(map[string]interface{}{})
	if findQ.RowsAffected == 0 {
		var user models.User
		var admin models.Admin
		user.ID = uuid.New()

		user.UserName = "admin"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {

			panic("user not created")
		}
		user.UserPassword = string(hashedPassword)
		createAdmin := Connect.Create(&user)

		if createAdmin.Error != nil {
			panic("user not created")
		}
		if createAdmin.RowsAffected == 0 {
			panic("user not created")
		}

		admin.UserID = user.ID
		admin.Role = "ADMIN"
		admin.ID = uuid.New()

		Connect.Create(&admin)
	}

}
