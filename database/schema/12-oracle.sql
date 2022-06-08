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
    id           SERIAL NOT NULL PRIMARY KEY,
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
    id              SERIAL NOT NULL PRIMARY KEY,
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
    id               SERIAl PRIMARY KEY,
    oracle_script_id INT REFERENCES oracle_script (id),
    calldata         TEXT,
    ask_count        INT,
    min_count        INT,
    client_id        TEXT,
    fee_limit        COIN[],
    prepare_gas      INT,
    execute_gas      INT,
    sender           TEXT NOT NULL REFERENCES account (address),
    tx_hash          TEXT,
    timestamp        TIMESTAMP WITHOUT TIME ZONE,
    report_timestamp TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMESTAMP 'epoch',
    reports_count    INT DEFAULT 0
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
    id         BIGSERIAL PRIMARY KEY,
    date       TIMESTAMP UNIQUE,
    requests_number BIGINT
);
CREATE INDEX requests_per_date_date_index ON requests_per_date (date);