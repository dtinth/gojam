package main

import (
	"bytes"
	"encoding/binary"
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

	if *pcmout != "" {
		installPCMOut(client, *pcmout, *vad)
	}

	if *apiserver != "" {
		installAPIServer(client, *apiserver)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("exiting")
	client.Close()
}

func installPCMOut(client *gojam.Client, pcmout string, vad bool) {
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
				client.Close()
				panic(err)
			}
		}
	}()
	hp := 0
	client.HandlePCM = func(pcm []int16) {
		var buf bytes.Buffer
		activityDetected := false
		for _, sample := range pcm {
			binary.Write(&buf, binary.LittleEndian, sample)
			if sample != 0 {
				activityDetected = true
			}
		}
		if vad {
			if activityDetected {
				hp = 100
			} else {
				hp--
				if hp < 0 {
					hp = 0
				}
				if hp == 0 {
					return
				}
			}
		}
		outChan <- buf.Bytes()
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

func installAPIServer(client *gojam.Client, apiserver string) {
	http.HandleFunc("/channel-info", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			info := client.GetChannelInfo()
			skillLevelId := uint8(info.SkillLevel)
			instrumentId := uint32(info.Instrument)
			json.NewEncoder(w).Encode(apiChannelInfo{
				Name:       &info.Name,
				City:       &info.City,
				Country:    &info.Country,
				SkillLevel: &skillLevelId,
				Instrument: &instrumentId,
			})
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
	}

	go func() {
		err := http.ListenAndServe(apiserver, nil)
		if err != nil {
			panic(err)
		}
	}()
}
