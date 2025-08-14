-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subscription_plans (
    id UUID PRIMARY KEY DEFAULT generate_uuidv7(),
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp,
    enabled BOOL NOT NULL DEFAULT true,
    plan_name TEXT NOT NULL,
    price INTEGER NOT NULL,
    currency TEXT NOT NULL,
    payment_processor_plan_id TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT generate_uuidv7(),
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp,
    renews_on TIMESTAMP NOT NULL,
    canceled_at TIMESTAMP,
    status TEXT NOT NULL DEFAULT 'trial',
    subscription_plan_id UUID NOT NULL,
    organization_id UUID NOT NULL,
    payment_processor_id TEXT NOT NULL,
    FOREIGN KEY (subscription_plan_id) REFERENCES subscription_plans (id),
    FOREIGN KEY (organization_id) REFERENCES organizations (id)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS subscription_plans;
-- +goose StatementEnd
