ALTER TABLE orders
DROP CONSTRAINT orders_pool_id_fkey;

ALTER TABLE orders
DROP COLUMN pool_id,
    ADD COLUMN pool_ids VARCHAR(42)[] NOT NULL,
    ADD COLUMN slippage DOUBLE PRECISION,
    ADD COLUMN twap_interval_seconds INT,
    ADD COLUMN twap_executed_times INT,
    ADD COLUMN twap_current_executed_times INT,
    ADD COLUMN twap_min_price NUMERIC(78,18),
    ADD COLUMN twap_max_price NUMERIC(78,18),
    ADD COLUMN deadline TIMESTAMP,
    ADD COLUMN rejected_at TIMESTAMP,
    ADD COLUMN signature VARCHAR(256);