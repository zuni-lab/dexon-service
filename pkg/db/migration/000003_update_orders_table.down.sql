ALTER TABLE orders
DROP COLUMN pool_ids,
    ADD COLUMN pool_id VARCHAR(42) NOT NULL,
    DROP COLUMN slippage,
    DROP COLUMN twap_interval_seconds,
    DROP COLUMN twap_executed_times,
    DROP COLUMN twap_current_executed_times,
    DROP COLUMN twap_min_price,
    DROP COLUMN twap_max_price,
    DROP COLUMN deadline,
    DROP COLUMN rejected_at;
    DROP COLUMN signature;

ALTER TABLE orders
    ADD CONSTRAINT orders_pool_id_fkey FOREIGN KEY (pool_id) REFERENCES pools(id) ON DELETE CASCADE;