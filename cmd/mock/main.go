package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type User struct {
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic"`
	Address    string `json:"address"`
}

func main() {
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			passportSerie := r.URL.Query().Get("passportSerie")
			passportNumber := r.URL.Query().Get("passportNumber")

			if passportSerie == "" || passportNumber == "" {
				http.Error(w, "passportSerie and passportNumber are required", http.StatusBadRequest)
				return
			}

			user1 := User{
				Surname:    "Иванов",
				Name:       "Иван",
				Patronymic: "Иванович",
				Address:    "г. Москва, ул. Ленина, д.5, кв.1",
			}

			user2 := User{
				Surname:    "Петров",
				Name:       "Петр",
				Patronymic: "Петрович",
				Address:    "г. Москва, ул. Брежнева, д.1, кв.5",
			}

			if strings.HasSuffix(passportSerie, "1") {
				response, err := json.Marshal(user1)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
			} else {
				response, err := json.Marshal(user2)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
			}
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	if err := http.ListenAndServe(":8221", nil); err != nil {
		log.Fatal(err)
	}
}
