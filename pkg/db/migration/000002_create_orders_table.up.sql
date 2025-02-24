CREATE TYPE ORDER_STATUS AS ENUM ('PENDING', 'PARTIAL_FILLED' ,'FILLED', 'REJECTED', 'CANCELED');
CREATE TYPE ORDER_SIDE AS ENUM ('BUY', 'SELL');
CREATE TYPE ORDER_CONDITION AS ENUM ('LIMIT', 'STOP', 'TWAP');

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT,
    wallet VARCHAR(42),
    from_token VARCHAR(10) NOT NULL,
    to_token VARCHAR(10) NOT NULL,
    status ORDER_STATUS NOT NULL DEFAULT 'PENDING',
    side ORDER_SIDE NOT NULL,
    condition ORDER_CONDITION NOT NULL,
    price NUMERIC(78,18) NOT NULL,
    amount NUMERIC(78,18) NOT NULL,
    twap_total_time INT,
    filled_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (parent_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (from_token) REFERENCES tokens(id) ON DELETE CASCADE,
    FOREIGN KEY (to_token) REFERENCES tokens(id) ON DELETE CASCADE
);
