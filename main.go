package main

import (
	"fmt"
	"io/ioutil"
	"linksparser/config"
	"linksparser/wordpress"
	"log"
	"strconv"
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
		task := MYSQL.GetFreeTask(937675)
		wp := wordpress.Base{}
		wp.Connect(`http://philli.beget.tech`, task.Login, task.Password, 1)
		//log.Fatal(wp.UploadFile("https://ru.ex-rate.com/wa-data/public/crcy/images/USD/UAH.jpg", 33, nil, false))
		//task.SetTimeout(2)
		//
		//
		//log.Fatal("")
		links := []string{
			//"https://www.investopedia.com/terms/v/virtual-reality.asp",
			//"https://www.iberdrola.com/innovation/virtual-reality",
			//"https://arvr.google.com/vr/",
			"https://www.pcmag.com/picks/the-best-vr-headsets",
			//"http://45.67.59.191/home/data",
		}

		go func() {
			job := JobHandler{}
			job.IsStart = true
			if job.Browser.Init() {
				for i := 0; i < len(links); i++ {
					link := links[i]
					buf, err := job.Browser.ScreenShot(link)
					if err != nil {
						fmt.Println("ERR.JobHandler.Run.Screenshot", err)
					} else {
						err = ioutil.WriteFile(CONF.ImgPath+"/"+strconv.Itoa(task.Id)+"-"+strconv.Itoa(i)+".jpg", *buf, 0644)
						if err != nil {
							fmt.Println("ERR.JobHandler.Run.Screenshot.2", err)
						}
						file, err := wp.UploadFile("", 0, buf, false)
						if err != nil {
							fmt.Println("ERR.JobHandler.Run.Screenshot.3", err)
						}
						log.Fatal(file)
						//fmt.Print(job.ExtractSimilarWebData(link))
						//time.Sleep(time.Second * time.Duration(rand.Intn(10)))
					}
				}
			}
		}()
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