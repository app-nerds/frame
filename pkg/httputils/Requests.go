package httputils

import "net/http"

/*
RealIP attempts to return the IP address of the caller. The result
will default to the RemoteAddr from http.Request. It will also
check the request headers for an "X-Forwarded-For" value and use
that. This is useful for when requests come through proxies or other
non-direct means.
*/
func RealIP(r *http.Request) string {
	result := r.RemoteAddr
	xForwardedFor := r.Header.Get("X-Forwarded-For")

	if xForwardedFor != "" {
		result = xForwardedFor
	}

	return result
}
