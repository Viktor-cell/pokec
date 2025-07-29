package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/rs/cors"
)

type User struct {
	Name string `json:"name"`
}

func (u *User) decodeFromRequestBody(b io.Reader) error {
	err := json.NewDecoder(b).Decode(u)
	return err
}

type Message struct {
	User    User
	Message string `json:"msg"`
}

func (m *Message) decodeFromRequestBody(b io.Reader) error {
	err := json.NewDecoder(b).Decode(m)
	return err
}

type Tuple[A, B any] struct {
	Fst A `json:"fst"`
	Scn B `json:"scn"`
}

type LoginResponse struct {
	IsOk     bool   `json:"ok"`
	Message  string `json:"message"`
	Redirect string `json:"redirect"`
}

func (lr LoginResponse) sendResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lr)
}

var usersCache []User
var chatCache []Tuple[Message, time.Time]
var mutex sync.Mutex

func main() {
	usersCache = make([]User, 0)
	chatCache = make([]Tuple[Message, time.Time], 0)

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static/"))

	mux.Handle("/", fs)
	mux.HandleFunc("/login", handleLogin)

	mux.HandleFunc("POST /message", mesageHandler)
	mux.HandleFunc("GET /getMessages", mesageGetHandler)

	fmt.Println("server at localhost:8080")
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8080", handler)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var user User
	err := user.decodeFromRequestBody(r.Body)

	if err != nil {
		panic(err)
	}

	var res LoginResponse

	if !slices.Contains(usersCache, user) {
		res.IsOk = true
		res.Redirect = "/chat.html"

		mutex.Lock()
		usersCache = append(usersCache, user)
		mutex.Unlock()
	} else {
		res.IsOk = false
		res.Message = "uz taky je buchto, daj ine meno"
	}

	res.sendResponse(w)
}

func mesageHandler(w http.ResponseWriter, r *http.Request) {
	var msg Message
	err := msg.decodeFromRequestBody(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	mutex.Lock()
	chatCache = append(chatCache, Tuple[Message, time.Time]{Fst: msg, Scn: time.Now()})
	mutex.Unlock()

	w.WriteHeader(http.StatusAccepted)
}

func mesageGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	val, err := json.Marshal(chatCache)

	if err != nil {
		panic(err)
	}

	w.Write(val)
}
