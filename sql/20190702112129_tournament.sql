-- +goose Up
-- +goose StatementBegin
CREATE TABLE tournaments (
     id INT NOT NULL AUTO_INCREMENT,
     name VARCHAR(20) NOT NULL,
     deposit INT(10) UNSIGNED NOT NULL,
     prize INT(10) UNSIGNED NOT NULL DEFAULT 0,
     finished BOOL DEFAULT false,
     winner INT,
     FOREIGN KEY(winner) REFERENCES users(id),
     PRIMARY KEY(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tournaments;
-- +goose StatementEnd
