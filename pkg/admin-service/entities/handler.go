package entities

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	// "go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"

	"introme-api/pkg/shared/database"
	"introme-api/pkg/shared/helper"
)

var updateOpts = options.Update().SetUpsert(true)

var fileUploadPath = ""
var ctx = context.Background()

func PostDocHandler(c *fiber.Ctx) error {
	// Extract the organization ID from the request headers
	orgId := c.Get("OrgId")
	if orgId == "" {
		// If the organization ID is missing, return a bad request response
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Organization Id missing",
		})
	}

	// Get the collection name from the URL parameters

	if c.Params("model_name") == "model_config" || c.Params("model_name") == "data_model" || c.Params("model_name") == "screen" || c.Params("model_name") == "user" { // collectionName == "user"  || collectionName == "group"
		// If it is one of these special collection names, call a function to handle it
		return helper.PostDataModelConfig(c)
	} else if c.Params("model_name") == "organisation" {
		return CloneAndInsertData(c)
	} else if c.Params("model_name") == "role" {
		return Clonedatabasedrolecollection(c)
	}

	//struct validation and Insert
	inputData, errmsg := helper.InsertValidateInDatamodel(c.Params("model_name"), string(c.Body()), orgId)
	var errmsgs []string
	if errmsg != nil {

		for _, values := range errmsg {
			errmsgs = append(errmsgs, values)
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": errmsg})
	}

	//*WITHOUT STRUCT
	// var inputData map[string]interface{}
	// if err := c.BodyParser(&inputData); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).SendString("Error parsing request body")
	// }
	collectionName := CollectionNameGet(c.Params("model_name"), orgId)
	// Insert the data into the database
	res, err := database.GetConnection().Collection(collectionName).InsertOne(ctx, inputData)
	if err != nil {
		return helper.BadRequest("Failed to insert data into the database: " + err.Error())
	}

	return helper.SuccessResponse(c, res)
}

func GetDocByIdHandler(c *fiber.Ctx) error {
	// orgId := c.Get("OrgId")
	// if orgId == "" {
	// 	return helper.BadRequest("Organization Id missing")
	// }
	filter := helper.DocIdFilter(c.Params("id"))
	collectionName := c.Params("collectionName")
	response, err := helper.GetQueryResult(collectionName, filter, int64(0), int64(1), nil)
	if err != nil {
		return helper.BadRequest(err.Error())
	}
	return helper.SuccessResponse(c, response)
}

func getDocsHandler(c *fiber.Ctx) error {
	// orgId := c.Get("OrgId")
	// if orgId == "" {
	// 	return helper.BadRequest("Organization Id missing")
	// }
	//  collectionName := c.Params("collectionName")
	var requestBody helper.PaginationRequest

	if err := c.BodyParser(&requestBody); err != nil {
		return nil
	}

	var pipeline []primitive.M
	pipeline = helper.MasterAggregationPipeline(requestBody, c)

	PagiantionPipeline := helper.PagiantionPipeline(requestBody.Start, requestBody.End)
	pipeline = append(pipeline, PagiantionPipeline)
	Response, err := helper.GetAggregateQueryResult(c.Params("collectionName"), pipeline)

	if err != nil {
		if cmdErr, ok := err.(mongo.CommandError); ok {
			return helper.BadRequest(cmdErr.Message)
		}
	}

	return helper.SuccessResponse(c, Response)
}

