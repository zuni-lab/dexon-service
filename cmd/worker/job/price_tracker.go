package job

import (
	"context"
	"fmt"
	"github.com/zuni-lab/dexon-service/internal/orders/services"
	"github.com/zuni-lab/dexon-service/pkg/utils"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/zuni-lab/dexon-service/config"
	"github.com/zuni-lab/dexon-service/pkg/evm"
)

type PriceTracker struct {
	client      *ethclient.Client
	backoff     *backoff.ExponentialBackOff
	maxAttempts uint64
}

func NewPriceTracker() *PriceTracker {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 1 * time.Second
	b.MaxInterval = 1 * time.Minute
	b.MaxElapsedTime = 30 * time.Minute

	return &PriceTracker{
		backoff:     b,
		maxAttempts: 100,
	}
}

func (p *PriceTracker) Start(ctx context.Context) error {
	for {
		if err := p.connect(); err != nil {
			log.Error().Err(err).Msg("Failed to connect to Ethereum client")
			continue
		}

		if err := p.WatchPools(ctx); err != nil {
			log.Error().Err(err).Msg("Error watching pools")
			if p.client != nil {
				p.client.Close()
				p.client = nil
			}
			continue
		}

		return nil
	}
}

func (p *PriceTracker) connect() error {
	url := strings.Replace(config.Env.AlchemyUrl, "https", "wss", 1)

	var client *ethclient.Client
	var err error

	operation := func() error {
		log.Info().Msg("Attempting to connect to Ethereum client...")
		client, err = ethclient.Dial(url)
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		return nil
	}

	if err := backoff.Retry(operation, p.backoff); err != nil {
		return err
	}

	p.client = client
	return nil
}

func (p *PriceTracker) WatchPools(ctx context.Context) error {
	pools := []common.Address{
		common.HexToAddress("0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640"),
		common.HexToAddress("0xc7bbec68d12a0d1830360f8ec58fa599ba1b0e9b"),
	}

	errChan := make(chan error, len(pools))
	var wg sync.WaitGroup

	for _, pool := range pools {
		wg.Add(1)
		go func(pool common.Address) {
			defer wg.Done()
			if err := p.WatchPool(ctx, pool); err != nil {
				select {
				case errChan <- fmt.Errorf("pool %s: %w", pool.Hex(), err):
				default:
				}
			}
		}(pool)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

func (p *PriceTracker) WatchPool(ctx context.Context, pool common.Address) error {
	for {
		if err := p.watchPoolWithRetry(ctx, pool); err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			log.Error().
				Err(err).
				Str("pool", pool.Hex()).
				Msg("Error watching pool, will retry")
			continue
		}
		return nil
	}
}

func (p *PriceTracker) watchPoolWithRetry(ctx context.Context, pool common.Address) error {
	log.Info().Msgf("Watching pool %s", pool.Hex())
	attempt := 0
RETRY:
	for attempt < int(p.maxAttempts) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		contract, err := evm.NewUniswapV3(pool, p.client)
		if err != nil {
			return fmt.Errorf("failed to create contract instance: %w", err)
		}

		watchOpts := &bind.WatchOpts{Context: ctx}
		sink := make(chan *evm.UniswapV3Swap)

		sub, err := contract.WatchSwap(watchOpts, sink, nil, nil)
		if err != nil {
			attempt++
			nextBackoff := p.backoff.NextBackOff()
			if nextBackoff == backoff.Stop {
				return fmt.Errorf("max elapsed time reached after %d attempts", attempt)
			}
			log.Warn().
				Err(err).
				Int("attempt", attempt).
				Dur("backoff", nextBackoff).
				Str("pool", pool.Hex()).
				Msg("Failed to watch swaps, retrying...")
			time.Sleep(nextBackoff)
			continue
		}
		defer sub.Unsubscribe()

		attempt = 0
		p.backoff.Reset()

		log.Info().Msg("🚀 Ready to watch pool")

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-sub.Err():
				if err != nil {
					log.Warn().
						Err(err).
						Str("pool", pool.Hex()).
						Msg("Subscription error, reconnecting...")
					goto RETRY
				}
			case event := <-sink:
				if err := p.handleEvent(event); err != nil {
					log.Error().
						Err(err).
						Str("pool", pool.Hex()).
						Msg("Error handling event")
				}
			}
		}
	}

	return fmt.Errorf("failed to maintain subscription after %d attempts", p.maxAttempts)
}

func (p *PriceTracker) handleEvent(event *evm.UniswapV3Swap) error {
	price := utils.CalculatePrice(nil, 0, 0, false)
	_, err := services.MatchOrder(context.Background(), price.String())
	if err != nil {
		log.Info().Any("event", event).Err(err).Msgf("[PriceTracker] [HandleEvent] failed to match order")
	}

	log.Info().Any("event", event).Msgf("[PriceTracker] [HandleEvent] handled %s event", event.Raw.Address.Hex())
	return nil
}
