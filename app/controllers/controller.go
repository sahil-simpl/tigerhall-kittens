package controllers

import (
	"net/http"
	"tigerhallKittens/app/lib/web"
)

type Controller func(request *web.Request, w http.ResponseWriter) (*web.JSONResponse, web.ErrorInterface)
