package controllers

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sandipbera35/jwt_authservice/database"
	"github.com/sandipbera35/jwt_authservice/models"
)

func AddUploadProfilePic(c *fiber.Ctx) error {
	profile, ErrC := GetUserFromToken(c.Get("Authorization"))

	if ErrC != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrUnauthorized.Code,
			"message": "Unauthorized",
			"data":    nil,
		})
		return nil
	}
	fh, e := c.FormFile("profile_pic")

	if e != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrBadRequest,
			"message": "Invalid Image Format",
			"data":    nil,
		})
		return nil
	}

	f, errF := fh.Open()
	if errF != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to upload profile picture",
			"data":    nil,
		})
		return nil
	}

	filePath := os.Getenv("STORE_PATH") + fh.Filename
	fc, errfc := os.Create(filePath)
	_, _ = io.Copy(fc, f)
	if errfc != nil {
		c.JSON("could not save file data")
		return nil
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
		return nil
	}

	mimeType := fh.Header.Get("Content-Type")
	if !strings.HasPrefix(mimeType, "image/") {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "File must be an image",
			"data":    nil,
		})
		return nil
	}

	var fileQ models.ProfileImage

	fQ := database.Connect.Where("user_id = ?", profile.ID).Find(&fileQ)

	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	minioAccessID := os.Getenv("MINIO_ACCESSID")
	minioAccessPass := os.Getenv("MINIO_ACCESSPASS")
	// minioUseSSL := os.Getenv("MINIO_USESSL")
	minioBucket := os.Getenv("MINIO_BUCKET")

	store := models.Store{
		EndPoint:   minioEndpoint,
		AccessId:   minioAccessID,
		AccessPass: minioAccessPass,
		UseSSL:     false,
	}
	storepath := ""
	var fileguid uuid.UUID
	var file models.ProfileImage

	if fQ.RowsAffected > 0 {
		fileguid = fileQ.ID
		file.CreatedAt = fileQ.CreatedAt
		storepath = "profile" + "/" + profile.ID.String() + "/" + fileQ.ID.String() + "/" + fileguid.String() + filepath.Ext(fh.Filename)

	} else {

		fileguid = uuid.New()
		file.ID = fileguid
		file.CreatedAt = time.Now().UTC()
		storepath = "profile" + "/" + profile.ID.String() + "/" + fileguid.String() + "/" + fileguid.String() + filepath.Ext(fh.Filename)

	}

	file.UpdatedAt = time.Now().UTC()
	file.FileName = fh.Filename
	file.Size = fh.Size
	file.Extension = filepath.Ext(fh.Filename)
	file.MimeType = mimeType
	file.UserID = profile.ID
	file.Path = storepath

	file.IsPublic = fileQ.IsPublic
	fmt.Printf("minioBucket: %v\n", minioBucket)
	errMFU := store.Upload(minioBucket, storepath, filePath, mimeType)

	if errMFU {
		fmt.Println("Error in File server : ", errMFU)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to upload profile picture",
			"data":    nil,
		})
		return nil
	}
	if fQ.RowsAffected > 0 {
		if err := database.Connect.Where("user_id = ?", profile.ID).Updates(&file).Error; err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  fiber.StatusInternalServerError,
				"message": "Failed to upload profile picture",
				"data":    nil,
			})
			return nil
		}
	} else {
		if err := database.Connect.Create(&file).Error; err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  fiber.StatusInternalServerError,
				"message": "Failed to upload profile picture",
				"data":    nil,
			})
			return nil
		}

	}

	database.Connect.Model(models.ProfileImage{}).Where("id = ?", file.ID).Find(&file)
	file.ID = fileguid
	// profile.Files = append(profile.Files, file)

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Profile picture uploaded successfully",
		"data":    file,
	})
	return nil
}
func AddUploadCoverPic(c *fiber.Ctx) error {
	profile, ErrC := GetUserFromToken(c.Get("Authorization"))

	if ErrC != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrUnauthorized.Code,
			"message": "Unauthorized",
			"data":    nil,
		})
		return nil
	}
	fh, e := c.FormFile("cover_pic")

	if e != nil {
		fmt.Printf("e.Error(): %v\n", e.Error())
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to upload cover picture",
			"data":    nil,
		})
		return nil
	}

	f, errF := fh.Open()
	if errF != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to upload cover picture",
			"data":    nil,
		})
		return nil
	}

	filePath := os.Getenv("STORE_PATH") + fh.Filename
	fc, errfc := os.Create(filePath)
	_, _ = io.Copy(fc, f)
	if errfc != nil {
		c.JSON("could not save file data")
		return nil
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
		return nil
	}

	mimeType := fh.Header.Get("Content-Type")
	if !strings.HasPrefix(mimeType, "image/") {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "File must be an image",
			"data":    nil,
		})
		return nil
	}

	var fileQ models.CoverImage

	fQ := database.Connect.Where("user_id = ?", profile.ID).Find(&fileQ)

	store := models.Store{
		EndPoint:   "localhost:9000",
		AccessId:   "minioadmin",
		AccessPass: "minioadmin",
		UseSSL:     false,
	}
	storepath := ""
	var fileguid uuid.UUID
	var file models.CoverImage

	if fQ.RowsAffected > 0 {
		fileguid = fileQ.ID
		file.CreatedAt = fileQ.CreatedAt
		storepath = "cover" + "/" + profile.ID.String() + "/" + fileQ.ID.String() + "/" + fileguid.String() + filepath.Ext(fh.Filename)

	} else {

		fileguid = uuid.New()
		file.ID = fileguid
		file.CreatedAt = time.Now().UTC()
		storepath = "cover" + "/" + profile.ID.String() + "/" + fileguid.String() + "/" + fileguid.String() + filepath.Ext(fh.Filename)

	}

	file.UpdatedAt = time.Now().UTC()
	file.FileName = fh.Filename
	file.Size = fh.Size
	file.Extension = filepath.Ext(fh.Filename)
	file.MimeType = mimeType
	file.UserID = profile.ID
	file.Path = storepath

	file.IsPublic = fileQ.IsPublic
	fmt.Printf("os.Getenv(\"MINIO_BUCKET\"): %v\n", os.Getenv("MINIO_BUCKET"))

	errMFU := store.Upload(os.Getenv("MINIO_BUCKET"), storepath, filePath, mimeType)

	if errMFU {
		fmt.Println("Error in File server")
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to upload profile picture",
			"data":    nil,
		})
		return nil
	}
	if fQ.RowsAffected > 0 {
		if err := database.Connect.Where("user_id = ?", profile.ID).Updates(&file).Error; err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  fiber.StatusInternalServerError,
				"message": "Failed to upload profile picture",
				"data":    nil,
			})
			return nil
		}
	} else {
		if err := database.Connect.Create(&file).Error; err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  fiber.StatusInternalServerError,
				"message": "Failed to upload profile picture",
				"data":    nil,
			})
			return nil
		}

	}

	database.Connect.Model(models.ProfileImage{}).Where("id = ?", file.ID).Find(&file)
	file.ID = fileguid
	// profile.Files = append(profile.Files, file)

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Profile picture uploaded successfully",
		"data":    file,
	})
	return nil
}

