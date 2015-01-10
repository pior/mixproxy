package main

import (
	// "runtime"
	"log"
	"net/http"
	"net/url"
)

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU())

	// Raise the number of idle connection to keep
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 50

	mixpanelHost, err := url.Parse("http://api.mixpanel.com")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/health", healthHandlerFunc)
	http.Handle("/", NewMixpanelProxy(mixpanelHost))

	log.Println("Listening...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func healthHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Health", "OK")
	w.Write([]byte("OK"))
}
