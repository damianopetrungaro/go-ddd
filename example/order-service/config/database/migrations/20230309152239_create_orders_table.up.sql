CREATE TABLE orders
(
    id           UUID PRIMARY KEY,
    number       VARCHAR(255) NOT NULL,
    status       VARCHAR(255) NOT NULL,
    placed_by    UUID         NOT NULL,
    placed_at    TIMESTAMP    NOT NULL,
    shipped_at   TIMESTAMP,
    delivered_at TIMESTAMP
);
