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

func (apiCfg *apiConfig) getSingleEvent(w http.ResponseWriter, r *http.Request, user database.User) {
	//event_name OR event_id work to find the event
	type params struct {
		EventID   uuid.UUID `json:"event_id"`
		EventName string    `json:"event_name"`
	}

	parameters, err := decodeParams[params](r.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding request body: %s", err))
		return
	}

	if parameters.EventID == uuid.Nil && parameters.EventName == "" {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("no event data parsed: %s", err))
		return
	}

	var actualEvent database.Event

	//If event id not provided
	if parameters.EventID == uuid.Nil {
		event, err := apiCfg.DB.GetEventByName(r.Context(), parameters.EventName)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error getting event in DB: %s", err))
			return
		}

		actualEvent = event
	} else {
		event, err := apiCfg.DB.GetEventByID(r.Context(), parameters.EventID)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error getting event in DB: %s", err))
			return
		}

		actualEvent = event
	}

	//Works anyway, as actualEventID is assigned via logic
	event_user_pair, err := apiCfg.DB.GetEventConcerns(r.Context(), database.GetEventConcernsParams{
		EventID: actualEvent.EventID,
		UserID:  user.UserID,
	})

	if event_user_pair == *new(database.EventUser) {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("error, authed user not related to event in any way: %s", err))
		return
	}

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error verifying relation between authed user and event: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, DBEventToLocalEvent(actualEvent))
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

func (apiCfg *apiConfig) addUsersToEvent(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		EventID uuid.UUID   `json:"event_id"`
		ToAdd   []uuid.UUID `json:"to_remove"`
	}

	parameters, err := decodeParams[params](r.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding request body: %s", err))
		return
	}

	event, err := apiCfg.DB.GetEventByID(r.Context(), parameters.EventID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error getting event from DB: %s", err))
		return
	}

	concerns, err := apiCfg.DB.GetUserIDsForEvent(r.Context(), parameters.EventID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error finding users in DB for event given: %s", err))
		return
	}

	if event.AuthorID != user.UserID {
		respondWithError(w, http.StatusUnauthorized, "only able to add users from event if authed user is author of event")
		return
	}

	for _, toAddUserID := range parameters.ToAdd {
		//If user is already concerned, skip
		if contains(toAddUserID, concerns) {
			return
		}

		_, err = apiCfg.DB.AddUserToEvent(r.Context(), database.AddUserToEventParams{
			EventID: parameters.EventID,
			UserID:  toAddUserID,
		})

		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error adding user to event: %s", err))
			break
		}
	}

	if err != nil {
		return
	}

	userIDArr, err := apiCfg.DB.GetUserIDsForEvent(r.Context(), event.EventID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting user id array from DB: %s", err))
		return
	}

	users := []User{}

	for _, userID := range userIDArr {
		user, err := apiCfg.DB.GetUserByID(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error getting users's event from DB: %s", err))
			return
		}
		users = append(users, DBUserToLocalUser(user))
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (apiCfg *apiConfig) getUsersForEvent(w http.ResponseWriter, r *http.Request, user database.User) {
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
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error getting event from DB: %s", err))
		return
	}

	if event.AuthorID != user.UserID {
		respondWithError(w, http.StatusUnauthorized, "only able to get users from event if authed user is author of event")
		return
	}

	concerns, err := apiCfg.DB.GetUserIDsForEvent(r.Context(), parameters.EventID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error finding users in DB for event given: %s", err))
		return
	}

	userArr, err := apiCfg.resGetUsersFromArr(r.Context(), concerns)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error getting concerned users from DB: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, userArr)
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
		respondWithError(w, http.StatusUnauthorized, "cannot delete event if not the author")
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

	concerned, err := apiCfg.DB.GetUserIDsForEvent(r.Context(), parameters.EventID)

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

func (apiCfg *apiConfig) removeUsersFromEvent(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		EventID  uuid.UUID   `json:"event_id"`
		ToRemove []uuid.UUID `json:"to_remove"`
	}

	parameters, err := decodeParams[params](r.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding request body: %s", err))
		return
	}

	event, err := apiCfg.DB.GetEventByID(r.Context(), parameters.EventID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error getting event from DB: %s", err))
		return
	}

	if event.AuthorID != user.UserID {
		respondWithError(w, http.StatusUnauthorized, "only able to remove users from event if authed user is author of event")
		return
	}

	concerns, err := apiCfg.DB.GetUserIDsForEvent(r.Context(), parameters.EventID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error finding users in DB for event given: %s", err))
		return
	}

	for _, userID := range concerns {
		//If the userID in the event's concerns is in the array of userIDs to remove, then don't append it to the list
		if contains(userID, parameters.ToRemove) {
			err = apiCfg.DB.DeleteUserFromEvent(r.Context(), database.DeleteUserFromEventParams{
				EventID: parameters.EventID,
				UserID:  userID,
			})

			if err != nil {
				respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error deleting user from event in DB: %s", err))
				break
			}
		}
	}

	if err != nil {
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}
