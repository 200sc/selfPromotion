package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/nlopes/slack/slackevents"
)

func raffleChallenge() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: os.Getenv("SLACK_API_TOKEN")}))
		if e != nil {
			fmt.Println(e, body)
			w.WriteHeader(http.StatusInternalServerError)
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}
		fmt.Println("URL:", r.RequestURI)
		fmt.Println("Type:", eventsAPIEvent.Type)
	}
}

func raffleStart() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: os.Getenv("SLACK_API_TOKEN")}))
		if e != nil {
			fmt.Println(e, body)
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Println("URL:", r.RequestURI)
		fmt.Println("Type:", eventsAPIEvent.Type)
	}
}

func raffleOptin() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: os.Getenv("SLACK_API_TOKEN")}))
		if e != nil {
			fmt.Println(e, body)
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Println("URL:", r.RequestURI)
		fmt.Println("Type:", eventsAPIEvent.Type)
	}
}

func raffleOptout() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: os.Getenv("SLACK_API_TOKEN")}))
		if e != nil {
			fmt.Println(e, body)
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Println("URL:", r.RequestURI)
		fmt.Println("Type:", eventsAPIEvent.Type)
	}
}

func raffleOptinAll() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: os.Getenv("SLACK_API_TOKEN")}))
		if e != nil {
			fmt.Println(e, body)
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Println("URL:", r.RequestURI)
		fmt.Println("Type:", eventsAPIEvent.Type)
	}
}

func raffleOptoutAll() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: os.Getenv("SLACK_API_TOKEN")}))
		if e != nil {
			fmt.Println(e, body)
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Println("URL:", r.RequestURI)
		fmt.Println("Type:", eventsAPIEvent.Type)
	}
}

func raffleWhosIn() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: os.Getenv("SLACK_API_TOKEN")}))
		if e != nil {
			fmt.Println(e, body)
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Println("URL:", r.RequestURI)
		fmt.Println("Type:", eventsAPIEvent.Type)
	}
}

func raffleDraw() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: os.Getenv("SLACK_API_TOKEN")}))
		if e != nil {
			fmt.Println(e, body)
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Println("URL:", r.RequestURI)
		fmt.Println("Type:", eventsAPIEvent.Type)
	}
}
