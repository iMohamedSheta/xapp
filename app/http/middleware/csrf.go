package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/iMohamedSheta/xerr"
	"github.com/imohamedsheta/xapp/app/domain/enums"
	"github.com/imohamedsheta/xapp/app/domain/utils"
	"github.com/imohamedsheta/xapp/app/x"
)

var (
	csrfMd   func(http.Handler) http.Handler
	onceCSRF sync.Once
)

func CSRF() gin.HandlerFunc {
	httpsEnabled := utils.IsHttpsEnabled()

	onceCSRF.Do(func() {
		csrfMd = csrf.Protect(
			[]byte(utils.GetAppSecret()),
			csrf.CookieName("X-CSRF-TOKEN"),
			csrf.RequestHeader("X-XSRF-TOKEN"),
			csrf.MaxAge(0),
			csrf.Secure(httpsEnabled), // set to true if you're using https (should be true in production)
			csrf.SameSite(csrf.SameSiteLaxMode),
			csrf.TrustedOrigins(x.Config().GetArrayOfStrings("cors.origin", []string{"*"})),
			csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				RequestErrorLogger(w, r, xerr.New("invalid request CSRF token invalid", enums.XErrForbiddenError, nil))
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte(`{"message": "FORBIDDEN - invalid CSRF token"}`))
			})),
		)
	})

	middleware := func(next http.Handler) http.Handler {
		// Wrap next with csrfMd, then set the token cookie after the CSRF middleware runs.
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			csrfMd(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Now token is available in the request context
				if r.Method == http.MethodGet || r.Method == http.MethodHead {
					token := csrf.Token(r)

					tokenHash := sha256.Sum256([]byte(token))
					utils.Dump("CSRF token verified for request " + r.Method + " " + r.URL.Path + " token hash: " + hex.EncodeToString(tokenHash[:]))

					http.SetCookie(w, &http.Cookie{
						Name:     "XSRF-TOKEN",
						Value:    token,
						Path:     "/",
						HttpOnly: false,
						Secure:   httpsEnabled,
						SameSite: http.SameSiteLaxMode,
						MaxAge:   0,
					})
				}

				next.ServeHTTP(w, r)
			})).ServeHTTP(w, r)
		})
	}

	return adapter.Wrap(middleware)
}
