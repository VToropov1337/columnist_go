package middleware

import (
	"fmt"
	"net/http"
	"columnist_go/sessions"
)

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessions.Store.Get(r, "session")
		_, ok := session.Values["username"]
		fmt.Println("authrequired ===>>>>>",session.Values, session.Name())
		if !ok {
			http.Redirect(w,r,"/login",302)
			return
		}
		handler.ServeHTTP(w,r)
	}
}
