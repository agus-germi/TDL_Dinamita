package services

import (
    "errors"
    "TDL_Dinamita/internal/models" // Cambia "myapp" al nombre de tu módulo
)

var tables []models.Table // Simulación de base de datos

func ReserveTable(userID, tableID int, date, time string) (models.Reservation, error) {
    // Aquí iría la lógica para reservar una mesa
    for _, table := range tables {
        if table.ID == tableID && table.IsAvailable {
            table.IsAvailable = false
            return models.Reservation{ID: 1, TableID: tableID, UserID: userID, Date: date, Time: time}, nil
        }
    }
    return models.Reservation{}, errors.New("mesa no disponible")
}
