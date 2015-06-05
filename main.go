package main

import (
	"encoding/json"
	"fmt"
	. "github.com/shaalx/echo/oauth2"
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	usage = []byte(`<h1>This is a app from@zhuulx</h1><a href="https://github.com/zhuulx/oauth" ><h1>https://github.com/zhuulx/oauth</h1></a>` + "\n" + `
		<a href="/signin" ><h1>GITHUB OAUTH2</h1></a>`)
	OA *OAGithub
)

func init() {
	OA = NewOAGithub("8ba2991113e68b4805c1", "b551e8a640d53904d82f95ae0d84915ba4dc0571", "user")
}

func main() {
	log.Println("ready...")
	http.HandleFunc("/", site)
	http.HandleFunc("/signin", signin)
	http.HandleFunc("/site", site)
	http.HandleFunc("/callback", callback)
	err := http.ListenAndServe(":80", nil)
	if check_err(err) {
		return
	}
	log.Println("go...")
}

func echo(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("[ECHO]"))
	q := req.URL.Query()
	b, err := json.Marshal(q)
	if check_err(err) {
		rw.Write([]byte("echo ： error"))
		return
	}
	rw.Write(b)
}

func root(rw http.ResponseWriter, req *http.Request) {
	log.Println(req.URL)
	rw.Write([]byte("[ROOT]" + time.Now().String()))
}

func site(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte(usage))
}

func signin(rw http.ResponseWriter, req *http.Request) {
	http.Redirect(rw, req, OA.AuthURL(), 302)
}

func callback(rw http.ResponseWriter, req *http.Request) {
	b := OA.NextStep(req)
	var ret map[string]interface{}
	err := json.Unmarshal(b, &ret)
	if nil == err {
		avatar, ok := ret["avatar_url"]
		if ok {
			avatar_url := fmt.Sprintf("%v", avatar)
			t := template.New("index.tpl")
			t, err := t.ParseFiles("index.tpl")
			if nil != err {
				rw.Write(b)
				return
			}
			rw.Write(b)
			t.Execute(rw, avatar_url)
		}
	}
}

func check_err(err error) bool {
	if nil != err {
		return true
	}
	return false
}
