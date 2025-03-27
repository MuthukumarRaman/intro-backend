package helper

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
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
