package authentication

import (
	"github.com/gofiber/fiber/v2"

	"introme-api/pkg/shared/helper"
)

func SetupRoutes(app *fiber.App) {
	//without JWT Token validation (without auth)

	auth := app.Group("/auth")
	auth.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Auth APIs")
	})
	auth.Post("/login", LoginHandler)
	auth.Get("/config", OrgConfigHandler)
	auth.Post("/register", RegisterUser)
	auth.Post("/sso-register", RegisterUserWithSSO)
	// auth.Post("/login/send-otp", MobileOtpGenerate)
	// auth.Post("/login/validate-otp", MobileOtpValidation)
	auth.Use(helper.JWTMiddleware())
	// Restricted Routes
	// auth.Post("/reset-password", ResetPasswordHandler)
	// auth.Post("/change-password", ChangePasswordHandler)
}
