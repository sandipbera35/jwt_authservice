package controllers

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	"github.com/sandipbera35/jwt_authservice/database"
	"github.com/sandipbera35/jwt_authservice/models"
)

func AddUploadProfilePic(c *fiber.Ctx) {
	profile, ErrC := GetUserFromToken(c.Get("Authorization"))

	if ErrC != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrUnauthorized.Code,
			"message": "Unauthorized",
			"data":    nil,
		})
		return
	}
	fh, e := c.FormFile("profile_pic")

	if e != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to upload profile picture",
			"data":    nil,
		})
		return
	}

	f, errF := fh.Open()
	if errF != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to upload profile picture",
			"data":    nil,
		})
		return
	}

	filePath := os.Getenv("STORE_PATH") + fh.Filename
	fc, errfc := os.Create(filePath)
	_, _ = io.Copy(fc, f)
	if errfc != nil {
		c.JSON("could not save file data")
		return
	}

	defer func() {
		f.Close()
		fc.Close()
		os.Remove(filePath)
	}()

	if fh.Size < 200 {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "File size must be greater than 200 bytes",
			"data":    nil,
		})
		return
	}

	mimeType := fh.Header.Get("Content-Type")
	if !strings.HasPrefix(mimeType, "image/") {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "File must be an image",
			"data":    nil,
		})
		return
	}

	var fileQ models.File

	fQ := database.Connect.Where("user_id = ?", profile.ID).Where("type = ?", "profile").Find(&fileQ)

	store := models.Store{
		EndPoint:   "localhost:9000",
		AccessId:   "minioadmin",
		AccessPass: "minioadmin",
		UseSSL:     false,
	}
	storepath := ""
	var fileguid uuid.UUID
	var file models.File

	if fQ.RowsAffected > 0 {
		fileguid = fileQ.ID
		storepath = "profile" + "/" + profile.ID.String() + "/" + fileQ.ID.String() + "/" + fileguid.String() + filepath.Ext(fh.Filename)

	} else {

		fileguid = uuid.New()
		file.ID = fileguid
		file.CreatedAt = time.Now().UTC()
		storepath = "profile" + "/" + profile.ID.String() + "/" + fileguid.String() + "/" + fileguid.String() + filepath.Ext(fh.Filename)

	}

	file.Type = "profile"
	file.UpdatedAt = time.Now().UTC()
	file.FileName = fh.Filename
	file.Size = fh.Size
	file.MimeType = mimeType
	file.UserID = profile.ID
	file.Path = storepath

	file.IsPublic = fileQ.IsPublic

	errMFU := store.Upload(os.Getenv("MINIO_BUCKET"), storepath, filePath, mimeType)

	if errMFU {
		fmt.Println("Error in File server")
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to upload profile picture",
			"data":    nil,
		})
		return
	}
	if fQ.RowsAffected > 0 {
		if err := database.Connect.Where("user_id = ?", profile.ID).Updates(&file).Error; err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  fiber.StatusInternalServerError,
				"message": "Failed to upload profile picture",
				"data":    nil,
			})
			return
		}
	} else {
		if err := database.Connect.Create(&file).Error; err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  fiber.StatusInternalServerError,
				"message": "Failed to upload profile picture",
				"data":    nil,
			})
			return
		}

	}

	database.Connect.Model(models.File{}).Where("id = ?", file.ID).Where("type = ?", "profile").Find(&file)
	file.ID = fileguid
	// profile.Files = append(profile.Files, file)

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Profile picture uploaded successfully",
		"data":    file,
	})
}

