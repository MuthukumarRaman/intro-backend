package authentication

import (
	"context"
	"fmt"

	// "fmt"
	"log"
	"time"

	"introme-api/pkg/shared/database"
	"introme-api/pkg/shared/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

func LoginHandler(c *fiber.Ctx) error {

	// OrgId := helper.GetOrgIdFromHeader(c)
	loginRequest := new(LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		return helper.BadRequest("Invalid params")
	}
	userFilter := bson.M{
		"$or": []bson.M{
			{"email_id": loginRequest.Id},
			{"mobile_number": loginRequest.Id},
		},
	}
	result := database.GetConnection().Collection("user").FindOne(ctx, userFilter)
	var user bson.M
	result.Decode(&user)

	if user == nil {
		return helper.BadRequest("Invalid User Id / Password")

	}

	if user["pwd"] == nil {
		return helper.BadRequest("Update the Password")
	}

	if !helper.CheckPasswordHashs(loginRequest.Password, user["pwd"].(string)) {
		return helper.BadRequest("Invalid User Id / Password")
	}

	// fmt.Println(loginRequest.Password, user["pwd"].(string), "password")
	// if !helper.ValidatePassword(loginRequest.Password, user["pwd"].(string)) {
	// 	return helper.BadRequest("Invalid User Id / Password")
	// }

	// If the password is valid, generate a JWT token
	claims := helper.GetNewJWTClaim()
	claims["id"] = user["_id"]
	claims["role"] = user["role"]
	// claims["acl"] = user["org_id"].(string)

	// claims["org_group"] = OrgId
	userName := user["name"]
	if userName == nil {
		userName = user["first_name"]
	}

	token := helper.GenerateJWTToken(claims, 24*60) //24*60
	// token := helper.GenerateJWTToken(claims, 1)
	response := &LoginResponse{
		Name:        userName.(string),
		UserRole:    user["role"].(string),
		UserProfile: user,
		Token:       token,
		Status:      200,
	}

	// Create a response message
	message := "Login Successfully"

	// Combine the response message and the JSON response
	responseWithMessage := struct {
		Message string `json:"message"`
		*LoginResponse
	}{
		Message:       message,
		LoginResponse: response,
	}

	return c.JSON(responseWithMessage)
}

func MobileOtpGenerate(c *fiber.Ctx) error {
	var req bson.M
	otpInfo := make(map[string]interface{})
	resp := make(map[string]string)
	orgId := c.Get("OrgId")
	if orgId == "" {
		return helper.BadRequest("Organization Id missing")
	}
	err := c.BodyParser(&req)
	_, isMobileNumExist := req["mobile"]
	if !isMobileNumExist {
		return helper.BadRequest("Invalid request, Unable to parse Mobile number")
	}
	mobile := req["mobile"].(string)
	result := database.GetConnection().Collection("user").FindOne(ctx,
		bson.M{
			"mobile":        req["mobile"].(string),
			"mobile_access": "Y",
			"status":        "A",
		})
	var user bson.M
	err = result.Decode(&user)
	if err == mongo.ErrNoDocuments {
		return helper.BadRequest("User Id not available")
	}
	if err != nil {
		return helper.BadRequest("Internal server Error")
	}
	id := uuid.New().String()
	otp := helper.GetOTPValue()
	helper.SmsInitOTP(req["mobile"].(string), otp)
	otpInfo["_id"] = id
	otpInfo["otp"] = otp
	otpInfo["otp_expired"] = false
	otpInfo["otp_verified"] = false
	if req["device_info"] != nil {
		otpInfo["device_info"] = req["device_info"]
	}
	otpInfo["created_by"] = req["mobile"].(string)
	otpInfo["created_on"] = time.Now()
	_, err = database.GetConnection().Collection("user").UpdateOne(
		ctx,
		bson.M{"mobile": mobile},
		bson.M{
			"$addToSet": bson.M{
				"otp_info": otpInfo,
			},
			//"$set": res,
		}, options.Update().SetUpsert(false))
	if err != nil {
		log.Print(err.Error())
	}
	//_, err = database.GetConnection(orgId).Collection("user_device").InsertOne(ctx, req)
	if err != nil {
		return helper.BadRequest(err.Error())
	}
	//resp = OTP{AuthKey: id, Otp: otp}
	resp["auth_key"] = id
	return helper.SuccessResponse(c, resp)
}

func MobileOtpValidation(c *fiber.Ctx) error {
	var req OTP
	orgId := c.Get("OrgId")
	if orgId == "" {
		return helper.BadRequest("Organization Id missing")
	}
	//ctx := context.Background()
	err := c.BodyParser(&req)
	if err != nil || req.Otp == 0 || req.AuthKey == "" {
		return helper.BadRequest("Invalid request, Unable to parse OTP or Auth Key")
	}
	filter := bson.M{
		"otp_info": bson.M{
			"$elemMatch": bson.M{
				"otp_expired":  false,
				"otp_verified": false,
				"_id":          req.AuthKey,
				"otp":          req.Otp,
				"created_on": bson.M{
					"$gte": time.Now().Add(-5 * time.Minute),
					"$lt":  time.Now(),
				},
			},
		},
	}

	// Run the query and retrieve the matching document
	var result bson.M
	err = database.GetConnection().Collection("user").FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return helper.BadRequest("Invalid OTP")
	}
	if err != nil {
		return helper.BadRequest("Internal server Error")
	}
	updateDoc := bson.M{
		"$set": bson.M{
			"otp_info.$[].otp_expired":      true,
			"otp_info.$[elem].otp_verified": true,
			"otp_info.$[elem].updated_by":   result["mobile"].(string),
			"otp_info.$[elem].updated_on":   time.Now(),
		},
	}

	// Define the filter to match the document containing the array
	updateFilter := bson.M{"_id": result["_id"].(string)}

	// Define the array element positional operator
	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"elem._id": req.AuthKey}},
	})
	_, err = database.GetConnection().Collection("user").UpdateOne(ctx, updateFilter, updateDoc, arrayFilters)
	if err != nil {
		log.Print(err.Error())
	}
	claims := helper.GetNewJWTClaim()
	claims["id"] = result["_id"]
	claims["role"] = result["role"]
	claims["org_id"] = orgId
	// claims["org_group"] = orgId
	userName := result["email"]
	if userName == nil {
		userName = result["name"]
	}
	token := helper.GenerateJWTToken(claims, 365*10)
	response := OTPResponse{token, result["_id"].(string)}
	return helper.SuccessResponse(c, response)
}

