package service

import (
	"errors"

	"github.com/sandipbera35/jwt_authservice/database"
	"github.com/sandipbera35/jwt_authservice/models"
	"gorm.io/gorm"
)

type AdminService struct {
}

func (s *AdminService) GetAdminByUserId(userId string) (adminReturn models.Admin, IsFound bool, erRtnr error) {
	var admin models.Admin
	err := database.Connect.Where("user_id = ?", userId).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Admin{}, false, nil // Return empty admin if not found
		}
		return models.Admin{}, false, err // Return error for other issues
	}
	return admin, true, nil
}

func (s *AdminService) GetRoleByUserId(userId string) (string, bool, error) {
	var admin models.Admin
	err := database.Connect.Where("user_id = ?", userId).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", false, nil // Return empty role if not found
		}
		return "", false, err // Return error for other issues
	}
	return admin.Role, true, nil
}
