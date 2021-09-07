package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"

	"github.com/gofrs/uuid"
)

func raffleChallenge() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		event, body, err := apiEvent(r)
		if err != nil {
			fmt.Println("Error parsing event", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		if event.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}
		fmt.Println("URL:", r.RequestURI)
		fmt.Println("Type:", event.Type)
	}
}

var (
	raffleLock sync.Mutex
	// map from channel to raffle
	ongoingRaffles = map[string]*Raffle{}
	allInUsers     = map[string]struct{}{}
	userNames      = map[string]string{}
	channelNames   = map[string]string{}
)

type Raffle struct {
	description string
	starterID   string
	in          map[string]struct{}
}

func teeVerifier(r *http.Request) (*http.Request, error) {
	verifier, err := slack.NewSecretsVerifier(r.Header, os.Getenv("SLACK_API_TOKEN"))
	if err != nil {
		return nil, err
	}
	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	return r, nil
}

func reply(w http.ResponseWriter, message string) {
	params := &slack.Msg{Text: message}
	b, err := json.Marshal(params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func replyToChannel(cl *slack.Client, channelID, message string) {
	s1, s2, err := cl.PostMessage(channelID, slack.MsgOptionText(message, false))
	fmt.Println("reply response", s1, s2, err)
}

func slashCommand(cmd func(slack.SlashCommand, http.ResponseWriter, *http.Request)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r, err := teeVerifier(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		slash, err := slack.SlashCommandParse(r)
		if err != nil {
			fmt.Println("Error parsing event", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("URL:", r.RequestURI)
		fmt.Println("Token:", slash.Token)
		fmt.Println("Channel:", slash.ChannelName)
		fmt.Println("Command:", slash.Command)
		fmt.Println("User:", slash.UserName)
		fmt.Println("Text:", slash.Text)
		raffleLock.Lock()

		userNames[slash.UserID] = slash.UserName
		channelNames[slash.ChannelID] = slash.ChannelName

		cmd(slash, w, r)

		raffleLock.Unlock()
	}
}

func raffleStart(cl *slack.Client) func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {
	return func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {
		if _, ok := ongoingRaffles[slash.ChannelID]; ok {
			reply(w, "Please stop the current raffle before starting a new one.")
			return
		}
		in := make(map[string]struct{}, len(allInUsers))
		for k, v := range allInUsers {
			in[k] = v
		}
		ongoingRaffles[slash.ChannelID] = &Raffle{
			description: slash.Text,
			starterID:   slash.UserName,
			in:          in,
		}
		replyToChannel(cl, slash.ChannelID, "Raffle started! -- "+slash.Text)
	}
}

func raffleOptin(_ *slack.Client) func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {
	return func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {
		raff, ok := ongoingRaffles[slash.ChannelID]
		if !ok {
			reply(w, "Cannot opt-in. No raffle is ongoing in this channel.")
			return
		}
		if _, ok := raff.in[slash.UserID]; ok {
			reply(w, slash.UserName+" is already in the raffle.")
		}
		raff.in[slash.UserID] = struct{}{}
		reply(w, slash.UserName+" added to raffle.")
	}
}

func raffleOptout(_ *slack.Client) func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {
	return func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {
		raff, ok := ongoingRaffles[slash.ChannelID]
		if !ok {
			reply(w, "Cannot opt-out. No raffle is ongoing in this channel.")
			return
		}
		if _, ok := raff.in[slash.UserID]; !ok {
			reply(w, slash.UserName+" is not in the raffle.")
		}
		delete(raff.in, slash.UserID)
		reply(w, slash.UserName+" removed from raffle.")
	}
}

func raffleOptinAll(_ *slack.Client) func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {
	return func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {
		if _, ok := allInUsers[slash.UserID]; ok {
			reply(w, slash.UserName+" is already opted in to all raffles.")
		}
		allInUsers[slash.UserID] = struct{}{}
		for _, raff := range ongoingRaffles {
			raff.in[slash.UserID] = struct{}{}
		}
		reply(w, slash.UserName+" added to all ongoing and future raffles")
	}
}

func raffleOptoutAll(_ *slack.Client) func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {
	return func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {
		delete(allInUsers, slash.UserID)
		for _, raff := range ongoingRaffles {
			delete(raff.in, slash.UserID)
		}
		reply(w, slash.UserName+" removed from all ongoing and future raffles")
	}
}

func raffleWhosIn(cl *slack.Client) func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {
	return func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {
		raff, ok := ongoingRaffles[slash.ChannelID]
		if !ok {
			reply(w, "There is no ongoing raffle in this channel.")
			return
		}
		in := make([]string, 0, len(raff.in))
		for k := range raff.in {
			in = append(in, userNames[k])
		}
		replyToChannel(cl, slash.ChannelID, "Who's In: "+strings.Join(in, ", "))
	}
}

const validateUsersAreReal = false

func raffleSetUsers(cl *slack.Client) func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {
	return func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {

		if len(slash.Text) == 0 {
			reply(w, "Usage: /raffleset april,bill,clara,danthro")
			return
		}
		toAdd := strings.Split(slash.Text, ",")
		toAddIDs := make([]string, len(toAdd))
		allUsers, err := cl.GetUsers()
		if err != nil {
			fmt.Println("Error getting users", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i, userName := range toAdd {
			for _, u := range allUsers {
				if u.Name == userName {
					userNames[u.ID] = u.Name
					toAddIDs[i] = u.ID
					break
				}
			}
			if toAddIDs[i] == "" {
				if validateUsersAreReal {
					reply(w, "Couldn't find user: "+userName)
					return
				}
				// Support 'fake' users, like 'engineering'
				id, err := uuid.NewV4()
				if err != nil {
					reply(w, "Couldn't create uuid: "+err.Error())
					return
				}
				userNames[id.String()] = userName
				toAddIDs[i] = id.String()
			}
		}
		fmt.Println("Got IDs for user names:", toAddIDs, toAdd)

		// secret flag, set opt in all or opt out all
		// if len(split) != 1 {
		// 	if split[1] == "optinall" {
		// 		reply(w, "Opting in all secretly")
		// 		for _, id := range toAddIDs {
		// 			allInUsers[id] = struct{}{}
		// 		}
		// 	} else if split[1] == "optoutall" {
		// 		reply(w, "Opting out all secretly")
		// 		for _, id := range toAddIDs {
		// 			delete(allInUsers, id)
		// 		}
		// 	}
		// 	return
		// }
		raff, ok := ongoingRaffles[slash.ChannelID]
		if !ok {
			reply(w, "There is no ongoing raffle in this channel.")
			return
		}
		raff.in = map[string]struct{}{}
		for _, s := range toAddIDs {
			raff.in[s] = struct{}{}
		}

		in := make([]string, 0, len(raff.in))
		for k := range raff.in {
			in = append(in, userNames[k])
		}
		replyToChannel(cl, slash.ChannelID, "Who's In set to: "+strings.Join(in, ", "))
	}
}

func raffleDraw(cl *slack.Client) func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {
	return func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {
		raff, ok := ongoingRaffles[slash.ChannelID]
		if !ok {
			reply(w, "There is no ongoing raffle in this channel.")
			return
		}
		if len(raff.in) == 0 {
			reply(w, "No one has opted-in to the current raffle.")
			return
		}

		winnerCount := 1
		if slash.Text != "" {
			var err error
			winnerCount, err = strconv.Atoi(strings.TrimSpace(slash.Text))
			if err != nil {
				reply(w, "Unable to parse number to draw: "+slash.Text+" : "+err.Error())
				return
			}
		}

		if winnerCount > len(raff.in) {
			reply(w, fmt.Sprintf("Only %d potential winners have opted in, but %d were drawn.", len(raff.in), winnerCount))
			return
		}

		rand.Seed(time.Now().UnixNano())

		out := make([]int, winnerCount)
		weights := make([]float64, len(raff.in))
		inSlice := make([]string, len(raff.in))
		i := 0
		for k := range raff.in {
			weights[i] = 1
			inSlice[i] = userNames[k]
			i++
		}
		stwh := newSTWHeap(weights)
		for i := 0; i < winnerCount; i++ {
			out[i] = stwh.Pop()
		}

		winners := make([]string, len(out))
		for i, v := range out {
			winners[i] = inSlice[v]
		}
		print := "Names drawn: " + strings.Join(winners, ", ")
		if len(winners) == 1 {
			print = "Name drawn: " + winners[0]
		}
		replyToChannel(cl, slash.ChannelID, print)
		w.WriteHeader(200)
	}
}

func raffleStop(cl *slack.Client) func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {
	return func(slash slack.SlashCommand, w http.ResponseWriter, r *http.Request) {
		_, ok := ongoingRaffles[slash.ChannelID]
		if !ok {
			reply(w, "There is no ongoing raffle in this channel.")
			return
		}
		// Todo: admins, limited stopping
		//if raff.starterID != slash.UserID {
		//
		//}
		delete(ongoingRaffles, slash.ChannelID)
		replyToChannel(cl, slash.ChannelID, "Ended raffle for channel "+channelNames[slash.ChannelID])
		w.WriteHeader(200)
	}
}

func apiEvent(r *http.Request) (slackevents.EventsAPIEvent, string, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	ev, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: os.Getenv("SLACK_API_TOKEN")}))
	return ev, body, err
}
