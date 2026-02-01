package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS creates a CORS middleware with specified origins
func CORS(origins []string) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     joinOrigins(origins),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	})
}

func joinOrigins(origins []string) string {
	if len(origins) == 0 {
		return "*"
	}
	result := ""
	for i, origin := range origins {
		if i > 0 {
			result += ","
		}
		result += origin
	}
	return result
}
