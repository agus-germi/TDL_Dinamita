package main

import (
	"context"
	"fmt"

	"github.com/agus-germi/TDL_Dinamita/database"
	"github.com/agus-germi/TDL_Dinamita/internal/api"
	"github.com/agus-germi/TDL_Dinamita/internal/repository"
	"github.com/agus-germi/TDL_Dinamita/internal/service"
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

func configureLifeCycleHooks(lc fx.Lifecycle, api *api.API, e *echo.Echo) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				fmt.Println("Starting application...")
				// El valor de "address" podemos leerlo de la variable de entorno API_PORT
				go api.Start(e, ":8080") // api.Start(e, address)

				return nil
			},

			OnStop: func(ctx context.Context) error {
				fmt.Println("Shuting down application...")
				return nil
			},
		},
	)

}
