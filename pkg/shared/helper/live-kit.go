package helper

import (
	"context"
	"fmt"
	"introme-api/pkg/shared/database"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/livekit/protocol/auth"
	"go.mongodb.org/mongo-driver/bson"

	// "github.com/livekit/protocol/livekit"

	livekit "github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

// Generate Token for liveKIt
func LiveKitGetToken(c *fiber.Ctx) error {
	var LK_API_KEY string
	var LK_API_SECRET string
	var HourDuration time.Duration
	roomName := c.Params("roomName")
	userName := c.Params("userName")
	// chatId := c.Params("chatId")

	LK_API_KEY = os.Getenv("LIVE_KIT_LICENSE_KEY")
	LK_API_SECRET = os.Getenv("LIVE_KIT_SECERET_KEY")
	TOKEN_HOURS := os.Getenv("LIVE_KIT_TTL_TOKEN_HOURS")

	h, _ := time.ParseDuration(TOKEN_HOURS + "h")
	HourDuration = h
	canPublish := true
	canSubscribe := true
	canPublishData := true
	canPublishSources := []string{"camera", "microphone", "screen_share", "screen_share_audio"}
	if roomName == "" || userName == "" {
		return BadRequest("Invalid Room")
	}
	at := auth.NewAccessToken(LK_API_KEY, LK_API_SECRET)

	grant := &auth.VideoGrant{
		RoomJoin:          true,
		Room:              roomName,
		CanPublish:        &canPublish,
		CanSubscribe:      &canSubscribe,
		CanPublishData:    &canPublishData,
		CanPublishSources: canPublishSources,
		RoomRecord:        true,
	}

	at.SetVideoGrant(grant).
		SetIdentity(userName).
		SetValidFor(HourDuration)

	response, err := at.ToJWT()
	if err != nil {
		return BadRequest(err.Error())
	}

	database.GetConnection().Collection("user_matched").UpdateOne(context.Background(), bson.M{"_id": roomName}, bson.M{"$set": bson.M{"room_token": response}})
	fmt.Println(response)
	return SuccessResponse(c, response)
}

func RecordingCalls(c *fiber.Ctx) error {

	egressClient := lksdk.NewEgressClient(
		"http://10.0.0.125:7880",
		"KriyaTec",
		"LVAIG43bc26tqlCeWjK9KKlqHHAsoag2lL",
	)

	fileRequest := &livekit.RoomCompositeEgressRequest{
		RoomName: "THI-ADD782",
		Layout:   "speaker",
		Output: &livekit.RoomCompositeEgressRequest_File{
			File: &livekit.EncodedFileOutput{
				FileType:        livekit.EncodedFileType_MP4,
				Filepath:        "livekit-demo/my-room-test.mp4",
				DisableManifest: true,
				Output: &livekit.EncodedFileOutput_S3{
					S3: &livekit.S3Upload{
						AccessKey:          "AKIA54HQYBDJOJRN5RGG",
						Secret:             "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",
						Bucket:             "fe-app-assets",
						Endpoint:           "https://s3.amazonaws.com",
						ContentDisposition: "attachment",
						Metadata: map[string]string{
							"ACL": "public-read",
							"acl": "public-read",
						},
					},
				},
			},
		},
	}

	Info, err := egressClient.StartRoomCompositeEgress(ctx, fileRequest)

	if err != nil {
		// Handle error appropriately
		return Unexpected(err.Error())
	}

	// Return success response
	return SuccessResponse(c, Info)
}

// Send a private message
func sendPrivateMessage(sender string, receiver string, message string) {
	// Connect as a participant
	LK_API_KEY := os.Getenv("LIVE_KIT_LICENSE_KEY")
	LK_API_SECRET := os.Getenv("LIVE_KIT_SECERET_KEY")
	LK_URL := os.Getenv("LIVE_KIT_URL")

	room, err := lksdk.ConnectToRoom(LK_URL, lksdk.ConnectInfo{
		APIKey:              LK_API_KEY,
		APISecret:           LK_API_SECRET,
		RoomName:            "sample",
		ParticipantIdentity: sender,
	}, nil) // Pass nil for RoomCallback

	if err != nil {
		log.Fatalf("Failed to connect to LiveKit: %v", err)
	}
	defer room.Disconnect()

	// Send Data Packet to a specific participant
	err = room.LocalParticipant.PublishData([]byte(message), nil)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Printf("Message sent from %s to %s: %s\n", sender, receiver, message)
}
