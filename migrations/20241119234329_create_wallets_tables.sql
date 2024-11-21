-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets (
                         id UUID PRIMARY KEY,
                         balance INT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallets;
-- +goose StatementEnd
