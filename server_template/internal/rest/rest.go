package rest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	c "github.com/robitx/inceptus/route/ctx"
)

// NotFound ..
func NotFound(file string) func(w http.ResponseWriter, r *http.Request) {
	workDir, _ := os.Getwd()
	file = filepath.Join(workDir, file)
	response, _ := ioutil.ReadFile(file)

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf(`{"fuck": "you"}`)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(404)
		w.Write(response)
		return
	}

}

// Echo ..
func Echo(w http.ResponseWriter, r *http.Request) {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(r.Body)
	w.Write(buffer.Bytes())
}

// HiUser ..
func HiUser(w http.ResponseWriter, r *http.Request) {
	userID := c.GetUserID(r)
	w.Write([]byte(fmt.Sprintf("Hi user: %v\n", userID)))
}

// Ping ..
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

// Panic ..
func Panic(w http.ResponseWriter, r *http.Request) {
	panic("for testing..")
}
