package main

import (
	"fmt"
	"io/ioutil"
	"linksparser/config"
	"linksparser/wordpress"
	"log"
	"net/http"
	"net/url"
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

	//link := "http://google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=&cad=rja&uact=8&ved=2ahUKEwjhqbmWtITwAhUUtqQKHekGDdEQFjANegQIIhAD&url=https%3A%2F%2Fwww.slideshare.net%2Fsaishanesarikar%2Fvirtual-reality-ppt-80531390&usg=AOvVaw15fAgjHzHd6bDTqlRAEgOo"
	//parsedUrl, _ := url.Parse(link)
	//if parsedUrl != nil {
	//	originUrl := parsedUrl.Query()["url"]
	//	if len(originUrl) > 0 {
	//
	//	}
	//}
	//log.Fatal()
	testGetReq()

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

func testGetReq(){
	//url_proxy, _ := url_i.Parse("phillip:I2n9BeJ@45.151.68.227:45785")

	//creating the proxyURL
	proxyStr := "http://phillip:I2n9BeJ@5.188.52.65:45785"
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Println(err)
	}

	//creating the URL to be loaded through the proxy
	urlStr := "https://www.google.com/search?q=virtual+reality"
	//urlStr := "http://httpbin.org/get"
	url, err := url.Parse(urlStr)
	if err != nil {
		log.Println(err)
	}

	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
	}

	//generating the HTTP GET request
	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		log.Println(err)
	}

	request.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	request.Header.Set("connection", "keep-alive")
	request.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")
	request.Header.Set("accept-encoding", "gzip, deflate, br")
	request.Header.Set("accept-language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	request.Header.Set("cache-control", "no-cache")
	request.Header.Set("cookie", "CGIC=IocBdGV4dC9odG1sLGFwcGxpY2F0aW9uL3hodG1sK3htbCxhcHBsaWNhdGlvbi94bWw7cT0wLjksaW1hZ2UvYXZpZixpbWFnZS93ZWJwLGltYWdlL2FwbmcsKi8qO3E9MC44LGFwcGxpY2F0aW9uL3NpZ25lZC1leGNoYW5nZTt2PWIzO3E9MC45; CONSENT=YES+UA.ru+20161009-18-0; ANID=AHWqTUkvwxouUV4f1TaH9JR9YsMu4IeIBkeAycj-NbJAafTK-AoTzFFUxU-pNLle; SEARCH_SAMESITE=CgQI-5EB; HSID=AM7F52QcRUd6Vwv23; SSID=ASeSpmCKqkdIYLooz; APISID=evpD-xcogtiDoUIF/AkN184Apyfp8Zge-k; SAPISID=ouhzu1bNBIOH01Ln/A-maxl0J6Ra4aPG7D; __Secure-3PAPISID=ouhzu1bNBIOH01Ln/A-maxl0J6Ra4aPG7D; OTZ=5943417_44_48_123900_44_436380; SID=8wfD0Wdb0ajAn5iVvpdE8cdHKsYV2I-WqYr_Jh0BkpSBQjJeNjADuABim4ct8Pvflll2Og.; __Secure-3PSID=8wfD0Wdb0ajAn5iVvpdE8cdHKsYV2I-WqYr_Jh0BkpSBQjJe4Q_qjLWwn2ZQ5fJfQaLDAQ.; 1P_JAR=2021-05-04-12; NID=214=l4gZSJpwlRlodLkNKsKdxtAJaMe65OYcxpCNb1Yzf0AqSkEsfEBZML-xpHavadbDjakvA9f8dc8FX5FlMHXB4DXys4KaR84OAwNiWTSh0-em6EOUX9u54xKpkYDfZQPP8P8lFwe9sDOtgeaAMwbuOpgPIhZ1Fdko2Svrqse_fgZxHIf8B1Zc9fyOy1L3Jr6-pOOSLl6AUq_PGjMppTntCkV-hFEvBWXY0OPeqg_xODd3hykk2W1TTJjzGbwJMA4fVEXCKKeSgKkYoEJ0DfAP901haH9rj7QGlRF6ddDZmk-pLVJBxj_dfKUfChtb_Hm_06wLsWuqhpKYykxlBIDqMbev6z7BpqLyr1WMM3Fhjehyv7UlMSQkbQeMp1t_9uc; SIDCC=AJi4QfEeTgtjUP6I38tgP4YDMpDpS6PflJytIdzBVPCQ2IbgRksy8GKE4qsifauy1jj0t-3n-MM; __Secure-3PSIDCC=AJi4QfF-k7TRzKo3ht1Pf7S6Yqg2wvkkxzKE5MttdIdVN4EE2jBnuVJLZXiffibyJHIw8WZrgJuL")
	request.Header.Set("upgrade-insecure-requests", "1")
	request.Header.Set("pragma", "no-cache")
	request.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"`)
	request.Header.Set("sec-ch-ua-mobile", "?0")
	request.Header.Set("sec-fetch-dest", "document")
	request.Header.Set("sec-fetch-mode", "navigate")
	request.Header.Set("sec-fetch-site", "same-origin")
	request.Header.Set("sec-fetch-user", "?1")
//	sec-ch-ua-mobile: ?0
//	sec-fetch-dest: document
//	sec-fetch-mode: navigate
//	sec-fetch-site: same-origin
//	sec-fetch-user: ?1
	//'Accept': '*/*', 'Connection': 'keep-alive', 'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.75 Safari/537.36', 'Accept-Encoding': 'gzip, deflate, br', 'Accept-Language': 'en-US;q=0.5,en;q=0.3', 'Cache-Control': 'max-age=0', 'DNT': '1', 'Upgrade-Insecure-Requests': '1'

	//calling the URL
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	//getting the response
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	//printing the response
	log.Fatal(string(data))
	//fmt.Println(client.Get("https://www.google.com/search?q=virtual+reality"))
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