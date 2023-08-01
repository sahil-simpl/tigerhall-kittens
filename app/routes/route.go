package routes

import (
	"github.com/newrelic/go-agent/v3/integrations/nrhttprouter"
)

func Init(router *nrhttprouter.Router) {
	router.SaveMatchedRoutePath = true

}
