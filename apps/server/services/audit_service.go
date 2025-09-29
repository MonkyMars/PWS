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
	logger := config.SetupLogger()
	return &AuditService{Logger: logger}
}

func (as *AuditService) GetLogs() (*[]types.AuditLog, error) {
	query := Query().
		SetOperation("select").
		SetTable(lib.TableAuditLogs).
		SetSelect(database.PrefixQuery(lib.TableAuditLogs, []string{"id", "timestamp", "level", "message", "attrs", "entry_hash"})).
		AddOrder(fmt.Sprintf("%s.timestamp DESC", lib.TableAuditLogs))
	result, err := database.ExecuteQuery[types.AuditLog](query)
	if err != nil {
		as.Logger.AuditError("Failed to retrieve audit logs", "error", err)
		return nil, err
	}

	return &result.Data, nil
}
