package main

import (
	"linksparser/config"
	"time"

	//"os"
	//"time"
	"linksparser/mysql"
	"os"
)

var (
	MYSQL   mysql.Instance
	CONF    config.Configuration
	STREAMS Streams
)

func main() {
	path, _ := os.Getwd()

	CONF.Create(path + "/config.json")

	// Connect to MysqlDB
	MYSQL = mysql.CreateConnection(CONF.MysqlHost, CONF.MysqlDb, CONF.MysqlLogin, CONF.MysqlPass)

	if CONF.Env == "local" {
		//task := MYSQL.GetFreeTask(937675)
		//task.SetTimeout(2)
		//
		//browser := Browser{}
		//browser.Init()
		//buf, _ := browser.ScreenShot("https://www.thesaurus.com/browse/redirect#:~:text=divert,change")
		//err := ioutil.WriteFile("testImage.jpg", *buf, 0644)
		//if err != nil {
		//	fmt.Println("ERR.JobHandler.Run.Screenshot.2", err)
		//}
		//
		//log.Fatal("")

		//go func() {
		//	job := JobHandler{}
		//	job.IsStart = true
		//	if job.Browser.Init() {
		//		job.Run(2)
		//	}
		//}()
	}

	//else if MYSQL.CountWorkingTasks() > 0 {
	//	conf := MYSQL.GetConfig()
	//	extra := conf.GetExtra()
	//	if extra.CountStreams > 0 {
	//		STREAMS.StartLoop(extra.CountStreams, extra.LimitStreams, extra.CmdStreams)
	//	}
	//}


	routes := Routes{}
	routes.Run()

	time.Sleep(time.Minute)
}