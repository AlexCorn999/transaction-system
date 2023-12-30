-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    withdrawals (
        wallet_number VARCHAR(255) NOT NULL,
        currency VARCHAR(5) NOT NULL,
        uploaded_at TIMESTAMPTZ NOT NULL,
        amount numeric
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE IF EXISTS withdrawals;

-- +goose StatementEnd