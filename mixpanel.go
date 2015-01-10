package main

type EventPropertiesPayload struct {
	Token      string `json:"token"`
	DistinctId string `json:"distinct_id"`
}

type EventPayload struct {
	Event      string                 `json:"event"`
	Properties EventPropertiesPayload `json:"properties"`
}

type PeoplePayload struct {
	Token      string `json:"$token"`
	DistinctId string `json:"$distinct_id"`
}
