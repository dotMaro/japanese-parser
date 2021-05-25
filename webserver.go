package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func StartServer(dict Dictionary) {
	server := server{
		dict: dict,
	}
	http.Handle("/parse", server)
	fmt.Println("Server up and running at :8080")
	_ = http.ListenAndServe(":8080", server)
}

type server struct {
	dict Dictionary
}

type input struct {
	Sentence string `json:"sentence"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (s server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("got req")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if req.Method == http.MethodOptions {
		w.WriteHeader(200)
		return
	}
	if req.Method != http.MethodPost {
		writeError(w, "only method POST is allowed", 400)
		return
	}

	body := io.LimitReader(req.Body, 1000*1000) // 1 MB
	dec := json.NewDecoder(body)
	var input input
	err := dec.Decode(&input)
	if err != nil {
		writeError(w, "invalid body", 400)
		return
	}

	sentence := s.dict.ParseSentence(input.Sentence)
	enc := json.NewEncoder(w)
	enc.Encode(sentence)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	resp, _ := json.Marshal(errorResponse{Error: msg})
	w.WriteHeader(status)
	w.Write(resp)
}
