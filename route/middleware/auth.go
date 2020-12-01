package middleware

import (
	"net/http"

	"context"

	firebase "firebase.google.com/go"
	auth "firebase.google.com/go/auth"

	option "google.golang.org/api/option"

	c "github.com/robitx/inceptus/route/ctx"
)

func setupFirebase(credentialsFile string) *auth.Client {
opt := option.WithCredentialsFile(credentialsFile)
//Firebase admin SDK initialization
 app, err := firebase.NewApp(context.Background(), nil, opt)
 if err != nil {
  panic("Firebase load error")
 }
//Firebase Auth
 auth, err := app.Auth(context.Background())
 if err != nil {
  panic("Firebase load error")
 }
 return auth
}


// Auth middleware which uses Firebase auth tokens
// and puts userID into ctx for handlers to decide what
// can be accessed by given user
func Auth(
  realm string,
  authHeader string,
  credentialsFile string,
  ) func(next http.Handler) http.Handler {

  auther := setupFirebase(credentialsFile)

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
      rawToken := r.Header.Get(authHeader)
      token, err := auther.VerifyIDToken(r.Context(), rawToken)
      var uid string
      if err == nil {
        uid = token.UID
      }
      // if err != nil {
      //   w.Header().Add("WWW-Authenticate",
      //     fmt.Sprintf(`Bearer realm="%s"`, realm))
      //   w.WriteHeader(http.StatusUnauthorized)
      //   w.Write([]byte("You're not authorized"))
      //   return
      // }

      r = c.SetUserID(uid, r)

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}