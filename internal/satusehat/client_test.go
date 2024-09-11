package satusehat

import (
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/jasoet/fhir-worker/pkg/util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSetExpired(t *testing.T) {
	tests := []struct {
		name string
		td   *TokenDetail
	}{
		{
			name: "set expired time for token",
			td:   &TokenDetail{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.td.SetExpired()
			assert.True(t, tt.td.IsExpired())
		})
	}
}

func TestIsExpired(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		td   TokenDetail
		want bool
	}{
		{
			name: "NotExpired",
			td: TokenDetail{
				ExpiresIn: now.Add(10 * time.Minute),
			},
			want: false,
		},
		{
			name: "Expired",
			td: TokenDetail{
				ExpiresIn: now.Add(-10 * time.Minute),
			},
			want: true,
		},
		{
			name: "OneSecondToExpire",
			td: TokenDetail{
				ExpiresIn: now.Add(4*time.Minute + 59*time.Second),
			},
			want: true,
		},
		{
			name: "ExactlyFiveMinutesToExpire",
			td: TokenDetail{
				ExpiresIn: now.Add(5 * time.Minute),
			},
			want: true,
		},
		{
			name: "SixMinutesToExpire",
			td: TokenDetail{
				ExpiresIn: now.Add(6 * time.Minute),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.td.IsExpired(); got != tt.want {
				t.Errorf("TokenDetail.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		options  []ClientOption
		expected *Client
	}{
		{
			name:    "Default",
			options: []ClientOption{},
			expected: &Client{
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
				restClient: resty.New(),
			},
		},
		{
			name:    "WithCredentials",
			options: []ClientOption{WithCredential(Credential{AuthUrl: "username", BaseUrl: "password"})},
			expected: &Client{
				restConfig: defaultRestConfig(),
				credential: &Credential{AuthUrl: "username", BaseUrl: "password"},
				restClient: resty.New(),
			},
		},
		{
			name: "WithRestConfig",
			options: []ClientOption{WithRestConfig(RestConfig{
				RetryCount:       5,
				RetryWaitTime:    10 * time.Millisecond,
				RetryMaxWaitTime: 1000 * time.Millisecond,
				Timeout:          2000 * time.Millisecond,
			})},
			expected: &Client{
				restConfig: &RestConfig{
					RetryCount:       5,
					RetryWaitTime:    10 * time.Millisecond,
					RetryMaxWaitTime: 1000 * time.Millisecond,
					Timeout:          2000 * time.Millisecond,
				},
				credential: defaultCredentials(),
				restClient: resty.New(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := NewClient(test.options...)
			assert.Equal(t, test.expected.restConfig, result.restConfig)
			assert.Equal(t, test.expected.credential, result.credential)
			assert.Equal(t, test.expected.token, result.token)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	issuedAt := time.Now()
	expiredIn := 3600

	response := map[string]any{
		"access_token": "access_token",
		"expires_in":   expiredIn,
		"issued_at":    issuedAt.UnixMilli(),
		"token_type":   "Bearer",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path[:4] {
		case "/200":
			w.WriteHeader(http.StatusOK)
			body, err := json.Marshal(response)
			if err != nil {
				panic(err)
			}
			_, _ = w.Write(body)
		case "/401":
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Unauthorized"))
		case "/500":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Internal Server Error"))
		case "/404":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Not Found"))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Not Found"))
		}
	}))
	defer server.Close()

	tests := []struct {
		name        string
		endpoint    string
		client      *Client
		expected    *TokenDetail
		expectErr   bool
		targetError error
	}{
		{
			name:     "RefreshToken-200",
			endpoint: "/200",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
			},
			expected: &TokenDetail{
				AccessToken: "access_token",
				TokenType:   "Bearer",
				ExpiresIn:   issuedAt.Add(time.Duration(expiredIn) * time.Second),
			},
			expectErr:   false,
			targetError: nil,
		},
		{
			name:     "RefreshToken-401",
			endpoint: "/401",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
			},
			expected:    &TokenDetail{},
			expectErr:   true,
			targetError: &UnauthorizedError{},
		},
		{
			name:     "RefreshToken-500",
			endpoint: "/500",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
			},
			expected:    &TokenDetail{},
			expectErr:   true,
			targetError: &ServerError{},
		},
		{
			name:     "RefreshToken-404",
			endpoint: "/404",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
			},
			expected:    &TokenDetail{},
			expectErr:   true,
			targetError: &ResponseError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.client.credential.AuthUrl = server.URL + tt.endpoint
			_, err := tt.client.RefreshToken(context.Background())
			if tt.expectErr {
				assert.Error(t, err)
				assert.True(t, util.IsSameType(err, tt.targetError))
			} else {
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.expected.AccessToken, tt.client.token.AccessToken)
				assert.True(t, tt.client.token.ExpiresIn.After(time.Now()))
			}
		})
	}
}