func RoleBasedData(c *fiber.Ctx) bson.A {
	users := helper.GetUserTokenValue(c)
	var Pipeline bson.A

	collectionName := c.Params("collectionName")
	Pipeline = bson.A{
		bson.D{{"$unwind", bson.D{{"path", "$response"}}}},
		bson.D{
			{"$match",
				bson.D{
					{"response.org_id", users.Org_name},
					{"response.name", users.UserRole},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "role_data_acl"},
					{"localField", "response.org_id"},
					{"foreignField", "org_id"},
					{"as", "result"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$result"},
					{"preserveNullAndEmptyArrays", true},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"result.model_ref_id", 1},
					{"result.org_id", 1},
					{"_id", 0},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "data_model"},
					{"let", bson.D{{"localFieldIds", bson.D{{"$toObjectId", "$result.model_ref_id"}}}}},
					{"pipeline",
						bson.A{
							bson.D{
								{"$match",
									bson.D{
										{"$expr",
											bson.D{
												{"$eq",
													bson.A{
														"$_id",
														"$$localFieldIds",
													},
												},
											},
										},
									},
								},
							},
						},
					},
					{"as", "datamodel"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$datamodel"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$set",
				bson.D{
					{"frmm", "$datamodel.model_name"},
					{"localfields", "$datamodel.description"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", collectionName},
					{"localField", "result.org_id"},
					{"foreignField", "_id"},
					{"as", "lookupResult"},
				},
			},
		},
		bson.D{
			{"$unset",
				bson.A{
					"datamodel",
					"result",
				},
			},
		},
		bson.D{{"$unwind", "$lookupResult"}},
		bson.D{
			{"$group",
				bson.D{
					{"_id", bson.D{}},
					{"localfields", bson.D{{"$addToSet", "$localfields"}}},
					{"commonFields",
						bson.D{
							{"$push",
								bson.D{
									{"$map",
										bson.D{
											{"input", bson.D{{"$objectToArray", "$lookupResult"}}},
											{"as", "field"},
											{"in",
												bson.D{
													{"k", "$$field.k"},
													{"v", "$$field.v"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"finalvalues",
						bson.D{
							{"$arrayElemAt",
								bson.A{
									"$commonFields",
									0,
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$unset",
				bson.A{
					"commonFields",
					"_id",
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"localfields", 1},
					{"finalvalues",
						bson.D{
							{"$filter",
								bson.D{
									{"input", "$finalvalues"},
									{"as", "item"},
									{"cond",
										bson.D{
											{"$in",
												bson.A{
													"$$item.k",
													"$localfields",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 0},
					{"finalvalues", bson.D{{"$arrayToObject", "$finalvalues"}}},
				},
			},
		},
		bson.D{{"$replaceRoot", bson.D{{"newRoot", "$finalvalues"}}}},
	}

	return Pipeline
}

// func GetTesting (c *fiber.Ctx)error{

// 	results, err := helper.GetAggregateQueryResult(orgId, collectionName, pipeline)
// 	if err != nil {
// 		return helper.BadRequest(err.Error())
// 	}

// 	return helper.SuccessResponse(c, results)
// }

// func childepipeline(collectionName string, id string) (bson.A, bool) {
// 	var childpipeline bson.A
// 	var flag bool = false
// 		fmt.Println(collectionName,id)
// 	if collectionName == "facility" && id == "SA" {
// 		flag = true
// 		fmt.Println("hi")
// 		childpipeline = bson.A{
// 			bson.D{
// 				{"$unwind",
// 					bson.D{
// 						{"path", "$response"},
// 						{"preserveNullAndEmptyArrays", true},
// 					},
// 				},
// 			},
// 			bson.D{{"$match", bson.D{{"response.org_id", id}}}},
// 			bson.D{
// 				{"$lookup",
// 					bson.D{
// 						{"from", "device"},
// 						{"localField", "response._id"},
// 						{"foreignField", "facility_id"},
// 						{"as", "device"},
// 					},
// 				},
// 			},
// 			bson.D{{"$match", bson.D{{"device.status", "Active"}}}},
// 		}

// 	}
// 	return childpipeline, flag
// }

func DeleteById(c *fiber.Ctx) error {
	orgId := c.Get("OrgId")
	if orgId == "" {
		return helper.BadRequest("Organization Id missing")
	}
	collectionName := c.Params("collectionName")

	filter := helper.DocIdFilter(c.Params("id"))

	if collectionName == "user_files" {
		return helper.DeleteFileIns3(c)
	}

	_, err := database.GetConnection().Collection(collectionName).DeleteOne(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error deleting document"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Document successfully deleted"})
}

func DeleteByAll(c *fiber.Ctx) error {
	orgId := c.Get("OrgId")
	if orgId == "" {
		return helper.BadRequest("Organization Id missing")
	}
	collectionName := c.Params("collectionName")

	filter := bson.M{}
	_, err := database.GetConnection().Collection(collectionName).DeleteMany(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error deleting documents"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Documents successfully deleted"})
}

// todo it worked extra for two collection only it act upset method based
func putDocByIDHandlers(c *fiber.Ctx) error {

	collectionName := c.Params("collectionName")

	// Define variables for filter and update
	var filter interface{}
	var update interface{}
	//org_data_acl
	if collectionName == "role_data_acl" || collectionName == "org_type_data_acl" || collectionName == "org_data_acl" {
		var configData map[string]interface{}

		// Parse the request body into the configData map
		if err := c.BodyParser(&configData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Error parsing request body",
				"error":   err.Error(),
			})
		}

		// model_ref_id := c.Params("model_ref_id")
		// role := c.Params("role")
		// filter is matched update if not matched upset method act this
		filter = bson.M{
			"model_ref_id": configData["model_ref_id"],
			// "role":         role,
		}

		// Define the update operation
		update = bson.M{
			"$set": configData,
		}

	} else if collectionName == "model_config" || collectionName == "data_model" || collectionName == "screen" { //|| collectionName == "user"
		response := helper.Updateformodel(c)
		return response
	} else {
		//with struct Update
		// For other collection names, use a document ID filter
		// filter = helper.DocIdFilter(c.Params("id"))
		var data map[string]interface{}
		err := c.BodyParser(&data)
		if err != nil {
			return helper.BadRequest(err.Error())
		}

		return helper.UpdateData(c, c.Params("collectionName"), helper.DocIdFilter(c.Params("id")), data, c.Get("OrgId"))
	}

	// Assuming 'ctx' is defined elsewhere in your code
	res, err := database.GetConnection().Collection(collectionName).UpdateOne(ctx, filter, update, updateOpts)
	if err != nil {
		// Handle the error appropriately, e.g., return an error response
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating document",
			"error":   err.Error(),
		})
	}

	// Return a success response with the result
	return c.Status(fiber.StatusOK).JSON(res)
}

func CloneAndInsertData(c *fiber.Ctx) error {
	orgId := c.Get("OrgId")
	if orgId == "" {
		return helper.BadRequest("Organization Id missing")
	}
	start := time.Now()
	dataMap, errmsg := helper.InsertValidateInDatamodel("organisation", string(c.Body()), orgId)
	fmt.Println(errmsg)
	// var errmsgs []string
	if errmsg != nil {
		// for _, values := range errmsg {
		// 	errmsgs = append(errmsgs, values)
		// }
		// var inputData map[string]interface{}
		// if err := c.BodyParser(&inputData); err != nil {
		// 	return c.Status(fiber.StatusBadRequest).SendString("Error parsing request body")
		// }

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": errmsg})
	}

	// Define the aggregation pipeline to match and set data
	pipeline := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"org_type", dataMap["org_type"]},
					// {"acl", bson.D{{"$ne", "N"}}}, //!undone
				},
			},
		},
		bson.D{{"$unset", "_id"}},
		bson.D{{"$set", bson.D{{"org_id", dataMap["_id"]}}}},
	}

	//check the filter to return the data
	orgDataArray, err := helper.GetAggregateQueryResult("org_type_data_acl", pipeline)
	if err != nil {
		return helper.BadRequest(err.Error())
	}

	_, err = database.GetConnection().Collection("organisation").InsertOne(ctx, dataMap)
	if err != nil {
		return helper.BadRequest("Failed to insert data into the database: " + err.Error())
	}

	//result Came from org_type_data_acl collection
	_, err = database.GetConnection().Collection("org_data_acl").InsertOne(ctx, orgDataArray[0])
	if err != nil {
		return helper.BadRequest("Failed to insert data into the database: " + err.Error())
	}

	//Once organisation create by default inisde the role collection
	var names = fmt.Sprintf("AD-%s", dataMap["_id"])

	var RolecollectionData = map[string]interface{}{
		"org_id": dataMap["_id"],
		"_id":    names,
		"status": "A",
		"name":   "Admin",
	}

	//todo inbuild struct
	_, err = database.GetConnection().Collection("role").InsertOne(ctx, RolecollectionData)
	if err != nil {

	}
	filter :=
		bson.A{
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "org_data_acl"},
						{"localField", "org_id"},
						{"foreignField", "org_id"},
						{"as", "result"},
					},
				},
			},
			bson.D{{"$unwind", bson.D{{"path", "$result"}}}},
			bson.D{{"$set", bson.D{{"result.role", "$_id"}}}},
			bson.D{{"$replaceRoot", bson.D{{"newRoot", "$result"}}}},
			bson.D{{"$unset", "_id"}},
			bson.D{{"$match", bson.D{{"acl", bson.D{{"$ne", "N"}}}}}}, //only role
		}

	roleDataArray, err := helper.GetAggregateQueryResult("role", filter)
	if err != nil {
		return helper.BadRequest(err.Error())
	}
	_, err = database.GetConnection().Collection("role_data_acl").InsertOne(ctx, roleDataArray[0])
	if err != nil {
		return helper.BadRequest("Failed to insert data into the database: " + err.Error())
	}
	fmt.Println("End Time", time.Since(start))

	return helper.SuccessResponse(c, "Successfully Data Added")
}

func Clonedatabasedrolecollection(c *fiber.Ctx) error {
	orgId := c.Get("OrgId")
	if orgId == "" {
		return helper.BadRequest("Organization Id missing")
	}
	collectionName := c.Params("collectionName") //role collection

	var inputData map[string]interface{}
	if err := c.BodyParser(&inputData); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing request body")
	}

	_, err := database.GetConnection().Collection(collectionName).InsertOne(ctx, inputData) //first insert the data in the role collection
	if err != nil {
		return helper.BadRequest("Failed to insert data into the database: " + err.Error())
	}

	filter :=
		bson.A{
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "org_data_acl"},
						{"localField", "org_id"},
						{"foreignField", "org_id"},
						{"as", "result"},
					},
				},
			},
			bson.D{{"$unwind", bson.D{{"path", "$result"}}}},
			bson.D{{"$set", bson.D{{"result.role", "$_id"}}}},
			bson.D{{"$replaceRoot", bson.D{{"newRoot", "$result"}}}},
			bson.D{{"$unset", "_id"}},
			bson.D{{"$match", bson.D{{"acl", bson.D{{"$ne", "N"}}}}}}, //only role
		}

	results, err := helper.GetAggregateQueryResult(collectionName, filter)
	if err != nil {
		return helper.BadRequest(err.Error())
	}

	for _, result := range results {

		_, err = database.GetConnection().Collection("role_data_acl").InsertOne(ctx, result)
		if err != nil {
			return helper.BadRequest("Failed to insert data into the database: " + err.Error())
		}
	}

	return helper.SuccessResponse(c, "Successfully Data Added")
}

