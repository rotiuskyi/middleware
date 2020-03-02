package middleware

import (
	"net/http"
	"strings"
)

// Middleware is a composable http.HandlerFunc wrapper
type Middleware func(next http.HandlerFunc) http.HandlerFunc

// Compose http.HandlerFunc from Middleware array
func Compose(middlewareArr ...Middleware) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if len(middlewareArr) == 0 {
			return
		}

		hf := func(http.ResponseWriter, *http.Request) { /* noop */ }

		lastI := len(middlewareArr) - 1
		last, body := middlewareArr[lastI], middlewareArr[:lastI]

		hf = last(hf)

		for i := len(body) - 1; i >= 0; i-- {
			hf = body[i](hf)
		}

		hf(rw, req)
	}
}

// Post method filter
func Post(hf http.HandlerFunc) Middleware {
	return filterMethod("POST", hf)
}

// Get method filter
func Get(hf http.HandlerFunc) Middleware {
	return filterMethod("GET", hf)
}

// Put method filter
func Put(hf http.HandlerFunc) Middleware {
	return filterMethod("PUT", hf)
}

// Delete method filter
func Delete(hf http.HandlerFunc) Middleware {
	return filterMethod("DELETE", hf)
}

func filterMethod(method string, hf http.HandlerFunc) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			if strings.ToUpper(req.Method) == method {
				hf(rw, req)
			}
			next(rw, req)
		}
	}
}