func TestPostBundle(t *testing.T) {
	response := map[string]any{
		"result":  "success",
		"message": "Bundle posted successfully",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path[:4] {
		case "/200":
			w.WriteHeader(http.StatusOK)
			body, err := json.Marshal(response)
			if err != nil {
				panic(err)
			}
			_, _ = w.Write(body)
		case "/401":
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Unauthorized"))
		case "/500":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Internal Server Error"))
		case "/404":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Not Found"))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Not Found"))
		}
	}))
	defer server.Close()

	nowUtc := time.Now().UTC()

	tests := []struct {
		name        string
		endpoint    string
		client      *Client
		expectErr   bool
		targetError error
	}{
		{
			name:     "PostBundle-200",
			endpoint: "/200",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
				token: TokenDetail{
					ExpiresIn:   nowUtc.Add(1 * time.Hour),
					IssuedAt:    nowUtc,
					AccessToken: "access token",
				},
			},
			expectErr:   false,
			targetError: nil,
		},
		{
			name:     "PostBundle-500",
			endpoint: "/500",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
				token: TokenDetail{
					ExpiresIn:   nowUtc.Add(1 * time.Hour),
					IssuedAt:    nowUtc,
					AccessToken: "access token",
				},
			},
			expectErr:   true,
			targetError: &ServerError{},
		},
		{
			name:     "PostBundle-404",
			endpoint: "/404",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
				token: TokenDetail{
					ExpiresIn:   nowUtc.Add(1 * time.Hour),
					IssuedAt:    nowUtc,
					AccessToken: "access token",
				},
			},
			expectErr:   true,
			targetError: &ResponseError{},
		},
		{
			name:     "PostBundle-401",
			endpoint: "/401",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
				token: TokenDetail{
					ExpiresIn:   nowUtc.Add(1 * time.Hour),
					IssuedAt:    nowUtc,
					AccessToken: "access token",
				},
			},
			expectErr:   true,
			targetError: &UnauthorizedError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.client.credential.BaseUrl = server.URL + tt.endpoint
			body, err := tt.client.PostBundle(context.Background(), "body")
			if tt.expectErr {
				assert.Error(t, err)
				assert.True(t, util.IsSameType(err, tt.targetError))

				// if UnauthorizedError returned token must be expired
				if util.IsSameType(err, &UnauthorizedError{}) {
					assert.True(t, tt.client.token.IsExpired())
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				assert.NotEmpty(t, body)
			}
		})
	}
}

func TestGetPatientId(t *testing.T) {
	resp := map[string]any{
		"entry": []map[string]any{
			{
				"fullUrl": "https://api-satusehat-stg.dto.kemkes.go.id/fhir-r4/v1/Patient/P02478375538",
				"resource": map[string]any{
					"id": "P02478375538",
					"name": []map[string]string{
						{
							"text": "pat***",
							"use":  "official",
						},
					},
				},
			},
		},
		"resourceType": "Bundle",
		"total":        1,
		"type":         "searchset",
	}

	respEmpty := map[string]any{
		"entry":        []map[string]any{},
		"resourceType": "Bundle",
		"total":        0,
		"type":         "searchset",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path[:4] {
		case "/200":
			w.WriteHeader(http.StatusOK)
			body, err := json.Marshal(resp)
			if err != nil {
				panic(err)
			}
			_, _ = w.Write(body)
		case "/204": //simulate response 200 but the data is empty
			w.WriteHeader(http.StatusOK)
			body, err := json.Marshal(respEmpty)
			if err != nil {
				panic(err)
			}
			_, _ = w.Write(body)
		case "/401":
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Unauthorized"))
		case "/500":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Internal Server Error"))
		case "/404":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Not Found"))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Not Found"))
		}
	}))
	defer server.Close()

	nowUtc := time.Now().UTC()

	tests := []struct {
		name        string
		endpoint    string
		client      *Client
		expectErr   bool
		targetError error
	}{
		{
			name:     "GetPatientId-200",
			endpoint: "/200",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
				token: TokenDetail{
					ExpiresIn:   nowUtc.Add(1 * time.Hour),
					IssuedAt:    nowUtc,
					AccessToken: "access token",
				},
			},
			expectErr:   false,
			targetError: nil,
		},
		{
			name:     "GetPatientId-204",
			endpoint: "/204",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
				token: TokenDetail{
					ExpiresIn:   nowUtc.Add(1 * time.Hour),
					IssuedAt:    nowUtc,
					AccessToken: "access token",
				},
			},
			expectErr:   true,
			targetError: &ResourceNotFoundError{},
		},
		{
			name:     "GetPatientId-500",
			endpoint: "/500",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
				token: TokenDetail{
					ExpiresIn:   nowUtc.Add(1 * time.Hour),
					IssuedAt:    nowUtc,
					AccessToken: "access token",
				},
			},
			expectErr:   true,
			targetError: &ServerError{},
		},
		{
			name:     "GetPatientId-404",
			endpoint: "/404",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
				token: TokenDetail{
					ExpiresIn:   nowUtc.Add(1 * time.Hour),
					IssuedAt:    nowUtc,
					AccessToken: "access token",
				},
			},
			expectErr:   true,
			targetError: &ResponseError{},
		},
		{
			name:     "GetPatientId-401",
			endpoint: "/401",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
				token: TokenDetail{
					ExpiresIn:   nowUtc.Add(1 * time.Hour),
					IssuedAt:    nowUtc,
					AccessToken: "access token",
				},
			},
			expectErr:   true,
			targetError: &UnauthorizedError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.client.credential.BaseUrl = server.URL + tt.endpoint
			id, err := tt.client.GetPatientId(context.Background(), "patientNik")
			if tt.expectErr {
				assert.Error(t, err)
				assert.True(t, util.IsSameType(err, tt.targetError))

				if util.IsSameType(err, &UnauthorizedError{}) {
					assert.True(t, tt.client.token.IsExpired())
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, "P02478375538", id)
			}
		})
	}
}

