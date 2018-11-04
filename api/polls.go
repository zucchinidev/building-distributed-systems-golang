package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type poll struct {
	ID      bson.ObjectId  `bson:"_id" json:"id"`
	Title   string         `json:"title"`
	Options []string       `json:"options"`
	Results map[string]int `json:"results,omitempty"`
	APIKey  string         `json:"apiKey"`
}

type Server struct {
	db *mgo.Session
}

func (s *Server) handlePolls(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		s.handlePollsGet(writer, request)
		return
	case "POST":
		s.handlePollsPost(writer, request)
		return
	case "DELETE":
		s.handlePollsDelete(writer, request)
		return

	case "OPTIONS":
		writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		respond(writer, request, http.StatusOK, nil)
		return
	}

	respondHTTPErr(writer, request, http.StatusNotFound)
}

func (s *Server) handlePollsGet(writer http.ResponseWriter, request *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	var query *mgo.Query
	var result []*poll

	collection := session.DB("ballots").C("polls")
	path := NewPath(request.URL.Path)

	if path.HasID() {
		query = collection.FindId(bson.ObjectIdHex(path.ID))
	} else {
		query = collection.Find(nil)
	}

	if err := query.All(&result); err != nil {
		respondErr(writer, request, http.StatusInternalServerError, err)
		return
	}

	respond(writer, request, http.StatusOK, &result)
}

func (s *Server) handlePollsPost(writer http.ResponseWriter, request *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	var p poll

	collection := session.DB("ballots").C("polls")

	if err := decodeBody(request, &p); err != nil {
		respondErr(writer, request, http.StatusBadRequest, "failed to read poll from request", err)
	}

	apiKey, ok := APIKey(request.Context())
	if ok {
		p.APIKey = apiKey
	}

	p.ID = bson.NewObjectId()
	if err := collection.Insert(p); err != nil {
		respondErr(writer, request, http.StatusInternalServerError, "failed to insert poll", err)
		return
	}

	writer.Header().Set("Location", "polls/"+p.ID.Hex())
	respond(writer, request, http.StatusCreated, nil)
}

func (s *Server) handlePollsDelete(writer http.ResponseWriter, request *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	collection := session.DB("ballots").C("polls")
	p := NewPath(request.URL.Path)
	if !p.HasID() {
		respondErr(writer, request, http.StatusMethodNotAllowed, "Cannot delete all polls.")
		return
	}

	if err := collection.RemoveId(bson.ObjectIdHex(p.ID)); err != nil {
		respondErr(writer, request, http.StatusInternalServerError, "failed to delete poll", err)
		return
	}

	respond(writer, request, http.StatusOK, nil)
}
