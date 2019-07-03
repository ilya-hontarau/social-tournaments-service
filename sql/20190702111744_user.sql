-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id      INT              NOT NULL AUTO_INCREMENT,
    name    VARCHAR(20)      NOT NULL,
    balance INT(10) UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
