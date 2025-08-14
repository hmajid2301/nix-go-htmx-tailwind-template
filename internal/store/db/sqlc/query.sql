-- name: AddUser :one
insert into users (email) values ($1) returning *;

-- name: AddOrganization :one
insert into organizations (display_name, slug) values ($1, $2) returning *;

-- name: GetUserAndOrganizationsByEmail :one
select
    u.id as user_id,
    u.email,
    u.created_at as user_created_at,
    u.updated_at as user_updated_at,
    o.display_name,
    o.slug,
    o.created_at as organization_created_at,
    o.updated_at as organization_updated_at,
    ou.created_at as organization_joined_at
from users as u
inner join organizations_users as ou on u.id = ou.user_id
inner join organizations as o on ou.organization_slug = o.slug
where u.email = $1;

-- name: GetUsers :many
select
    u.id as user_id,
    u.email,
    u.created_at,
    u.updated_at
from users as u;


-- name: UpdateOrganizationAvatar :one
update organizations
set avatar_url = $1 returning *;

-- name: UpdateOrganizationDescription :one
update organizations
set description = $1 returning *;

-- name: UpdateOrganizationProjectLink :one
update organizations
set project_link = $1 returning *;

-- name: UpdateOrganizationName :one
update organizations
set display_name = $1, updated_at = current_timestamp
where slug = $2 returning *;

-- name: GetOrganizationBySlug :one
select * from organizations
where slug = $1 and deleted_at is null;

-- name: SoftDeleteOrganization :one
update organizations
set deleted_at = current_timestamp, updated_at = current_timestamp
where slug = $1 and deleted_at is null returning *;

-- name: AddUsersOrganization :one
insert into organizations_users (user_id, organization_slug) values (
    $1, $2
) returning *;

-- name: GetActiveSubscriptionByEmail :one
select
    sp.*,
    s.*,
    o.slug as organization_slug
from users as u
inner join organizations_users as ou on u.id = ou.user_id
inner join organizations as o on ou.organization_slug = o.slug
inner join subscriptions as s on o.slug = s.organization_slug
inner join subscription_plans as sp on s.subscription_plan_id = sp.id
where
    u.email = $1
    and s.status = 'active'
    and s.canceled_at is null;

-- name: AddSubscription :one
insert into subscriptions (
    renews_on,
    status,
    subscription_plan_id,
    organization_slug,
    payment_processor_id
) values ($1, $2, $3, $4, $5) returning *;

-- name: ActivateSubscriptionUsingOrganizationSlug :one
update subscriptions
set
    status = 'active',
    updated_at = current_timestamp
where
    organization_slug = $1
returning *;

-- name: AddSubscriptionTrial :one
insert into subscription_trials (
    trial_start_at, trial_end_at, subscription_id
) values ($1, $2, $3) returning *;

-- name: GetSubscriptionNotificationByEventID :one
select *
from subscription_notifications
where event_id = $1 and event_type = $2;

-- name: AddSubscriptionNotification :one
insert into subscription_notifications (event_type, event_id, payload) values (
    $1, $2, $3
) returning *;

-- name: AddSubscriptionPlan :one
insert into subscription_plans (
    plan_name,
    monthly_price,
    annual_price,
    currency,
    payment_processor_monthly_plan_id,
    payment_processor_annual_plan_id
) values (
    $1, $2, $3, $4, $5, $6
) returning *;

-- name: GetSubscriptionPlanByName :one
select *
from subscription_plans
where plan_name = $1;

-- name: AddPlanFeature :one
insert into plan_features (
    subscription_plan_id,
    feature_key,
    feature_text,
    limit_value
) values (
    $1, $2, $3, $4
) returning *;

-- name: GetCurrentPlanFeature :one
select
    pf.feature_key,
    pf.feature_text,
    pf.limit_value,
    sp.plan_name
from plan_features pf
inner join subscription_plans sp on pf.subscription_plan_id = sp.id
where pf.feature_key = $1;

-- name: GetSubscriptionTiers :many
select
    sp.id,
    sp.plan_name,
    sp.monthly_price,
    sp.annual_price,
    sp.currency,
    sp.payment_processor_monthly_plan_id,
    sp.payment_processor_annual_plan_id,
    sp.enabled,
    sp.created_at,
    sp.updated_at,
    coalesce(
        json_agg(
            json_build_object(
                'feature_key',
                pf.feature_key,
                'feature_text',
                pf.feature_text,
                'limit_value',
                pf.limit_value::integer,
                'created_at',
                pf.created_at,
                'updated_at',
                pf.updated_at
            )
        ) filter (where pf.feature_key is not null),
        '[]'::json
    ) as features  /* : json */
from subscription_plans as sp
left join plan_features as pf on sp.id = pf.subscription_plan_id
where sp.enabled = true
group by sp.id;
