package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"linksparser/mysql"
	"linksparser/services"
	"log"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/webjohny/chromedp"
)

type JobHandler struct {
	SearchHtml string
	IsStart bool

	taskId int
	config mysql.Config
	task mysql.FreeTask
	proxy Proxy

	Browser Browser
	ctx context.Context
	isFinished chan bool
}

type WpPost struct {
	Title string
	Content string
	Url string
	AskedBy string
	Text string
	CatId int
	PhotoId int
}

type LinkResult struct {
	Link string
	Image string
	Author string
	Title string
	Description string
	GlobalRank int64
	PageViews int64
	CountryCode string
	CountryName string
}

type TopCountryShare struct {
	Value float64 `json:"Value"`
	Country int32 `json:"Country"`
}

type GlobalRank struct {
	Rank int32 `json:"Rank"`
}

type CountryRank struct {
	Country int32 `json:"Country"`
	Rank int32 `json:"Rank"`
}

type CategoryRank struct {
	Category string `json:"Category"`
	Rank string `json:"Rank"`
}

type TrafficSources struct {
	Social float64 `json:"Social"`
	PaidReferrals float64 `json:"Paid Referrals"`
	Mail float64 `json:"Mail"`
	Referrals float64 `json:"Referrals"`
	Search float64 `json:"Search"`
	Direct float64 `json:"Direct"`
}

type SimilarWebResp struct {
	SiteName string `json:"SiteName"`
	Title string `json:"Title"`
	Category string `json:"Category"`
	CategoryRank CategoryRank `json:"CategoryRank"`
	LargeScreenshot string `json:"LargeScreenshot"`
	EstimatedMonthlyVisits map[string]int `json:"EstimatedMonthlyVisits"`
	Description string `json:"Description"`
	TopCountryShares []TopCountryShare `json:"Description"`
	GlobalRank GlobalRank `json:"GlobalRank"`
	CountryRank CountryRank `json:"CountryRank"`
	IsSmall bool `json:"IsSmall"`
	TrafficSources TrafficSources `json:"TrafficSources"`
}

