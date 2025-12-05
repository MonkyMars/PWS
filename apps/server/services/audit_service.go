package services

import (
	"fmt"

	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
)

type AuditService struct {
	Logger *config.Logger
}

func NewAuditService() *AuditService {
	return &AuditService{
		Logger: config.SetupLogger(),
	}
}

func (as *AuditService) GetLogs() (*[]types.AuditLog, error) {
	query := Query().
		SetOperation("select").
		SetTable(lib.TableAuditLogs).
		SetSelect([]string{"id", "timestamp", "level", "message", "attrs", "entry_hash"}).
		AddOrder(fmt.Sprintf("%s.timestamp DESC", lib.TableAuditLogs))
	result, err := database.ExecuteQuery[types.AuditLog](query)
	if err != nil {
		as.Logger.AuditError("Failed to retrieve audit logs", "error", err)
		return &[]types.AuditLog{}, err
	}

	if len(result.Data) == 0 {
		as.Logger.AuditError("No audit logs found")
		return &[]types.AuditLog{}, nil
	}

	return &result.Data, nil
}

type AuditServiceInterface interface {
	GetLogs() (*[]types.AuditLog, error)
}
