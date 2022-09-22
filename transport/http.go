package transport

import (
	"encoding/json"
	"go-nsq/application/entrypoint"
	"go-nsq/model"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	entryPointService entrypoint.IEntryPointService
}

func NewHTTPServer(
	entryPointService entrypoint.IEntryPointService,
) *mux.Router {
	router := mux.NewRouter()
	server := server{
		entryPointService: entryPointService,
	}
	router.HandleFunc("/sendDocument", server.SendDocument).Methods(http.MethodPost)

	return router
}

func httpWriteResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (s *server) SendDocument(w http.ResponseWriter, r *http.Request) {
	err := s.entryPointService.SendData()

	if err != nil {
		log.Println("Error sending data")
		httpWriteResponse(w, &model.ServerResponse{
			Message: "Error Sending Data",
		})
	}
	log.Println("Upload Document Success")
	httpWriteResponse(w, &model.ServerResponse{
		Message: "Success",
	})
}
