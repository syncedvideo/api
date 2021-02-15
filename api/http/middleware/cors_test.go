package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorsMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		return
	})

	t.Run("add required headers", func(t *testing.T) {
		rr := httptest.NewRecorder()
		CorsMiddleware(handler).ServeHTTP(rr, newRequest(http.MethodGet))
		for key, val := range requiredHeaders() {
			want := val
			got := rr.HeaderMap.Get(key)
			if want != got {
				t.Errorf(`%s: want "%s", got "%s"`, key, want, got)
			}
		}
	})

	t.Run("options request returns 200", func(t *testing.T) {
		rr := httptest.NewRecorder()
		CorsMiddleware(handler).ServeHTTP(rr, newRequest(http.MethodOptions))
		want := http.StatusOK
		got := rr.Result().StatusCode
		if want != got {
			t.Errorf("want %v, got %v", want, got)
		}
	})
}

func newRequest(method string) *http.Request {
	req, _ := http.NewRequest(method, "/test", nil)
	req.Header.Set("Origin", "test")
	return req
}

func requiredHeaders() map[string]string {
	h := make(map[string]string)
	h[allowOriginHeaderKey] = "test"
	h[allowMethodsHeaderKey] = allowedMethods
	h[allowHeadersHeaderKey] = allowedHeaders
	h[allowCredentialsHeaderKey] = "true"
	return h
}
