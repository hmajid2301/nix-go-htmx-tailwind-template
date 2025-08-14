-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subscription_transactions (
    id UUID PRIMARY KEY DEFAULT generate_uuidv7(),
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp,
    payment_proccessor_transaction_id TEXT NOT NULL UNIQUE,
    occurred_at TIMESTAMP NOT NULL,
    subscription_id UUID NOT NULL,
    FOREIGN KEY (subscription_id) REFERENCES subscriptions (id)
);

CREATE TABLE IF NOT EXISTS subscription_notifications (
    id UUID PRIMARY KEY DEFAULT generate_uuidv7(),
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp,
    event_type TEXT NOT NULL,
    event_id TEXT NOT NULL UNIQUE,
    payload JSONB NOT NULL
);

CREATE TABLE IF NOT EXISTS subscription_trials (
    id UUID PRIMARY KEY DEFAULT generate_uuidv7(),
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp,
    trial_start_at TIMESTAMP NOT NULL,
    trial_end_at TIMESTAMP NOT NULL,
    subscription_id UUID NOT NULL,
    FOREIGN KEY (subscription_id) REFERENCES subscriptions (id)
);

ALTER TABLE subscription_plans RENAME COLUMN price TO monthly_price;

ALTER TABLE subscription_plans ADD COLUMN annual_price INTEGER;

ALTER TABLE subscription_plans
ADD CONSTRAINT valid_pricing CHECK (
    (monthly_price IS NOT NULL AND monthly_price > 0)
    OR (annual_price IS NOT NULL AND annual_price > 0)
);

ALTER TABLE subscription_plans RENAME COLUMN payment_processor_plan_id TO payment_processor_monthly_plan_id;

ALTER TABLE subscription_plans ADD COLUMN payment_processor_annual_plan_id TEXT NOT NULL;

ALTER TABLE subscription_plans ADD CONSTRAINT plan_name_unique UNIQUE (
    plan_name
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE subscription_plans DROP CONSTRAINT plan_name_unique;

ALTER TABLE subscription_plans DROP CONSTRAINT valid_pricing;

ALTER TABLE subscription_plans DROP COLUMN annual_price;

ALTER TABLE subscription_plans RENAME COLUMN monthly_price TO price;

DROP TABLE IF EXISTS subscription_transactions;
DROP TABLE IF EXISTS subscription_notifications;
DROP TABLE IF EXISTS subscription_trials;
-- +goose StatementEnd
