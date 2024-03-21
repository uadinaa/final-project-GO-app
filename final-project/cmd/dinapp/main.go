package main

import (
	"errors"
	"os"

	"database/sql"
	"flag"
	"log"
	"net/http"

	"final-project/pkg/dinapp/model"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

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

// var movies []model.Movies

func main() {

	var confg config
	flag.StringVar(&confg.port, "port", ":8082", "API server port")
	flag.StringVar(&confg.env, "env", "development", "Environment(development|staging|production)")
	//  postgres://<username>:<password>@localhost:<port>/<db_name>?sslmode=disable   dinaabitova
	flag.StringVar(&confg.db.dsn, "db-dsn", "postgres://postgres:dinaisthebest@localhost:5434/postgres?sslmode=disable", "PostgreSQL DSN")
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

	router := mux.NewRouter()

	r := router.PathPrefix("/api/main").Subrouter()

	r.HandleFunc("/movie/{movieId:[0-9]+}", app.getMoviesHandler).Methods("GET")
	r.HandleFunc("/movie", app.createMoviesHandler).Methods("POST")
	r.HandleFunc("/movie/{movieId:[0-9]+}", app.updateMovieHandler).Methods("PUT")
	r.HandleFunc("/api/movie/{movieId:[0-9]+}", app.deleteMovieHandler).Methods("DELETE")

	r.HandleFunc("/genre/{movieId:[0-9]+}", app.getGenresHandler).Methods("GET")
	r.HandleFunc("/genre", app.createGenresHandler).Methods("POST")
	r.HandleFunc("/genre/{movieId:[0-9]+}", app.updateGenreHandler).Methods("PUT")
	r.HandleFunc("/genre/{movieId:[0-9]+}", app.deleteGenreHandler).Methods("DELETE")

	log.Printf("Starting server on %s\n", app.config.port)
	errorr := http.ListenAndServe(app.config.port, r)
	log.Fatal(errorr)

}

func openDB(confg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", confg.db.dsn)
	// db, err := sql.Open("postgres", "postgres://postgres:dinaisthebest@localhost:5434/postgres?sslmode=disable")
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// router := mux.NewRouter()

// musics = append(musics, MusicMax{ID: "1", Isbn: "1234", Name: "Tadow", Artist: &Author{FirstName: "Masaego", LastName: "FKJ"}})
// musics = append(musics, MusicMax{ID: "2", Isbn: "5678", Name: "Далада", Artist: &Author{FirstName: "PRiNCE", LastName: "Папа"}})
// musics = append(musics, MusicMax{ID: "3", Isbn: "9012", Name: "35+34", Artist: &Author{FirstName: "Ariana", LastName: "Grande"}})

// router.HandleFunc("/api/musics", getMusics).Methods("GET") //это ссылканын сондары биз жасап жатырмыз
// router.HandleFunc("/api/musics/{id}", getMusic).Methods("GET")
// router.HandleFunc("/api/musics", createMusic).Methods("POST")
// router.HandleFunc("/api/musics/{id}", updateMusics).Methods("PUT")
// router.HandleFunc("/api/musics/{id}", deleteMusics).Methods("DELETE") // а тут можно айдига не только намберс но и стринги запихнуть можно
// //r.HandleFunc("/restaurants/{id:[0-9]+}", restaurant) прикол это типа онли фор диджитс

// // const PORT = ":8080" можно и так
// // log.Fatal(http.ListenAndServe(":8000", router)) BEFORE
// const port = ":8000"
// log.Fatal(http.ListenAndServe(port, router))
// log.Printf("starting server on %s \n", port)
// // 	log.Printf("Starting server on %s\n", PORT)

// errorr := http.ListenAndServe(port, router) //корейк не шыгады екеен, но пон что еррор хаха
// log.Fatal(errorr)

// router := r.PathPrefix("/api/main").Subrouter()

// movies = append(movies, model.Movies{Title: "flmf", Description: "hi", YearOfProduction: 2300, GenreId: "1"})

// router.HandleFunc("/api/musics", getMusics).Methods("GET") //это ссылканын сондары биз жасап жатырмыз

// "Users/dinaabitova/code/golan/final-project/pkg/dinapp/model"
// "github.com/uadinaa/final-project/pkg/dinapp/model"

// "main.go/pkg/dinapp/model"

// "github.com/uadinaa/final-project-GO-app/tree/main/model"
// model "command-line-arguments/Users/dinaabitova/code/golan/final-project/pkg/dinapp/model/movies.go"

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
