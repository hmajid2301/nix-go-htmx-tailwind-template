-- +goose Up
-- +goose StatementBegin
ALTER TABLE plan_features
ADD COLUMN created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now();

UPDATE subscription_plans SET annual_price = 0
WHERE annual_price IS NULL;

ALTER TABLE subscription_plans ALTER COLUMN annual_price SET NOT NULL;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE plan_features
DROP COLUMN updated_at,
DROP COLUMN created_at;

ALTER TABLE subscription_plans ALTER COLUMN annual_price DROP NOT NULL;
-- +goose StatementEnd
