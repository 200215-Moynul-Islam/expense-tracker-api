package middlewares

import (
	"net/http"
	"strconv"

	beegoCtx "github.com/beego/beego/v2/server/web/context"
)

func RequireAuthentication(ctx *beegoCtx.Context) {
	userIDStr := ctx.Request.Header.Get("X-User-ID")

	if userIDStr == "" {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]interface{}{
			"success": false,
			"message": "Unauthorized.",
		}, false, false)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]interface{}{
			"success": false,
			"message": "Unauthorized.",
		}, false, false)
		return
	}

	ctx.Input.SetData("userID", userID)
}