package main

import (
	"context"
	"fmt"

	"github.com/agus-germi/TDL_Dinamita/database"
	"github.com/agus-germi/TDL_Dinamita/internal/api"
	"github.com/agus-germi/TDL_Dinamita/internal/repository"
	"github.com/agus-germi/TDL_Dinamita/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
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
			echo.New,
		),
		fx.Invoke(
			configureLifeCycleHooks,
		),
	).Run()
}

func configureLifeCycleHooks(lc fx.Lifecycle, api *api.API, e *echo.Echo, dbPool *pgxpool.Pool) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				fmt.Println("Starting application...")
				// El valor de "address" podemos leerlo de la variable de entorno API_PORT (o APP_PORT)
				//api.SetStaticFiles(e)
				//api.SetRoutes(e)

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
