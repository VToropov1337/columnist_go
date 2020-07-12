package sessions

import "github.com/gorilla/sessions"

var Store = sessions.NewCookieStore([]byte("SESSION_KEY_s3cr3t"))

