package api

import (
	"io"
	"net/http"
)

// MockAuthenticatedRequestForTest creates a new authenticated request for testing purposes.
func MockAuthenticatedRequestForTest(method, target, userID string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, target, body)
	req.Header.Set("X-Test-User-ID", userID)
	return req
}
