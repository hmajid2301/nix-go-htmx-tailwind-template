package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jackc/pgx/v5/pgtype"

	"{{project_slug}}/internal/store/db"
)

type PricingTier struct {
	Name                 string
	MonthlyPrice         int
	AnnualPrice          int
	AnnualSavings        int
	Features             []Feature
	IsMostPopular        bool
	IsFree               bool
	CustomDataJSON       string
	AnnualDataItemsJSON  string
	MonthlyDataItemsJSON string
}

type Feature struct {
	FeatureKey  string
	FeatureText string
	LimitValue  int
}

type SubscriptionStore interface {
	GetUserAndOrganizationsByEmail(ctx context.Context, email string) (db.GetUserAndOrganizationsByEmailRow, error)
	GetActiveSubscriptionByEmail(ctx context.Context, email string) (db.GetActiveSubscriptionByEmailRow, error)
	GetSubscriptionNotificationByEventID(ctx context.Context, arg db.GetSubscriptionNotificationByEventIDParams) (db.SubscriptionNotification, error)
	GetSubscriptionTiers(ctx context.Context) ([]db.GetSubscriptionTiersRow, error)
	TransactionWithRetry(ctx context.Context, fn func(db.Querier) error) error
}

type SubscriptionService struct {
	store SubscriptionStore
}

func NewSubscription(store SubscriptionStore) SubscriptionService {
	return SubscriptionService{
		store: store,
	}
}

// TODO: cache the response this will not change very often.
func (s *SubscriptionService) GetSubscriptionTiers(ctx context.Context, email string, organizationSlug string) ([]PricingTier, error) {
	tiers := []PricingTier{}
	subPlans, err := s.store.GetSubscriptionTiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription tier info: %w", err)
	}

	for i, plan := range subPlans {
		annualDataItems, err := getDataItems(plan.PaymentProcessorAnnualPlanID)
		if err != nil {
			return nil, fmt.Errorf("failed marshal AnnualPriceID: %w", err)
		}

		monthlyDataItems, err := getDataItems(plan.PaymentProcessorMonthlyPlanID)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal monthly price ID: %w", err)
		}
		customData, err := getCustomData(email, organizationSlug)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal custom data ID: %w", err)
		}

		features, err := convertFeatures(plan.Features)
		if err != nil {
			return nil, fmt.Errorf("failed to convert from interface{} to []features: %w", err)
		}

		isFree := false
		if plan.MonthlyPrice == 0 {
			isFree = true
		}

		isMostPopular := false
		// INFO: At the moment we expect to only show three plans so the middle plan will always be marked as most
		// popular. We should probably make this dynamic. But for now this is good enough.
		if i == 1 {
			isMostPopular = true
		}

		tier := PricingTier{
			Name:                 plan.PlanName,
			AnnualDataItemsJSON:  annualDataItems,
			AnnualPrice:          int(plan.AnnualPrice),
			MonthlyDataItemsJSON: monthlyDataItems,
			MonthlyPrice:         int(plan.MonthlyPrice),
			CustomDataJSON:       customData,
			AnnualSavings:        (int(plan.MonthlyPrice) * 12) - int(plan.AnnualPrice),
			Features:             features,
			IsFree:               isFree,
			IsMostPopular:        isMostPopular,
		}

		tiers = append(tiers, tier)
	}

	return tiers, nil
}

func getDataItems(priceID string) (string, error) {
	d := []map[string]any{{"priceId": priceID, "quantity": 1}}
	dataItem, err := json.Marshal(d)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data items: %w", err)
	}
	return string(dataItem), nil
}

func getCustomData(email string, organizationSlug string) (string, error) {
	d := map[string]string{"user_email": email, "organization_slug": organizationSlug}
	customData, err := json.Marshal(d)
	if err != nil {
		return "", fmt.Errorf("failed to custom items: %w", err)
	}
	return string(customData), nil
}

func convertFeatures(input interface{}) ([]Feature, error) {
	rawSlice, ok := input.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected features type: %T, expected []interface{}", input)
	}

	features := make([]Feature, 0, len(rawSlice))
	for idx, item := range rawSlice {
		rf, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected element type at index %d: %T, expected map[string]interface{}", idx, item)
		}

		f := Feature{
			FeatureKey:  getString(rf, "feature_key"),
			FeatureText: getString(rf, "feature_text"),
		}

		if lv, exists := rf["limit_value"]; exists && lv != nil {
			if fv, ok := lv.(float64); ok {
				f.LimitValue = int(fv)
			} else {
				return nil, fmt.Errorf("unexpected limit_value type at index %d: %T, expected number", idx, lv)
			}
		}

		features = append(features, f)
	}

	return features, nil
}

func getString(m map[string]interface{}, key string) string {
	v, ok := m[key].(string)
	if !ok {
		return ""
	}
	return v
}

type Subscription struct {
	Name             string
	OrganizationSlug string
	RenewsOn         time.Time
}

