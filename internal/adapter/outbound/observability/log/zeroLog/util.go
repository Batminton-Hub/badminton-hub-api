package zeroLog

import (
	"Badminton-Hub/internal/core/domain"
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ZeroLog struct{}

func NewZeroLog() *ZeroLog {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	return &ZeroLog{}
}

func (z *ZeroLog) Info(ctx context.Context, info domain.LogInfo) {
	logger := logInit(ctx)
	body := buildBody(info)
	logger.Info().Fields(body).Send()
}

func (z *ZeroLog) Error(ctx context.Context, info domain.LogError) {
	logger := logInit(ctx)
	body := buildBody(info)
	logger.Error().Fields(body).Send()
}

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
