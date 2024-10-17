package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sebstainsgit/calendar/internal/database"
)

func main() {
	godotenv.Load()

	portStr := os.Getenv("PORT")

	if portStr == "" {
		log.Fatal("Port not found in .env")
		return
	}

	dbURL := os.Getenv("DB_URL")

	conn, err := sql.Open("postgres", dbURL)

	db := database.New(conn)

	if err != nil {
		log.Printf("Unable to connect to database: %s", err)
	}

	apiCfg := apiConfig{
		DB: db,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           3000,
	}))

	//ADD ADMIN TABLE AND ADMINS SHOULD BE ABLE TO GET ALL USERS + EVENTS AND STUFF
	
	//Responds with 200
	router.Get("/ready", handlerReadiness)
	//Responds with error
	router.Get("/error", handlerError)

	userRouter := chi.NewRouter()

	userRouter.Get("/events", apiCfg.middlewareUserAuth(apiCfg.getUsersEvents))

	userRouter.Post("/events", apiCfg.middlewareUserAuth(apiCfg.createEvent))

	userRouter.Delete("/delete_event", apiCfg.middlewareUserAuth(apiCfg.deleteEvent))

	userRouter.Post("/users", apiCfg.createUser)

	userRouter.Get("/users", apiCfg.getAllUsers)

	userRouter.Delete("/users", apiCfg.middlewareUserAuth(apiCfg.deleteUserSelf))
	
	router.Mount("/user", userRouter)

	adminRouter := chi.NewRouter()

	router.Mount("/admin", adminRouter)
	
	srv := &http.Server{
		Addr:    ":" + portStr,
		Handler: router,
	}

	log.Println("Server running on port", portStr)

	log.Fatal(srv.ListenAndServe())
}
