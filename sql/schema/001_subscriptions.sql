-- +goose Up
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(100) NOT NULL,
    price INTEGER NOT NULL,  -- в рублях, без копеек
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,  -- формат MM-YYYY
    end_date DATE NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE subscriptions;