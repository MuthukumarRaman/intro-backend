package entities

import (
	openaiservice "introme-api/pkg/openai-service"
	"introme-api/pkg/shared/helper"

	"github.com/gofiber/fiber/v2"
)

func SetupAllRoutes(app *fiber.App) {
	SetupCRUDRoutes(app)
	SetupLookupRoutes(app)
	SetupQueryRoutes(app)
	SetupaccessUser(app)
	SetupDownloadRoutes(app)
	SetupBulkUploadRoutes(app)
	SetupDatasets(app)
	SetupOdooApis(app)
	SetupFCMRoutes(app)
	SetupLiveKitRoutes(app)
	SetupLocationRoutes(app)
	app.Static("/image", fileUploadPath)
}

// Without token access
func SetupaccessUser(app *fiber.App) {
	r := app.Group("/activation-api")
	r.Put("/generate-pwd/:access_key", helper.UpdateUserPasswordAndDeleteTempData)
	r.Get("/:access_key", helper.GetTemporaryUserDataByAccessKey) //for angular
}

// Basic Crud
func SetupCRUDRoutes(app *fiber.App) {
	r := helper.CreateRouteGroup(app, "/entities/", "REST API")
	r.Post("/:model_name", PostDocHandler)
	r.Put("/:collectionName/:id?/:model_ref_id?/:role?", putDocByIDHandlers)
	r.Get("/:collectionName/:id", GetDocByIdHandler)
	r.Delete("/:collectionName/:id", DeleteById)
	r.Delete("/:collectionName", DeleteByAll)
	r.Post("/filter/:collectionName", getDocsHandler)
}

// func SetTesting
func SetupLookupRoutes(app *fiber.App) {
	r := helper.CreateRouteGroup(app, "/group", "Data Lookup API")
	r.Get("/:groupname", helper.GroupDataBasedOnRules)
	r.Get("/testing/:modelName", helper.Testing)
}

func SetupDatasets(app *fiber.App) {
	r := helper.CreateRouteGroup(app, "/dataset", "Data Sets")
	r.Post("/config/:options?", helper.DatasetsConfig)
	r.Post("/data/:datasetname", helper.DatasetsRetrieve)
	r.Put("/:datasetname", helper.UpdateDataset)
}

func SetupQueryRoutes(app *fiber.App) {
	r := helper.CreateRouteGroup(app, "/query", "Raw Query API")
	// r.Post("/:type/:collectionName", helper.RawQueryHandler)  // currently removed
	r.Get("/:organisation", helper.GetDeviceDataByOrganizationID)
}

// S3 File Upload
func SetupDownloadRoutes(app *fiber.App) {
	r := helper.CreateRouteGroup(app, "/file", "Upload APIs")
	r.Post("/:folder/:refId", helper.FileUpload)
	r.Get("/all/:folder/:status/:page?/:limit?", helper.GetAllFileDetails)
	r.Get("/:folder/:refId", helper.GetFileDetails)
	// r.Delete("/:collectionName/:id",helper.DeleteFileIns3)
}

func SetupBulkUploadRoutes(app *fiber.App) {
	r := helper.CreateRouteGroup(app, "/upload_bulk", "Bulk Api")
	r.Get("/", helper.UploadbulkData) //todo pending
}

func SetupLocationRoutes(app *fiber.App) {
	r := helper.CreateRouteGroup(app, "/location", "Bulk Api")
	r.Post("/near", openaiservice.MatchUserProfileById) //todo pending
}

func SetupOdooApis(app *fiber.App) {
	r := helper.CreateRouteGroup(app, "/odoo", "Data Sets")
	r.Post("/test", ParnterCreateAndAddSubscriptions)
}

func SetupLiveKitRoutes(app *fiber.App) {
	r := helper.CreateRouteGroup(app, "/livekit", "Data Sets")
	r.Get("/gettoken/:roomName/:userName", helper.LiveKitGetToken)
	r.Post("/get_unread", GetUnreadMesssgaeData)
	r.Put("/update_chat_status", UpdateAllChatHandler)
	r.Post("/chats", PostChatHandler)

}

func SetupFCMRoutes(app *fiber.App) {
	r := helper.CreateRouteGroup(app, "/fcm", "Data Sets")
	r.Post("/send", SendNewFcmMessage)
}
