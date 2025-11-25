package main

import (
	"AvitoProject/internal/config"
	db "AvitoProject/internal/connections"
	"AvitoProject/internal/delivery/http/teamhandler"
	"AvitoProject/internal/delivery/http/userhandler"
	"AvitoProject/internal/delivery/http/prhandler"
	"AvitoProject/internal/repository/postgres"
	"AvitoProject/internal/usecase"
	"AvitoProject/utils"

	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 5 * time.Second
	idleTimeout  = 60 * time.Second
)

func main() {
	cfg := config.GetConfig() 

	dbConn, err := db.New(cfg)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer dbConn.CloseAll()

	repoCfg := postgres.New(dbConn)

	useCaseCfg := usecase.New(repoCfg)

	teamHandler := teamhandler.New(useCaseCfg.TeamUsecase)
	userHandler := userhandler.New(useCaseCfg.UserUsecase)
	prHandler := prhandler.New(useCaseCfg.PrUsecase)

	router := mux.NewRouter()

	router.HandleFunc("/team/add", teamHandler.CreateTeam).Methods(http.MethodPost)
	router.HandleFunc("/team/get", teamHandler.GetTeam).Methods(http.MethodGet)

	router.HandleFunc("/users/setIsActive", userHandler.SetIsActive).Methods(http.MethodPost)
	router.HandleFunc("/users/getReview", userHandler.GetReview).Methods(http.MethodGet)

	router.HandleFunc("/pullRequest/create", prHandler.Create).Methods(http.MethodPost)
	router.HandleFunc("/pullRequest/merge", prHandler.Merge).Methods(http.MethodPost)
	router.HandleFunc("/pullRequest/reassign", prHandler.Reassign).Methods(http.MethodPost)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, "OK", "service is running", http.StatusOK) // или просто {"status":"ok"}
	}).Methods(http.MethodGet)

	addr := ":8080" 
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	log.Println("🚀 Server starting on http://localhost" + addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("server failed:", err)
	}
}