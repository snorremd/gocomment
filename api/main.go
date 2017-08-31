package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/snorremd/gocomment/api/db"
	"github.com/snorremd/gocomment/api/model"
	"github.com/snorremd/gocomment/api/router"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// server return new http server with comment routes
func server(hostAddress string, router *router.Router) error {
	muxRouter := router.Router()

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	credentialsOk := handlers.AllowCredentials()

	return http.ListenAndServe(hostAddress, handlers.CORS(originsOk, headersOk, methodsOk, credentialsOk)(muxRouter))
}

func main() {

	db, err := db.DB(os.Getenv("DB"))

	if err != nil {
		log.Fatal("Could not connect to database", err)
	}

	defer db.Close()

	if err := model.Migrate(db); err != nil {
		log.Fatal("Could not migrate database", err)
	}

	router := &router.Router{
		Commenter: model.SqliteCommentStore{
			DB: db,
		},
	}

	if err := server(os.Getenv("HOST"), router); err != nil {
		log.Fatal(err)
	}

	log.Println("App successfully ran")

}
