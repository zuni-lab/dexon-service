CREATE TABLE IF NOT EXISTS tokens (
    id VARCHAR(42) PRIMARY KEY,  -- Ethereum address (0x + 40 hex chars)
    name VARCHAR(255) NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    decimals INTEGER NOT NULL,
    is_stable BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pools (
    id VARCHAR(42) PRIMARY KEY,  -- Ethereum address (0x + 40 hex chars)
    token0_id VARCHAR(42) NOT NULL REFERENCES tokens(id),
    token1_id VARCHAR(42) NOT NULL REFERENCES tokens(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_pools_token0_id ON pools(token0_id);
CREATE INDEX idx_pools_token1_id ON pools(token1_id);


--- Pool Histories ---
CREATE TABLE IF NOT EXISTS prices (
    id BIGSERIAL,
    pool_id VARCHAR(42) NOT NULL REFERENCES pools(id),
    block_number BIGINT NOT NULL,
    block_timestamp BIGINT NOT NULL,
    sender VARCHAR(42) NOT NULL,
    recipient VARCHAR(42) NOT NULL,
    amount0 BIGINT NOT NULL,
    amount1 BIGINT NOT NULL,
    sqrt_price_x96 BIGINT NOT NULL,
    liquidity BIGINT NOT NULL,
    tick INTEGER NOT NULL,
    price_usd NUMERIC(78,18) NOT NULL,
    timestamp_utc TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, timestamp_utc)
);


CREATE TYPE ORDER_STATUS AS ENUM ('PENDING', 'PARTIAL_FILLED' ,'FILLED', 'REJECTED', 'CANCELLED');
CREATE TYPE ORDER_SIDE AS ENUM ('BUY', 'SELL');
CREATE TYPE ORDER_TYPE AS ENUM ('MARKET', 'LIMIT', 'STOP', 'TWAP');

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    pool_ids VARCHAR(42)[] NOT NULL,
    paths VARCHAR(256) NOT NULL,

    wallet VARCHAR(42) NOT NULL,
    status ORDER_STATUS NOT NULL DEFAULT 'PENDING',
    
    side ORDER_SIDE NOT NULL,
    type ORDER_TYPE NOT NULL,
    
    price NUMERIC(78,18) NOT NULL,
    -- Actual USD amount after swap
    actual_amount NUMERIC(78,18),
    amount NUMERIC(78,18) NOT NULL,
    slippage DOUBLE PRECISION,
    nonce BIGINT NOT NULL,
    signature VARCHAR(255) NOT NULL,
    tx_hash VARCHAR(255),

    parent_id BIGINT,
    twap_interval_seconds INT,
    twap_executed_times INT,
    twap_current_executed_times INT,
    twap_min_price NUMERIC(78,18),
    twap_max_price NUMERIC(78,18),
    twap_started_at TIMESTAMP,

    deadline TIMESTAMP,
    partial_filled_at TIMESTAMP,
    filled_at TIMESTAMP,
    rejected_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (parent_id) REFERENCES orders(id) ON DELETE CASCADE
);

--- Seed data ---
INSERT INTO tokens (id, name, symbol, decimals, is_stable) VALUES
(LOWER('0x9f6006523bbe9d719e83a9f050108dd5463f269d'), 'USD Coin', 'USDC', 6, TRUE),
(LOWER('0xbcb4d4effb4820abe4ab77f4349605dc2ebae551'), 'Wrapped BTC', 'WBTC', 8, FALSE),
(LOWER('0x951dbc0e23228a5b5a40f4b845da75e5658ba3e4'), 'Wrapped ETH', 'WETH', 18, FALSE),
(LOWER('0xe6ae5d42b0952c5a885538ec0aceb8f5c0c3857d'), 'Wrapped SOL', 'WSOL', 18, FALSE);


INSERT INTO pools (id, token0_id, token1_id) VALUES
(LOWER('0x70429da31815168EcCDd6898DA18D44d5641540d'), 
LOWER('0x9f6006523bbe9d719e83a9f050108dd5463f269d'), 
LOWER('0xbcb4d4effb4820abe4ab77f4349605dc2ebae551'));

INSERT INTO pools (id, token0_id, token1_id) VALUES
(LOWER('0xb8bd80BA7aFA32006Ae4cF7D1dA2Ecb8bBCa9Bf8'), 
LOWER('0x951dbc0e23228a5b5a40f4b845da75e5658ba3e4'), 
LOWER('0x9f6006523bbe9d719e83a9f050108dd5463f269d'));

INSERT INTO pools (id, token0_id, token1_id) VALUES
(LOWER('0x1BeCb7209e86A3d7D4631E6dC3bc59E897F54aF5'), 
LOWER('0x9f6006523bbe9d719e83a9f050108dd5463f269d'), 
LOWER('0xe6ae5d42b0952c5a885538ec0aceb8f5c0c3857d'));



-- -- For tracking processed blocks
CREATE TABLE block_processing_state (
    pool_address VARCHAR(42) NOT NULL,
    last_processed_block BIGINT NOT NULL,
    is_backfill BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (pool_address, is_backfill)
);


--- For user chat ---
CREATE TABLE IF NOT EXISTS chat_threads (
    id BIGSERIAL PRIMARY KEY,
    thread_id VARCHAR(256) NOT NULL,
    user_address VARCHAR(42) NOT NULL,
    thread_name VARCHAR(256) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_chat_threads_thread_id_user_address ON chat_threads(thread_id, user_address) WHERE NOT is_deleted;