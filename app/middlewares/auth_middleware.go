package middlewares

import (
	"fmt"
	"log"
	"strings"

	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"net/http"

	"tigerhallKittens/app"
	"tigerhallKittens/app/controllers"
	"tigerhallKittens/app/lib/web"
	"tigerhallKittens/app/utils"
)

func AuthenticateWithHmacDigestMiddleware(next controllers.Controller) controllers.Controller {
	return func(request *web.Request, w http.ResponseWriter) (*web.JSONResponse, web.ErrorInterface) {
		serviceID := getReqHeader(request, serviceID)
		nonce := getReqHeader(request, serviceNonce)
		serviceSignature := getReqHeader(request, serviceSignature)

		if !validateHmacDigest(serviceID, nonce, serviceSignature, utils.GetMapEnvConfig(app.Env.AllowedOrigins)) {
			return nil, web.NewError(web.UnauthorizedRequest, "Failed to authenticate",
				fmt.Sprintf("serviceID:%s serviceID is not whitelisted", serviceID), http.StatusUnauthorized)
		}
		log.Println(request.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		return next(request, w)
	}
}

func getReqHeader(r *web.Request, key string) string {
	if len(r.Header[key]) < 1 {
		return ""
	}

	return r.Header[key][0]
}

func validateHmacDigest(serviceID, nonce, serviceSignature string, serviceConfig map[string]string) bool {
	if serviceID == "" || nonce == "" || serviceSignature == "" {
		return false
	}

	message := strings.Join([]string{nonce, serviceID}, "-")
	digest := CalculateHmacSignature(serviceConfig[serviceID], message)

	return digest == serviceSignature
}

func CalculateHmacSignature(key, data string) string {
	hash := hmac.New(sha1.New, []byte(key))
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
