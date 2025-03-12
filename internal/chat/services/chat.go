package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/zuni-lab/dexon-service/pkg/db"
	"github.com/zuni-lab/dexon-service/pkg/openai"
	"github.com/zuni-lab/dexon-service/pkg/swap"
	"github.com/zuni-lab/dexon-service/pkg/utils"
)

type ChatParams struct {
	Message     string `json:"message"`
	ThreadID    string `json:"thread_id,omitempty"`
	UserAddress string `json:"user_address" validate:"eth_addr"`
}

type ChatResponse struct {
	Message  string `json:"message"`
	ThreadID string `json:"thread_id"`
}

func ChatDex(ctx context.Context, input ChatParams, w http.ResponseWriter) error {
	var threadID string

	cryptoData, err := swap.PoolInfo.GetPrices(ctx)
	if err != nil {
		return fmt.Errorf("failed to get crypto data: %w", err)
	}

	log.Info().Any("cryptoData", cryptoData).Msg("crypto data")
	// if len(cryptoData) == 0 {
	// 	return echo.NewHTTPError(http.StatusServiceUnavailable, "no cryptocurrency data is currently available, please try again later.")
	// }

	for _, crypto := range cryptoData {
		if crypto.Price == "" {
			return fmt.Errorf("the price of %s is not available, please try again later", crypto.TokenName)
		}
	}

	log.Debug().Any("cryptoData", cryptoData).Msg("crypto data")

	cryptoDataString, errJson := json.Marshal(cryptoData)
	if errJson != nil {
		return fmt.Errorf("failed to marshal crypto data: %w", errJson)
	}

	// Create a new thread if threadId is not provided
	if input.ThreadID == "" {
		threadID, _ = openai.CreateThreadWithData(ctx, string(cryptoDataString))
	} else {
		threadID = input.ThreadID
	}

	log.Debug().Msgf("threadID: %s", threadID)

	// Stream first
	if err := openai.Streaming(ctx, threadID, input.Message, w); err != nil {
		return fmt.Errorf("failed to stream chat: %w", err)
	}

	// Insert chat thread after successful streaming
	_, err = db.DB.UpsertChatThread(ctx, db.UpsertChatThreadParams{
		ThreadID:    threadID,
		UserAddress: utils.NormalizeAddress(input.UserAddress),
		ThreadName:  input.Message,
	})
	if err != nil {
		return fmt.Errorf("failed to insert chat thread: %w", err)
	}

	return nil
}
