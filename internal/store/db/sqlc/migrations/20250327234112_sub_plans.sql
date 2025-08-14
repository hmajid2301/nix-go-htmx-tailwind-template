-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS plan_features (
    subscription_plan_id UUID NOT NULL,
    feature_key TEXT NOT NULL,
    feature_text TEXT NOT NULL,
    limit_value INT NULL,
    PRIMARY KEY (subscription_plan_id, feature_key),
    FOREIGN KEY (subscription_plan_id) REFERENCES subscription_plans (id)
);

ALTER TABLE subscription_plans DROP CONSTRAINT IF EXISTS valid_pricing;

ALTER TABLE subscription_plans
ADD CONSTRAINT valid_pricing CHECK (
    (monthly_price IS NOT NULL AND monthly_price >= 0)
    OR (annual_price IS NOT NULL AND annual_price >= 0)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE subscription_plans
DROP CONSTRAINT IF EXISTS valid_pricing;

ALTER TABLE subscription_plans
ADD CONSTRAINT valid_pricing CHECK (
    (monthly_price IS NOT NULL AND monthly_price > 0)
    OR (annual_price IS NOT NULL AND annual_price > 0)
);

DROP TABLE IF EXISTS plan_features;
-- +goose StatementEnd
