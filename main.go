package main

import (
	"context"
	"fmt"

	"github.com/agus-germi/TDL_Dinamita/database"
	"github.com/agus-germi/TDL_Dinamita/internal/api"
	"github.com/agus-germi/TDL_Dinamita/internal/repository"
	"github.com/agus-germi/TDL_Dinamita/internal/service"
	"github.com/agus-germi/TDL_Dinamita/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx" // fx es un framework que sirve para inyectar dependencias.
)

func main() {
	fx.New(
		fx.Provide(
			context.Background,
			database.CreateConnection,
			repository.New,
			service.New,
			api.New,
			api.NewValidator,
			echo.New,
			logger.New,
			// logger.InitAppLoggerAdapter,
			// logger.InitEchoLoggerAdapter,
		),
		fx.Invoke(
			configureLifeCycleHooks,
		),
	).Run()
}

func configureLifeCycleHooks(lc fx.Lifecycle, api *api.API, e *echo.Echo, dbPool *pgxpool.Pool, log logger.Logger) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				fmt.Println("Starting application...")

				err := godotenv.Load("/usr/src/app/.env")
				if err != nil {
					log.Fatalf("'.env' file couldn't be loaded: %v", err)
				}
				log.Info("'.env' file loaded successfully.")

				// El valor de "address" podemos leerlo de la variable de entorno API_PORT (o APP_PORT)

				// Echo middleware logger configuration --> transfering echo http logs management to logrus
				e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
					Output: log.Writer(),
				}))

				go func() {
					e.Logger.Fatal(api.Start(e, ":8080")) // api.Start(e, address)
				}()

				return nil
			},

			OnStop: func(ctx context.Context) error {
				fmt.Println("Shutting down application...")
				e.Shutdown(ctx)
				dbPool.Close()
				return nil
			},
		},
	)

}
