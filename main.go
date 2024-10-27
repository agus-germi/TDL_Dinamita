package main

import (
	"context"

	"github.com/agus-germi/TDL_Dinamita/database"
	"github.com/agus-germi/TDL_Dinamita/internal/repository"
	"github.com/agus-germi/TDL_Dinamita/internal/service"
	"go.uber.org/fx" // fx es un framework que sirve para inyectar dependencias.
)

func main() {
	app := fx.New(
		fx.Provide(
			context.Background,
			database.New,
			repository.New,
			service.New,
		),
		fx.Invoke(
		// func(db *sqlx.DB) {
		// 	_, err := db.Query("SELECT * FROM users")
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// },
		),
	)

	app.Run()
}
