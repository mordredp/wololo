package auth

import (
	"net/http"
	"time"
)

// Clean removes expired sessions
func Clean(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if time.Now().After(lastCleanup.Add(maxSessionLength / 2)) {

			for key, s := range sessions {
				if s.isExpired() {
					delete(sessions, key)
				}
			}
			lastCleanup = time.Now()
		}

		next.ServeHTTP(w, r)
	})
}
