package main

import (
	"log"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/noonyuu/comparison/backend/internal/db"
	"github.com/noonyuu/comparison/backend/internal/graphql"

	"github.com/noonyuu/comparison/backend/internal/rest"
)

func main() {
	// データベース接続
	database, err := db.Connect(os.Getenv("MONGODB_URI"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Disconnect()

	// REST API サーバーの設定
	restRouter := rest.NewRouter(database)
	go func() {
		log.Println("REST server starting on port 8080...")
		if err := http.ListenAndServe(":8080", restRouter); err != nil {
			log.Fatalf("Failed to start REST server: %v", err)
		}
	}()

	// GraphQL サーバーの設定
	graphqlHandler := graphql.NewHandler(database)
	http.Handle("/graphql", graphqlHandler)
	go func() {
		log.Println("GraphQL server starting on port 8090...")
		if err := http.ListenAndServe(":8090", nil); err != nil {
			log.Fatalf("Failed to start GraphQL server: %v", err)
		}
	}()

	// グレースフルシャットダウンの設定
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	// REST と GraphQL サーバーのシャットダウン処理をここに追加

	log.Println("Servers stopped")
}
