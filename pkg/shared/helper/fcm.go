package helper

import (
	"context"
	"fmt"
	"introme-api/pkg/shared/database"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/api/option"
)

func SendNewFCMNotification(deviceToken, title, body string, data map[string]string) error {
	// Initialize Firebase App with credentials
	//opt := option.WithCredentialsFile("service-account.json") // Replace with actual path
	opt := option.WithCredentialsFile("introme-firebase-service.json")
	// 🔹 Explicitly set Project ID
	config := &firebase.Config{
		ProjectID: "introme-847cf", // Replace with actual Firebase Project ID
	}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return fmt.Errorf("error initializing Firebase app: %v", err)
	}

	// Get Firebase Messaging Client
	client, err := app.Messaging(context.Background())
	if err != nil {
		return fmt.Errorf("error getting Messaging client: %v", err)
	}

	// Define FCM Message
	message := &messaging.Message{
		Token: deviceToken,
		Data:  data,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	// Send FCM Message
	response, err := client.Send(context.Background(), message)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	fmt.Println("FCM Notification Sent Successfully! Response:", response)
	return nil
}

// func SendFCMNotification(deviceToken, title, body string) error {
// 	// Replace with your Firebase Server Key
// 	serverKey := "BI7fuE1Qd9rGXMuLmzXoWPQN2vdpZtY9WEXAZBxHlHWEUJKW-IP7zGB-r17lzDKIiBWifrdevHzLB60xiryu__Q"

// 	// Create FCM client
// 	client := fcm.NewFcmClient(serverKey)

// 	// Set notification payload
// 	data := map[string]interface{}{
// 		"title": title,
// 		"body":  body,
// 	}

// 	// Send to a single device
// 	client.NewFcmRegIdsMsg([]string{deviceToken}, data)

// 	// Send message
// 	status, err := client.Send()
// 	if err != nil {
// 		return fmt.Errorf("error sending message: %v", err)
// 	}
// 	fmt.Println(status)
// 	return nil
// }

func SendUserNotification(fromUser string, toUsers []bson.M, notificationType string, msg string, chatId string) error {

	var body string
	var title string
	var data map[string]string
	// var profileImage string
	var fromUserProfileImage string
	var fromUserData map[string]interface{}
	err := database.GetConnection().Collection("user").FindOne(context.Background(), bson.M{"_id": fromUser}).Decode(&fromUserData)
	if err != nil {
		return err
	}

	if fromUserData["profile_image"] != nil {
		fromUserProfileImage = fromUserData["profile_image"].(string)
	}

	for _, user := range toUsers {
		firstName := user["first_name"].(string)
		lastName := user["last_name"].(string)
		fullName := firstName + " " + lastName
		fcm_token := user["fcm_token"].(string)
		toUser := user["_id"].(string)

		// if user["profile_image"] != nil {
		// 	profileImage = user["profile_image"].(string)
		// }

		if notificationType == "FRIEND" {
			title = "Friend Request"
			body = fullName + " has sent you a friend request"
			data = map[string]string{
				"type":      notificationType,
				"date_time": time.Now().Format(time.RFC3339),
				"image":     fromUserProfileImage,
			}

		} else if notificationType == "FRIENDACCEPTED" {
			title = "Friend Request Accepted"
			body = fullName + " has accepted your friend request"
			data = map[string]string{
				"type":      notificationType,
				"date_time": time.Now().Format(time.RFC3339),
				"image":     fromUserProfileImage,
			}
		} else if notificationType == "CHAT" {
			title = fullName
			body = msg
			data = map[string]string{
				"type":      notificationType,
				"date_time": time.Now().Format(time.RFC3339),
				"image":     fromUserProfileImage,
				"from_id":   fromUser,
				"to":        toUser,
				"chat_id":   chatId,
				"message":   msg,
				"status":    "sent",
			}
		}

		err := SendNewFCMNotification(fcm_token, title, body, data)
		if err != nil {
			fmt.Println(err.Error())
		}

		if err == nil {
			if notificationType == "FRIEND" {
				res := map[string]interface{}{
					"_id":        uuid.New().String(),
					"from":       fromUser,
					"req_status": "pending",
					"to":         toUser,
					"type":       "FRIEND",
					"created_on": time.Now(),
					"status":     "sent",
				}
				database.GetConnection().Collection("notifications").InsertOne(context.Background(), res)
			} else if notificationType == "FRIENDACCEPTED" {
				var userIds []string
				userIds = append(userIds, fromUser)
				userIds = append(userIds, toUser)
				// GetNextSeqNumber("FRIENDSHIP")
				res := map[string]interface{}{
					"_id":         uuid.New().String(),
					"user_ids":    userIds,
					"type":        "FRIENDACCEPTED",
					"created_on":  time.Now(),
					"accepted_by": fromUser,
					"status":      "Accepted",
				}

				database.GetConnection().Collection("user_matched").InsertOne(context.Background(), res)
			} else if notificationType == "CHAT" {

				res := map[string]interface{}{
					"_id":       uuid.New().String(),
					"from":      fromUser,
					"to":        toUser,
					"chat_id":   chatId,
					"message":   msg,
					"date_time": time.Now(),
					"status":    "sent",
				}

				database.GetConnection().Collection("chats").InsertOne(context.Background(), res)
			}

		}

	}
	return nil

}
