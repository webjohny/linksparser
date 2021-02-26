package main

import (
	"linksparser/config"
	//"os"
	//"time"
	"linksparser/mysql"
	"os"
)

var (
	UTILS   Utils
	MYSQL   mysql.Instance
	CONF    config.Configuration
	//STREAMS Streams
)

func main() {
	path, _ := os.Getwd()

	CONF.Create(path + "/config.json")

	// Connect to MysqlDB
	MYSQL = mysql.CreateConnection(CONF.MysqlHost, CONF.MysqlDb, CONF.MysqlLogin, CONF.MysqlPass)

	// Run routes

	//if CONF.Env == "local" {
		//task := MYSQL.GetFreeTask(564805)
		//task.SetTimeout(2)

		//go func() {
			//job := JobHandler{}
			//job.IsStart = true
			//if job.Browser.Init() {
				//job.Run(2)
				//job.Run(1)
				//job.Run(1)
			//}
		//}()

		//time.Sleep(100)

		//go func() {
		//	job := JobHandler{}
		//	job.IsStart = true
		//	if job.Browser.Init() {
		//		job.Run(1)
		//		job.Run(0)
		//		job.Run(2)
		//	}
		//}()
		//
		//time.Sleep(100)
		//
		//go func() {
		//	job := JobHandler{}
		//	job.IsStart = true
		//	if job.Browser.Init() {
		//		job.Run(0)
		//		job.Run(2)
		//		job.Run(1)
		//	}
		//}()
	//}else if MYSQL.CountWorkingTasks() > 0 {
		//conf := MYSQL.GetConfig()
		//extra := conf.GetExtra()
		//if extra.CountStreams > 0 {
			//STREAMS.StartLoop(extra.CountStreams, extra.LimitStreams, //extra.CmdStreams)
		//}
	//}

	//routes := Routes{}
	//routes.Run()

	//time.Sleep(time.Minute)
}