package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var startedAt = time.Now()

func main() {
	http.HandleFunc("/healthz", Healtz)
	http.HandleFunc("/secret", Secret)
	http.HandleFunc("/configmap", ConfigMap)
	http.HandleFunc("/", Hello)
	http.ListenAndServe(":80", nil)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	age := os.Getenv("AGE")

	fmt.Fprintf(w, "Hello, I`m %s. I`m %s.", name, age)
}

func ConfigMap(w http.ResponseWriter, r *http.Request) {

	// LER ARQUIVO
	data, err := ioutil.ReadFile("myfamily/family.txt")

	if err != nil {
		log.Fatalf("Error reading file", err)
	}

	fmt.Fprintf(w, "My Family: %s.", string(data))
}

func Secret(w http.ResponseWriter, r *http.Request) {

	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")

	fmt.Fprintf(w, "User: %s. Password: %s", user, password)
}

func Healtz(w http.ResponseWriter, r *http.Request) {

	duration := time.Since(startedAt)

	// UTILIZADO PARA TESTAR O LIVENESS...
	// if duration.Seconds() > 25 {
	// 	w.WriteHeader(500)
	// 	w.Write([]byte(fmt.Sprintf("Duration: %v", duration.Seconds())))
	// } else {
	// 	w.WriteHeader(200)
	// 	w.Write([]byte("ok"))
	// }

	// UTILIZADO PARA TESTAR O READINESS
	if duration.Seconds() < 10 || duration.Seconds() > 30 {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Duration: %v", duration.Seconds())))
	} else {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}
}
