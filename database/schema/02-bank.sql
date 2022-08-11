/* ---- SUPPLY ---- */

CREATE TABLE supply
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    coins      COIN[] NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX supply_height_index ON supply (height);

CREATE TABLE account_balance
(
    address         TEXT NOT NULL UNIQUE PRIMARY KEY,
    loki_balance    BIGINT,
    mgeo_balance    BIGINT,
    all_balances    COIN[],
    height          BIGINT NOT NULL
);
CREATE INDEX account_balance_address_index ON account_balance (address);