package zeroLog

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func buildBody(logInfo any) map[string]interface{} {
	body := map[string]interface{}{}
	byteBody, _ := json.Marshal(logInfo)
	json.Unmarshal(byteBody, &body)
	return body
}

func logInit(ctx context.Context) *zerolog.Logger {
	logInit := log.With().Str("path", "test").Logger()
	ctx = logInit.WithContext(ctx)
	logger := log.Ctx(ctx)
	return logger
}
