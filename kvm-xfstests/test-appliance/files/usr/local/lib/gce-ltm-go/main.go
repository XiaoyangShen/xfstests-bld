package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/gce-xfstests", runTests)
	err := http.ListenAndServeTLS(":443", CertPath, SecretPath, nil)
	if err != nil {
		log.Fatal("ListenandServer: ", err)
	}
}
