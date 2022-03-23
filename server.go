package main

import (
    "github.com/Duclmict/go-backend/routes"
    "github.com/Duclmict/go-backend/database"
    "github.com/Duclmict/go-backend/config"
)

func main() {
    config.LoadENV()
    database.Init()
    routes.InitRouter()
}