func GetProfilePic(c *fiber.Ctx) error {

	token := c.Get("Authorization")

	if token == "" {
		token = c.Query("token")
		if token == "" {
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  fiber.ErrUnauthorized.Code,
				"message": "Unauthorized",
				"data":    nil,
			})
			return nil
		}
	}

	profile, ErrC := GetUserFromToken(token)
	if ErrC != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrUnauthorized.Code,
			"message": "Unauthorized",
			"data":    nil,
		})
		return nil
	}

	if profile.ProfileImage == nil {

		f, _ := os.ReadFile("./defaultprofile.png")

		c.Set("Content-Disposition", "inline; filename=defaultprofile.png")
		c.Set("Content-Type", "image/png")
		c.Set("Content-Length", fmt.Sprintf("%v", len(f)))

		c.Status(fiber.StatusOK).SendFile("./defaultprofile.png", true)
		return nil

	}

	var file models.ProfileImage

	fileQ := database.Connect.Where("user_id = ?", profile.ID).Find(&file)

	if fileQ.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile picture not found",
			"data":    nil,
		})
		return nil
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
		return nil
	}

	// stream file with fiber
	objectInfo, err := obj.Stat()
	if err != nil {
		log.Println("Error getting object info:", err)
		c.Status(fiber.StatusInternalServerError).JSON("Error getting object info")
		return nil
	}

	// Set headers for the response
	c.Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", file.FileName))
	c.Set("Content-Type", objectInfo.ContentType)
	c.Set("Content-Length", fmt.Sprintf("%d", objectInfo.Size))

	// Stream the object to the client
	c.SendStream(obj, int(file.Size))
	return nil

}
func GetCoverPic(c *fiber.Ctx) error {

	token := c.Get("Authorization")

	if token == "" {
		token = c.Query("token")
		if token == "" {
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  fiber.ErrUnauthorized.Code,
				"message": "Unauthorized",
				"data":    nil,
			})
			return nil
		}
	}

	profile, ErrC := GetUserFromToken(token)
	if ErrC != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.ErrUnauthorized.Code,
			"message": "Unauthorized",
			"data":    nil,
		})
		return nil
	}

	if profile.CoverImage == nil {

		f, _ := os.ReadFile("./default_cover.jpg")

		c.Set("Content-Disposition", "inline; filename=default_cover.jpg")
		c.Set("Content-Type", "image/jpeg")
		c.Set("Content-Length", fmt.Sprintf("%v", len(f)))

		c.Status(fiber.StatusOK).SendFile("./default_cover.jpg", true)
		return nil

	}

	var file models.CoverImage

	fileQ := database.Connect.Where("user_id = ?", profile.ID).Find(&file)

	if fileQ.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Cover picture not found",
			"data":    nil,
		})
		return nil
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
			"message": "Failed to get cover picture",
			"data":    nil,
		})
		return nil
	}

	// stream file with fiber
	objectInfo, err := obj.Stat()
	if err != nil {
		log.Println("Error getting object info:", err)
		c.Status(fiber.StatusInternalServerError).JSON("Error getting object info")
		return nil
	}

	// Set headers for the response
	c.Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", file.FileName))
	c.Set("Content-Type", objectInfo.ContentType)
	c.Set("Content-Length", fmt.Sprintf("%d", objectInfo.Size))

	// Stream the object to the client
	c.SendStream(obj, int(file.Size))
	return nil

}

