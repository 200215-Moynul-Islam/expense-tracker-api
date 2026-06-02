package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	beegoCtx "github.com/beego/beego/v2/server/web/context"
)

func TestHealthController_Get(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	ctx := beegoCtx.NewContext()
	ctx.Reset(rec, req)

	controller := &HealthController{}
	controller.Init(ctx, "HealthController", "", controller)

	controller.Get()

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	if rec.Body == nil {
		t.Fatal("expected response body, got nil")
	}
}
