package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
	"flag"
	"log"

	"final-project/pkg/dinapp/mailer"
	"final-project/pkg/dinapp/model"
	"final-project/pkg/jsonlog"

	_ "github.com/lib/pq"
	"github.com/peterbourgon/ff/v3"
)

type config struct {
	port int
	env  string
	db   struct {
		port       int
		env        string
		fill       bool
		migrations string
		db         struct {
			dsn string
		}
	}

	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config config
	models model.Models
	logger *jsonlog.Logger
	wg     sync.WaitGroup
	mailer mailer.Mailer
}

// var movies []model.Movies

func main() {

	var confg config
	fs := flag.NewFlagSet("demo-app-demo-app", flag.ContinueOnError)

	log.Println("Starting API server")
	flag.IntVar(&confg.port, "port", 8082, "API server port")
	flag.StringVar(&confg.env, "env", "development", "Environment(development|staging|production)")

	var (
		// confg      config
		// fill       = fs.Bool("fill", false, "Fill database with dummy data")
		// migrations = fs.String("migrations", "", "Path to migration files folder. If not provided, migrations do not applied")
		port  = fs.Int("port", 8082, "API server port")
		env   = fs.String("env", "development", "Environment (development|staging|production)")
		//dbDsn = fs.String("dsn", "postgres://postgres:dinaisthebest@localhost:5434/postgres?sslmode=disable", "PostgreSQL DSN")
	)
	//  postgres://<username>:<password>@localhost:<port>/<db_name>?sslmode=disable   dinaabitova
	// flag.Float64Var(&confg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	// flag.IntVar(&confg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	// flag.BoolVar(&confg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&confg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&confg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&confg.smtp.username, "smtp-username", "0abf276416b183", "SMTP username")
	flag.StringVar(&confg.smtp.password, "smtp-password", "d8672aa2264bb5", "SMTP password")
	flag.StringVar(&confg.smtp.sender, "smtp-sender", "Greenlight <no-reply@greenlight.alexedwards.net>", "SMTP sender")
//"host=db port=5432 user=postgres dbname=postgres password=dinaisthebest sslmode=disable"
	flag.StringVar(&confg.db.db.dsn, "db-dsn", "host=db port=5432 user=postgres dbname=postgres password=dinaisthebest sslmode=disable", "PostgreSQL DSN")

	flag.Parse()

	// logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVars()); err != nil {
		logger.PrintFatal(err, nil)
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	confg.port = *port
	confg.env = *env
	// confg.fill = *fill
	//confg.db.db.dsn = *dbDsn
	// confg.migrations = *migrations

	logger.PrintInfo("starting application with configuration", map[string]string{
		"port": fmt.Sprintf("%d", confg.port),
		// "fill": fmt.Sprintf("%t", confg.fill),
		"env": confg.env,
		"db":  confg.db.db.dsn,
	})

	db, err := openDB(confg)

	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.PrintFatal(err, nil)
		}
	}()

	// if err != nil {
	// 	var pgErr *os.SyscallError
	// 	if errors.As(err, &pgErr) {
	// 		log.Fatalf("Error opening database: %s\n", pgErr.Error())
	// 		return
	// 	}
	// 	log.Fatal(err)
	// 	return
	// }
	// defer db.Close()

	app := &application{
		config: confg,
		models: model.NewModels(db),
		logger: logger,
	}

	// srv := &http.Server{
	// 	Addr:         fmt.Sprintf(":%d", confg.port),
	// 	Handler:      app.routes(),
	// 	IdleTimeout:  time.Minute,
	// 	ReadTimeout:  10 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// }

	// logger.PrintInfo("starting server", map[string]string{
	// 	"addr": srv.Addr,
	// 	"env":  confg.env,
	// })
	// err = srv.ListenAndServe()

	// logger.PrintFatal(err, nil)

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sqlx.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config // struct.
	db, err := sqlx.Open("postgres", cfg.db.db.dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
