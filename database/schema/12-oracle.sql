/* ---- ORACLE PARAMS ---- */

CREATE TABLE oracle_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

CREATE TABLE data_source
(
    id           INT NOT NULL PRIMARY KEY,--SERIAL NOT NULL PRIMARY KEY,
    create_block BIGINT NOT NULL,
    edit_block   BIGINT,
    name         TEXT NOT NULL,
    description  TEXT,
    executable   TEXT,
    fee          COIN[],
    owner        TEXT NOT NULL REFERENCES account (address),
    sender       TEXT,
    timestamp    TIMESTAMP WITHOUT TIME ZONE
);
CREATE INDEX data_source_id_index ON data_source (id);

CREATE TABLE oracle_script
(
    id              INT NOT NULL PRIMARY KEY,--SERIAL NOT NULL PRIMARY KEY,
    create_block    BIGINT NOT NULL,
    edit_block      BIGINT,
    name            TEXT NOT NULL,
    description     TEXT,
    schema          TEXT,
    source_code_url TEXT,
    owner           TEXT NOT NULL REFERENCES account (address),
    sender          TEXT,
    timestamp       TIMESTAMP WITHOUT TIME ZONE
);
CREATE INDEX oracle_script_id_index ON oracle_script (id);

CREATE TABLE request
(
    id                INT PRIMARY KEY,--SERIAl PRIMARY KEY,
    height            BIGINT,
    oracle_script_id  INT REFERENCES oracle_script (id),
    --data_source_id    INT REFERENCES data_source (id), REMOVED
    calldata          TEXT,
    ask_count         INT,
    min_count         INT,
    client_id         TEXT,
    fee_limit         COIN[],
    prepare_gas       INT,
    execute_gas       INT,
    sender            TEXT NOT NULL REFERENCES account (address),
    tx_hash           TEXT,
    timestamp         TIMESTAMP WITHOUT TIME ZONE,
    resolve_timestamp TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMESTAMP 'epoch',
    reports_count     INT DEFAULT 0
);
CREATE INDEX request_id_index ON request (id);

CREATE TABLE report
(
    id               BIGSERIAL PRIMARY KEY,
    validator        TEXT,
    oracle_script_id INT REFERENCES oracle_script (id),
    tx_hash          TEXT
);
CREATE INDEX report_id_index ON report (id);

CREATE TABLE requests_per_date
(
    id              BIGSERIAL PRIMARY KEY,
    date            TIMESTAMP UNIQUE,
    requests_number BIGINT
);
CREATE INDEX requests_per_date_date_index ON requests_per_date (date);

CREATE TABLE data_providers_pool
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    coins      COIN[],
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX data_providers_pool_height_index ON data_providers_pool (height);

CREATE TABLE total_requests
(
    id                    BIGSERIAL PRIMARY KEY,
    date                  TIMESTAMP UNIQUE,
    total_requests_number BIGINT
);
CREATE INDEX total_requests_date_index ON total_requests (date);

CREATE TABLE request_data_source
(
    request_id        INT REFERENCES request (id),
    data_source_id    INT REFERENCES data_source (id),
);
CREATE INDEX request_source_request_index ON request_source (request_id);
CREATE INDEX request_source_data_source_index ON request_source (data_source_id);
