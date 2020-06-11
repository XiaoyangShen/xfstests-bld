package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	ServerLogPath = "/var/log/lgtm/lgtm.log"
	TestLogPath   = "/var/log"
	SecretPath    = "/etc/lighttpd/server.pem"
	CertPath      = "/root/xfstests_bld/kvm-xfstests/.gce_xfstests_cert.pem"
)

type Opts struct {
	no_region_shard string
}

type Respond struct {
	Status bool `json:"status"`
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("LTM test page"))
	fmt.Println("Hello World")
}

func login(w http.ResponseWriter, r *http.Request) {
	stat := Respond{true}
	js, err := json.Marshal(stat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	fmt.Println("log in here", string(js))
}

func runTests(w http.ResponseWriter, r *http.Request) {
	status := Respond{true}
	js, _ := json.Marshal(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func test() {
	arg := os.Args[1]
	branch := "master"
	if len(os.Args) > 2 {
		branch = os.Args[2]
	}
	switch arg {
	case "clone":
		gitClone("https://github.com/XiaoyangShen/spinner_test.git", "test")
	case "commit":
		id := gitCommit("test", branch)
		fmt.Println(id)
	case "pull":
		id := gitPull("test", branch)
		fmt.Println(id)
	case "watch":
		gitWatch("https://github.com/XiaoyangShen/spinner_test.git", "test", "master")
	}
}

// func init() {
// 	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
// }

func main() {
	test()
	return
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/gce-xfstests", runTests)
	err := http.ListenAndServeTLS(":443", CertPath, SecretPath, nil)
	if err != nil {
		log.Fatal("ListenandServer: ", err)
	}
}
