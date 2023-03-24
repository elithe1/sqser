package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"net/http"
	"os"
	"sqser/app"
	_ "sqser/plugins/enrichers/all"
	_ "sqser/plugins/filters/all"
	_ "sqser/plugins/inputs/all"
	_ "sqser/plugins/outputs/all"
)

func health(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	//todo switch to contextual logger with reqids for each request https://medium.com/@gosamv/using-gos-context-library-for-logging-4a8feea26690
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
	err := godotenv.Load(".env.local")
	if err != nil {
		logger.Sugar().Error(err)
	}

	port, ok := os.LookupEnv("PORT")

	if !ok {
		zap.S().Error("Couldn't retrieve server port")
		return
	}

	a := app.NewApp()
	fmt.Sprint(len(a.Config.Config.Inputs))
	logger.Info("Starting server at port " + port)
	http.HandleFunc("/health", health)

	http.HandleFunc("/get-item", a.GetItem)
	http.HandleFunc("/delete-item", a.DeleteItem)
	http.HandleFunc("/move-items", a.MoveItems)
	http.HandleFunc("/list-items", a.ListItems)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Sugar().Error(err)
	}
}
