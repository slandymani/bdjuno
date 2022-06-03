/* ---- ORACLE PARAMS ---- */

CREATE TABLE oracle_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

/*TODO create indexes, add .yaml*/
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
    report_timestamp TIMESTAMP WITHOUT TIME ZONE,
    reports_count    INT DEFAULT 0
);