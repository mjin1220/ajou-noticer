package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	CACertPath = "/cert_storage/Ajou-Noticer-Root-CA/root-ca.cert.pem"
	CRLPath    = "/cert_storage/Ajou-Noticer-Root-CA/crl/root-ca.crl"
	// CACertPath = "./root-ca.cert.pem"
	// CRLPath    = "./root-ca.crl"
)

type ReqJSON interface{}

func StartServer(port string) {
	router := httprouter.New()
	router.GET("/", IndexHandler)

	router.GET("/release/AjouNoticer.crt", CertReleaseHandler)
	router.GET("/release/AjouNoticer.crl", CRLReleaseHandler)

	router.GET("/webhook", WebhookGetHandler)
	router.POST("/webhook", WebhookPostHandler)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func IndexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logRequest(r)

	w.Write([]byte("<h1>Index Page</h1>"))
}

func CertReleaseHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logRequest(r)

	data, err := ioutil.ReadFile(CACertPath)
	if err != nil {
		log.Fatal(w, err)
	}
	http.ServeContent(w, r, "AjouNoticer.crt", time.Now(), bytes.NewReader(data))
}

func CRLReleaseHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logRequest(r)

	data, err := ioutil.ReadFile(CRLPath)
	if err != nil {
		log.Fatal(w, err)
	}
	http.ServeContent(w, r, "AjouNoticer.crl", time.Now(), bytes.NewReader(data))
}

func WebhookGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logRequest(r)

	keys, ok := r.URL.Query()["hub.verify_token"]

	if !ok || len(keys) < 1 {
		log.Println("Url Param 'hub.verify_token' is missing")
		return
	}

	key := keys[0]

	if key == "mjin1220" {
		keys, ok = r.URL.Query()["hub.challenge"]
		if !ok || len(keys) < 1 {
			log.Println("Url Param 'hub.challenge' is missing")
			return
		}
		key = keys[0]
		w.Write([]byte(key))
		return
	}

	log.Println("Error, wrong validation token")
	w.Write([]byte("Error, wrong validation token"))
}

func WebhookPostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logRequest(r)

}

func logRequest(r *http.Request) {
	logStr := fmt.Sprintf("%v %v\n", r.Method, r.URL.Path)
	fd, err := os.OpenFile("./request.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer fd.Close()

	if _, err = fd.WriteString(logStr); err != nil {
		log.Fatal(err)
		return
	}

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			return
		}
		if _, err = fd.WriteString(string(body) + "\n"); err != nil {
			log.Fatal(err)
			return
		}
	}

}
