package main

import (
	"encoding/json"
	"fmt"
	"linksparser/services"
	"log"
	"net/http"
	"strconv"
)

func (rt *Routes) startLoopStreams(w http.ResponseWriter, r *http.Request) {
	count := services.ToInt(r.FormValue("count"))
	limit := services.ToInt(r.FormValue("limit"))
	cmd := r.FormValue("cmd")

	if limit < 1 {
		limit = 10
	}

	config := MYSQL.GetConfig()
	extra := config.GetExtra()
	extra.CmdStreams = cmd
	extra.LimitStreams = limit
	extra.CountStreams = count
	_ = MYSQL.SetExtra(extra)

	fmt.Println(count, limit, cmd)
	STREAMS.StartLoop(count, limit, cmd)

	err := json.NewEncoder(w).Encode(map[string]bool{
		"status": true,
	})
	if err != nil {
		log.Println("Routes.StartLoopStreams.HasError", err)
	}
}

func (rt *Routes) stopLoopStreams(w http.ResponseWriter, r *http.Request) {
	STREAMS.isStarted = false
	go STREAMS.StopAll()

	err := json.NewEncoder(w).Encode(map[string]bool{
		"status": true,
	})
	if err != nil {
		log.Println("Routes.StopLoopStreams.HasError", err)
	}
}

func (rt *Routes) countLoopStreams(w http.ResponseWriter, r *http.Request) {
	count := STREAMS.Count()

	_, err := fmt.Fprintln(w, strconv.Itoa(count))
	if err != nil {
		log.Println("Routes.CountLoopStreams.HasError", err)
	}
}
