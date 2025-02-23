package main

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/zuni-lab/dexon-service/cmd/worker/handlers"
	"github.com/zuni-lab/dexon-service/config"
	"github.com/zuni-lab/dexon-service/pkg/db"
	"github.com/zuni-lab/dexon-service/pkg/evm"
	"github.com/zuni-lab/dexon-service/pkg/openobserve"
)

func main() {
	loadConfig()

	ctx := context.Background()

	loadSvcs(ctx)

	mgr := evm.NewManager([]common.Address{
		common.HexToAddress("0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640"),
		common.HexToAddress("0xc7bbec68d12a0d1830360f8ec58fa599ba1b0e9b"),
	})

	for {
		if err := mgr.Connect(); err != nil {
			log.Error().Err(err).Msg("Failed to connect to Ethereum client")
			continue
		}

		if err := mgr.WatchPools(ctx, handlers.HandleSwap); err != nil {
			log.Error().Err(err).Msg("Error watching pools")
			mgr.Close()
			continue
		}
	}
}

func loadConfig() {
	config.LoadEnv()

	appName := config.Env.AppName
	if config.Env.IsDev {
		appName = appName + "-dev"
	}
	openobserve.Init(openobserve.OpenObserveConfig{
		Endpoint:    config.Env.OpenObserveEndpoint,
		Credential:  config.Env.OpenObserveCredential,
		ServiceName: appName,
		Env:         config.Env.Env,
	})

	config.InitLogger()
}

func loadSvcs(ctx context.Context) {
	db.Init(ctx, config.Env.PostgresUrl, config.Env.MigrationUrl)
}
