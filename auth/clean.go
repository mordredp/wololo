package auth

import (
	"log"
	"net/http"
	"time"
)

// Clean removes expired sessions
func Clean(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if time.Now().After(lastCleanup.Add(maxSessionLength / 2)) {

			deletedCount := 0
			for id, session := range sessions {
				if session.isExpired() {
					delete(sessions, id)
					deletedCount++
				}
			}
			lastCleanup = time.Now()

			log.Printf("removed %d expired sessions", deletedCount)
		}

		next.ServeHTTP(w, r)
	})
}
