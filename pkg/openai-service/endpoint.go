package openaiservice

import (
	"introme-api/pkg/shared/helper"

	"github.com/gofiber/fiber/v2"
)

func SetupAiRoutes(app *fiber.App) {
	SetupProfileApis(app)
}

func SetupProfileApis(app *fiber.App) {

	r := helper.CreateRouteGroup(app, "/profile", "PROFILE API'S")

	r.Post("/update/:profileId", GetUserProfile)

	r.Get("/match", MatchUserProfileById)

	r.Get("/get/onboarding-question/:userId?", GetUserOnboardingController)

	// r.Get("/match/:userId", MatchUserProfileById)

}
