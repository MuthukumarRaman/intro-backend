package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/api/option"
)

func UploadFileToGoogleCloudStorage(bucketName, objectName, filePath string) error {
	ctx := context.Background()

	// Create a new storage client
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Upload to GCS
	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if _, err = wc.Write([]byte("This is a test file upload.")); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	fmt.Printf("File %s uploaded to bucket %s as %s\n", filePath, bucketName, objectName)
	return nil

}

func NewUploadFileToGoogleCloudStorage(bucketName, objectName, filePath string) error {
	ctx := context.Background()

	// Create a new storage client
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("introme-firebase-service.json"))
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Upload to GCS
	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)

	// Copy file content to writer
	if _, err = io.Copy(wc, file); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	fmt.Printf("File %s uploaded to bucket %s as %s\n", filePath, bucketName, objectName)
	return nil
}

func FileUploadToGoogle(c *fiber.Ctx) error {
	fileCategory := "userFiles"
	token := GetUserTokenValue(c)
	refId := token.UserId
	folderName := path.Join(fileCategory, refId)

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": "Failed to parse multipart form: " + err.Error(),
		})
	}

	files := form.File["file"]
	if len(files) == 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": "No file uploaded under 'file' key",
		})
	}

	var result []interface{}

	for _, file := range files {
		// Generate a timestamped filename
		ext := path.Ext(file.Filename)
		name := strings.TrimSuffix(file.Filename, ext)
		timestampedName := fmt.Sprintf("%s__%s%s", name, time.Now().Format("2006-01-02-15-04-05"), ext)

		// Save file temporarily
		tmpPath := path.Join(os.TempDir(), timestampedName)
		if err := c.SaveFile(file, tmpPath); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"errors": "Failed to save uploaded file: " + err.Error(),
			})
		}

		// Upload to GCS
		objectPath := path.Join(folderName, timestampedName)
		fmt.Println(objectPath, tmpPath)
		isError := NewUploadFileToGoogleCloudStorage(
			"introme-webapp-user-storage", // bucket name
			objectPath,                    // object name
			tmpPath,                       // local file path
		)

		if isError != nil {
			// log.Println(errMsg)
			return c.Status(422).JSON(fiber.Map{"errors": isError.Error()})
		}

		// Store file metadata
		id := uuid.New().String()
		apiResponse := bson.M{
			"_id":          id,
			"ref_id":       refId,
			"uploaded_by":  refId,
			"folder":       fileCategory,
			"file_name":    file.Filename,
			"storage_name": objectPath,
			"size":         file.Size,
		}
		InsertData(c, "user_files", apiResponse)
		result = append(result, apiResponse)

		// Clean up temp file
		os.Remove(tmpPath)
	}

	return SuccessResponse(c, result)
}

func GetSignedURL(c *fiber.Ctx) error {
	filename := c.Query("filename")
	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing filename query parameter",
		})
	}

	url, err := generateSignedURL(filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "Could not generate signed URL",
			"detail": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"url": url})
}

func generateSignedURL(objectName string) (string, error) {
	// ctx := context.Background()

	// Load service account credentials from JSON file
	credsJSON, err := os.ReadFile("introme-firebase-service.json")
	if err != nil {
		return "", fmt.Errorf("failed to read service account file: %w", err)
	}

	// Parse JSON to extract required fields
	var creds struct {
		ClientEmail string `json:"client_email"`
		PrivateKey  string `json:"private_key"`
	}
	if err := json.Unmarshal(credsJSON, &creds); err != nil {
		return "", fmt.Errorf("failed to unmarshal service account: %w", err)
	}

	opts := &storage.SignedURLOptions{
		Method:         "GET",
		Expires:        time.Now().Add(8 * time.Hour),
		Scheme:         storage.SigningSchemeV4,
		GoogleAccessID: creds.ClientEmail,
		PrivateKey:     []byte(creds.PrivateKey),
	}

	// Generate signed URL
	url, err := storage.SignedURL("introme-webapp-user-storage", objectName, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}

func DeleteFileFromGCS(objectName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := storage.NewClient(ctx, option.WithCredentialsFile("introme-firebase-service.json"))
	if err != nil {
		return err
	}
	defer client.Close()

	object := client.Bucket("introme-webapp-user-storage").Object(objectName)
	return object.Delete(ctx)
}
