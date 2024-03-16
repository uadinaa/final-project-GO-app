package main

// import (
// 	"database/sql"
// 	"flag"
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/mux"
// 	_ "github.com/lib/pq"
// 	"main.go/pkg/dinapp/model"
// )

import (
	"errors"
	"os"

	// "github.com/uadinaa/final-project-GO-app/tree/main/model"
	// model "command-line-arguments/Users/dinaabitova/code/golan/final-project/pkg/dinapp/model/movies.go"
	"database/sql"
	"flag"
	"log"
	"net/http"

	// "Users/dinaabitova/code/golan/final-project/pkg/dinapp/model"
	"github.com/uadinaa/final-project/pkg/dinapp/model"

	// "main.go/pkg/dinapp/model"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// type Server struct {
// 	httpServer *http.Server
// }

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
}

// func main() {
// 	// For windows use postgresql:// instead of postgres:// in the connection string first part
// 	db, err := sql.Open("postgres", "postgres://postgres:dinaisthebest@localhost/postgres?sslmode=disable")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Connected to DB!")
// }

func main() {

	var confg config
	flag.StringVar(&confg.port, "port", ":8082", "API server port")
	flag.StringVar(&confg.env, "env", "development", "Environment(development|staging|production)")
	//  postgres://<username>:<password>@localhost:<port>/<db_name>?sslmode=disable
	flag.StringVar(&confg.db.dsn, "db-dsn", "postgres://postgres:dinaisthebest@localhost/postgres?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	db, err := openDB(confg)
	if err != nil {
		var pgErr *os.SyscallError
		if errors.As(err, &pgErr) {
			log.Fatalf("Error opening database: %s\n", pgErr.Error())
			return
		}
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &application{
		config: confg,
		models: model.NewModels(db),
	}
	app.run()
}

func (app *application) run() {
	r := mux.NewRouter()

	router := r.PathPrefix("/api/main").Subrouter()

	// router.HandleFunc("/api/musics", getMusics).Methods("GET") //это ссылканын сондары биз жасап жатырмыз
	router.HandleFunc("/movieGet/{movieId:[0-9]+}", app.getMoviesHandler).Methods("GET")
	router.HandleFunc("/moviePost", app.createMoviesHandler).Methods("POST")

	router.HandleFunc("/movieUpdate/{movieId:[0-9]+}", app.updateMovieHandler).Methods("PUT")
	router.HandleFunc("/movieDelete/{movieId:[0-9]+}", app.deleteMovieHandler).Methods("DELETE")

	log.Printf("Starting server on %s\n", app.config.port)
	errorr := http.ListenAndServe(app.config.port, r)
	log.Fatal(errorr)

}

func openDB(confg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", confg.db.dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
