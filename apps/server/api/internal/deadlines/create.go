package deadlines

import (
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/services"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CreateDealineRequest struct {
	SubjectID   uuid.UUID `json:"subject_id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     string    `json:"due_date"`
	CreatedAt   string    `json:"created_at"`
}

func CreateDeadline(c fiber.Ctx) error {
	logger := config.SetupLogger()

	body, err := middleware.GetValidatedRequest[CreateDealineRequest](c)
	if err != nil {
		logger.Error("Failed to get validated request", "error", err)
		return lib.HandleValidationError(c, err, "request")
	}

	if body == nil {
		return response.NotFound(c, "Data not found")
	}

	return response.Accepted(c, "Deadline creation accepted")
}
