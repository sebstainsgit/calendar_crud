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

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           3000,
	}))
	
	apiRouter := chi.NewRouter()
	//Responds with 200
	apiRouter.Get("/ready", handlerReadiness)
	//Responds with error
	apiRouter.Get("/error", handlerError)
	//Works accross whole API, not specific to elevation
	apiRouter.Post("/login", apiCfg.loginUser)
	
	r.Mount("/api", apiRouter)
	
	userRouter := chi.NewRouter()

	userRouter.Get("/events", apiCfg.middlewareUserAuth(apiCfg.getUsersEvents))

	userRouter.Post("/events", apiCfg.middlewareUserAuth(apiCfg.createSelfEvent))

	userRouter.Post("/group_event", apiCfg.middlewareUserAuth(apiCfg.createGroupEvent))

	userRouter.Get("/refresh", apiCfg.makeJWTfromRefrToken)

	userRouter.Post("/update_event", apiCfg.middlewareUserAuth(apiCfg.updateEvent))

	userRouter.Delete("/delete_event", apiCfg.middlewareUserAuth(apiCfg.deleteEvent))

	userRouter.Delete("/remove_from_event", apiCfg.middlewareUserAuth(apiCfg.removeSelfFromEvent))

	userRouter.Post("/users", apiCfg.createUser)

	userRouter.Post("/update_self", apiCfg.middlewareUserAuth(apiCfg.updateUserInfo))

	userRouter.Delete("/delete_self", apiCfg.middlewareUserAuth(apiCfg.deleteUserSelf))

	r.Mount("/user", userRouter)

	adminRouter := chi.NewRouter()

	adminRouter.Get("/remove_expired_tokens", apiCfg.middlewareAdminAuth(apiCfg.removeOldRefrTokens))

	adminRouter.Post("/admins", apiCfg.createAdmin)

	adminRouter.Delete("/delete_user", apiCfg.middlewareAdminAuthWithUser(apiCfg.deleteUser))

	adminRouter.Get("/users", apiCfg.middlewareAdminAuth(apiCfg.getAllUsers))

	r.Mount("/admin", adminRouter)

	srv := &http.Server{
		Addr:    ":" + portStr,
		Handler: r,
	}

	log.Println("Server running on port", portStr)

	log.Fatal(srv.ListenAndServe())
}
