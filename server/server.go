package main

import (
	"fmt"
	"net/http"
	"os"
	"server/config"
	"server/database"
	"server/routes"
)

func main() {

	err := os.MkdirAll("./static/image", os.ModePerm)
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Println("目录已创建或已存在:", "./static/image")
	}

	err = os.MkdirAll("./static/file", os.ModePerm)
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Println("目录已创建或已存在:", "./static/file")
	}

	config.InitConfig()

	database.InitMySQL()
	database.InitRedis()

	r := routes.InitRouter()
	r.StaticFS("/images", http.Dir("./static/image"))
	r.StaticFS("/file", http.Dir("./static/file"))

	_ = r.Run(":" + config.Global.Application.Port)
}