func OrgConfigHandler(c *fiber.Ctx) error {
	org, exists := helper.GetOrg(c)
	if !exists {
		//send error
		return helper.BadRequest("Org not found")
	}
	return helper.SuccessResponse(c, org)
}

func LoginWithSSO(c *fiber.Ctx) error {

	// OrgId := helper.GetOrgIdFromHeader(c)
	loginRequest := new(LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		return helper.BadRequest("Invalid params")
	}
	userFilter := bson.M{
		"$or": []bson.M{
			{"email_id": loginRequest.Id},
			{"mobile_number": loginRequest.Id},
		},
	}
	result := database.GetConnection().Collection("user").FindOne(ctx, userFilter)
	var user bson.M
	result.Decode(&user)

	if user == nil {
		return helper.BadRequest("Invalid User Id / Password")

	}

	// if user["pwd"] == nil {
	// 	return helper.BadRequest("Update the Password")
	// }

	// if !helper.CheckPasswordHashs(loginRequest.Password, user["pwd"].(string)) {
	// 	return helper.BadRequest("Invalid User Id / Password")
	// }

	// fmt.Println(loginRequest.Password, user["pwd"].(string), "password")
	// if !helper.ValidatePassword(loginRequest.Password, user["pwd"].(string)) {
	// 	return helper.BadRequest("Invalid User Id / Password")
	// }

	// If the password is valid, generate a JWT token
	claims := helper.GetNewJWTClaim()
	claims["id"] = user["_id"]
	claims["role"] = user["role"]
	// claims["acl"] = user["org_id"].(string)

	// claims["org_group"] = OrgId
	userName := user["name"]
	if userName == nil {
		userName = user["first_name"]
	}

	token := helper.GenerateJWTToken(claims, 24*60) //24*60
	// token := helper.GenerateJWTToken(claims, 1)
	response := &LoginResponse{
		Name:        userName.(string),
		UserRole:    user["role"].(string),
		UserProfile: user,
		Token:       token,
		Status:      200,
	}

	// Create a response message
	message := "Login Successfully"

	// Combine the response message and the JSON response
	responseWithMessage := struct {
		Message string `json:"message"`
		*LoginResponse
	}{
		Message:       message,
		LoginResponse: response,
	}

	return c.JSON(responseWithMessage)
}

