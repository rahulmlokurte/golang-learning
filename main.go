package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var users = make(map[string]struct{})

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := users[req.Username]; ok {
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	} else {
		users[req.Username] = struct{}{}
	}
}

func hasTimedOut(err error) bool {
	switch err := err.(type) {
	case *url.Error:
		if err, ok := err.Err.(net.Error); ok && err.Timeout() {
			return true
		}

	case net.Error:
		if err.Timeout() {
			return true
		}

	case *net.OpError:
		if err.Timeout() {
			return true
		}
	}

	errTxt := "use of closed connection network"

	if err != nil && strings.Contains(err.Error(), errTxt) {
		return true
	}
	return false
}

func download(location string, file *os.File, retries int64) error {
	req, err := http.NewRequest(http.MethodGet, location, nil)
	fi, err := file.Stat()
	current := fi.Size()
	if current > 0 {
		start := strconv.FormatInt(current, 10)
		req.Header.Set("Range", "bytes="+start+"-")
	}

	cc := &http.Client{Timeout: 10}
	res, err := cc.Do(req)
	if err != nil && hasTimedOut(err) {
		if retries > 0 {
			return download(location, file, retries-1)
		}
	} else if err != nil {
		return err
	}

	_, err = io.Copy(file, res.Body)
	if err != nil && hasTimedOut(err) {
		if retries > 0 {
			return download(location, file, retries-1)
		}
		return err
	} else if err != nil {
		return err
	}
	return nil
}

type Error struct {
	HTTPCode int    `json:"-"`
	Code     int    `json:"code,omitempty"`
	Message  string `json:"message"`
}

func displayError(w http.ResponseWriter, r *http.Request) {
	e := Error{
		HTTPCode: http.StatusForbidden,
		Code:     123,
		Message:  "An Error Occured",
	}
	JSONError(w, e)
}

func JSONError(w http.ResponseWriter, e Error) {
	data := struct {
		Err Error `json:"error"`
	}{e}
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.HTTPCode)
	fmt.Fprint(w, string(b))
}

func main() {
	http.HandleFunc("/register", handleRegister)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
