package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/snorremd/gocomment/api/db"
	"github.com/snorremd/gocomment/api/model"
	"github.com/snorremd/gocomment/api/router"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// server starts a go http server and returns any error encountered
func server(hostAddress string, router *router.Router) error {
	muxRouter := router.Router()

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	credentialsOk := handlers.AllowCredentials()

	return http.ListenAndServe(hostAddress, handlers.CORS(originsOk, headersOk, methodsOk, credentialsOk)(muxRouter))
}

// serveCmd represents the serve command which starts the api server
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts gocomment http server",
	Long: `Starts the gocomment http server on the selected host and port
using the specified database.

If the specified database does not exist it will automatically create it.`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db, err := db.DB(viper.GetString("db"))

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

		listen := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))

		if err := server(listen, router); err != nil {
			log.Fatal(err)
		}

		log.Println("App successfully ran")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.
	serveCmd.PersistentFlags().String("host", "", "host to listen to, defaults to localhost")
	serveCmd.PersistentFlags().Uint("port", 0, "port to bind to, defaults to 8080")
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", "8080")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
