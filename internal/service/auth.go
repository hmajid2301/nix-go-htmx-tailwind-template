package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"{{project_slug}}/internal/adapter/wristband"
	"{{project_slug}}/internal/store/db"
)

type AuthStore interface {
	GetUserAndOrganizationsByEmail(ctx context.Context, email string) (db.GetUserAndOrganizationsByEmailRow, error)
	GetSubscriptionPlanByName(ctx context.Context, planName string) (db.SubscriptionPlan, error)
	TransactionWithRetry(ctx context.Context, fn func(db.Querier) error) error
}

type AuthAdapter interface {
	AddToWaitlist(ctx context.Context, email string, refID string) error
	Login(ctx context.Context, tenant string) (string, string, error)
	Callback(ctx context.Context, code string, verifier string) (*wristband.TokenResponse, error)
	Revoke(ctx context.Context, refreshToken string, tenant string) (string, error)
	GetUserInfo(ctx context.Context, token string) (*wristband.AuthUser, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*wristband.TokenResponse, error)
	InviteUser(ctx context.Context, token string, inviteDetails wristband.InviteDetails) error
	GetInvitedUsers(ctx context.Context, token string, tenantID string) ([]wristband.Invitation, error)
}

type AuthConfig struct {
	DefaultPlanName string
}

type AuthService struct {
	adapter AuthAdapter
	store   AuthStore
	config  AuthConfig
}

type User struct {
	Email    string
	Avatar   string
	JoinedAt time.Time
}

func NewAuth(adapter AuthAdapter, store AuthStore, config AuthConfig) AuthService {
	return AuthService{
		adapter: adapter,
		store:   store,
		config:  config,
	}
}

func (a *AuthService) AddToWaitlist(ctx context.Context, email string, refID string) error {
	return a.adapter.AddToWaitlist(ctx, email, refID)
}

func (a *AuthService) Login(ctx context.Context, tenant string) (string, string, error) {
	return a.adapter.Login(ctx, tenant)
}

type Token struct {
	OrganizationSlug string
	AccessToken      string
	RefreshToken     string
	ExpiresAt        time.Time
}

func (a *AuthService) Callback(ctx context.Context, code string, verifier string, tenantName string) (*Token, error) {
	token, err := a.adapter.Callback(ctx, code, verifier)
	if err != nil {
		return nil, err
	}

	userInfo, err := a.adapter.GetUserInfo(ctx, token.AccessToken)
	if err != nil {
		return nil, err
	}

	var orgSlug string
	org, err := a.store.GetUserAndOrganizationsByEmail(ctx, userInfo.Email)
	if errors.Is(err, sql.ErrNoRows) {
		org, err := a.createUserAndOrg(ctx, userInfo.Email, tenantName)
		if err != nil {
			return nil, fmt.Errorf("failed to create user and organization: %w", err)
		}
		orgSlug = org.Slug
	} else if err != nil {
		return nil, err
	}
	if orgSlug == "" {
		orgSlug = org.Slug
	}

	state := &Token{
		AccessToken:      token.AccessToken,
		RefreshToken:     token.RefreshToken,
		ExpiresAt:        token.ExpiresAt,
		OrganizationSlug: orgSlug,
	}

	return state, nil
}

func (a *AuthService) createUserAndOrg(ctx context.Context, email string, tenantName string) (db.Organization, error) {
	subPlan, err := a.store.GetSubscriptionPlanByName(ctx, a.config.DefaultPlanName)
	if err != nil {
		return db.Organization{}, fmt.Errorf("failed to get free subscription plan: %w", err)
	}

	var org db.Organization
	err = a.store.TransactionWithRetry(ctx, func(q db.Querier) error {
		user, err := q.AddUser(ctx, email)
		if err != nil {
			return err
		}

		org, err = q.AddOrganization(ctx, db.AddOrganizationParams{
			DisplayName: tenantName,
			Slug:        tenantName,
		})
		if err != nil {
			return err
		}

		_, err = q.AddUsersOrganization(ctx, db.AddUsersOrganizationParams{
			UserID:           user.ID,
			OrganizationSlug: org.Slug,
		})
		if err != nil {
			return err
		}

		_, err = q.AddSubscription(ctx, db.AddSubscriptionParams{
			OrganizationSlug:   org.Slug,
			Status:             "active",
			SubscriptionPlanID: subPlan.ID,
			PaymentProcessorID: "",
		})
		if err != nil {
			return err
		}

		return nil
	})

	return org, err
}

func (a *AuthService) GetUserInfo(ctx context.Context, token string) (*wristband.AuthUser, error) {
	return a.adapter.GetUserInfo(ctx, token)
}

func (a *AuthService) Logout(ctx context.Context, refreshToken string, tenant string) (string, error) {
	return a.adapter.Revoke(ctx, refreshToken, tenant)
}

func (a *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (*Token, error) {
	token, err := a.adapter.RefreshAccessToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh access token: %w", err)
	}

	state := &Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.ExpiresAt,
	}
	return state, nil
}

func (a *AuthService) InviteUser(ctx context.Context, token string, email string) error {
	userInfo, err := a.adapter.GetUserInfo(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	invite := wristband.InviteDetails{
		TenantID: userInfo.TenantID,
		Email:    email,
	}
	return a.adapter.InviteUser(ctx, token, invite)
}

func (a *AuthService) GetInvitedUsers(ctx context.Context, token string) ([]wristband.Invitation, error) {
	userInfo, err := a.adapter.GetUserInfo(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return a.adapter.GetInvitedUsers(ctx, token, userInfo.TenantID)
}