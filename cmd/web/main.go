package main

import (
	"log"
	"net/http"
	"text/template"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	log.Print("home handler called!")

	ts, err := template.ParseFiles("./ui/html/index.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	log.Print(ts.Name())
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func sayBye(w http.ResponseWriter, r *http.Request) {
	log.Print("satbye handler called!")
	w.Write([]byte("Come back again!"))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Print("healthCheck handler called!")
	w.Write([]byte("API working good!"))
}

func contactMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
	}
	log.Print("contact handler invoked")
	w.Write([]byte("Send me a message!"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/healthz", healthCheck)
	mux.HandleFunc("GET /bye", sayBye)
	mux.HandleFunc("POST /contact", contactMe)

	log.Print("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
