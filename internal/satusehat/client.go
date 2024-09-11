package satusehat

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/jasoet/fhir-worker/pkg/util"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
	"strings"
	"time"
)

type TokenDetail struct {
	OrganizationName string
	DeveloperEmail   string
	TokenType        string
	IssuedAt         time.Time
	ClientId         string
	AccessToken      string
	ApplicationName  string
	ExpiresIn        time.Time
	Status           string
}

func (td *TokenDetail) IsExpired() bool {
	fvMinBefore := td.ExpiresIn.Add(-5 * time.Minute)
	return time.Now().UTC().After(fvMinBefore)
}

func (td *TokenDetail) SetExpired() {
	td.ExpiresIn = time.Time{}
}

type Client struct {
	restClient *resty.Client
	credential *Credential
	restConfig *RestConfig
	token      TokenDetail
}

type ClientOption func(*Client)

func WithCredential(credential Credential) ClientOption {
	return func(client *Client) {
		client.credential = &credential
	}
}

func WithRestConfig(restConfig RestConfig) ClientOption {
	return func(client *Client) {
		client.restConfig = &restConfig
	}
}

func NewClient(options ...ClientOption) *Client {
	client := &Client{
		restConfig: defaultRestConfig(),
		credential: defaultCredentials(),
	}

	for _, option := range options {
		option(client)
	}

	httpClient := resty.New()
	httpClient.
		SetRetryCount(client.restConfig.RetryCount).
		SetRetryWaitTime(client.restConfig.RetryWaitTime).
		SetRetryMaxWaitTime(client.restConfig.RetryMaxWaitTime).
		SetTimeout(client.restConfig.Timeout)

	client.restClient = httpClient

	return client
}

func (t *Client) RefreshToken(ctx context.Context) (*resty.Response, error) {
	config := t.credential
	getTokenURL := fmt.Sprintf("%s%s", config.AuthUrl, "/accesstoken?grant_type=client_credentials")

	_log := log.With().Ctx(ctx).Str("function", "RefreshTokenFn").Str("url", getTokenURL).Logger()

	params := map[string]string{
		"client_id":     config.ClientId,
		"client_secret": config.ClientSecret,
	}

	response, err := t.restClient.R().
		SetContext(ctx).
		SetFormData(params).
		Post(getTokenURL)

	if err != nil {
		_log.Error().Err(err).Msg("Failed to refresh token")
		return nil, NewExecutionError("Failed to refresh token", err)
	}

	if util.IsUnauthorized(response) {
		_log.Error().Int("statusCode", response.StatusCode()).Str("body", response.String()).Msg("unauthorized")
		return nil, NewUnauthorizedError(response.StatusCode(), "Unauthorized access", response.String())
	}

	if util.IsServerError(response) {
		_log.Error().Int("statusCode", response.StatusCode()).Str("body", response.String()).Msg("produce server error")
		return response, NewServerError(response.StatusCode(), "server error", response.String())
	}

	if response.IsError() {
		_log.Error().Int("statusCode", response.StatusCode()).Str("body", response.String()).Msg("response error")
		return response, NewResponseError(response.StatusCode(), "response error", response.String())
	}

	gjsonResult := gjson.ParseBytes(response.Body())
	issuedAtEpoc := gjsonResult.Get("issued_at").Int()
	expiresInSec := gjsonResult.Get("expires_in").Int()

	issueAt := time.Unix(issuedAtEpoc/1000, 0)
	expiresIn := time.Duration(expiresInSec) * time.Second

	tokenDetail := TokenDetail{
		IssuedAt:         issueAt,
		ExpiresIn:        issueAt.Add(expiresIn),
		OrganizationName: gjsonResult.Get("organization_name").String(),
		DeveloperEmail:   gjsonResult.Get("developer\\.email").String(),
		TokenType:        gjsonResult.Get("issued_at").String(),
		ClientId:         gjsonResult.Get("client_id").String(),
		AccessToken:      gjsonResult.Get("access_token").String(),
		ApplicationName:  gjsonResult.Get("application_name").String(),
		Status:           gjsonResult.Get("status").String(),
	}

	t.token = tokenDetail

	return response, nil
}

