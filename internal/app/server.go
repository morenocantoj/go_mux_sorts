package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/morenocantoj/go_mux_sorts/internal/app/types/responses"
)

func home(writter http.ResponseWriter, request *http.Request) {
	fmt.Println("GET /")
	response := responses.HomeJSON{Message: "Welcome to sort algorithms server"}
	responseJSON, err := json.Marshal(&response)
	if err != nil {
		return
	}
	writter.Write(responseJSON)
}

func defineRoutes(router *mux.Router) {
	router.HandleFunc("/", home)
}

func launchServer(address string, port string) {
	listeningAddress := address + port
	// suscripción SIGINT
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	router := mux.NewRouter()
	defineRoutes(router)

	server := &http.Server{
		Addr:         listeningAddress,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second}

	go func() {
		fmt.Printf("Server listening at %s...\n", listeningAddress)
		err := server.ListenAndServe()
		if err != nil {
			return
		}
	}()

	<-stopChan // Wait SIGINT
	log.Println("Shutting down server ...")

	// Shut down server safely
	ctx, fnc := context.WithTimeout(context.Background(), 5*time.Second)
	fnc()
	server.Shutdown(ctx)

	log.Println("Bye!")
}

// Start : Starts the server
func Start() {
	fmt.Println("### Sorting Algorithms Server ###")
	fmt.Println("Server launching...")
	launchServer("127.0.0.1", ":9080")
}