func GetPublicProfilePicById(c *fiber.Ctx) error {

	fileid := c.Query("file_id")

	if fileid == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "file_id is required",
			"data":    nil,
		})
		return nil
	}

	var file models.ProfileImage

	fileQ := database.Connect.Where("id = ?", fileid).Where("is_public = ?", true).Find(&file)

	if fileQ.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile picture not found",
			"data":    nil,
		})
		return nil
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
		return nil
	}

	// stream file with fiber
	objectInfo, err := obj.Stat()
	if err != nil {
		log.Println("Error getting object info:", err)
		c.Status(fiber.StatusInternalServerError).JSON("Error getting object info")
		return nil
	}

	// Set headers for the response
	c.Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", file.FileName))
	c.Set("Content-Type", objectInfo.ContentType)
	c.Set("Content-Length", fmt.Sprintf("%d", objectInfo.Size))

	// Stream the object to the client
	c.SendStream(obj, int(file.Size))
	return nil

}
func GetPublicCoverPicById(c *fiber.Ctx) error {

	fileid := c.Query("file_id")

	if fileid == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "file_id is required",
			"data":    nil,
		})
		return nil
	}

	var file models.CoverImage

	fileQ := database.Connect.Where("id = ?", fileid).Where("is_public = ?", true).Find(&file)

	if fileQ.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Cover picture not found",
			"data":    nil,
		})
		return nil
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
			"message": "Failed to get cover picture",
			"data":    nil,
		})
		return nil
	}

	// stream file with fiber
	objectInfo, err := obj.Stat()
	if err != nil {
		log.Println("Error getting object info:", err)
		c.Status(fiber.StatusInternalServerError).JSON("Error getting object info")
		return nil
	}

	// Set headers for the response
	c.Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", file.FileName))
	c.Set("Content-Type", objectInfo.ContentType)
	c.Set("Content-Length", fmt.Sprintf("%d", objectInfo.Size))

	// Stream the object to the client
	c.SendStream(obj, int(file.Size))
	return nil

}

func DeleteProfilePic(c *fiber.Ctx) error {

	// fileid := c.Query("file_id")
	// if fileid == "" {
	// 	c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"status":  fiber.StatusBadRequest,
	// 		"message": "file_id is required",
	// 		"data":    nil,
	// 	})
	// 	return
	// }
	profile, errC := GetUserFromToken(c.Get("Authorization"))

	if errC != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.ErrUnauthorized.Code,
			"message": "Unauthorized",
			"data":    nil,
		})
		return nil
	}
	if profile.ProfileImage == nil {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile picture not found",
			"data":    nil,
		})
		return nil
	}
	fileQ := database.Connect.Where("id = ?", profile.ProfileImage.ID).Delete(&models.ProfileImage{})

	if fileQ.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to delete profile picture",
			"data":    nil,
		})
		return nil
	}

	if fileQ.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Profile picture not found",
			"data":    nil,
		})
		return nil
	}

	store := models.Store{
		EndPoint:   "localhost:9000",
		AccessId:   "minioadmin",
		AccessPass: "minioadmin",
		UseSSL:     false,
	}

	store.Delete(os.Getenv("MINIO_BUCKET"), profile.ProfileImage.Path)

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Profile picture deleted successfully",
		"data":    nil,
	})
	return nil

}
func DeleteCoverPic(c *fiber.Ctx) error {

	// fileid := c.Query("file_id")
	// if fileid == "" {
	// 	c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"status":  fiber.StatusBadRequest,
	// 		"message": "file_id is required",
	// 		"data":    nil,
	// 	})
	// 	return
	// }
	profile, errC := GetUserFromToken(c.Get("Authorization"))

	if errC != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.ErrUnauthorized.Code,
			"message": "Unauthorized",
			"data":    nil,
		})
		return nil
	}
	if profile.CoverImage == nil {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Cover picture not found",
			"data":    nil,
		})
		return nil
	}
	fileQ := database.Connect.Where("id = ?", profile.CoverImage.ID).Delete(&models.CoverImage{})

	if fileQ.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to delete cover picture",
			"data":    nil,
		})
		return nil
	}

	if fileQ.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Cover picture not found",
			"data":    nil,
		})
		return nil
	}
	store := models.Store{
		EndPoint:   "localhost:9000",
		AccessId:   "minioadmin",
		AccessPass: "minioadmin",
		UseSSL:     false,
	}

	store.Delete(os.Getenv("MINIO_BUCKET"), profile.CoverImage.Path)

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Cover picture deleted successfully",
		"data":    nil,
	})
	return nil

}