func (t *Client) PostBundle(ctx context.Context, body string) (string, error) {
	if t.token.IsExpired() {
		_, err := t.RefreshToken(ctx)
		if err != nil {
			return "", err
		}
	}

	credential := t.credential

	requestUrl := fmt.Sprintf("%s", credential.BaseUrl)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", t.token.AccessToken),
	}

	_log := log.With().Ctx(ctx).Str("function", "PostBundleFn").Str("url", requestUrl).Logger()

	response, err := t.restClient.R().
		SetHeaders(headers).
		SetContext(ctx).
		EnableTrace().
		SetBody(body).
		Post(requestUrl)

	if err != nil {
		_log.Error().Err(err).Msg("Failed to post data")
		return "", NewExecutionError("Failed to post data", err)
	}

	responseBody := response.String()

	if util.IsUnauthorized(response) {
		t.token.SetExpired()
		return responseBody, NewUnauthorizedError(response.StatusCode(), "Unauthorized access", responseBody)
	}

	if util.IsServerError(response) {
		_log.Error().Int("statusCode", response.StatusCode()).Str("body", response.String()).Msg("produce server error")
		return responseBody, NewServerError(response.StatusCode(), "server error", responseBody)
	}

	if response.IsError() {
		_log.Error().Int("statusCode", response.StatusCode()).Str("body", response.String()).Msg("response error")
		return responseBody, NewResponseError(response.StatusCode(), "response error", responseBody)
	}

	return responseBody, err
}

func (t *Client) GetPatientId(ctx context.Context, nik string) (string, error) {
	if t.token.IsExpired() {
		_, err := t.RefreshToken(ctx)
		if err != nil {
			return "", err
		}
	}

	credential := t.credential
	requestUrl := fmt.Sprintf("%s%s%s", credential.BaseUrl, "/Patient?identifier=https://fhir.kemkes.go.id/id/nik|", nik)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", t.token.AccessToken),
	}

	_log := log.With().Ctx(ctx).Str("function", "GetPatientId").Str("url", requestUrl).Logger()

	response, err := t.restClient.R().
		SetContext(ctx).
		SetHeaders(headers).
		EnableTrace().
		Get(requestUrl)

	if err != nil {
		return "", err
	}

	if err != nil {
		_log.Error().Err(err).Msg("Failed to refresh token")
		return "", NewExecutionError("Failed to refresh token", err)
	}

	if util.IsUnauthorized(response) {
		t.token.SetExpired()
		return "", NewUnauthorizedError(response.StatusCode(), "Unauthorized access", response.String())
	}

	if util.IsServerError(response) {
		_log.Error().Int("statusCode", response.StatusCode()).Str("body", response.String()).Msg("produce server error")
		return "", NewServerError(response.StatusCode(), "server error", response.String())
	}

	if response.IsError() {
		_log.Error().Int("statusCode", response.StatusCode()).Str("body", response.String()).Msg("response error")
		return "", NewResponseError(response.StatusCode(), "response error", response.String())
	}

	idPath := "entry.#(resource.id%\"P*\").resource.id"
	respJson := gjson.ParseBytes(response.Body())
	id := respJson.Get(idPath).String()

	if strings.TrimSpace(id) == "" {
		return "", NewResourceNotFoundError(response.StatusCode(), "Patient ID not found", response.String())
	}

	return id, nil

}

func (t *Client) GetPractitionerId(ctx context.Context, nik string) (string, error) {
	if t.token.IsExpired() {
		_, err := t.RefreshToken(ctx)
		if err != nil {
			return "", err
		}
	}

	credential := t.credential
	requestUrl := fmt.Sprintf("%s%s%s", credential.BaseUrl, "/Practitioner?identifier=https://fhir.kemkes.go.id/id/nik|", nik)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", t.token.AccessToken),
	}

	_log := log.With().Ctx(ctx).Str("function", "GetPractitionerId").Str("url", requestUrl).Logger()

	response, err := t.restClient.R().
		SetContext(ctx).
		SetHeaders(headers).
		EnableTrace().
		Get(requestUrl)

	if err != nil {
		return "", err
	}

	if err != nil {
		_log.Error().Err(err).Msg("Failed to refresh token")
		return "", NewExecutionError("Failed to refresh token", err)
	}

	if util.IsUnauthorized(response) {
		t.token.SetExpired()
		return "", NewUnauthorizedError(response.StatusCode(), "Unauthorized access", response.String())
	}

	if util.IsServerError(response) {
		_log.Error().Int("statusCode", response.StatusCode()).Str("body", response.String()).Msg("produce server error")
		return "", NewServerError(response.StatusCode(), "server error", response.String())
	}

	if response.IsError() {
		_log.Error().Int("statusCode", response.StatusCode()).Str("body", response.String()).Msg("response error")
		return "", NewResponseError(response.StatusCode(), "response error", response.String())
	}

	idPath := "entry.0.resource.id"
	respJson := gjson.ParseBytes(response.Body())
	id := respJson.Get(idPath).String()

	if strings.TrimSpace(id) == "" {
		return "", NewResourceNotFoundError(response.StatusCode(), "Practitioner ID not found", response.String())
	}

	return id, nil

}
