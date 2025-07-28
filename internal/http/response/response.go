package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// Пример для 400 Bad Request
type BadRequestError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad Request"`
}

// Пример для 500 Internal Server Error
type InternalServerError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Internal Server Error"`
}

// ErrSubscriptionExists — ошибка при попытке создать дубликат подписки
type ErrSubscriptionExists struct {
	Code    int    `json:"code" example:"409"`
	Message string `json:"message" example:"Subscription already exists"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code >= http.StatusInternalServerError {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	RespondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
