package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/gorilla/websocket"
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

var usersNameCache []User
var chatCache []Tuple[Message, time.Time]
var mutex sync.Mutex

var clientConnection []*websocket.Conn
var broadcast chan []Tuple[Message, time.Time]

var upgrader = websocket.Upgrader{}

func main() {
	usersNameCache = make([]User, 0)
	chatCache = make([]Tuple[Message, time.Time], 0)
	clientConnection = make([]*websocket.Conn, 0)
	broadcast = make(chan []Tuple[Message, time.Time])

	go sendToAll()

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static/"))

	mux.Handle("/", fs)
	mux.HandleFunc("/login", handleLogin)

	mux.HandleFunc("/message", mesageHandler)

	fmt.Println("server at localhost:8080")
	handler := cors.Default().Handler(mux)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Println(err)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var user User
	err := user.decodeFromRequestBody(r.Body)

	if err != nil {
		panic(err)
	}

	var res LoginResponse

	if !slices.Contains(usersNameCache, user) {
		res.IsOk = true
		res.Redirect = "/chat.html"

		mutex.Lock()
		usersNameCache = append(usersNameCache, user)
		mutex.Unlock()
	} else {
		res.IsOk = false
		res.Message = "uz taky je buchto, daj ine meno"
	}

	res.sendResponse(w)
}

func sendToAll() {
	for {
		msgs := <-broadcast

		mutex.Lock()
		for _, client := range clientConnection {
			client.WriteJSON(msgs)
		}
		mutex.Unlock()
	}
}

func mesageHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	broadcast <- chatCache

	var msg Message

	mutex.Lock()
	clientConnection = append(clientConnection, conn)
	mutex.Unlock()

	for {
		_, inMsg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := json.Unmarshal(inMsg, &msg); err != nil {
			fmt.Println(err)
		}

		mutex.Lock()
		chatCache = append(chatCache, Tuple[Message, time.Time]{Fst: msg, Scn: time.Now()})
		mutex.Unlock()

		fmt.Println(msg)

		broadcast <- chatCache
	}
}
