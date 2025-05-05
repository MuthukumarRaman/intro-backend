package helper

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"introme-api/pkg/shared/database"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func initS3() (*s3.S3, string) {
	var api_key = GetenvStr("S3_API_KEY")
	var secret = GetenvStr("S3_API_SECRET")
	var endpoint = GetenvStr("S3_ENDPOINT")
	var region = GetenvStr("S3_REGION")
	var bucket = GetenvStr("S3_BUCKET")

	var s3Config = &aws.Config{
		Credentials: credentials.NewStaticCredentials(api_key, secret, ""),
		Endpoint:    aws.String(endpoint),
		Region:      aws.String(region),
	}

	// var newSession = session.New(s3Config)

	// Create a new session using NewSession
	var newSession = session.Must(session.NewSession(s3Config))
	var s3Client = s3.New(newSession)
	return s3Client, bucket
}

func UploadFile(fileIn *multipart.FileHeader, key string) (bool, string) {
	s3Client, bucket := initS3()
	var errContent string
	var isErrExist bool
	file, err := fileIn.Open()
	if err != nil {
		isErrExist = true
		errContent = err.Error()
		return isErrExist, errContent
	}
	defer file.Close()
	buf := bytes.NewBuffer(nil)
	_, err = buf.ReadFrom(file)
	if err != nil {
		isErrExist = true
		errContent = err.Error()
		return isErrExist, errContent
	}
	_, err = s3Client.PutObjectWithContext(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
		ACL:    aws.String("public-read"), ContentType: aws.String(filepath.Ext(fileIn.Filename)),
	})
	if err != nil {
		isErrExist = true
		errContent = err.Error()
		return isErrExist, errContent
	}
	return isErrExist, errContent
}

func FileUpload(c *fiber.Ctx) error {

	fileCategory := c.Params("folder") // c.Params("category")

	// request, err := c.MultipartForm()
	// if err != nil {
	// 	return c.Status(422).JSON(fiber.Map{"errors": err.Error()})
	// }

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": "Failed to parse multipart form: " + err.Error()})
	}
	// Get the files under the "file" key
	files, ok := form.File["file"]
	if !ok || len(files) == 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": "No file uploaded under 'file' key"})
	}
	token := GetUserTokenValue(c)
	refId := token.UserId
	folderName := fileCategory + "/" + refId

	var result []interface{}
	for _, file := range files {
		// fmt.Println(file)
		fileNew := file
		fileExtn := filepath.Ext(fileNew.Filename)
		fileName := strings.TrimSuffix(fileNew.Filename, fileExtn)
		fileName = fileName + "__" + time.Now().Format("2006-01-02-15-04-05") + fileExtn
		isErrorExist, errContent := UploadFile(fileNew, folderName+"/"+fileName)
		if isErrorExist {
			log.Print(errContent)
			return c.Status(422).JSON(fiber.Map{"errors": errContent})
		}
		//Save file name to the DB
		id := uuid.New().String()
		storageName := folderName + "/" + fileName
		apiResponse := bson.M{"_id": id, "ref_id": refId, "uploaded_by": token.UserId, "folder": fileCategory, "file_name": fileNew.Filename, "storage_name": storageName, "size": fileNew.Size} // "extn": filepath.Ext(fileName),

		//S3
		InsertData(c, "user_files", apiResponse) //Without Struct
		result = append(result, apiResponse)

		// fmt.Printf("User Id %s, File Name:%s, Size:%d", "test", fileName, file[0].Size)
	}
	return SuccessResponse(c, result)
}

// *S3
func InsertData(c *fiber.Ctx, collectionName string, data interface{}) error {
	response, err := database.GetConnection().Collection(collectionName).InsertOne(ctx, data)
	if err != nil {
		return BadRequest(err.Error())
	}
	return SuccessResponse(c, response)
}

func GetAllFileDetails(c *fiber.Ctx) error {
	orgId := c.Get("OrgId")
	if orgId == "" {
		return BadRequest("Organization Id missing")
	}
	fileCategory := c.Params("folder")
	//status := c.Params("status")

	page := c.Params("page")
	if page == "" {
		page = "0"
	}
	limit := c.Params("limit")
	if limit == "" {
		limit = "25"
	}
	query := bson.M{"folder": fileCategory}
	response, err := GetQueryResult("user_files", query, Page(page), Limit(limit), nil)
	if err != nil {
		return BadRequest(err.Error())
	}
	return SuccessResponse(c, response)
}

func GetFileDetails(c *fiber.Ctx) error {
	orgId := c.Get("OrgId")
	if orgId == "" {
		return BadRequest("Organization Id missing")
	}
	fileCategory := c.Params("folder")
	refId := c.Params("refId")
	//	token := GetUserTokenValue(c)
	query := bson.M{"ref_id": refId, "folder": fileCategory}
	response, err := GetQueryResult("user_files", query, int64(0), int64(200), nil)
	if err != nil {
		return BadRequest(err.Error())
	}
	return SuccessResponse(c, response)
}

func DeleteFileIns3(c *fiber.Ctx) error {

	// s3Client, bucket := initS3()

	ID := c.Params("id")
	collectionName := c.Params("collectionName")
	filter := bson.M{"_id": ID}

	// Define a MongoDB aggregation pipeline to retrieve file metadata
	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"_id", ID}}}},
		bson.D{{"$unset", "_id"}},
		bson.D{{"$project", bson.D{{"storage_name", 1}}}},
	}

	// Retrieve file metadata from MongoDB
	res, err := GetAggregateQueryResult(collectionName, pipeline)
	// GetQueryResult(orgId, "user_files", pipeline, int64(0), int64(200), nil)
	if err != nil {
		return BadRequest(err.Error())
	}

	// Delete the file metadata from MongoDB
	_, err = database.GetConnection().Collection(collectionName).DeleteOne(ctx, filter)
	if err != nil {

	}

	for _, obj := range res {
		storageName, found := obj["storage_name"].(string)
		// if found {

		// 	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(storageName)})
		// 	if err != nil {
		// 	}
		// 	// Wait until the object is deleted in S3
		// 	err = s3Client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		// 		Bucket: aws.String(bucket),
		// 		Key:    aws.String(storageName),
		// 	})

		// 	if err != nil {
		// 	}
		// }
		if found {
			DeleteFileFromGCS(storageName)
		}
	}

	return SuccessResponse(c, nil)
}

func UploadbulkData(c *fiber.Ctx) error {

	f, err := excelize.OpenFile("Parthiban.xlsx")
	if err != nil {
		// fmt.Println(err)
		return nil
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			// fmt.Println(err)
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		// fmt.Println(err)
		return nil

	}
	for _, row := range rows {
		for _, colCell := range row {

			fmt.Print(colCell, "\t")
		}

	}

	return nil
}
