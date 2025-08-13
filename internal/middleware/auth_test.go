package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/auth"
	"github.com/espennoreng/go-http-rental-server/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {

	expectedIdentity := auth.Identity{
		UserID: "a53e4b0c-9d6c-4f7f-8c3b-5a1e2f3g4h5i",
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		identity, err := auth.FromContext(r.Context())

		assert.NoError(t, err, "auth.FromContext should not return an error")
		assert.Equal(t, expectedIdentity, identity, "The identity in the context should match the expected one")

		w.WriteHeader(http.StatusOK)
	})

	handlerToTest := middleware.AuthMiddleware(testHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	rr := httptest.NewRecorder()

	handlerToTest.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "The final handler should be called, returning a 200 OK status")
}
