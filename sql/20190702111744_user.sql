-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
    CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(20) NOT NULL,
    balance INT(10) UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (id)
);
-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
