// @title           Expense Tracker API
// @version         1.0
// @description     REST API for managing personal expenses with JWT authentication.
//
// @contact.name   API Support
// @contact.email  support@expensetracker.io
//
// @license.name  MIT
//
// @host      localhost:8080
// @BasePath  /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token obtained from POST /auth/login or /auth/register. Enter value as: **Bearer &lt;token&gt;**

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/luisfucros/expense-tracker-app/docs"
	"github.com/luisfucros/expense-tracker-app/internal/config"
	"github.com/luisfucros/expense-tracker-app/internal/router"
	"github.com/luisfucros/expense-tracker-app/internal/server"
	"github.com/luisfucros/expense-tracker-app/migrations"
	"github.com/luisfucros/expense-tracker-app/internal/repository"
	"github.com/luisfucros/expense-tracker-app/internal/service"
	"github.com/luisfucros/expense-tracker-app/internal/handler"
	dbrepo "github.com/luisfucros/expense-tracker-app/internal/repository/db"
	applogger "github.com/luisfucros/expense-tracker-app/pkg/logger"
)

// Ensure interfaces are satisfied at compile time.
var _ repository.UserRepository = (*dbrepo.UserRepo)(nil)
var _ repository.ExpenseRepository = (*dbrepo.ExpenseRepo)(nil)

func main() {
	// Load .env file if present (ignore error in production/Docker)
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := applogger.New(cfg.Env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync() //nolint:errcheck

	// Connect to MySQL via GORM
	db, err := dbrepo.Connect(cfg.DSN())
	if err != nil {
		log.Sugar().Fatalf("failed to connect to database: %v", err)
	}

	// Get underlying *sql.DB for migrations
	sqlDB, err := db.DB()
	if err != nil {
		log.Sugar().Fatalf("failed to get sql.DB: %v", err)
	}

	// Resolve migrations path relative to the binary or source
	migrationsPath := resolveMigrationsPath()

	// Run database migrations
	if err := migrations.Run(sqlDB, migrationsPath, log); err != nil {
		log.Sugar().Fatalf("failed to run migrations: %v", err)
	}

	// Wire repositories
	userRepo := dbrepo.NewUserRepository(db)
	expenseRepo := dbrepo.NewExpenseRepository(db)

	// Wire services
	authSvc := service.NewAuthService(userRepo, cfg)
	expenseSvc := service.NewExpenseService(expenseRepo)

	// Wire handler
	h := handler.NewHandler(authSvc, expenseSvc, log)

	// Build router
	engine := router.New(cfg, h, log)

	// Create HTTP server
	addr := fmt.Sprintf(":%s", cfg.Port)
	srv := server.New(addr, engine)

	// Start server in background
	go func() {
		log.Sugar().Infof("starting server on %s", addr)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Sugar().Fatalf("server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Sugar().Errorf("server shutdown error: %v", err)
	}

	log.Info("server stopped")
}

// resolveMigrationsPath returns an absolute path to the migrations directory.
func resolveMigrationsPath() string {
	// When running via `go run ./cmd/api`, __file__ is not available in Go,
	// but we can use the executable path or a relative path.
	// Try relative to working directory first.
	if _, err := os.Stat("migrations"); err == nil {
		abs, _ := filepath.Abs("migrations")
		return abs
	}

	// Fallback: use the source file location (useful during development)
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		// filename is .../cmd/api/main.go, go up two levels to project root
		root := filepath.Join(filepath.Dir(filename), "..", "..")
		path := filepath.Join(root, "migrations")
		if _, err := os.Stat(path); err == nil {
			abs, _ := filepath.Abs(path)
			return abs
		}
	}

	// Default fallback
	return "migrations"
}
