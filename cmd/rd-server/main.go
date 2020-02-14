package main

import (
	"context"
	"flag"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/cigarframework/reconciled/cmd/rd-server/app"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func main() {
	pConfig := flag.String("config", "", "config file location")
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err != nil {
			logger.Fatal("server closed", zap.Error(err))
		}
		logger.Sync()
	}()

	var file []byte
	file, err = ioutil.ReadFile(*pConfig)
	if err != nil {
		return
	}

	rand.Seed(time.Now().UnixNano())
	config := &app.Config{}
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return
	}

	var server *app.Server
	server, err = app.New(logger, config)
	if err != nil {
		return
	}
	exit := make(chan struct{}, 1)
	err = server.Start(context.TODO())
	if err != nil {
		return
	}
	logger.Info("server started")
	<-exit
}