func GetProfilePic(c *fiber.Ctx) {

	token := c.Get("Authorization")

	if token == "" {
		token = c.Query("token")
		if token == "" {
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  fiber.ErrUnauthorized.Code,
				"message": "Unauthorized",
				"data":    nil,
			})
			return
		}
	}

	profile, ErrC := GetUserFromToken(token)
	if ErrC != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrUnauthorized.Code,
			"message": "Unauthorized",
			"data":    nil,
		})
		return
	}

	var file models.File

	fileQ := database.Connect.Where("user_id = ?", profile.ID).Where("type = ?", "profile").First(&file)

	if fileQ.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile picture not found",
			"data":    nil,
		})
		return
	}

	store := models.Store{
		EndPoint:   "localhost:9000",
		AccessId:   "minioadmin",
		AccessPass: "minioadmin",
		UseSSL:     false,
	}

	obj, errMFS := store.Stream(os.Getenv("MINIO_BUCKET"), file.Path)

	if errMFS {
		fmt.Println("Error in File server")
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to get profile picture",
			"data":    nil,
		})
		return
	}

	// stream file with fiber
	objectInfo, err := obj.Stat()
	if err != nil {
		log.Println("Error getting object info:", err)
		c.Status(fiber.StatusInternalServerError).JSON("Error getting object info")
		return
	}

	// Set headers for the response
	c.Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", file.FileName))
	c.Set("Content-Type", objectInfo.ContentType)
	c.Set("Content-Length", fmt.Sprintf("%d", objectInfo.Size))

	// Stream the object to the client
	c.SendStream(obj, int(file.Size))

}

func UpdateProfileDetails(c *fiber.Ctx) {

	token := c.Get("Authorization")
	// fmt.Printf("token: %v\n", token)
	claims, err := VerifyJWT(token)
	if err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
		return
	}
	var profile models.User
	// fmt.Printf("profile: %v\n", claims)
	findQ := database.Connect.Where("id = ?", claims.UserId).First(&profile)

	if findQ.Error != nil {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile not found",
			"data":    nil,
		})
		return
	}

	var usermodel models.UserUiModel
	if err := c.BodyParser(&usermodel); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
		return
	}

	fmt.Printf("usermodel: %v\n", usermodel)

	profile.FirstName = usermodel.FirstName
	profile.LastName = usermodel.LastName
	profile.Gender = usermodel.Gender
	profile.BirthDate = usermodel.BirthDate

	saveQ := database.Connect.Where("id = ?", claims.UserId).Save(&profile)

	if saveQ.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to update profile details",
			"data":    nil,
		})
		return
	}

	var file models.File
	fileQ := database.Connect.Where("user_id = ?", claims.UserId).Where("type = ?", "profile").First(&file)
	if fileQ.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile picture not found",
			"data":    nil,
		})
		return
	}

	var ispub bool = false

	if usermodel.ProfilePicStatus {
		ispub = true
	}

	// file.IsPublic = ispub
	fileSaveQ := database.Connect.Model(models.File{}).Where("id = ?", file.ID).UpdateColumn("is_public", ispub)

	if fileSaveQ.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to update profile details",
			"data":    nil,
		})
		return
	}

	fileQQ := database.Connect.Where("user_id = ?", claims.UserId).Where("type = ?", "profile").First(&file)
	if fileQQ.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Internal server errors",
			"data":    nil,
		})
		return
	}
	profile.Files = append(profile.Files, file)
	c.Status(fiber.StatusOK).JSON(profile)
}
func GetPublicProfilePicById(c *fiber.Ctx) {

	fileid := c.Query("file_id")

	if fileid == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "file_id is required",
			"data":    nil,
		})
		return
	}

	var file models.File

	fileQ := database.Connect.Where("id = ?", fileid).Where("type = ?", "profile").Where("is_public = ?", true).Find(&file)

	if fileQ.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile picture not found",
			"data":    nil,
		})
		return
	}

	store := models.Store{
		EndPoint:   "localhost:9000",
		AccessId:   "minioadmin",
		AccessPass: "minioadmin",
		UseSSL:     false,
	}

	obj, errMFS := store.Stream(os.Getenv("MINIO_BUCKET"), file.Path)

	if errMFS {
		fmt.Println("Error in File server")
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to get profile picture",
			"data":    nil,
		})
		return
	}

	// stream file with fiber
	objectInfo, err := obj.Stat()
	if err != nil {
		log.Println("Error getting object info:", err)
		c.Status(fiber.StatusInternalServerError).JSON("Error getting object info")
		return
	}

	// Set headers for the response
	c.Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", file.FileName))
	c.Set("Content-Type", objectInfo.ContentType)
	c.Set("Content-Length", fmt.Sprintf("%d", objectInfo.Size))

	// Stream the object to the client
	c.SendStream(obj, int(file.Size))

}
