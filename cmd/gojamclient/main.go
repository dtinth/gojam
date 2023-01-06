package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dtinth/gojam/pkg/gojam"
	"github.com/dtinth/gojam/pkg/jamulusprotocol"
	"github.com/google/uuid"
)

func main() {
	// Parse command line arguments
	server := flag.String("server", "127.0.0.1:22124", "server to connect to")
	pcmout := flag.String("pcmout", "", "server to pipe PCM data to")
	apiserver := flag.String("apiserver", "", "server to listen for API requests")
	name := flag.String("name", "", "musician name")
	vad := flag.Bool("vad", false, "do not send audio when there is no activity")
	mp3 := flag.Bool("mp3", false, "encode mp3 and exposes via API server (requires ffmpeg)")

	flag.Parse()

	fmt.Println("gojam jamulus client")
	fmt.Println("server:", *server)

	client, err := gojam.NewClient(*server)
	if err != nil {
		panic(err)
	}
	fmt.Println("client:", client)

	info := client.GetChannelInfo()
	if *name != "" {
		info.Name = *name
		client.UpdateChannelInfo(info)
	}

	soundPipe := newSoundPipe(*vad)
	client.HandlePCM = func(pcm []int16) {
		soundPipe.WritePCM(pcm)
	}

	if *pcmout != "" {
		installPCMOut(soundPipe, *pcmout)
	}

	if *apiserver != "" {
		installAPIServer(client, soundPipe, *apiserver, *mp3)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("exiting")
	client.Close()
}

func installPCMOut(pipe *soundPipe, pcmout string) {
	conn, err := net.Dial("tcp", pcmout)
	if err != nil {
		panic(err)
	}
	outChan := make(chan []byte, 100)
	go func() {
		for {
			data := <-outChan
			_, err := conn.Write(data)
			if err != nil {
				conn.Close()
				panic(err)
			}
		}
	}()
	pipe.HandlePCMBytes = func(pcmBytes []byte) {
		outChan <- pcmBytes
	}
}

type apiChannelInfo struct {
	Name       *string `json:"name"`
	City       *string `json:"city"`
	Country    *uint16 `json:"country"`
	SkillLevel *uint8  `json:"skillLevel"`
	Instrument *uint32 `json:"instrument"`
}

type chatHistoryEntry struct {
	Id        string `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type apiChatInput struct {
	Message string `json:"message"`
}

type apiEvent struct {
	Clients        *[]apiChannelInfo `json:"clients,omitempty"`
	Levels         *[]int            `json:"levels,omitempty"`
	NewChatMessage *chatHistoryEntry `json:"newChatMessage,omitempty"`
}

func installAPIServer(client *gojam.Client, pipe *soundPipe, apiserver string, mp3 bool) {
	http.HandleFunc("/channel-info", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			info := client.GetChannelInfo()
			json.NewEncoder(w).Encode(serializeChannelInfo(&info))
		case "PATCH":
			var patch apiChannelInfo
			err := json.NewDecoder(r.Body).Decode(&patch)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			info := client.GetChannelInfo()
			if patch.Name != nil {
				info.Name = *patch.Name
			}
			if patch.City != nil {
				info.City = *patch.City
			}
			if patch.Country != nil {
				info.Country = *patch.Country
			}
			if patch.SkillLevel != nil {
				info.SkillLevel = jamulusprotocol.SkillLevelId(*patch.SkillLevel)
			}
			if patch.Instrument != nil {
				info.Instrument = jamulusprotocol.InstrumentId(*patch.Instrument)
			}
			client.UpdateChannelInfo(info)
			json.NewEncoder(w).Encode(map[string]string{
				"status": "ok",
			})
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "method not allowed",
			})
		}
	})

	chatHistory := []chatHistoryEntry{}
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			var input apiChatInput
			err := json.NewDecoder(r.Body).Decode(&input)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			client.SendChatMessage(input.Message)
			json.NewEncoder(w).Encode(map[string]string{
				"status": "ok",
			})
		case "GET":
			json.NewEncoder(w).Encode(chatHistory)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	broadcaster := newEventBroadcaster()
	welcomeEvent := apiEvent{}
	broadcaster.GetWelcomeMessage = func() string {
		return jsonMarshal(welcomeEvent)
	}
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		f, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		ctx := r.Context()
		unregister := broadcaster.Register(r.RemoteAddr, func(data string) {
			go func() {
				fmt.Fprintf(w, "data: %s\n\n", data)
				f.Flush()
			}()
		})

		<-ctx.Done()
		unregister()
	})

	if mp3 {
		station := newMp3Broadcaster()
		oldHandler := pipe.HandlePCMBytes
		pipe.HandlePCMBytes = func(pcmBytes []byte) {
			if oldHandler != nil {
				oldHandler(pcmBytes)
			}
			station.WritePCMBytes(pcmBytes)
		}
		http.HandleFunc("/mp3", func(w http.ResponseWriter, r *http.Request) {
			f, ok := w.(http.Flusher)
			if !ok {
				http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "audio/mpeg")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Access-Control-Allow-Origin", "*")

			outChan := make(chan []byte, 100)
			unregister := station.Register(r.RemoteAddr, outChan)
			defer unregister()
			for {
				select {
				case <-r.Context().Done():
					return
				case data := <-outChan:
					w.Write(data)
					f.Flush()
				}
			}
		})
	}

	client.HandleChatMessage = func(message string) {
		chatHistory = append(chatHistory, chatHistoryEntry{
			Id:        uuid.New().String(),
			Message:   message,
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		})

		// Trim chat history
		if len(chatHistory) > 100 {
			chatHistory = chatHistory[len(chatHistory)-100:]
		}

		// Broadcast
		event := apiEvent{}
		event.NewChatMessage = &chatHistory[len(chatHistory)-1]
		broadcaster.Broadcast(jsonMarshal(event))
	}
	client.HandleSoundLevels = func(list []uint8) {
		event := apiEvent{}
		levels := make([]int, len(list))
		for i, level := range list {
			levels[i] = int(level)
		}
		event.Levels = &levels
		welcomeEvent.Levels = &levels
		broadcaster.Broadcast(jsonMarshal(event))
	}
	client.HandleClientList = func(list []*jamulusprotocol.ClientListItem) {
		event := apiEvent{}
		clients := make([]apiChannelInfo, len(list))
		for i, item := range list {
			clients[i] = serializeClientListItem(item)
		}
		event.Clients = &clients
		welcomeEvent.Clients = &clients
		broadcaster.Broadcast(jsonMarshal(event))
	}
	go func() {
		err := http.ListenAndServe(apiserver, nil)
		if err != nil {
			panic(err)
		}
	}()
}

func jsonMarshal(event apiEvent) string {
	data, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func serializeChannelInfo(info *jamulusprotocol.ChannelInfo) apiChannelInfo {
	skillLevelId := uint8(info.SkillLevel)
	instrumentId := uint32(info.Instrument)
	return apiChannelInfo{
		Name:       &info.Name,
		City:       &info.City,
		Country:    &info.Country,
		SkillLevel: &skillLevelId,
		Instrument: &instrumentId,
	}
}

func serializeClientListItem(item *jamulusprotocol.ClientListItem) apiChannelInfo {
	skillLevelId := uint8(item.SkillLevel)
	instrumentId := uint32(item.InstrumentId)
	return apiChannelInfo{
		Name:       &item.Name,
		City:       &item.City,
		Country:    &item.CountryId,
		SkillLevel: &skillLevelId,
		Instrument: &instrumentId,
	}
}