// func ResetPasswordHandler(c *fiber.Ctx) error {
// 	orgId := c.Get("OrgId")
// 	if orgId == "" {
// 		return helper.BadRequest("Organization Id missing")
// 	}
// 	userToken := helper.GetUserTokenValue(c)
// 	// ctx := context.Background()
// 	req := new(ResetPasswordRequest)
// 	err := c.BodyParser(req)
// 	if err != nil {
// 		return helper.BadRequest(err.Error())
// 	}
// 	if req.Id == "" {
// 		req.Id = userToken.UserId
// 	}

// 	result := database.GetConnection(orgId).Collection("user").FindOne(ctx, bson.M{
// 		"_id": req.Id,
// 	})
// 	var user bson.M
// 	err = result.Decode(&user)
// 	if err == mongo.ErrNoDocuments {
// 		return helper.BadRequest("User Id not available")
// 	}
// 	if err != nil {
// 		return helper.BadRequest("Internal server Error")
// 	}
// 	if userToken.UserRole == "SA" {
// 		//Check the old password
// 		// if !helper.CheckPasswordHash(req.OldPwd, user["pwd"].(primitive.Binary)) {
// 		// 	return helper.BadRequest("Given user id and old password mismated")
// 		// }
// 		if !helper.CheckPassword(req.OldPwd, primitive.Binary(user["pwd"].(primitive.Binary)).Data) {
// 			return helper.BadRequest("Given user id and old password mismated")

// 		}
// 	}
// 	// TODO set random string - hard coded for now
// 	passwordHash, _ := helper.GeneratePasswordHash(req.NewPwd)

// 	_, err = database.GetConnection(orgId).Collection("user").UpdateByID(ctx,
// 		req.Id,
// 		bson.M{"$set": bson.M{"pwd": passwordHash, "password_hash": passwordHash}},
// 	)
// 	if err != nil {
// 		return helper.BadRequest(err.Error())
// 	}
// 	return c.JSON("Password Updated")
// 	// automatically return 200 success (http.StatusOK) - no need to send explictly
// }

// func ChangePasswordHandler(c *fiber.Ctx) error {
// 	orgId := c.Get("OrgId")
// 	if orgId == "" {
// 		return helper.BadRequest("Organization Id missing")
// 	}
// 	userToken := helper.GetUserTokenValue(c)
// 	// ctx := context.Background()
// 	var req ResetPasswordRequest
// 	err := c.BodyParser(&req)
// 	if err != nil {
// 		return helper.BadRequest(err.Error())
// 	}
// 	if req.Id == "" {
// 		req.Id = userToken.UserId
// 	}
// 	var user bson.M
// 	err = database.GetConnection(orgId).Collection("user").FindOne(ctx, bson.M{
// 		"_id": req.Id,
// 	}).Decode(&user)
// 	if err == mongo.ErrNoDocuments {
// 		return helper.BadRequest("User Id not available")
// 	}
// 	if err != nil {
// 		return helper.BadRequest("Internal server Error")
// 	}
// 	//Check given old password is right or not?
// 	if req.OldPwd != "" {
// 		if !helper.ValidatePassword(req.OldPwd, user["pwd"].(string)) {
// 			return helper.BadRequest("Your Old password is Wrong!")
// 		}
// 	}
// 	//update new password hash to the table
// 	passwordHash := helper.PasswordHash(req.NewPwd)
// 	_, err = database.GetConnection(orgId).Collection("user").UpdateByID(ctx,
// 		req.Id,
// 		bson.M{"$set": bson.M{"pwd": passwordHash}},
// 	)
// 	if err != nil {
// 		return helper.BadRequest(err.Error())
// 	}
// 	return c.JSON("Password Updated")
// 	// automatically return 200 success (http.StatusOK) - no need to send explictly
// }

func RegisterUser(c *fiber.Ctx) error {

	// ctx := context.Background()
	var req UserRegister
	err := c.BodyParser(&req)
	if err != nil {
		return helper.BadRequest(err.Error())
	}
	userFilter := bson.M{
		"$or": []bson.M{
			{"email_id": req.EmailId},
			{"mobile_number": req.MobileNumber},
		},
	}

	var user bson.M
	database.GetConnection().Collection("user").FindOne(ctx, userFilter).Decode(&user)

	if user != nil {
		return helper.BadRequest("User Already Exists")
	}

	id := helper.ToString(helper.GetNextSeqNumber("USR"))
	//update new password hash to the table
	passwordHash, err := helper.GeneratePasswordHash(req.Pwd)
	if err != nil {
		return helper.Unexpected("Cannot Hash Password")
	}
	req.Status = "A"
	req.Id = "USR" + id
	req.Pwd = string(passwordHash)
	req.CreatedOn = time.Now()
	_, err = database.GetConnection().Collection("user").InsertOne(ctx, req)
	if err != nil {
		return helper.BadRequest(err.Error())
	}
	res := map[string]interface{}{
		"status": 200, "message": "User Registered Successfully", "Data": req,
	}
	return c.JSON(res)

}

