package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type MixpanelProxy struct {
	proxy *httputil.ReverseProxy
}

func NewMixpanelProxy(target *url.URL) *MixpanelProxy {
	log.Printf("Proxying to %s://%s\n", target.Scheme, target.Host)

	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
	}

	proxy := &httputil.ReverseProxy{Director: director}

	return &MixpanelProxy{
		proxy: proxy,
	}
}

func (p *MixpanelProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	setpoint := handleMixpanelRequest(req)

	if setpoint.enabled {
		p.proxy.ServeHTTP(rw, req)
	} else {

	}
}

type MixpanelRequestSetpoint struct {
	enabled bool
}

func (s *MixpanelRequestSetpoint) setForEvent(event EventPayload) {
	s.enabled = false
}

func handleMixpanelRequest(req *http.Request) (setpoint MixpanelRequestSetpoint) {
	setpoint.enabled = true // Enabled by default

	form, err := dumpRequestForm(req)
	if err != nil {
		log.Println("error: dumpRequestForm: ", err)
		return
	}
	data, err := base64.StdEncoding.DecodeString(form.Get("data"))
	if err != nil {
		log.Println("error: base64: ", err)
	}
	// data := getData(req)
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
