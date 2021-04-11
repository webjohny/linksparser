package main

import (
	"fmt"
	"linksparser/config"
	"linksparser/wordpress"
	"log"
	"net/http"
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

		//postWP()
		go func() {
			job := JobHandler{}
			job.IsStart = true
			if job.Browser.Init() {
				//byt, _ := job.Browser.ScreenShot("https://www.investopedia.com/terms/v/virtual-reality.asp")
				//postWP()
				fmt.Println(job.Run(0))
			}
		}()
	} else if MYSQL.CountWorkingTasks() > 0 {
		conf := MYSQL.GetConfig()
		extra := conf.GetExtra()
		if extra.CountStreams > 0 {
			STREAMS.StartLoop(extra.CountStreams, extra.LimitStreams, extra.CmdStreams)
		}
	}


	routes := Routes{}
	routes.Run()

	time.Sleep(time.Minute)
}

func postWP() {
	//task := MYSQL.GetFreeTask(1097080)
	client := wordpress.NewClient(&wordpress.Options{
		BaseAPIURL: `http://philli.beget.tech/wp-json/wp/v2`, // example: `http://192.168.99.100:32777/wp-json/wp/v2`
		//Username:   "proger",
		//Password:   "qwerty12345",
		Username:   "Jekyll1911",
		Password:   "ghjcnjgfhjkm",
	})

	// for eg, to get current user (GET /users/me)
	_, resp, body, _ := client.Users().Me(map[string]int{})
	if resp.StatusCode != http.StatusOK {
		// handle error
	}
	fmt.Println(string(body))


	cats, _, body, _ := client.Categories().List(map[string]string{
		"slug": "qa",
	})
	log.Println(cats)
	log.Fatal(string(body))

	//filee, resp, body, _ := client.Media().Create(&wordpress.MediaUploadOptions{
	//	Filename:    "test-image.jpg",
	//	ContentType: "image/jpeg",
	//	Data:        image,
	//})
	//log.Println(string(body))
	//log.Fatal(filee)

	//slugName := "test-posts-create-2"
	//posts, _, _, _ := client.Posts().List("slug=" + slugName)
	//
	//if posts != nil && len(posts) > 0 {
	//	post := posts[0]
	//	post.Categories = []int{2}
	//	post.Content.Raw = "Hello world!"
	//
	//	_, _, _, _ = client.Posts().Update(2699, &post)
	//}
	//log.Fatal(string(body))
}