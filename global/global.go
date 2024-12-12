package global

import (
	"go.uber.org/zap"
	"wallet_sdk/models"
)

var (
	CONFIG *models.Server
	LOG    *zap.Logger
)

var (
	UtxoSpendPath       string
	UtxoUnSpendPath     string
	UtxoUserSpendPath   string
	UtxoUserUnSpendPath string
)
