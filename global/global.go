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
	ChainName string

	UtxoBlockHeightPath   string
	UtxoSpendPath         string
	UtxoUnSpendPath       string
	UtxoUserSpendPath     string
	UtxoUserUnSpendPath   string
	UtxoBlockHeightByUser string
)
