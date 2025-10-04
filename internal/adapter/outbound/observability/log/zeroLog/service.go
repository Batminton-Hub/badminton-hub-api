package zeroLog

import (
	"Badminton-Hub/internal/core/domain"
	"context"

	"github.com/rs/zerolog"
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
