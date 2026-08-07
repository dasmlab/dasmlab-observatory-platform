package auth

import (
	"net/http"
	"os"
	"strings"
)

// Middleware optionally gates API routes when DPO_AUTH=oidc (Track B stub).
// Default is open for pilot until Keycloak client `dpo` is wired (ADR-0007).
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mode := strings.ToLower(strings.TrimSpace(os.Getenv("DPO_AUTH")))
		if mode != "oidc" {
			next.ServeHTTP(w, r)
			return
		}
		// Stub: reject unauthenticated API until real OIDC validation lands.
		if strings.HasPrefix(r.URL.Path, "/api/") && r.Header.Get("Authorization") == "" {
			http.Error(w, "oidc required (DPO_AUTH=oidc stub)", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
