package test

import (
	"errors"
	"net/http"

	"github.com/supabase-community/gotrue-go"
	"github.com/supabase-community/gotrue-go/types"
)

type MockAuth struct {
	gotrue.Client
}

func (m *MockAuth) AdminAudit(req types.AdminAuditRequest) (*types.AdminAuditResponse, error) {
	return &types.AdminAuditResponse{}, nil
}

func (m *MockAuth) AdminCreateSSOProvider(req types.AdminCreateSSOProviderRequest) (*types.AdminCreateSSOProviderResponse, error) {
	return &types.AdminCreateSSOProviderResponse{}, nil
}

func (m *MockAuth) AdminCreateUser(req types.AdminCreateUserRequest) (*types.AdminCreateUserResponse, error) {
	return &types.AdminCreateUserResponse{}, nil
}

func (m *MockAuth) AdminDeleteSSOProvider(req types.AdminDeleteSSOProviderRequest) (*types.AdminDeleteSSOProviderResponse, error) {
	return &types.AdminDeleteSSOProviderResponse{}, nil
}

func (m *MockAuth) AdminDeleteUser(req types.AdminDeleteUserRequest) error {
	return nil
}

func (m *MockAuth) AdminDeleteUserFactor(req types.AdminDeleteUserFactorRequest) error {
	return nil
}

func (m *MockAuth) AdminGenerateLink(req types.AdminGenerateLinkRequest) (*types.AdminGenerateLinkResponse, error) {
	return &types.AdminGenerateLinkResponse{}, nil
}

func (m *MockAuth) AdminGetSSOProvider(req types.AdminGetSSOProviderRequest) (*types.AdminGetSSOProviderResponse, error) {
	return &types.AdminGetSSOProviderResponse{}, nil
}

func (m *MockAuth) AdminGetUser(req types.AdminGetUserRequest) (*types.AdminGetUserResponse, error) {
	return &types.AdminGetUserResponse{}, nil
}

func (m *MockAuth) AdminListSSOProviders() (*types.AdminListSSOProvidersResponse, error) {
	return &types.AdminListSSOProvidersResponse{}, nil
}

func (m *MockAuth) AdminListUserFactors(req types.AdminListUserFactorsRequest) (*types.AdminListUserFactorsResponse, error) {
	return &types.AdminListUserFactorsResponse{}, nil
}

func (m *MockAuth) AdminListUsers() (*types.AdminListUsersResponse, error) {
	return &types.AdminListUsersResponse{}, nil
}

func (m *MockAuth) AdminUpdateSSOProvider(req types.AdminUpdateSSOProviderRequest) (*types.AdminUpdateSSOProviderResponse, error) {
	return &types.AdminUpdateSSOProviderResponse{}, nil
}

func (m *MockAuth) AdminUpdateUser(req types.AdminUpdateUserRequest) (*types.AdminUpdateUserResponse, error) {
	return &types.AdminUpdateUserResponse{}, nil
}

func (m *MockAuth) AdminUpdateUserFactor(req types.AdminUpdateUserFactorRequest) (*types.AdminUpdateUserFactorResponse, error) {
	return &types.AdminUpdateUserFactorResponse{}, nil
}

func (m *MockAuth) Authorize(req types.AuthorizeRequest) (*types.AuthorizeResponse, error) {
	return &types.AuthorizeResponse{}, nil
}

func (m *MockAuth) ChallengeFactor(req types.ChallengeFactorRequest) (*types.ChallengeFactorResponse, error) {
	return &types.ChallengeFactorResponse{}, nil
}

func (m *MockAuth) EnrollFactor(req types.EnrollFactorRequest) (*types.EnrollFactorResponse, error) {
	return &types.EnrollFactorResponse{}, nil
}

func (m *MockAuth) GetSettings() (*types.SettingsResponse, error) {
	return &types.SettingsResponse{}, nil
}

func (m *MockAuth) GetUser() (*types.UserResponse, error) {
	return &types.UserResponse{}, nil
}

func (m *MockAuth) HealthCheck() (*types.HealthCheckResponse, error) {
	return &types.HealthCheckResponse{}, nil
}

func (m *MockAuth) Invite(req types.InviteRequest) (*types.InviteResponse, error) {
	return &types.InviteResponse{}, nil
}

func (m *MockAuth) Logout() error {
	return nil
}

func (m *MockAuth) Magiclink(req types.MagiclinkRequest) error {
	return nil
}

func (m *MockAuth) OTP(req types.OTPRequest) error {
	return nil
}

func (m *MockAuth) Reauthenticate() error {
	return nil
}

func (m *MockAuth) Recover(req types.RecoverRequest) error {
	return nil
}

func (m *MockAuth) RefreshToken(token string) (*types.TokenResponse, error) {
	return &types.TokenResponse{}, nil
}

func (m *MockAuth) SAMLACS(req *http.Request) (*http.Response, error) {
	return &http.Response{}, nil
}

func (m *MockAuth) SAMLMetadata() ([]byte, error) {
	return []byte{}, nil
}

func (m *MockAuth) SSO(req types.SSORequest) (*types.SSOResponse, error) {
	return &types.SSOResponse{}, nil
}

func (m *MockAuth) SignInWithEmailPassword(email, password string) (*types.TokenResponse, error) {
	return &types.TokenResponse{}, nil
}

func (m *MockAuth) SignInWithPhonePassword(phone, password string) (*types.TokenResponse, error) {
	return &types.TokenResponse{}, nil
}

func (m *MockAuth) Signup(req types.SignupRequest) (*types.SignupResponse, error) {
	return &types.SignupResponse{
		User: types.User{
			Email: req.Email,
		},
	}, nil
}

func (m *MockAuth) Token(req types.TokenRequest) (*types.TokenResponse, error) {
	return &types.TokenResponse{}, nil
}

func (m *MockAuth) UnenrollFactor(req types.UnenrollFactorRequest) (*types.UnenrollFactorResponse, error) {
	return &types.UnenrollFactorResponse{}, nil
}

func (m *MockAuth) UpdateUser(req types.UpdateUserRequest) (*types.UpdateUserResponse, error) {
	return &types.UpdateUserResponse{}, nil
}

func (m *MockAuth) Verify(req types.VerifyRequest) (*types.VerifyResponse, error) {
	return &types.VerifyResponse{}, nil
}

func (m *MockAuth) VerifyFactor(req types.VerifyFactorRequest) (*types.VerifyFactorResponse, error) {
	return &types.VerifyFactorResponse{}, nil
}

func (m *MockAuth) VerifyForUser(req types.VerifyForUserRequest) (*types.VerifyForUserResponse, error) {
	return &types.VerifyForUserResponse{}, nil
}

func (m *MockAuth) WithClient(client http.Client) gotrue.Client {
	return gotrue.New("test", "test")
}

func (m *MockAuth) WithCustomGoTrueURL(url string) gotrue.Client {
	return gotrue.New("test", "test")
}

func (m *MockAuth) WithToken(token string) gotrue.Client {
	return gotrue.New("test", "test")
}

type MockFailedAuth struct {
	MockAuth
}

func (m *MockFailedAuth) SignInWithEmailPassword(email, password string) (*types.TokenResponse, error) {
	return nil, errors.New("failed to sign in")
}
