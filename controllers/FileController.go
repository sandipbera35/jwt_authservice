package controllers

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	"github.com/sandipbera35/jwt_authservice/database"
	"github.com/sandipbera35/jwt_authservice/models"
)

func AddUploadProfilePic(c *fiber.Ctx) error {
	profile, ErrC := GetUserFromToken(c.Get("Authorization"))

	if ErrC != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrUnauthorized.Code,
			"message": "Unauthorized",
			"data":    nil,
		})
	}
	fh, e := c.FormFile("profile_pic")

	if e != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to upload profile picture",
			"data":    nil,
		})
	}

	f, errF := fh.Open()
	if errF != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to upload profile picture",
			"data":    nil,
		})
	}

	filePath := os.Getenv("STORE_PATH") + fh.Filename
	fc, errfc := os.Create(filePath)
	_, _ = io.Copy(fc, f)
	if errfc != nil {
		return fmt.Errorf("could not save file data: %v", errfc)
	}

	defer func() {
		f.Close()
		fc.Close()
		os.Remove(filePath)
	}()

	if fh.Size < 200 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "File size must be greater than 200 bytes",
			"data":    nil,
		})
	}

	mimeType := fh.Header.Get("Content-Type")
	if !strings.HasPrefix(mimeType, "image/") {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "File must be an image",
			"data":    nil,
		})
	}

	var fileQ models.File

	fQ := database.Connect.Where("user_id = ?", profile.ID).First(&fileQ)

	store := models.Store{
		EndPoint:   "localhost:9000",
		AccessId:   "minioadmin",
		AccessPass: "minioadmin",
		UseSSL:     false,
	}
	storepath := ""
	var fileguid uuid.UUID
	if fQ.RowsAffected > 0 {
		storepath = "profile" + "/" + profile.ID.String() + "/" + fileQ.ID.String() + "/" + fh.Filename
		fileguid = fileQ.ID

	} else {

		fileguid = uuid.New()
		storepath = "profile" + "/" + profile.ID.String() + "/" + fileguid.String() + "/" + fh.Filename

	}

	file := new(models.File)
	file.ID = fileguid
	file.Type = "profile"
	file.FileName = fh.Filename
	file.Size = fh.Size
	file.MimeType = mimeType
	file.FileGUID = fileguid.String()
	file.UserID = profile.ID
	file.Path = storepath

	errMFU := store.Upload(os.Getenv("MINIO_BUCKET"), storepath, filePath, mimeType)

	if errMFU {
		fmt.Println("Error in File server")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to upload profile picture",
			"data":    nil,
		})
	}

	if err := database.Connect.Create(&file).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to upload profile picture",
			"data":    nil,
		})
	}

	profile.Files = append(profile.Files, *file)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Profile picture uploaded successfully",
		"data":    profile,
	})
}

func GetProfilePic(c *fiber.Ctx) error {

	token := c.Get("Authorization")

	if token == "" {
		token = c.Query("token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  fiber.ErrUnauthorized.Code,
				"message": "Unauthorized",
				"data":    nil,
			})
		}
	}

	profile, ErrC := GetUserFromToken(token)
	if ErrC != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrUnauthorized.Code,
			"message": "Unauthorized",
			"data":    nil,
		})
	}

	var file models.File

	fileQ := database.Connect.Where("user_id = ?", profile.ID).Where("type = ?", "profile").Preload("User").First(&file)

	if fileQ.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile picture not found",
			"data":    nil,
		})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to get profile picture",
			"data":    nil,
		})
	}

	// stream file with fiber
	objectInfo, err := obj.Stat()
	if err != nil {
		log.Println("Error getting object info:", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error getting object info")
	}

	// Set headers for the response
	c.Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", file.FileName))
	c.Set("Content-Type", objectInfo.ContentType)
	c.Set("Content-Length", fmt.Sprintf("%d", objectInfo.Size))

	// Stream the object to the client
	c.SendStream(obj, int(file.Size))

	return nil

}

func UpdateProfileDetails(c *fiber.Ctx) error {

	token := c.Get("Authorization")
	// fmt.Printf("token: %v\n", token)
	claims, err := VerifyJWT(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	var profile models.User
	// fmt.Printf("profile: %v\n", claims)
	findQ := database.Connect.Where("id = ?", claims.UserId).First(&profile)

	if findQ.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile not found",
			"data":    nil,
		})
	}

	var usermodel models.UserUiModel
	if err := c.BodyParser(&usermodel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	fmt.Printf("usermodel: %v\n", usermodel)

	profile.FirstName = usermodel.FirstName
	profile.LastName = usermodel.LastName
	profile.Gender = usermodel.Gender
	profile.BirthDate = usermodel.BirthDate

	saveQ := database.Connect.Where("id = ?", claims.UserId).Save(&profile)

	if saveQ.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to update profile details",
			"data":    nil,
		})
	}

	var file models.File
	fileQ := database.Connect.Where("user_id = ?", claims.UserId).Where("type = ?", "profile").First(&file)
	if fileQ.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile picture not found",
			"data":    nil,
		})
	}

	var ispub bool = false

	if usermodel.ProfilePicStatus {
		ispub = true
	}

	// file.IsPublic = ispub
	fileSaveQ := database.Connect.Model(models.File{}).Where("id = ?", file.ID).UpdateColumn("is_public", ispub)

	if fileSaveQ.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to update profile details",
			"data":    nil,
		})
	}

	fileQQ := database.Connect.Where("user_id = ?", claims.UserId).Where("type = ?", "profile").First(&file)
	if fileQQ.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Internal server errors",
			"data":    nil,
		})
	}
	profile.Files = append(profile.Files, file)
	return c.Status(fiber.StatusOK).JSON(profile)
}
func GetPublicProfilePicById(c *fiber.Ctx) error {

	fileid := c.Query("file_id")

	if fileid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "file_id is required",
			"data":    nil,
		})
	}

	var file models.File

	fileQ := database.Connect.Where("id = ?", fileid).Where("type = ?", "profile").Where("is_public = ?", true).Preload("User").Preload("User").First(&file)

	if fileQ.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile picture not found",
			"data":    nil,
		})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to get profile picture",
			"data":    nil,
		})
	}

	// stream file with fiber
	objectInfo, err := obj.Stat()
	if err != nil {
		log.Println("Error getting object info:", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Error getting object info")
	}

	// Set headers for the response
	c.Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", file.FileName))
	c.Set("Content-Type", objectInfo.ContentType)
	c.Set("Content-Length", fmt.Sprintf("%d", objectInfo.Size))

	// Stream the object to the client
	c.SendStream(obj, int(file.Size))

	return nil

}
