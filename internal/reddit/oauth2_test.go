package reddit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewOAuth2Client(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, requestCount int)
		expectError    bool
		retryCount     int
	}{
		{
			name: "successful client creation",
			serverResponse: func(w http.ResponseWriter, _ int) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(
					[]byte(`{
					"access_token": "test_token",
					"token_type": "bearer",
					"expires_in": 3600
				}`),
				)
			},
			expectError: false,
			retryCount:  0,
		},
		{
			name: "rate limited then success",
			serverResponse: func(w http.ResponseWriter, requestCount int) {
				// Return rate limit on first attempt.
				if requestCount == 1 {
					w.Header().Set("Request-Count", "1")
					w.Header().Set("Retry-After", "1")
					w.WriteHeader(429)
					return
				}

				// Return success on second attempt.
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(
					[]byte(`{
					"access_token": "test_token",
					"token_type": "bearer",
					"expires_in": 3600
				}`),
				)
			},
			expectError: false,
			retryCount:  1,
		},
		{
			name: "authentication error",
			serverResponse: func(w http.ResponseWriter, _ int) {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error": "invalid_client"}`))
			},
			expectError: true,
			retryCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				requestCount := 0
				server := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							// Verify the User-Agent header is set.
							require.Equal(t, "TestUserAgent", r.Header.Get("User-Agent"))
							requestCount++
							tt.serverResponse(w, requestCount)
						},
					),
				)
				defer server.Close()

				cfg := Config{
					BaseURL:      server.URL,
					ClientID:     "test_client_id",
					ClientSecret: "test_client_secret",
					Username:     "test_username",
					Password:     "test_password",
					UserAgent:    "TestUserAgent",
				}

				client, err := NewOAuth2Client(context.Background(), cfg)

				if tt.expectError {
					require.Error(t, err)
					require.Nil(t, client)
				} else {
					require.NoError(t, err)
					require.NotNil(t, client)
				}

				// Verify the number of retries.
				require.Equal(t, tt.retryCount+1, requestCount)
			},
		)
	}
}

func TestUserAgentTransport(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Verify the User-Agent header was set correctly.
				require.Equal(t, "TestUserAgent", r.Header.Get("User-Agent"))
				w.WriteHeader(http.StatusOK)
			},
		),
	)
	defer server.Close()

	transport := &userAgentTransport{
		base:      http.DefaultTransport,
		userAgent: "TestUserAgent",
	}

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	resp, err := transport.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
