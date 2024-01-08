-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    invoices (
        wallet_number VARCHAR(255) NOT NULL,
        currency VARCHAR(5) NOT NULL,
        uploaded_at TIMESTAMPTZ NOT NULL,
        amount numeric,
        status VARCHAR(10) NOT NULL,
        user_id integer REFERENCES users (id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS invoices;
-- +goose StatementEnd