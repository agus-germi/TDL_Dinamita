package main

import (
	"context"
	"fmt"

	"github.com/agus-germi/TDL_Dinamita/database"
	"github.com/agus-germi/TDL_Dinamita/internal/repository"
	"github.com/agus-germi/TDL_Dinamita/internal/service"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx" // fx es un framework que sirve para inyectar dependencias.
)

func main() {
	fx.New(
		fx.Provide(
			context.Background,
			database.CreateConnection,
			repository.New,
			service.New,
		),
		fx.Invoke(
			configureLifeCycleHooks,
		),
	).Run()
}

func configureLifeCycleHooks(lc fx.Lifecycle, db *sqlx.DB, repo repository.Repository) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				fmt.Println("Starting application...")
				return nil
			},

			OnStop: func(ctx context.Context) error {
				fmt.Println("Shuting down application...")
				return nil
			},
		},
	)

}