func TestGetPractitionerId(t *testing.T) {
	resp := map[string]any{
		"entry": []map[string]any{
			{
				"fullUrl": "https://api-satusehat-stg.dto.kemkes.go.id/fhir-r4/v1/Practitioner/N10000004",
				"resource": map[string]any{
					"id":           "N10000004",
					"name":         []map[string]string{{"text": "Pamela Educator, RN", "use": "official"}},
					"resourceType": "Practitioner",
				},
			},
		},
		"resourceType": "Bundle",
		"total":        1,
		"type":         "searchset",
	}

	respEmpty := map[string]any{
		"entry":        []map[string]any{},
		"resourceType": "Bundle",
		"total":        0,
		"type":         "searchset",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path[:4] {
		case "/200":
			w.WriteHeader(http.StatusOK)
			body, err := json.Marshal(resp)
			if err != nil {
				panic(err)
			}
			_, _ = w.Write(body)
		case "/204": //simulate response 200 but the data is empty
			w.WriteHeader(http.StatusOK)
			body, err := json.Marshal(respEmpty)
			if err != nil {
				panic(err)
			}
			_, _ = w.Write(body)
		case "/401":
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Unauthorized"))
		case "/500":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Internal Server Error"))
		case "/404":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Not Found"))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Not Found"))
		}
	}))
	defer server.Close()

	nowUtc := time.Now().UTC()

	tests := []struct {
		name        string
		endpoint    string
		client      *Client
		expectErr   bool
		targetError error
	}{
		{
			name:     "GetPractitionerId-200",
			endpoint: "/200",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
				token: TokenDetail{
					ExpiresIn:   nowUtc.Add(1 * time.Hour),
					IssuedAt:    nowUtc,
					AccessToken: "access token",
				},
			},
			expectErr:   false,
			targetError: nil,
		},
		{
			name:     "GetPractitionerId-204",
			endpoint: "/204",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
				token: TokenDetail{
					ExpiresIn:   nowUtc.Add(1 * time.Hour),
					IssuedAt:    nowUtc,
					AccessToken: "access token",
				},
			},
			expectErr:   true,
			targetError: &ResourceNotFoundError{},
		},
		{
			name:     "GetPractitionerId-500",
			endpoint: "/500",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
				token: TokenDetail{
					ExpiresIn:   nowUtc.Add(1 * time.Hour),
					IssuedAt:    nowUtc,
					AccessToken: "access token",
				},
			},
			expectErr:   true,
			targetError: &ServerError{},
		},
		{
			name:     "GetPractitionerId-404",
			endpoint: "/404",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
				token: TokenDetail{
					ExpiresIn:   nowUtc.Add(1 * time.Hour),
					IssuedAt:    nowUtc,
					AccessToken: "access token",
				},
			},
			expectErr:   true,
			targetError: &ResponseError{},
		},
		{
			name:     "GetPractitionerId-401",
			endpoint: "/401",
			client: &Client{
				restClient: resty.New(),
				restConfig: defaultRestConfig(),
				credential: defaultCredentials(),
				token: TokenDetail{
					ExpiresIn:   nowUtc.Add(1 * time.Hour),
					IssuedAt:    nowUtc,
					AccessToken: "access token",
				},
			},
			expectErr:   true,
			targetError: &UnauthorizedError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.client.credential.BaseUrl = server.URL + tt.endpoint
			id, err := tt.client.GetPractitionerId(context.Background(), "patientNik")
			if tt.expectErr {
				assert.Error(t, err)
				assert.True(t, util.IsSameType(err, tt.targetError))

				if util.IsSameType(err, &UnauthorizedError{}) {
					assert.True(t, tt.client.token.IsExpired())
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, "N10000004", id)
			}
		})
	}
}
