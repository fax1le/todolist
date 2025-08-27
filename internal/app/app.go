package application

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"
	"todo/internal/config"
	"todo/internal/http/handlers"
	"todo/internal/http/handlers/auth"
	"todo/internal/http/handlers/todo"
	"todo/internal/log"
	"todo/internal/middleware"
	"todo/internal/storage/postgres"
	redis_ "todo/internal/storage/redis"

	"github.com/redis/go-redis/v9"
)

type App struct {
	Cfg    config.Config
	Server *http.Server
	DB     *sql.DB
	Cache  *redis.Client
	Logger *slog.Logger
}

func New(cfg config.Config) *App {
	return &App{
		Cfg: cfg,
		Server: &http.Server{
			Addr:         cfg.Addr,
			Handler:      nil,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (app *App) Init() {
	var err error

	app.Logger,_ = log.NewLogger(app.Cfg)
	app.DB, err = postgres.StartDB(app.Cfg, app.Logger)

	if err != nil {
		app.Logger.Error("postgres connection failed", "err", err)
		os.Exit(1)
	} else {
		app.Logger.Info("postgres connection established")
	}

	app.Cache, err = redis_.StartRedis(app.Cfg, app.Logger)

	if err != nil {
		app.Logger.Error("redis connection failed", "err", err)
		os.Exit(1)
	} else {
		app.Logger.Info("redis connection established")
	}

	authH := &auth.AuthHandler{DB: app.DB, Cache: app.Cache, Logger: app.Logger}

	tasksH := &todo.TasksHandler{DB: app.DB, Cache: app.Cache, Logger: app.Logger}

	mux := http.NewServeMux()
	base := &handlers.BaseHandler{AuthHandler: authH, TasksHandler: tasksH, Mux: mux}
	base.HandleRoutes()

	middleware := middleware.LoggingMiddleWare(
		middleware.AuthMiddleWare(base.Mux, app.Logger, app.Cache),
		app.Logger)

	app.Server.Handler = middleware
}

func (app *App) Run() {
	if app.Server.Handler == nil {
		logger := slog.Logger{}
		logger.Error("Application was not initialized")
		os.Exit(1)
	}

	app.Logger.Info("Application started on port " + app.Cfg.Addr)
	app.Server.ListenAndServe()
}
