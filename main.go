package main

import (
	"encoding/json"
	. "github.com/everfore/oauth/oauth2"
	"github.com/shaalx/membership/dbu"
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	usage = []byte(`<a href="https://github.com/everfore/oauth" ><h1>@oauth2</h1></a>` + "\n" + `
		<a href="/signin" ><h1>@sign in#github#</h1></a>`)
	OA *OAGithub

	MgoDB   = dbu.NewMgoDB("daocloud")
	usersC  = MgoDB.GetCollection([]string{"lEyTj8hYrUIKgMfi", "users"}...)
	onlineC = MgoDB.GetCollection([]string{"lEyTj8hYrUIKgMfi", "online"}...)
	vcountC = MgoDB.GetCollection([]string{"lEyTj8hYrUIKgMfi", "vcount"}...)
)

func init() {
	OA = NewOAGithub("8ba2991113e68b4805c1", "b551e8a640d53904d82f95ae0d84915ba4dc0571", "user")
}

func main() {
	log.Println("ready...")
	http.HandleFunc("/db", dbH)
	http.HandleFunc("/ddb", ddbH)
	http.HandleFunc("/", root)
	http.HandleFunc("/signin", signin)
	http.HandleFunc("/callback", callback)
	err := http.ListenAndServe(":80", nil)
	if check_err(err) {
		return
	}
}

func dbH(rw http.ResponseWriter, req *http.Request) {
	n := usersC.Count(nil)
	log.Println(n)
	v := usersC.ISelect(nil)
	log.Println(v)
}

func ddbH(rw http.ResponseWriter, req *http.Request) {
	usersC.C.DropCollection()
	onlineC.C.DropCollection()
	vcountC.C.DropCollection()
}

func root(rw http.ResponseWriter, req *http.Request) {
	rw.Write(usage)
}

func signin(rw http.ResponseWriter, req *http.Request) {
	http.Redirect(rw, req, OA.AuthURL(), 302)
}

func callback(rw http.ResponseWriter, req *http.Request) {
	log.Printf("%s\n", req.RemoteAddr)
	b, err := OA.NextStep(req)
	if nil != err {
		rw.Write([]byte(err.Error()))
		return
	}
	var ret map[string]interface{}
	err = json.Unmarshal(b, &ret)
	if nil == err {
		t := template.New("index.html")
		t, err := t.ParseFiles("index.html")
		if nil != err {
			return
		}
		now := time.Now().String()
		ret["now"] = now
		t.Execute(rw, ret)
	} else {
		rw.Write([]byte(err.Error()))
	}
}

func check_err(err error) bool {
	if nil != err {
		return true
	}
	return false
}
