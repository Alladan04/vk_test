package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authDelivery "github.com/Alladan04/vk_test/internal/pkg/auth/delivery/http"
	authRepo "github.com/Alladan04/vk_test/internal/pkg/auth/repo"
	authUsecase "github.com/Alladan04/vk_test/internal/pkg/auth/usecase"
	marketDelivery "github.com/Alladan04/vk_test/internal/pkg/market/delivery/http"
	marketRepo "github.com/Alladan04/vk_test/internal/pkg/market/repo"
	marketUsecase "github.com/Alladan04/vk_test/internal/pkg/market/usecase"
	"github.com/Alladan04/vk_test/internal/pkg/middleware"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
}
func main() {
	db, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	AuthRepo := authRepo.NewAuthRepo(db)
	AuthUsecase := authUsecase.NewAuthUsecase(AuthRepo)
	AuthDelivery := authDelivery.NewAuthHandler(AuthUsecase)

	MarketRepo := marketRepo.NewMarketRepo(db)
	MarketUsecase := marketUsecase.NewMarketUsecase(MarketRepo)
	MarketDelivery := marketDelivery.NewMarketHandler(MarketUsecase)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.Handle("/signup", http.HandlerFunc(AuthDelivery.SignUp)).Methods(http.MethodPost, http.MethodOptions)
		auth.Handle("/login", http.HandlerFunc(AuthDelivery.SignIn)).Methods(http.MethodPost, http.MethodOptions)
		auth.Handle("/logout", middleware.JwtMiddleware(http.HandlerFunc(AuthDelivery.LogOut))).Methods(http.MethodDelete, http.MethodOptions)

	}
	market := r.PathPrefix("/market").Subrouter()
	{
		market.Handle("/get", middleware.JwtMiddlewareCommon(http.HandlerFunc(MarketDelivery.GetAll))).Methods(http.MethodGet, http.MethodOptions)
		market.Handle("/add", middleware.JwtMiddleware(http.HandlerFunc(MarketDelivery.AddItem))).Methods(http.MethodPost, http.MethodOptions)

	}

	http.Handle("/", r)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	server := http.Server{
		Handler:           r,
		Addr:              ":8080",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("Server Stopped")
		}
	}()
	fmt.Println("Server started ")

	sig := <-signalCh
	fmt.Println("recieved signal: ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Server shutdown failed: " + err.Error())
	}
}
