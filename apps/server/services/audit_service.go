package services

import (
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/types"
)

type AuditService struct {
	Logger *config.Logger
}

func NewAuditService() *AuditService {
	logger := config.SetupLogger()
	return &AuditService{Logger: logger}
}

func (as *AuditService) GetLogs() ([]types.AuditLog, error) {
	query := Query().
		SetOperation("raw").
		SetRawSQL("SELECT id, timestamp, level, message, attrs, entry_hash FROM audit_logs ORDER BY timestamp DESC")

	result, err := database.ExecuteQuery[types.AuditLog](query)
	if err != nil {
		as.Logger.AuditError("Failed to retrieve audit logs", "error", err)
		return nil, err
	}

	return result.Data, nil
}
