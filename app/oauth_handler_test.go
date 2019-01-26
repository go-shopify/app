package app

import "net/http"

// Make sure handlerImpl implements http.Handler.
var _ http.Handler = oauthHandlerImpl{}
