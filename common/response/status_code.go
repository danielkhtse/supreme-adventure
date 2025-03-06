package response

import (
	"net/http"
)

// StatusCode represents HTTP status codes
type StatusCode int

const (
	StatusContinue                     StatusCode = http.StatusContinue                     // 100
	StatusSwitchingProtocols           StatusCode = http.StatusSwitchingProtocols           // 101
	StatusOK                           StatusCode = http.StatusOK                           // 200
	StatusCreated                      StatusCode = http.StatusCreated                      // 201
	StatusAccepted                     StatusCode = http.StatusAccepted                     // 202
	StatusNonAuthoritativeInfo         StatusCode = http.StatusNonAuthoritativeInfo         // 203
	StatusNoContent                    StatusCode = http.StatusNoContent                    // 204
	StatusResetContent                 StatusCode = http.StatusResetContent                 // 205
	StatusPartialContent               StatusCode = http.StatusPartialContent               // 206
	StatusMultipleChoices              StatusCode = http.StatusMultipleChoices              // 300
	StatusMovedPermanently             StatusCode = http.StatusMovedPermanently             // 301
	StatusFound                        StatusCode = http.StatusFound                        // 302
	StatusSeeOther                     StatusCode = http.StatusSeeOther                     // 303
	StatusNotModified                  StatusCode = http.StatusNotModified                  // 304
	StatusUseProxy                     StatusCode = http.StatusUseProxy                     // 305
	StatusTemporaryRedirect            StatusCode = http.StatusTemporaryRedirect            // 307
	StatusPermanentRedirect            StatusCode = http.StatusPermanentRedirect            // 308
	StatusBadRequest                   StatusCode = http.StatusBadRequest                   // 400
	StatusUnauthorized                 StatusCode = http.StatusUnauthorized                 // 401
	StatusPaymentRequired              StatusCode = http.StatusPaymentRequired              // 402
	StatusForbidden                    StatusCode = http.StatusForbidden                    // 403
	StatusNotFound                     StatusCode = http.StatusNotFound                     // 404
	StatusMethodNotAllowed             StatusCode = http.StatusMethodNotAllowed             // 405
	StatusNotAcceptable                StatusCode = http.StatusNotAcceptable                // 406
	StatusProxyAuthRequired            StatusCode = http.StatusProxyAuthRequired            // 407
	StatusRequestTimeout               StatusCode = http.StatusRequestTimeout               // 408
	StatusConflict                     StatusCode = http.StatusConflict                     // 409
	StatusGone                         StatusCode = http.StatusGone                         // 410
	StatusLengthRequired               StatusCode = http.StatusLengthRequired               // 411
	StatusPreconditionFailed           StatusCode = http.StatusPreconditionFailed           // 412
	StatusRequestEntityTooLarge        StatusCode = http.StatusRequestEntityTooLarge        // 413
	StatusRequestURITooLong            StatusCode = http.StatusRequestURITooLong            // 414
	StatusUnsupportedMediaType         StatusCode = http.StatusUnsupportedMediaType         // 415
	StatusRequestedRangeNotSatisfiable StatusCode = http.StatusRequestedRangeNotSatisfiable // 416
	StatusExpectationFailed            StatusCode = http.StatusExpectationFailed            // 417
	StatusTeapot                       StatusCode = http.StatusTeapot                       // 418
	StatusMisdirectedRequest           StatusCode = http.StatusMisdirectedRequest           // 421
	StatusUnprocessableEntity          StatusCode = http.StatusUnprocessableEntity          // 422
	StatusLocked                       StatusCode = http.StatusLocked                       // 423
	StatusFailedDependency             StatusCode = http.StatusFailedDependency             // 424
	StatusTooEarly                     StatusCode = http.StatusTooEarly                     // 425
	StatusUpgradeRequired              StatusCode = http.StatusUpgradeRequired              // 426
	StatusPreconditionRequired         StatusCode = http.StatusPreconditionRequired         // 428
	StatusTooManyRequests              StatusCode = http.StatusTooManyRequests              // 429
	StatusRequestHeaderFieldsTooLarge  StatusCode = http.StatusRequestHeaderFieldsTooLarge  // 431
	StatusUnavailableForLegalReasons   StatusCode = http.StatusUnavailableForLegalReasons   // 451
	StatusInternalServerError          StatusCode = http.StatusInternalServerError          // 500
	StatusNotImplemented               StatusCode = http.StatusNotImplemented               // 501
	StatusBadGateway                   StatusCode = http.StatusBadGateway                   // 502
	StatusServiceUnavailable           StatusCode = http.StatusServiceUnavailable           // 503
	StatusGatewayTimeout               StatusCode = http.StatusGatewayTimeout               // 504
	StatusHTTPVersionNotSupported      StatusCode = http.StatusHTTPVersionNotSupported      // 505
	StatusVariantAlsoNegotiates        StatusCode = http.StatusVariantAlsoNegotiates        // 506
	StatusInsufficientStorage          StatusCode = http.StatusInsufficientStorage          // 507
	StatusLoopDetected                 StatusCode = http.StatusLoopDetected                 // 508
	StatusNotExtended                  StatusCode = http.StatusNotExtended                  // 510
)
