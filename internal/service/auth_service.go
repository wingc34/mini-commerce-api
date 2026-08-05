package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wingc34/mini-commerce-api/internal/config"
	"github.com/wingc34/mini-commerce-api/internal/model"
	"github.com/wingc34/mini-commerce-api/internal/repository"
	pkgjwt "github.com/wingc34/mini-commerce-api/pkg/jwt"
	"golang.org/x/oauth2"
)

// googleUserInfo holds the profile data returned by Google's userinfo API.
type googleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type AuthService interface {
	GetGoogleAuthURL() string
	HandleGoogleCallback(code string) (string, error)
}

type authService struct {
	userRepo    repository.UserRepository
	oauthConfig *oauth2.Config
	jwtSecret   string
}

func NewAuthService(
	userRepo repository.UserRepository,
	oauthConfig *oauth2.Config,
	jwtSecret string,
) AuthService {
	return &authService{
		userRepo:    userRepo,
		oauthConfig: oauthConfig,
		jwtSecret:   jwtSecret,
	}
}

// GetGoogleAuthURL returns the Google OAuth consent screen URL.
// state is a fixed string here - in production use a random value per request to prevent CSRF.
func (s *authService) GetGoogleAuthURL() string {
	return s.oauthConfig.AuthCodeURL("state")
}

// HandleGoogleCallback exchanges the OAuth code for a JWT token.
// It fetches the user's Google profile, upserts the user in the DB, and returns a signed JWT.
func (s *authService) HandleGoogleCallback(code string) (string, error) {
	// 第一步：用 code 換 Google access token
	token, err := s.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code: %w", err)
	}

	// 第二步：用 access token 拿 Google profile
	userInfo, err := fetchGoogleUserInfo(token.AccessToken)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user info: %w", err)
	}

	// 第三步：Upsert user 進 DB
	user := &model.User{
		ID:    userInfo.ID,
		Email: userInfo.Email,
		Name:  &userInfo.Name,
		Image: &userInfo.Picture,
	}
	if err := s.userRepo.UpsertByGoogle(user); err != nil {
		return "", fmt.Errorf("failed to upsert user: %w", err)
	}

	// 第四步：簽發 JWT
	jwtToken, err := pkgjwt.Sign(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return jwtToken, nil
}

// fetchGoogleUserInfo calls Google's userinfo API to get the user's profile.
func fetchGoogleUserInfo(accessToken string) (*googleUserInfo, error) {
	resp, err := http.Get(config.GoogleUserInfoURL + "?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo googleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
