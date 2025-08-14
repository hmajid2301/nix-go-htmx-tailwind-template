package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"{{project_slug}}/internal/store/db"
)

type OrganizationStore interface {
	GetUsers(ctx context.Context) ([]db.GetUsersRow, error)
	UpdateOrganizationAvatar(ctx context.Context, avatarUrl pgtype.Text) (db.Organization, error)
	UpdateOrganizationDescription(ctx context.Context, description pgtype.Text) (db.Organization, error)
	UpdateOrganizationProjectLink(ctx context.Context, projectLink pgtype.Text) (db.Organization, error)
	UpdateOrganizationName(ctx context.Context, arg db.UpdateOrganizationNameParams) (db.Organization, error)
	GetOrganizationBySlug(ctx context.Context, slug string) (db.Organization, error)
	SoftDeleteOrganization(ctx context.Context, slug string) (db.Organization, error)
}

type OrganizationService struct {
	store OrganizationStore
}

func NewOrganization(store OrganizationStore) OrganizationService {
	return OrganizationService{
		store: store,
	}
}

func (o *OrganizationService) GetUsers(ctx context.Context) ([]User, error) {
	users, err := o.store.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	uu := []User{}
	for _, u := range users {
		user := User{
			Email:    u.Email,
			JoinedAt: u.CreatedAt.Time,
			// TODO: avatar set and fetch
		}
		uu = append(uu, user)
	}

	return uu, nil
}

type Organization struct {
	Name        string
	Slug        string
	URL         string
	Description string
	Avatar      string
}

func (o *OrganizationService) UpdateAvatar(ctx context.Context, avatar string) (Organization, error) {
	org, err := o.store.UpdateOrganizationAvatar(ctx, pgtype.Text{Valid: true, String: avatar})
	if err != nil {
		return Organization{}, fmt.Errorf("failed to update avatar: %w", err)
	}

	return Organization{
		Name:        org.DisplayName,
		Slug:        org.Slug,
		URL:         org.ProjectLink.String,
		Description: org.Description.String,
		Avatar:      org.AvatarUrl.String,
	}, nil
}

func (o *OrganizationService) UpdateDescription(ctx context.Context, description string) (Organization, error) {
	org, err := o.store.UpdateOrganizationDescription(ctx, pgtype.Text{Valid: true, String: description})
	if err != nil {
		return Organization{}, fmt.Errorf("failed to update description: %w", err)
	}

	return Organization{
		Name:        org.DisplayName,
		Slug:        org.Slug,
		URL:         org.ProjectLink.String,
		Description: org.Description.String,
		Avatar:      org.AvatarUrl.String,
	}, nil
}

func (o *OrganizationService) UpdateProjectLink(ctx context.Context, projectLink string) (Organization, error) {
	org, err := o.store.UpdateOrganizationProjectLink(ctx, pgtype.Text{Valid: true, String: projectLink})
	if err != nil {
		return Organization{}, fmt.Errorf("failed to update project link: %w", err)
	}

	return Organization{
		Name:        org.DisplayName,
		Slug:        org.Slug,
		URL:         org.ProjectLink.String,
		Description: org.Description.String,
		Avatar:      org.AvatarUrl.String,
	}, nil
}

func (o *OrganizationService) UpdateName(ctx context.Context, name, slug string) (Organization, error) {
	org, err := o.store.UpdateOrganizationName(ctx, db.UpdateOrganizationNameParams{
		DisplayName: name,
		Slug:        slug,
	})
	if err != nil {
		return Organization{}, fmt.Errorf("failed to update organization name: %w", err)
	}

	return Organization{
		Name:        org.DisplayName,
		Slug:        org.Slug,
		URL:         org.ProjectLink.String,
		Description: org.Description.String,
		Avatar:      org.AvatarUrl.String,
	}, nil
}

func (o *OrganizationService) GetBySlug(ctx context.Context, slug string) (Organization, error) {
	org, err := o.store.GetOrganizationBySlug(ctx, slug)
	if err != nil {
		return Organization{}, fmt.Errorf("failed to get organization: %w", err)
	}

	return Organization{
		Name:        org.DisplayName,
		Slug:        org.Slug,
		URL:         org.ProjectLink.String,
		Description: org.Description.String,
		Avatar:      org.AvatarUrl.String,
	}, nil
}

func (o *OrganizationService) SoftDelete(ctx context.Context, slug string) error {
	_, err := o.store.SoftDeleteOrganization(ctx, slug)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}
	return nil
}