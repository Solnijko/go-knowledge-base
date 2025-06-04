package auth

import "net/http"

func AuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/login", LoginHandler)
	mux.HandleFunc("/api/protected", ProtectedHandler)
}
