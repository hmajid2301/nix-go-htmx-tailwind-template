package wristband

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type AuthClient struct {
	httpClient *retryablehttp.Client
	wristband  WristbandConf
	waitlist   WaitlistConf
}

type WristbandConf struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
}

type WaitlistConf struct {
	BaseURL string
	ListID  string
}

func NewClient(waitlist WaitlistConf, wristband WristbandConf, maxRetries int) *AuthClient {
	client := retryablehttp.NewClient()
	client.RetryMax = maxRetries

	return &AuthClient{
		httpClient: client,
		waitlist:   waitlist,
		wristband:  wristband,
	}
}

type WaitlistRequest struct {
	Email        string `json:"email"`
	WaitlistID   string `json:"waitlist_id"`
	ReferralLink string `json:"referral_link,omitempty"`
}

func (a *AuthClient) AddToWaitlist(ctx context.Context, email string, refID string) error {
	body := WaitlistRequest{
		Email:        email,
		WaitlistID:   a.waitlist.ListID,
		ReferralLink: refID,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := retryablehttp.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/api/v1/signup", a.waitlist.BaseURL),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("waitlist sign up failed (%d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (a *AuthClient) Login(ctx context.Context, tenant string) (string, string, error) {
	randomStr := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, randomStr); err != nil {
		return "", "", err
	}

	verifier := base64.RawURLEncoding.EncodeToString(randomStr)
	hash := sha256.Sum256([]byte(verifier))
	code := base64.RawURLEncoding.EncodeToString(hash[:])

	parsedURL, err := url.Parse(a.wristband.BaseURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse URL %s: %w", a.wristband.BaseURL, err)
	}

	u := &url.URL{
		Scheme: parsedURL.Scheme,
		Host:   fmt.Sprintf("%s-%s", tenant, parsedURL.Host),
		Path:   "/api/v1/oauth2/authorize",
	}

	q := url.Values{}
	q.Add("client_id", a.wristband.ClientID)
	q.Add("response_type", "code")
	q.Add("code_challenge_method", "S256")
	q.Add("code_challenge", code)
	q.Add("scope", "openid offline_access email profile")
	u.RawQuery = q.Encode()

	return u.String(), verifier, nil
}

func (a *AuthClient) Callback(ctx context.Context, code string, verifier string, tenant string) (*TokenResponse, error) {
	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", code)
	formData.Add("scope", "openid offline_access email profile")
	formData.Set("code_verifier", verifier)

	token, err := a.getJWT(ctx, formData)
	if err != nil {
		return nil, fmt.Errorf("failed to get JWT: %w", err)
	}

	token.TenantDomain = tenant
	return token, nil
}

func (a *AuthClient) Revoke(ctx context.Context, refreshToken string, tenant string) (string, error) {
	u := fmt.Sprintf("%s/api/v1/oauth2/revoke", a.wristband.BaseURL)

	formData := url.Values{}
	formData.Set("token", refreshToken)
	formData.Set("client_id", a.wristband.ClientID)
	body := bytes.NewBufferString(formData.Encode())

	req, err := retryablehttp.NewRequestWithContext(ctx, "POST", u, body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(a.wristband.ClientID, a.wristband.ClientSecret)

	parsedURL, err := url.Parse(a.wristband.BaseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL %s: %w", a.wristband.BaseURL, err)
	}

	host := fmt.Sprintf("%s-%s", tenant, parsedURL.Host)
	req.Header.Set("Host", host)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("logout request failed with status (%d): %s", resp.StatusCode, string(respBody))
	}

	redirectURL := fmt.Sprintf("https://%s/api/v1/logout", host)
	return redirectURL, nil
}

type AuthUser struct {
	UserID    string `json:"sub"`
	Email     string `json:"email"`
	AvatarURL string `json:"picture"`
	TenantID  string `json:"tnt_id"`
}

func (a *AuthClient) GetUserInfo(ctx context.Context, token string) (*AuthUser, error) {
	u := fmt.Sprintf("%s/api/v1/oauth2/userinfo", a.wristband.BaseURL)

	req, err := retryablehttp.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making userinfo request: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo request failed with status (%d): %s", resp.StatusCode, string(respBody))
	}

	authUser := &AuthUser{}
	if err := json.Unmarshal(respBody, authUser); err != nil {
		return nil, fmt.Errorf("error decoding userinfo response: %w", err)
	}

	return authUser, nil
}

func (a *AuthClient) RefreshAccessToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("refreshToken", refreshToken)

	token, err := a.getJWT(ctx, formData)
	if err != nil {
		return nil, fmt.Errorf("failed to get JWT: %w", err)
	}

	return token, nil
}

type UserUpdate struct {
	TenantID string `json:"tenantId"`
}

// TODO: move somewhere common, used in cookie and here keep it generic?.
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	IdToken      string    `json:"id_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
	TenantDomain string    `json:"tenant_domain"`
	OrgName      string    `json:"org_name"`
}

func (a *AuthClient) getJWT(ctx context.Context, formData url.Values) (*TokenResponse, error) {
	body := bytes.NewBufferString(formData.Encode())

	u := fmt.Sprintf("%s/api/v1/oauth2/token", a.wristband.BaseURL)
	req, err := retryablehttp.NewRequestWithContext(ctx, "POST", u, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(a.wristband.ClientID, a.wristband.ClientSecret)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making token request: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed with status (%d): %s", resp.StatusCode, string(respBody))
	}

	tokenResp := TokenResponse{}
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("error decoding token response: %w", err)
	}
	now := time.Now().UTC()
	tokenResp.ExpiresAt = now.Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return &tokenResp, nil
}
