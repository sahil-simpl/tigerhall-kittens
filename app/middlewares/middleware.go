package middlewares

import (
	"net/http"
	"tigerhallKittens/app/controllers"
)

const (
	serviceID        = "Service-Id"
	serviceNonce     = "Service-Nonce"
	serviceSignature = "Service-Signature"
)

type Middleware func(nextHandler controllers.Controller) controllers.Controller

func EmptyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
