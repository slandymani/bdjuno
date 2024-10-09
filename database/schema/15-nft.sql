CREATE TABLE nft_class
(
    id          TEXT PRIMARY KEY,
    name        TEXT,
    symbol      TEXT,
    description TEXT,
    uri         TEXT,
    uri_hash    TEXT,
    data        TEXT NULL,
    owner       TEXT   NOT NULL,
    metadata    TEXT,
    height      BIGINT NOT NULL
);

CREATE TABLE nft
(
    id       TEXT   NOT NULL,
    class_id TEXT   NOT NULL REFERENCES nft_class (id),
    uri      TEXT,
    uri_hash TEXT,
    data     TEXT NULL,
    owner    TEXT   NOT NULL,
    metadata TEXT,
    height   BIGINT NOT NULL,
    PRIMARY KEY (id, class_id)
);
