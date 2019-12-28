package base_handler

import (
	"2019_2_Covenant/pkg/logger"
	"2019_2_Covenant/pkg/middlewares"
	"2019_2_Covenant/pkg/reader"
)

type BaseHandler struct {
	MManager  *middlewares.MiddlewareManager
	Logger    *logger.LogrusLogger
	ReqReader *reader.ReqReader
}
