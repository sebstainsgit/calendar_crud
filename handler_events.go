package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sebstainsgit/calendar/internal/database"
)

func (apiCfg *apiConfig) createSelfEvent(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		Event_Name string `json:"event_name"`
		Date       string `json:"date"`
	}
	//Date in format "2018-04-08 15:04:05"

	parameters, err := decodeParams[params](r.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding request body: %s", err))
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
		AuthorID:  user.UserID,
		Date:      dueTime,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error writing event to DB: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, DBEventToLocalEvent(event))
}

func (apiCfg *apiConfig) createGroupEvent(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		Event_Name string      `json:"event_name"`
		Date       string      `json:"date"`
		Concerns   []uuid.UUID `json:"concerns"`
	}
	//Date in format "2018-04-08 15:04:05"

	parameters, err := decodeParams[params](r.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding request body: %s", err))
		return
	}

	dueTime, err := time.Parse(time.DateTime, parameters.Date)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing time into correct format: %s", err))
		return
	}

	event, err := apiCfg.DB.CreateEvent(r.Context(), database.CreateEventParams{
		EventID:   uuid.New(),
		EventName: parameters.Event_Name,
		AuthorID:  user.UserID,
		Date:      dueTime,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error writing event to DB: %s", err))
		return
	}

	for _, DBUserID := range parameters.Concerns {
		apiCfg.DB.AddUserToEvent(r.Context(), database.AddUserToEventParams{
			EventID: event.EventID,
			UserID:  DBUserID,
		})
	}

	respondWithJSON(w, http.StatusCreated, DBEventToLocalEvent(event))
}

func (apiCfg *apiConfig) updateEvent(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		EventID    uuid.UUID `json:"event_id"`
		Event_Name string    `json:"event_name"`
		Date       time.Time `json:"date"`
	}

	parameters, err := decodeParams[params](r.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding request body: %s", err))
		return
	}

	event, err := apiCfg.DB.GetEventByID(r.Context(), parameters.EventID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("event not found: %s", err))
		return
	}

	if event.AuthorID != user.UserID {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("cannot delete event if not the author: %s", err))
		return
	}

	updatedEvent, err := apiCfg.DB.UpdateEvent(r.Context(), database.UpdateEventParams{
		Date:      parameters.Date,
		EventName: parameters.Event_Name,
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error updating event in DB: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, DBEventToLocalEvent(updatedEvent))
}

func (apiCfg *apiConfig) getUsersEvents(w http.ResponseWriter, r *http.Request, user database.User) {
	eventIDArr, err := apiCfg.DB.GetEventsIDsForUser(r.Context(), user.UserID)

	events := []Event{}

	for _, eventID := range eventIDArr {
		event, err := apiCfg.DB.GetEventByID(r.Context(), eventID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error getting users's event from DB: %s", err))
			return
		}
		events = append(events, DBEventToLocalEvent(event))
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error fetching all events from DB: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, events)
}

// Could make first part of deleteEvent + updateEvent into function but code is more readable and errors are more specific
func (apiCfg *apiConfig) deleteEvent(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		EventID uuid.UUID `json:"event_id"`
	}

	parameters, err := decodeParams[params](r.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding request body: %s", err))
		return
	}

	event, err := apiCfg.DB.GetEventByID(r.Context(), parameters.EventID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("event not found: %s", err))
		return
	}

	if event.AuthorID != user.UserID {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("cannot delete event if not the author: %s", err))
		return
	}

	err = apiCfg.DB.DeleteEvent(r.Context(), event.EventID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error deleting event in database: %s", err))
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}

func (apiCfg *apiConfig) removeSelfFromEvent(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		EventID uuid.UUID `json:"event_id"`
	}

	parameters, err := decodeParams[params](r.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding request body: %s", err))
		return
	}

	concerned, err := apiCfg.DB.GetUsersForEvent(r.Context(), parameters.EventID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error getting users for event from DB: %s", err))
		return
	}

	isConcerned := false

	for _, val := range concerned {
		if val == user.UserID {
			isConcerned = true
			break
		}
	}

	if !isConcerned {
		respondWithError(w, http.StatusUnauthorized, "authenticated user is not concerned with group event")
		return
	}

	err = apiCfg.DB.DeleteUserFromEvent(r.Context(), database.DeleteUserFromEventParams{
		EventID: parameters.EventID,
		UserID:  user.UserID,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error deleting user from event in DB: %s", err))
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}
