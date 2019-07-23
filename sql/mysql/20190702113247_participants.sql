-- +goose Up
-- +goose StatementBegin
CREATE TABLE participants (
    user_id INT NOT NULL,
    tournament_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (tournament_id) REFERENCES tournaments(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id,tournament_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE participants;
-- +goose StatementEnd
