package main

import (
	"fmt"
	"io/ioutil"
	"linksparser/config"
	"log"
	//"os"
	//"time"
	"linksparser/mysql"
	"os"
)

var (
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
	//log.Fatal(MYSQL.GetFreeTask(0))

	links := []string{
		"https://exprealty.com/",
		"https://join.exprealty.com/",
		"https://buildingbetteragents.com/about-exp-realty/",
		"https://www.linkedin.com/company/exp-realty",
		"https://www.facebook.com/eXpRealty/",
		"https://exprealty.ca/",
		"https://www.globenewswire.com/news-release/2020/12/14/2144576/0/en/eXp-Realty-Surpasses-40-000-Real-Estate-Agents-Globally-On-Its-Immersive-Cloud-Based-Platform.html",
		"https://www.indeed.com/cmp/Exp-Realty/reviews?fjobtitle=Realtor",
		"https://expcloud.com/",
		"https://www.inman.com/2018/05/30/i-toured-exp-realtys-virtual-reality-heres-what-its-like/",
		"https://www.fortunebuilders.com/exp-realty-the-4-pillars-explained/",
		"https://twitter.com/exprealty",
		"https://www.expaustralia.com.au/",
		"https://finance.yahoo.com/quote/EXPI/",
		"https://exprealty.com/",
		"https://join.exprealty.com/",
		"https://buildingbetteragents.com/about-exp-realty/",
		"https://www.linkedin.com/company/exp-realty",
		"https://www.facebook.com/eXpRealty/",
		"https://exprealty.ca/",
		"https://www.globenewswire.com/news-release/2020/12/14/2144576/0/en/eXp-Realty-Surpasses-40-000-Real-Estate-Agents-Globally-On-Its-Immersive-Cloud-Based-Platform.html",
		"https://www.indeed.com/cmp/Exp-Realty/reviews?fjobtitle=Realtor",
		"https://expcloud.com/",
		"https://www.inman.com/2018/05/30/i-toured-exp-realtys-virtual-reality-heres-what-its-like/",
		"https://www.fortunebuilders.com/exp-realty-the-4-pillars-explained/",
		"https://twitter.com/exprealty",
		"https://www.expaustralia.com.au/",
		"https://finance.yahoo.com/quote/EXPI/",
		"https://exprealty.com/",
		"https://join.exprealty.com/",
		"https://buildingbetteragents.com/about-exp-realty/",
		"https://www.linkedin.com/company/exp-realty",
		"https://www.facebook.com/eXpRealty/",
		"https://exprealty.ca/",
		"https://www.globenewswire.com/news-release/2020/12/14/2144576/0/en/eXp-Realty-Surpasses-40-000-Real-Estate-Agents-Globally-On-Its-Immersive-Cloud-Based-Platform.html",
		"https://www.indeed.com/cmp/Exp-Realty/reviews?fjobtitle=Realtor",
		"https://expcloud.com/",
		"https://www.inman.com/2018/05/30/i-toured-exp-realtys-virtual-reality-heres-what-its-like/",
		"https://www.fortunebuilders.com/exp-realty-the-4-pillars-explained/",
		"https://twitter.com/exprealty",
		"https://www.expaustralia.com.au/",
		"https://finance.yahoo.com/quote/EXPI/",
	}

	job := JobHandler{}
	proxy := NewProxy()
	if proxy == nil {
		log.Fatal("Need free proxies")
	}
	job.IsStart = true
	job.Browser.Init()
	job.proxy = *proxy
	var results []SimilarWebResp
	for i := 0; i < len(links); i++ {
		link := links[i]
		buf, err := job.Browser.ScreenShot(link)
		if err != nil {
			fmt.Println("ERR.JobHandler.Run.Screenshot", err)
		}


		//dsw, err := job.ExtractSimilarWebData(link)
		//if err != nil {
		//	fmt.Println("Iteration #" + strconv.Itoa(i), err)
		//	continue
		//}
		//if dsw != nil {
		//	results = append(results, *dsw)
		//	time.Sleep(time.Second * time.Duration(rand.Intn(20)+3))
		//}
	}
	log.Println(len(results))
	log.Fatal(results)
	job.IsStart = true
	if job.Browser.Init() {
		fmt.Println(job.Run(2))
		//job.Run(1)
		//job.Run(1)
	}
	//if CONF.Env == "local" {
	//	task := MYSQL.GetFreeTask(564805)
	//	task.SetTimeout(2)
	//
	//	go func() {
	//		job := JobHandler{}
	//		job.IsStart = true
	//		if job.Browser.Init() {
	//			job.Run(2)
	//			//job.Run(1)
	//			//job.Run(1)
	//		}
	//	}()
	//}

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