package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sebstainsgit/calendar/internal/database"
)

type apiConfig struct {
	DB         *database.Queries
	JWT_SECRET string
}

type errResponse struct {
	Error string `json:"error"`
}

type JWTResponse struct {
	JWT string `json:"access_token"`
}

type TokenResponse struct {
	JWT          string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// local event
type Event struct {
	EventID   uuid.UUID `json:"event_id"`
	EventName string    `json:"event_name"`
	AuthorID  uuid.UUID `json:"users_id"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Elevation string    `json:"elevation"`
}

type RefreshToken struct {
	RefrToken string    `json:"refresh_token"`
	UsersID   uuid.UUID `json:"users_id"`
	Expires   time.Time `json:"expires"`
}

func DBEventToLocalEvent(event database.Event) Event {
	return Event{
		EventID:   event.EventID,
		EventName: event.EventName,
		AuthorID:  event.AuthorID,
		Date:      event.Date,
		CreatedAt: event.CreatedAt,
		UpdatedAt: event.UpdatedAt,
	}
}

func DBUserToLocalUser(user database.User) User {
	return User{
		UserID:    user.UserID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Elevation: user.Elevation,
	}
}

func DBRefrTokenToLocalRefrToken(DBRefr_Token database.RefreshToken) RefreshToken {
	return RefreshToken{
		RefrToken: DBRefr_Token.RefrToken,
		UsersID:   DBRefr_Token.UsersID,
		Expires:   DBRefr_Token.Expires,
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
