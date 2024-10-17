package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sebstainsgit/calendar/internal/database"
)

type apiConfig struct {
	DB           *database.Queries
	Admin_Secret string
}

type errResponse struct {
	Error string `json:"error"`
}

type authedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

// local event
type Event struct {
	EventID   uuid.UUID `json:"event_id"`
	EventName string    `json:"event_name"`
	UsersID   uuid.UUID `json:"users_id"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	UserID    uuid.UUID `json:"user_id"`
	ApiKey    string    `json:"api_key"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DBEventToLocalEvent(event database.Event) Event {
	return Event{
		EventID:   event.EventID,
		EventName: event.EventName,
		UsersID:   event.UsersID,
		Date:      event.Date,
		CreatedAt: event.CreatedAt,
		UpdatedAt: event.UpdatedAt,
	}
}

func DBUserToLocalUser(user database.User) User {
	return User{
		UserID:    user.UserID,
		ApiKey:    user.ApiKey,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func DBEventsToLocalEvents(DBEventArr []database.Event) []Event {
	eventArr := make([]Event, len(DBEventArr))
	for i, dbEvent := range DBEventArr {
		eventArr[i] = DBEventToLocalEvent(dbEvent)
	}
	return eventArr
}

func DBUsersToLocalUsers(DBUsersArr []database.User) []User {
	localUserArr := make([]User, len(DBUsersArr))
	for i, dbUser := range DBUsersArr {
		localUserArr[i] = DBUserToLocalUser(dbUser)
	}
	return localUserArr
}