func RegisterUserWithSSO(c *fiber.Ctx) error {

	// ctx := context.Background()
	var req SSOUser
	var tokenId string
	err := c.BodyParser(&req)
	if err != nil {
		return helper.BadRequest(err.Error())
	}

	userFilter := bson.M{
		"email_id": req.EmailId,
	}

	var userExists bool
	var id string
	var user bson.M
	database.GetConnection().Collection("user").FindOne(ctx, userFilter).Decode(&user)
	fmt.Println(user)
	if user != nil {
		userExists = true
		// return helper.SuccessResponse(c, user)
	}
	if !userExists {

		id = helper.ToString(helper.GetNextSeqNumber("USR"))
		tokenId = "USR" + id
	} else {
		id = user["_id"].(string)
		tokenId = id
	}

	//update new password hash to the table
	if !userExists {
		req.Status = "A"
		req.Id = "USR" + id
		req.CreatedOn = time.Now()
		req.EmailVerified = true
		_, err = database.GetConnection().Collection("user").InsertOne(ctx, req)
		if err != nil {
			return helper.BadRequest(err.Error())
		}
		createWallet(req.Id)

	}

	// If the password is valid, generate a JWT token
	claims := helper.GetNewJWTClaim()
	claims["id"] = tokenId
	claims["role"] = "user"

	// claims["acl"] = user["org_id"].(string)

	// claims["org_group"] = OrgId
	userName := user["name"]
	if userName == nil {
		userName = req.FirstName
	}

	token := helper.GenerateJWTToken(claims, 24*60) //24*60
	// token := helper.GenerateJWTToken(claims, 1)
	val := getPoints("USR" + id)

	var response *LoginResponse
	if userExists {
		response = &LoginResponse{
			Name:        userName.(string),
			UserRole:    "user",
			UserProfile: user,
			Token:       token,
			Status:      200,
			Points:      val,
		}
	} else {
		response = &LoginResponse{
			Name:        userName.(string),
			UserRole:    "user",
			UserProfile: req,
			Token:       token,
			Status:      200,
			Points:      val,
		}
	}

	// Create a response message
	message := "Login Successfully"

	// Combine the response message and the JSON response
	responseWithMessage := struct {
		Message string `json:"message"`
		*LoginResponse
	}{
		Message:       message,
		LoginResponse: response,
	}
	return c.JSON(responseWithMessage)

}

func createWallet(user_id string) error {
	var wallet Wallet
	var id = uuid.New().String()
	if user_id != "" {
		wallet.ID = id
		wallet.User_ID = user_id
		wallet.Available_Credits = 600
		wallet.Status = "active"
		wallet.CreatedON = time.Now()
		wallet.UpdatedON = time.Now()

		_, err := database.GetConnection().Collection("wallet").InsertOne(ctx, wallet)
		if err != nil {
			return helper.BadRequest(err.Error())
		}
		transaction(wallet.ID)
	}
	return nil
}
func transaction(wallet_id string) error {
	var transaction Transaction
	var id = uuid.New().String()
	if wallet_id != "" {
		transaction.ID = id
		transaction.WalletId = wallet_id
		transaction.OpenAmount = 0
		transaction.Type = "CREDIT"
		transaction.TransactionAmount = 600
		transaction.AvailableAmount = 600
		transaction.TransactionTime = time.Now()
		_, err := database.GetConnection().Collection("transactions").InsertOne(ctx, transaction)
		if err != nil {
			return helper.BadRequest(err.Error())
		}
	}
	return nil
}

func getPoints(user_id string) float64 {
	userFilter := bson.M{
		"user_id": user_id,
	}
	var wallet Wallet
	database.GetConnection().Collection("wallet").FindOne(ctx, userFilter).Decode(&wallet)
	if wallet.ID != "" {
		return wallet.Available_Credits
	}
	return 0.0
}
