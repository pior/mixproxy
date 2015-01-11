package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	proxy *httputil.ReverseProxy
}

func NewProxy(target *url.URL) *Proxy {
	log.Printf("Proxying to %s://%s\n", target.Scheme, target.Host)

	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
	}

	proxy := &httputil.ReverseProxy{Director: director}

	return &Proxy{
		proxy: proxy,
	}
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	setpoint := handleMixpanelRequest(req)

	if setpoint.forward {
		p.proxy.ServeHTTP(rw, req)
	} else {
		serveDummyMixpanelResponse(rw, req)
	}
}

func serveDummyMixpanelResponse(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Printf("mixproxy: fake response error: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	redirect := req.Form.Get("redirect")
	if redirect != "" {
		rw.WriteHeader(302)
		rw.Header().Set("Location", redirect)
		rw.Header().Set("Cache-Control", "no-cache, no-store")
		return
	}

	verbose := req.Form.Get("verbose")
	if verbose == "1" {
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"error": "", "status": 1}`))
	} else {
		rw.Write([]byte("1"))
	}
}

type MixpanelRequestSetpoint struct {
	forward bool
}

func (s *MixpanelRequestSetpoint) setForEvent(event EventPayload) {
	s.forward = false
}

func handleMixpanelRequest(req *http.Request) (setpoint MixpanelRequestSetpoint) {
	setpoint.forward = true // forward by default

	form, err := dumpRequestForm(req)
	if err != nil {
		log.Println("error: dumpRequestForm: ", err)
		return
	}
	data, err := base64.StdEncoding.DecodeString(form.Get("data"))
	if err != nil {
		log.Println("error: base64: ", err)
	}
	if data == nil {
		log.Printf("debug: no data")
		return
	}

	switch req.URL.Path {
	case "/track", "/track/":
		log.Printf("debug: Event!")

		var event EventPayload
		err := json.Unmarshal(data, &event)
		if err != nil {
			log.Println("error: ", err)
			return
		}
		log.Printf("Token=%s", event.Properties.Token)
		log.Printf("DistinctId=%s", event.Properties.DistinctId)

		setpoint.setForEvent(event)

		// get settings for project
		// get settings for project/event
		// create setpoint accordingly

	case "/engage", "/engage/":
		log.Printf("debug: People!")

		var event PeoplePayload
		err := json.Unmarshal(data, &event)
		if err != nil {
			log.Println("error: ", err)
			return
		}
		log.Printf("Token=%s", event.Token)
		log.Printf("DistinctId=%s", event.DistinctId)
		// get settings for project
		// get settings for project/event
		// create setpoint accordingly

	default:
		log.Printf("debug: unkown endpoint: %s", req.URL.Path)
	}

	return
}
