package helper

import (
	"fmt"
	"log"
	"os"
	_ "strconv"
	"time"

	"introme-api/pkg/shared/database"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/mail.v2"
)

func SendSimpleEmailHandler(c *fiber.Ctx) error {
	orgId := c.Get("OrgId")
	if orgId == "" {
		return BadRequest("Organization Id missing")
	}
	var requestData map[string]interface{}
	err := c.BodyParser(&requestData)
	if err != nil {
		return BadRequest(err.Error())
	}

	email := requestData["_id"].(string)
	// Generate the 'decoding' value (replace this with your actual logic)
	decoding := GenerateAppaccesscode()

	// Generate the URL with parameters  it call the back end api

	link := fmt.Sprintf("http://localhost:4200/activate?accesskey=%s", decoding)

	body := createOnBoardtemplate(link)

	if err := SendEmailS(email, os.Getenv("CLIENT_EMAIL"), "Welcome to Amsort Onboarding", body); err == nil {
		// If email sending was successful
		if err := User_junked_files(email, decoding); err != nil {
			log.Println("Failed to insert user junked files:", err)
		} else {
			log.Println("Email sent successfully")
		}
	} else {
		log.Println("Email sending failed:", err)
	}

	requestData["created_on"] = time.Now()
	requestData["status"] = "A"
	//* instert the data to user collection after send the mail
	_, err = database.GetConnection().Collection("user").InsertOne(ctx, requestData)
	if err != nil {
		return BadRequest("Failed to insert data into the database: " + err.Error())
	}

	return SuccessResponse(c, requestData)
}

// USER ON BOARDING TEMPLATE  //todo
func createOnBoardtemplate(link string) string {

	body := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Welcome to Our Onboarding Process</title>
	</head>
	<body>
		<table cellpadding="0" cellspacing="0" width="100%" bgcolor="#f0f0f0">
			<tr>
				<td align="center">
					<table cellpadding="0" cellspacing="0" width="600" style="border-collapse: collapse;">
						<tr>
							<td align="center" bgcolor="#ffffff" style="padding: 40px 0 30px 0; border-top: 3px solid #007BFF;">
								<h1>Welcome to Our Onboarding Process</h1>
								<p>Thank you for choosing our services. We are excited to have you on board!</p>
								<p>Please follow the steps below to get started:</p>
								<ol>
									<div>Step 1: Complete your profile</div>
									<div>Step 2: Explore our platform</div>
									<div>Step 3: Contact our support team if you have any questions</div>
								</ol>
								<p>Enjoy your journey with us!</p>
								<p>
								<a href="` + link + `" style="background-color: #007BFF; color: #fff; padding: 10px 20px; text-decoration: none; display: inline-block; border-radius: 5px;">Activation Now</a>
								</p>
							</td>
						</tr>
					</table>
				</td>
			</tr>
		</table>
	</body>
	</html>`

	return body
}

func SendEmailS(recipientEmail string, senderEmail string, subject string, body string) error {
	email := mail.NewMessage()
	email.SetHeader("From", senderEmail)
	email.SetHeader("To", recipientEmail)

	email.SetHeader("Subject", subject)
	email.SetBody("text/html", body)

	sendinmail := mail.NewDialer("smtp.gmail.com", 587, senderEmail, os.Getenv("CLIENT_EMAIL_PASSWORD"))

	err := sendinmail.DialAndSend(email)
	if err != nil {
		return err
	}

	return nil
}

//todo not use
// func GenerateTimeLimitedLink(baseURL string) string {
// 	expirationTime := time.Now().Add(1 * time.Minute)
// 	expirationTimestamp := expirationTime.Unix() // Convert to Unix timestamp
// 	return baseURL + "?expires=" + strconv.FormatInt(expirationTimestamp, 10)
// }
