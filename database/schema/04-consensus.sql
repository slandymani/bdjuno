CREATE TABLE genesis
(
    one_row_id     BOOL      NOT NULL DEFAULT TRUE PRIMARY KEY,
    chain_id       TEXT      NOT NULL,
    time           TIMESTAMP NOT NULL,
    initial_height BIGINT    NOT NULL,
    CHECK (one_row_id)
);

CREATE TABLE average_block_time_per_minute
(
    one_row_id   BOOL    NOT NULL DEFAULT TRUE PRIMARY KEY,
    average_time DECIMAL NOT NULL,
    height       BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX average_block_time_per_minute_height_index ON average_block_time_per_minute (height);

CREATE TABLE average_block_time_per_hour
(
    one_row_id   BOOL    NOT NULL DEFAULT TRUE PRIMARY KEY,
    average_time DECIMAL NOT NULL,
    height       BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX average_block_time_per_hour_height_index ON average_block_time_per_hour (height);

CREATE TABLE average_block_time_per_day
(
    one_row_id   BOOL    NOT NULL DEFAULT TRUE PRIMARY KEY,
    average_time DECIMAL NOT NULL,
    height       BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX average_block_time_per_day_height_index ON average_block_time_per_day (height);

CREATE TABLE average_block_time_from_genesis
(
    one_row_id   BOOL    NOT NULL DEFAULT TRUE PRIMARY KEY,
    average_time DECIMAL NOT NULL,
    height       BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX average_block_time_from_genesis_height_index ON average_block_time_from_genesis (height);

CREATE TABLE average_block_size
(
    id            BIGSERIAL PRIMARY KEY,
    date          TIMESTAMP,
    blocks_number BIGINT,
    block_sizes   INT,
    average_size  INT
);
CREATE INDEX average_block_size_date_index ON average_block_size (date);

CREATE TABLE average_block_time
(
    id             BIGSERIAL PRIMARY KEY,
    date           TIMESTAMP,
    last_timestamp BIGINT,
    blocks_number  BIGINT,
    block_times    BIGINT,
    average_time   BIGINT
);
CREATE INDEX average_block_time_date_index ON average_block_time (date);

CREATE TABLE txs_per_date
(
    id         BIGSERIAL PRIMARY KEY,
    date       TIMESTAMP UNIQUE,
    txs_number BIGINT
);
CREATE INDEX txs_per_date_date_index ON txs_per_date (date);

CREATE TABLE average_block_fee
(
    id            BIGSERIAL PRIMARY KEY,
    date          TIMESTAMP UNIQUE,
    blocks_number BIGINT,
    block_fees    BIGINT,
    average_fee   BIGINT DEFAULT 0
);
CREATE INDEX average_block_fee_date_index ON average_block_fee (date);