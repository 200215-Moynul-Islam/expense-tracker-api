package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	beegoCtx "github.com/beego/beego/v2/server/web/context"
)

func newBaseCtx() (*BaseController, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	ctx := beegoCtx.NewContext()
	ctx.Reset(w, req)

	c := &BaseController{}
	c.Ctx = ctx
	c.Data = make(map[interface{}]interface{})

	return c, w
}

func TestBaseController_Success(t *testing.T) {
	c, w := newBaseCtx()

	c.Success(http.StatusOK, "ok message", map[string]string{
		"key": "value",
	})

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if body["success"] != true {
		t.Errorf("expected success true, got %v", body["success"])
	}

	if body["message"] != "ok message" {
		t.Errorf("unexpected message: %v", body["message"])
	}

	if body["data"] == nil {
		t.Errorf("expected data field")
	}
}

func TestBaseController_Error(t *testing.T) {
	c, w := newBaseCtx()

	c.Error(http.StatusBadRequest, "error message")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if body["success"] != false {
		t.Errorf("expected success false, got %v", body["success"])
	}

	if body["message"] != "error message" {
		t.Errorf("unexpected message: %v", body["message"])
	}
}