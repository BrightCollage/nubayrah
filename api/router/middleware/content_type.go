package middleware

import "net/http"

const (
	HeaderKeyContentType       = "Content-Type"
	HeaderValueContentTypeJSON = "application/json;charset=utf8"
	HeaderValueContentTypePNG  = "image/png"
)

func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderKeyContentType, HeaderValueContentTypeJSON)
		next.ServeHTTP(w, r)
	})
}

func ContentTypePNG(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderKeyContentType, HeaderValueContentTypePNG)
		next.ServeHTTP(w, r)
	})
}
