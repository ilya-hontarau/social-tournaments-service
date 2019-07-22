-- +goose Up
-- +goose StatementBegin
CREATE TABLE tournaments
(
    id       SERIAL,
    name     TEXT,
    deposit  BIGINT NOT NULL CHECK (deposit >= 0),
    prize    BIGINT DEFAULT 0 CHECK (prize >= 0),
    finished BOOL   DEFAULT FALSE,
    winner   INT,
    FOREIGN KEY (winner) REFERENCES users (id),
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tournaments;
-- +goose StatementEnd
