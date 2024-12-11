package global

import (
	"go.uber.org/zap"
	"wallet_sdk/models"
)

var (
	CONFIG *models.Server
	LOG    *zap.Logger
)
