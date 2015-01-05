package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/coreos/go-etcd/etcd"

	"github.com/igneous-systems/beta/shared"
)

const (
	USERPROP = "postgres/host_count/username"
	PWDPROP  = "postgres/host_count/password"
)

var (
	peers = []string{
		"http://localhost:4001",
		"http://etcd:4001", //for local dev case
	}
)

func post(w http.ResponseWriter, r *http.Request) {
	log.Printf("POST API call %+v", r)

	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var payload shared.ApiPayload

	if err := dec.Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payload.Username = strings.TrimSpace(payload.Username)
	payload.Password = strings.TrimSpace(payload.Password) //dangerous

	if err := WriteKV(USERPROP, payload.Username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := WriteKV(PWDPROP, payload.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//we return 200 from here, use {} to make the jquery api happy
	fmt.Fprint(w, "{}")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		post(w, r)
	} else {
		get(w, r)
	}
}

func ReadKV(name string) (string, error) {
	client := etcd.NewClient(peers)
	resp, err := client.Get(name, false, false)
	if err != nil {
		//special case not found
		if err.(*etcd.EtcdError).ErrorCode == 100 {
			return "", nil
		}
		return "", err
	}
	return resp.Node.Value, nil
}

func WriteKV(name string, value string) error {
	client := etcd.NewClient(peers)
	_, err := client.Set(name, value, 0)
	if err != nil {
		return err
	}
	return nil
}

func get(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET API call %+v", r)

	var result shared.ApiPayload
	resp, err := ReadKV(USERPROP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if resp == "" {
		http.NotFound(w, r)
		return
	}
	result.Username = resp

	resp, err = ReadKV(PWDPROP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if resp == "" {
		result.Password = ""
	} else {
		result.Password = resp
	}

	var buff bytes.Buffer
	enc := json.NewEncoder(&buff)
	if err := enc.Encode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buff.Bytes())
}

type staticFiles struct {
}

func (s staticFiles) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("serving static content: /%v", r.URL)
	http.FileServer(http.Dir("/static")).ServeHTTP(w, r)
}

func main() {
	log.Printf("binding to port 80 for beta application... try /index.html")
	http.HandleFunc("/api", apiHandler)
	http.Handle("/", http.StripPrefix("/", staticFiles{}))
	http.ListenAndServe(":80", nil)
}
