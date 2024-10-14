package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"disapp/internal/config"
	"disapp/internal/handlers/storage"

	"github.com/google/uuid"
)

func main() {
	projectBaseDir, err := filepath.Abs(filepath.Join("cmd", "disapp", "main.go"))
	if err != nil {
		panic(err)
	}
	projectBaseDir = filepath.Dir(filepath.Dir(filepath.Dir(projectBaseDir)))
	slog.Info(fmt.Sprintf("Project base dir: %s", projectBaseDir))

	defaultConfigPath := filepath.Join(projectBaseDir, "config", "local.yaml")

	configPath := flag.String("p", defaultConfigPath, "set config path")
	flag.Parse()

	config := config.MustLoad(*configPath)
	// тестово создал сообщение
	msgs := storage.New()
	dur, _ := time.ParseDuration("5m")
	testId := msgs.Add("ghljk23324", dur)
	slog.Info("New message was created: " + testId.String())

	router := http.NewServeMux()
	router.HandleFunc("/hi/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		w.Write([]byte("Hi, " + name))
	})
	router.HandleFunc("GET /m/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("uuid")
		parsedId, err := uuid.Parse(id)
		if err != nil {
			slog.Error("faild to pars uuid from URL", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			w.Write([]byte("Sorry, but your URL is incorrect"))
		}
		msg := msgs.Take(parsedId) // TODO: надо добавить проверку на наличие
		w.Write([]byte(msg.Body))
	})
	router.HandleFunc("POST /api/messages", func(w http.ResponseWriter, r *http.Request) {
		// decoder := json.NewDecoder(r.Body)
		// type testStruct struct {

		// }
		// var t testStruct
		// err := decoder.Decode(&t)
		body, _ := io.ReadAll(r.Body)
		slog.Info(string(body))
	})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>What`s that?</h1>"))
	})

	server := http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		IdleTimeout:  config.IdleTimeout,
	}

	// uu, err := uuid.NewRandom()
	// if err != nil {
	// 	slog.Error(err.Error())
	// }
	// dur, err := time.ParseDuration("5m")
	// if err != nil {
	// 	slog.Error(err.Error())
	// }
	// message.AllMessages[uu] = message.New("привет!", dur)
	// fmt.Print("AllMessages: ")
	// fmt.Println(message.AllMessages)
	// delete(message.AllMessages, uu)
	// fmt.Print("AllMessages after deletion: ")
	// fmt.Println(message.AllMessages)

	startingMsg := fmt.Sprintf("--- Starting server on %s", config.Address)
	slog.Info(startingMsg)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("failed to start server", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
		os.Exit(1)
	}
}
