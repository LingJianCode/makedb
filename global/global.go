package global

import (
	"makedb/conf"

	"go.uber.org/zap"
)

var (
	MAKEDB_LOG    *zap.Logger
	MAKEDB_CONFIG *conf.Config
)
