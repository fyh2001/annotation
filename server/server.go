package main

import (
	"server/database"
	"server/routes"
)

func main() {
	database.InitMySQL()
	database.InitRedis()

	r := routes.InitRouter()

	r.Run(":3030")
}