func (s *SubscriptionService) GetActiveSubscription(ctx context.Context, email string) (Subscription, error) {
	subscription, err := s.store.GetActiveSubscriptionByEmail(ctx, email)
	if err != nil {
		// INFO: Must be on free plan
		if errors.Is(err, sql.ErrNoRows) {
			org, err := s.store.GetUserAndOrganizationsByEmail(ctx, email)
			if err != nil {
				return Subscription{}, fmt.Errorf("failed to get organization by user email: %w", err)
			}
			return Subscription{Name: "Free", OrganizationSlug: org.Slug}, nil
		}
		return Subscription{}, fmt.Errorf("failed to get active subscription by email: %w", err)
	}

	return Subscription{
		Name:             subscription.PlanName,
		OrganizationSlug: subscription.OrganizationSlug,
		RenewsOn:         subscription.RenewsOn.Time,
	}, nil
}

type Event struct {
	EventID    string     `json:"event_id"`
	EventType  string     `json:"event_type"`
	CustomData CustomData `json:"custom_data"`
}

type CustomData struct {
	OrganizationSlug string `json:"organization_slug"`
	Email            string `json:"user_email"`
}

func (s *SubscriptionService) HandleSubscriptionEvent(msg *message.Message) ([]*message.Message, error) {
	e := Event{}
	err := json.Unmarshal(msg.Payload, &e)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal as JSON: %w", err)
	}

	ctx := msg.Context()

	_, err = s.store.GetSubscriptionNotificationByEventID(ctx, db.GetSubscriptionNotificationByEventIDParams{
		EventID:   e.EventID,
		EventType: e.EventType,
	})
	// INFO: subscription already exists means we have already handleded this event.
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	err = s.store.TransactionWithRetry(ctx, func(q db.Querier) error {
		_, err := q.AddSubscriptionNotification(ctx, db.AddSubscriptionNotificationParams{
			EventType: e.EventType,
			EventID:   e.EventID,
			Payload:   msg.Payload,
		})
		if err != nil {
			return err
		}

		switch e.EventType {
		case "subscription.created":
			return s.AddSubscription(ctx, msg.Payload, q)
		case "subscription.activated":
			return s.ActivateSubscription(ctx, msg.Payload, q)
		}

		return nil
	})

	return nil, err
}

type AddSubscriptionEvent struct {
	EventID string `json:"event_id"`
	Data    Data   `json:"data"`
}

type Data struct {
	ID           string     `json:"id"`
	NextBilledAt CustomTime `json:"next_billed_at"`
	Items        []Item     `json:"items"`
	CustomData   CustomData `json:"custom_data"`
}

type Item struct {
	Product    Product    `json:"product"`
	TrialDates TrialDates `json:"trial_dates"`
}

type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TrialDates struct {
	StartsAt CustomTime `json:"starts_at"`
	EndsAt   CustomTime `json:"ends_at"`
}

func (s *SubscriptionService) AddSubscription(ctx context.Context, payload []byte, q db.Querier) error {
	e := AddSubscriptionEvent{}
	err := json.Unmarshal(payload, &e)
	if err != nil {
		return fmt.Errorf("failed to unmarshal as JSON: %w", err)
	}

	if len(e.Data.Items) == 0 {
		return fmt.Errorf("failed to add subscription no items in subscription notification with event_id: %s", e.EventID)
	}

	// INFO: For now we only expect one item, as it'll be a single SAAS subscription i.e. Starter or Pro plan.
	item := e.Data.Items[0]
	plan, err := q.GetSubscriptionPlanByName(ctx, item.Product.Name)
	if err != nil {
		return err
	}

	org, err := q.GetUserAndOrganizationsByEmail(ctx, e.Data.CustomData.Email)
	if err != nil {
		return err
	}

	sub, err := q.AddSubscription(ctx, db.AddSubscriptionParams{
		RenewsOn: pgtype.Timestamp{
			Time:  e.Data.NextBilledAt.Time,
			Valid: !e.Data.NextBilledAt.IsZero(),
		},
		OrganizationSlug:   org.Slug,
		SubscriptionPlanID: plan.ID,
		Status:             "trial",
		PaymentProcessorID: e.Data.ID,
	})
	if err != nil {
		return err
	}

	_, err = q.AddSubscriptionTrial(ctx, db.AddSubscriptionTrialParams{
		TrialStartAt: pgtype.Timestamp{
			Time:  item.TrialDates.StartsAt.Time,
			Valid: !item.TrialDates.StartsAt.IsZero(),
		},
		TrialEndAt: pgtype.Timestamp{
			Time:  item.TrialDates.EndsAt.Time,
			Valid: !item.TrialDates.EndsAt.IsZero(),
		},
		SubscriptionID: sub.ID,
	})

	return err
}

type ActivateSubscriptionEvent struct {
	EventID string `json:"event_id"`
	Data    Data   `json:"data"`
}

func (s *SubscriptionService) ActivateSubscription(ctx context.Context, payload []byte, q db.Querier) error {
	e := ActivateSubscriptionEvent{}
	err := json.Unmarshal(payload, &e)
	if err != nil {
		return fmt.Errorf("failed to unmarshal as JSON: %w", err)
	}

	_, err = q.ActivateSubscriptionUsingOrganizationSlug(ctx, e.Data.CustomData.OrganizationSlug)
	return err
}

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}