func CollectionNameGet(model_name, orgId string) string {

	var collectionName string
	filter := bson.M{
		"model_name": model_name,
	}
	Response, err := helper.FindDocs(orgId, "model_config", filter)
	if err != nil {
		return ""
	}
	collectionName = Response["collection_name"].(string)
	return collectionName
}

type GeoNear struct {
	GeoNear []float64 `json:"geo_location" bson:"geo_location"`
	UserId  string    `json:"user_id" bson:"user_id"`
}

func GetNearByUser(c *fiber.Ctx) error {

	var inputData GeoNear
	if err := c.BodyParser(&inputData); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing request body")
	}

	pipeline := bson.A{
		bson.D{
			{"$geoNear",
				bson.D{
					{"near",
						bson.D{
							{"type", "Point"},
							{"coordinates",
								inputData.GeoNear,
							},
						},
					},
					{"distanceField", "string"},
					{"maxDistance", 50000},
					{"spherical", true},
				},
			},
		},
		// bson.D{{"$match", bson.D{{"_id", bson.D{{"$ne", inputData.UserId}}}}}},
	}

	results, err := helper.GetAggregateQueryResult("user", pipeline)
	if err != nil {
		return helper.BadRequest(err.Error())
	}

	return helper.SuccessResponse(c, results)

}
