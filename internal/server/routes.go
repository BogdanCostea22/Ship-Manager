package server

import (
	"encoding/json"
	"log"
	"net/http"

	// "Ship_Manager/cmd/web"
	"Ship_Manager/internal/handlers"
	"Ship_Manager/internal/repositories"
	"Ship_Manager/internal/services"
)

func (s *Server) RegisterRoutes() http.Handler {
	repository := repositories.NewPackageRepository()
	service := services.NewPackageService(repository)
	ph := handlers.NewPackageHandler(service)
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HelloWorldHandler)

	mux.HandleFunc("/calculator", ph.CalculatorIndex)
	mux.HandleFunc("/add-pack", ph.AddPack)
	mux.HandleFunc("/clear-packs", ph.ClearPacks)
	mux.HandleFunc("/calculate", ph.Calculate)
	mux.HandleFunc("/pack-sizes", ph.PackSizes)

	// fileServer := http.FileServer(http.FS(web.Files))
	// mux.Handle("/assets/", fileServer)
	// mux.Handle("/web", templ.Handler(web.HelloForm()))
	// mux.HandleFunc("/hello", web.HelloWebHandler)

	return mux
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
