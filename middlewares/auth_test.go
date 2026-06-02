package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	beegoCtx "github.com/beego/beego/v2/server/web/context"
)

// newContext builds a minimal beego context wrapping a real http.Request.
func newContext(t *testing.T, userIDHeader string) *beegoCtx.Context {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/expenses", nil)
	if userIDHeader != "" {
		req.Header.Set("X-User-ID", userIDHeader)
	}

	w := httptest.NewRecorder()

	ctx := beegoCtx.NewContext()
	ctx.Reset(w, req)

	return ctx
}

func TestRequireAuthentication(t *testing.T) {
	tests := []struct {
		name         string
		userIDHeader string
		wantStatus   int
		wantUserID   int
		wantBlocked  bool
	}{
		{
			name:         "missing header is rejected",
			userIDHeader: "",
			wantStatus:   http.StatusUnauthorized,
			wantBlocked:  true,
		},
		{
			name:         "non-numeric value is rejected",
			userIDHeader: "abc",
			wantStatus:   http.StatusUnauthorized,
			wantBlocked:  true,
		},
		{
			name:         "zero is rejected",
			userIDHeader: "0",
			wantStatus:   http.StatusUnauthorized,
			wantBlocked:  true,
		},
		{
			name:         "negative value is rejected",
			userIDHeader: "-1",
			wantStatus:   http.StatusUnauthorized,
			wantBlocked:  true,
		},
		{
			name:        "valid user id passes through",
			userIDHeader: "1",
			wantBlocked:  false,
			wantUserID:   1,
		},
		{
			name:        "large valid user id passes through",
			userIDHeader: "9999",
			wantBlocked:  false,
			wantUserID:   9999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.wantBlocked {
				w := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/api/v1/expenses", nil)

				if tt.userIDHeader != "" {
					req.Header.Set("X-User-ID", tt.userIDHeader)
				}

				ctx := beegoCtx.NewContext()
				ctx.Reset(w, req)

				RequireAuthentication(ctx)

				if w.Code != http.StatusUnauthorized {
					t.Errorf("expected 401, got %d", w.Code)
				}

				return
			}

			ctx := newContext(t, tt.userIDHeader)

			RequireAuthentication(ctx)

			userID, ok := ctx.Input.GetData("userID").(int)
			if !ok {
				t.Fatal("userID not set in context data after successful auth")
			}

			if userID != tt.wantUserID {
				t.Errorf("userID: got %d, want %d", userID, tt.wantUserID)
			}
		})
	}
}