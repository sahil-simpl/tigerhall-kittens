package middlewares

import (
	"bytes"
	"context"
	"fmt"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tuvistavie/securerandom"

	"tigerhallKittens/app/controllers"
	"tigerhallKittens/app/lib/logger"
	"tigerhallKittens/app/lib/web"
)

const (
	APIVersionV1 = 1
)

type ResponseBuilder func(data *web.JSONResponse, responseErr web.ErrorInterface) *web.JSONResponse

func ServeV1Endpoint(middleware Middleware, handler controllers.Controller) httprouter.Handle {
	return serve(buildResponseBuilder(APIVersionV1), middleware(handler))
}

func buildResponseBuilder(version int) ResponseBuilder {
	return func(data *web.JSONResponse, err web.ErrorInterface) *web.JSONResponse {
		if err == nil {
			return successResponse(version, data)
		} else {
			return errorResponse(version, err)
		}
	}
}

func RequestHeaderId(req *http.Request) string {
	var requestId string
	var err error

	requestId, err = securerandom.Uuid()
	if err != nil {
		logger.E(req.Context(), err, "endpoint/RequestHeaderId")
	}
	return requestId
}

func serve(responseBuilder ResponseBuilder, handler controllers.Controller) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		webReq := web.NewRequest(req)
		for i := range ps {
			webReq.SetPathParam(ps[i].Key, ps[i].Value)
		}

		_, decodeErr := readRequestBody(req)
		if decodeErr != nil {
			logger.I(req.Context(), "Error while decoding the request body",
				logger.Field("error", decodeErr.Error()))
		}

		defer func() {
			if recvr := recover(); recvr != nil {
				errorMessage := fmt.Sprintf("%v", recvr)
				err := web.NewError(web.InternalServerError, errorMessage,
					"", http.StatusInternalServerError)
				w.WriteHeader(err.HTTPStatusCode())
				writeResponse(req.Context(), w, responseBuilder(nil, err))
			}
		}()

		data, responseErr := handler(&webReq, w)
		responseCode := responseCode(responseErr)

		// setting response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(responseCode)

		writeResponse(req.Context(), w, responseBuilder(data, responseErr))

		if responseCode >= http.StatusInternalServerError {
			return
		}

		if responseCode >= http.StatusBadRequest {
			return
		}

	}
}

func writeResponse(ctx context.Context, w http.ResponseWriter, response *web.JSONResponse) {
	_, err := w.Write(response.ByteArray(ctx))
	if err != nil {
		logger.E(ctx, err, "error in writing response", logger.Field("error", err.Error()))
	}
}

func responseCode(err web.ErrorInterface) int {
	if err != nil {
		return err.HTTPStatusCode()
	}
	return http.StatusOK
}

func readRequestBody(req *http.Request) (map[string]interface{}, error) {
	var jsonPayload map[string]interface{}

	if req.Body != nil {
		body, _ := ioutil.ReadAll(req.Body)
		if len(body) != 0 {
			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			err := json.Unmarshal(body, &jsonPayload)
			if err != nil {
				return nil, err
			}
		}
	}

	return jsonPayload, nil
}

func successResponse(version int, data *web.JSONResponse) *web.JSONResponse {
	return &web.JSONResponse{
		"success":     true,
		"data":        data,
		"api_version": version,
	}
}

func errorResponse(version int, err web.ErrorInterface) *web.JSONResponse {
	return &web.JSONResponse{
		"success": false,
		"error": map[string]interface{}{
			"code":    err.Code(),
			"message": err.Description(),
		},
		"api_version": version,
	}
}
