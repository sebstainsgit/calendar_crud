package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sebstainsgit/calendar/internal/database"
)

func (apiCfg *apiConfig) createEvent(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		Event_Name string `json:"event_name"`
		Date       string `json:"date"`
	}

	var parameters params

	data, err := io.ReadAll(r.Body)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error reading request body: %s", err))
		return
	}

	err = json.Unmarshal(data, &parameters)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error unmarhsalling request body: %s", err))
		return
	}

	dueTime, err := time.Parse(time.RFC1123Z, parameters.Date)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing time into correct format: %s", err))
		return
	}

	event, err := apiCfg.DB.CreateEvent(r.Context(), database.CreateEventParams{
		EventID:   uuid.New(),
		EventName: parameters.Event_Name,
		UsersID: user.UserID,
		Date:      dueTime,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error writing event to DB: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, DBEventToLocalEvent(event))
}

func (apiCfg *apiConfig) getUsersEvents(w http.ResponseWriter, r *http.Request, user database.User) {
	events, err := apiCfg.DB.GetUsersEvents(r.Context(), user.UserID)
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error fetching all queries from DB: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, DBEventsToLocalEvents(events))
}

func (apiCfg *apiConfig) deleteEvent(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		EventID string `json:"event_id"`
	}

	var parameters params

	data, err := io.ReadAll(r.Body)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error reading request body: %s", err))
		return
	}

	err = json.Unmarshal(data, &parameters)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error unmarhsalling request body: %s", err))
		return
	}

	DBEventUUID, err := uuid.Parse(parameters.EventID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing uuid")
		return
	}
	
	event, err := apiCfg.DB.GetEventByID(r.Context(), DBEventUUID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "event not found")
		return
	}

	if event.UsersID != user.UserID {
		respondWithError(w, http.StatusUnauthorized, "cannot delete event if not the author")
		return
	}

	err = apiCfg.DB.DeleteEventByID(r.Context(), event.EventID)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error deleting event in database")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}