func (j *JobHandler) Run(parser int) (status bool, msg string) {
	if !j.IsStart {
		go j.Cancel()
		return false, "Задача закрыта"
	}

	fmt.Println("Start task")

	var taskId int

	//var fast QaFast

	// Берём свободную задачу в работу
	var task mysql.FreeTask
	if j.taskId < 1 {
		task = MYSQL.GetFreeTask(0)
	}else{
		task = MYSQL.GetFreeTask(j.taskId)
	}

	if task.Id < 1 {
		go j.Cancel()
		return false, "Свободных задач нет в наличии"
	}
	taskId = task.Id
	task.Domain = task.GetRandDomain()
	//task.SetLog("Задача #" + strconv.Itoa(taskId) + " с запросом (" + task.Keyword + ") взята в работу")

	j.task = task

	if j.CheckFinished() {
		task.FreeTask()
		return false, "Timeout"
	}

	if task.TryCount == 5 {
		task.FreeTask()
		go j.Cancel()
		return false, "Исключён после 5 попыток парсинга"
	}

	j.Browser.Proxy.setTimeout(parser, 5)
	task.SetLog("Подключаем прокси #" + strconv.Itoa(j.Browser.Proxy.Id) + " к браузеру (" + j.Browser.Proxy.LocalIp + ")")

	task.SetTimeout(parser)

	var searchHtml string
	var googleUrl string

	j.config = MYSQL.GetConfig()

	for i := 1; i < 2; i++ {
		if j.CheckFinished() {
			task.SetLog("Задача завершилась преждевременно из-за таймаута")
			return false, "Timeout"
		}

		// Запускаемся
		googleUrl = "https://www.google.com/search?q=" + url.QueryEscape(task.Keyword)
		task.SetLog("Открываем страницу (попытка №" + strconv.Itoa(i) + "): " + googleUrl)

		if j.Browser.ctx != nil {
			if err := chromedp.Run(j.Browser.ctx,
				// Устанавливаем страницу для парсинга
				//chromedp.Sleep(time.Second * 60),
				j.Browser.runWithTimeOut(20, false, chromedp.Tasks{
					chromedp.Navigate(googleUrl),
					chromedp.Sleep(time.Second*time.Duration(rand.Intn(5))),
					chromedp.WaitVisible("body", chromedp.ByQuery),
					// Вытащить html на проверку каптчи
					chromedp.OuterHTML("body", &searchHtml, chromedp.ByQuery),
				}),
			); err != nil {
				log.Println("JobHandler.Run.HasError", err)
				task.FreeTask()
				return false, "Not found page"
			}else if j.Browser.CheckCaptcha(searchHtml) {
				fmt.Println("Присутствует каптча")
				task.FreeTask()
				j.Cancel()
				return false, "Присутствует каптча"
			}
		}else{
			task.FreeTask()
			j.Cancel()
			return false, "Context undefined"
		}
	}

	if j.CheckFinished() {
		fmt.Println("Задача завершилась преждевременно")
		task.FreeTask()
		j.Cancel()
		return false, "Timeout"
	}

	if searchHtml == "" {
		fmt.Println("Контент не подгрузился, задачу закрываем")
		j.Cancel()
		task.SetLog("Контент не подгрузился, задачу закрываем")
		return
	}

	task.SetLog("Контент загружен")

	htmlReader := strings.NewReader(searchHtml)
	body, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		log.Println("JobHandler.SetFastAnswer.HasError", err)
	}

	task.SetLog("Парсинг ссылок из выдачи")

	var linkResults []LinkResult
	body.Find(".hlcw0c").Each(func(i int, hlcw0c *goquery.Selection) {
		hlcw0c.Find(".g").Each(func(y int, g *goquery.Selection) {
			var res LinkResult
			res.Title = g.Find("h3").Text()
			linkSel := g.Find(".yuRUbf").Find("a")
			if linkSel != nil {
				href, _ := linkSel.Attr("href")
				res.Link = href
			}

			linkResults = append(linkResults, res)
		})
	})

	task.SetLog("Обработка похожих запросов")

	var searchesRelated []string
	if body.Find(".k8XOCe").Length() > 0 {
		body.Find(".k8XOCe").Each(func(i int, k8XOCe *goquery.Selection) {
			searchesRelated = append(searchesRelated, k8XOCe.Text())
		})
	}

	task.SetLog("Извлечение информации по ссылкам из API Data.Similarweb.com")

	var results []SimilarWebResp
	for i := 0; i < len(linkResults); i++ {
		res := linkResults[i]
		dsw, err := j.ExtractSimilarWebData(res.Link)
		if err != nil {
			task.SetLog("Ошибка загрузки на ресурсе: " + res.Link)
			continue
		}

		if dsw != nil {
			buf, err := j.Browser.ScreenShot(res.Link)
			if err != nil {
				fmt.Println("ERR.JobHandler.Run.Screenshot", err)
			}
			err = ioutil.WriteFile("./" + strconv.Itoa(task.Id) + "-" + strconv.Itoa(i) + ".jpg", *buf, 0644)
			if err != nil {
				fmt.Println("ERR.JobHandler.Run.Screenshot.2", err)
			}
			results = append(results, *dsw)
		}
		time.Sleep(time.Second * time.Duration(rand.Intn(15)+3))
	}

	var WpPost WpPost
	WpPost.Title = task.Keyword
	WpPost.Content = "Looking for " + WpPost.Title + "? Get direct access to " + WpPost.Title + " through official links provided below."

	fmt.Println(linkResults)
	fmt.Println(searchesRelated)
	task.FreeTask()
	log.Fatal("STOP")

	if j.CheckFinished() {
		task.FreeTask()
		go j.Cancel()
		return false, "Timeout"
	}

	if task.ParseSearch4 < 1 {
		wp := services.Wordpress{}
		wp.Connect(`https://` + task.Domain + `/xmlrpc2.php`, task.Login, task.Password, 1)
		if !wp.CheckConn() {
			task.SetLog("Не получилось подключится к wp xmlrpc (https://" + task.Domain + "/xmlrpc2.php - " + task.Login + " / " + task.Password + ")")
			task.SetError(wp.GetError().Error())
			go j.Cancel()
			return false, "Не получилось подключится к wp xmlrpc (https://" + task.Domain + "/xmlrpc2.php - " + task.Login + " / " + task.Password + ")"
		}

		// Отправляем заметку на сайт
		postId := wp.NewPost(qaTotalPage.Title, qaTotalPage.Content, qaTotalPage.CatId, qaTotalPage.PhotoId)
		var fault bool
		if postId > 0 {
			post := wp.GetPost(postId)
			if post.Id > 0 {
				wp.EditPost(postId, qaTotalPage.Title, qaTotalPage.Content)
			}else{
				fault = true
			}
		}else{
			fault = true
		}

		if fault {
			task.SetLog("Не получилось разместить статью на сайте")
			task.SetError(wp.GetError().Error())
			go j.Cancel()
			return false, "Не получилось разместить статью на сайте"
		}

		task.SetLog("Статья размещена на сайте")
	}else{
		task.SetLog(`Данные сохранены в "Search for"`)
	}
	task.SetFinished(1, "")
	fmt.Println(taskId)
	go j.Cancel()
	return true, "Задача #" + strconv.Itoa(taskId) + " была успешно выполнена"
}

func (j *JobHandler) ExtractSimilarWebData(link string) (*SimilarWebResp, error) {
	var jsonResp string

	dswUrl := "https://data.similarweb.com/api/v1/data?domain=" + url.QueryEscape(link)
	if err := chromedp.Run(j.Browser.ctx,
		// Устанавливаем страницу для парсинга
		//chromedp.Sleep(time.Second * 60),
		j.Browser.runWithTimeOut(20, false, chromedp.Tasks{
			chromedp.Navigate(dswUrl),
			chromedp.WaitVisible("body", chromedp.ByQuery),
			// Вытащить html на проверку каптчи
			chromedp.InnerHTML("body", &jsonResp, chromedp.ByQuery),
		}),
	); err != nil {
		log.Println("ExtractSimilarWebData.Run.HasError", err)
		return nil, err
	}

	var re = regexp.MustCompile(`(?m)^(.*?)\{(.*)\}.*`)
	jsonResp = re.ReplaceAllString(jsonResp, "{$2}")
	if jsonResp == "{}" {
		return nil, nil
	}

	var obj SimilarWebResp
	err := json.Unmarshal([]byte(jsonResp), &obj)
	if err != nil {
		fmt.Println("ERR.ExtractSimilarWebData.3")
		return nil, err
	}

	return &obj, nil
}

func (j *JobHandler) CheckFinished() bool {
	select {
	case <-j.isFinished:
		return true
	default:
		return false
	}
}

func (j *JobHandler) Cancel() {
	if CONF.Env != "local" {
		j.isFinished <- true
	}
}

func (j *JobHandler) Reload() {
	j.Browser.Reload()
}