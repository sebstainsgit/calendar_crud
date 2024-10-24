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

	JWTString := os.Getenv("JWT_SECRET")

	apiCfg := apiConfig{
		DB:         db,
		JWT_SECRET: JWTString,
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
	//Responds with 200
	router.Get("/ready", handlerReadiness)
	//Responds with error
	router.Get("/error", handlerError)

	router.Post("/login", apiCfg.loginUser)

	userRouter := chi.NewRouter()

	userRouter.Get("/events", apiCfg.middlewareUserAuth(apiCfg.getUsersEvents))

	userRouter.Post("/events", apiCfg.middlewareUserAuth(apiCfg.createEvent))

	userRouter.Get("/refresh", apiCfg.makeJWTfromRefrToken)

	userRouter.Delete("/delete_event", apiCfg.middlewareUserAuth(apiCfg.deleteEvent))

	userRouter.Post("/users", apiCfg.createUser)

	userRouter.Post("/update_self", apiCfg.middlewareUserAuth(apiCfg.updateUserInfo))

	userRouter.Delete("/delete_self", apiCfg.middlewareUserAuth(apiCfg.deleteUserSelf))

	router.Mount("/user", userRouter)

	adminRouter := chi.NewRouter()

	adminRouter.Get("/remove_expired_tokens", apiCfg.middlewareAdminAuth(apiCfg.removeOldRefrTokens))

	adminRouter.Post("/admins", apiCfg.createAdmin)

	adminRouter.Delete("/delete_user", apiCfg.middlewareAdminAuthWithUser(apiCfg.deleteUser))

	adminRouter.Get("/users", apiCfg.middlewareAdminAuth(apiCfg.getAllUsers))

	router.Mount("/admin", adminRouter)

	srv := &http.Server{
		Addr:    ":" + portStr,
		Handler: router,
	}

	log.Println("Server running on port", portStr)

	log.Fatal(srv.ListenAndServe())
}
