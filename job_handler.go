package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bxcodec/faker"
	"io/ioutil"
	"linksparser/mysql"
	"linksparser/services"
	"linksparser/tmpl"
	"linksparser/wordpress_xmlrpc"
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

type Engagments struct {
	BounceRate string `json:"BounceRate"`
	Month string `json:"Month"`
	Year string `json:"Year"`
	PagePerVisit string `json:"PagePerVisit"`
	Visits string `json:"Visits"`
	TimeOnSite string `json:"TimeOnSite"`
}

type SimilarWebResp struct {
	SiteName string `json:"SiteName"`
	Description string `json:"Desc"`
	Title string `json:"Title"`
	Category string `json:"Category"`
	CategoryRank CategoryRank `json:"CategoryRank"`
	LargeScreenshot string `json:"LargeScreenshot"`
	EstimatedMonthlyVisits map[string]int `json:"EstimatedMonthlyVisits"`
	TopCountryShares []TopCountryShare `json:"Description"`
	GlobalRank GlobalRank `json:"GlobalRank"`
	CountryRank CountryRank `json:"CountryRank"`
	IsSmall bool `json:"IsSmall"`
	TrafficSources TrafficSources `json:"TrafficSources"`
	Engagments Engagments `json:"Engagments"`
}

func checkHostInArrayLinks(items []*tmpl.LinkResult, link string) bool {
	for _, v := range items {
		existUrl, err := url.Parse(v.Link)
		if err != nil {
			return false
		}
		compareUrl, err := url.Parse(link)
		if err != nil {
			return false
		}
		if existUrl.Host == compareUrl.Host {
			return true
		}
	}
	return false
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
	task.SetLog("Задача #" + strconv.Itoa(taskId) + " с запросом (" + task.Keyword + ") взята в работу")

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
	//task.FreeTask()

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

	task.SetLog(`Подключение к ` + task.Domain)

	wp := wordpress_xmlrpc.Base{}
	wp.Connect(`https://` + task.Domain, task.Login, task.Password, 1)
	if !wp.CheckConn() {
		task.SetLog("Не получилось подключится к wp xmlrpc (https://" + task.Domain + " - " + task.Login + " / " + task.Password + ")")
		task.SetLog("Пробуем http соединение")
		wp.Connect(`http://` + task.Domain, task.Login, task.Password, 1)
		if !wp.CheckConn() {
			task.SetLog("Не получилось подключится к wp xmlrpc")
			if wp.GetError() != nil {
				task.SetError(wp.GetError().Error())
			}
			go j.Cancel()
			return false, "Не получилось подключится к wp xmlrpc"
		}
	}

	task.SetLog("Парсинг ссылок из выдачи")

	var wpPost tmpl.WpPost
	var links []*tmpl.LinkResult
	body.Find(".hlcw0c").Each(func(i int, hlcw0c *goquery.Selection) {
		hlcw0c.Find(".g").Each(func(y int, g *goquery.Selection) {
			var res tmpl.LinkResult
			res.Title = g.Find("h3").Text()
			res.Description = g.Find(".aCOpRe").Find("span").Last().Text()
			linkSel := g.Find(".yuRUbf").Find("a")
			if linkSel != nil {
				href, _ := linkSel.Attr("href")
				if href != "" && !checkHostInArrayLinks(links, href) {
					res.Link = href
					links = append(links, &res)
				}
			}
		})
	})

	task.SetLog("Обработка похожих запросов")

	var searchesRelated []string
	if body.Find(".k8XOCe").Length() > 0 {
		body.Find(".k8XOCe").Each(func(i int, k8XOCe *goquery.Selection) {
			searchesRelated = append(searchesRelated, k8XOCe.Text())
		})
	}

	if len(searchesRelated) > 0 {
		for i := 0; i < len(searchesRelated); i++ {
			item := searchesRelated[i]

			if !MYSQL.GetTaskByKeyword(item).Id.Valid {
				if _, err := MYSQL.AddTask(map[string]interface{}{
					"site_id" : strconv.Itoa(task.SiteId),
					"cat_id" : strconv.Itoa(task.CatId),
					"parent_id" : strconv.Itoa(task.Id),
					"keyword" : item,
					"stream" : "",
					"error" : "",
				}); err != nil {
					log.Println("JobHandler.Run.6.HasError", err)
					task.SetLog("Не добавилась новая задача. (" + err.Error() + ")")
				}
			}
		}
	}

	task.SetLog("Извлечение информации по ссылкам из API Data.Similarweb.com")


	list, err := services.GetCountryList()
	if err != nil {
		fmt.Println(err)
	}

	params := map[string]string{
		"keyword": j.task.Keyword,
		"askedBy": wpPost.AskedBy,
	}
	configExtra := j.config.GetExtra()
	extra := j.task.Extra
	texts := extra.Texts
	if len(texts) < 1 {
		texts = configExtra.Texts
	}
	titles := extra.Titles
	if len(titles) < 1 {
		titles = configExtra.Titles
	}
	answers := extra.Answers
	if len(answers) < 1 {
		answers = configExtra.Answers
	}
	if len(texts) > 0 {
		content := services.ArrayRand(texts)
		wpPost.Content = services.SetTmpl(content, params)
	}
	if len(answers) > 0 {
		answerText := services.ArrayRand(answers)
		wpPost.Text = services.SetTmpl(answerText, params)
	}
	if len(titles) > 0 {
		titleText := services.ArrayRand(titles)
		wpPost.Title = services.SetTmpl(titleText, params)
	}
	wpPost.AskedBy = faker.FirstName() + " " + faker.LastName()

	if len(links) > 0 {
		for i := 0; i < 1; i++ {
		//for i := 0; i < len(links); i++ {
			res := links[i]
			dsw, err := j.ExtractSimilarWebData(res.Link)
			if err != nil {
				j.Reload()
				task.SetLog("Ошибка загрузки на ресурсе: " + res.Link)
				continue
			}
			task.SetLog("Данные получены по ресурсу: " + res.Link)

			if dsw != nil {
				buf, err := j.Browser.ScreenShot(res.Link)
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
					} else {
						res.Src = file.Url
					}
					res.Image = *buf
				}
				task.SetLog("Создался скриншот")
				res.GlobalRank = dsw.GlobalRank.Rank
				pageViews := strings.Split(dsw.Engagments.Visits, ".")
				res.PageViews = pageViews[0]
				//res.Title = strings.Title(dsw.Title)
				res.Description = strings.Title(dsw.Description)
				res.Author = faker.FirstName() + " " + faker.LastName()
				if list != nil && len(list.Country) > 0 {
					for _, country := range list.Country {
						iso, _ := strconv.Atoi(country.Iso)
						if int(dsw.CountryRank.Country) == iso {
							res.CountryName = country.English
							res.CountryCode = country.Iso
							res.CountryImg = strings.ToLower(country.Alpha2) + ".png"
						}
					}
				}
				wpPost.Links = append(wpPost.Links, res)
			}

			num := rand.Intn(10)+3
			task.SetLog("Данные собраны. Задержка " + strconv.Itoa(num) + " сек.")
			time.Sleep(time.Second * time.Duration(num))
		}
	}

	rendered := tmpl.CreateWpPostTmpl(wpPost)

	//f, _ := os.Create("data.txt")
	//
	//defer f.Close()
	//
	//f.WriteString(rendered)
	//
	//fmt.Println("done")
	//return

	jsLinks, _ := json.Marshal(wpPost.Links)

	_, err =  MYSQL.AddResult(map[string]interface{}{
		"domain": task.Domain,
		"site_id": task.SiteId,
		"cat_id": task.CatId,
		"task_id": task.Id,
		"keyword": wpPost.Title,
		"author": wpPost.AskedBy,
		"links": string(jsLinks),
		"text": wpPost.Text,
		"content": wpPost.Content,
	})
	if err != nil {
		fmt.Println("ERR.JobHandler.Run.AddResult", err)
	}
	task.SetLog("Добавлен результат в базу данных")

	// Отправляем заметку на сайт
	postId := wp.NewPost(wpPost.Title, rendered, wpPost.CatId, 12)
	var fault bool
	if postId > 0 {
		post := wp.GetPost(postId)
		if post.Id > 0 {
			wp.EditPost(postId, wpPost.Title, rendered)
		}else{
			fault = true
		}
	}else{
		fault = true
	}
	
	//fmt.Println(wpPost.Links)
	//fmt.Println(searchesRelated)
	//task.FreeTask()
	//log.Fatal("STOP")

	if j.CheckFinished() {
		task.FreeTask()
		go j.Cancel()
		return false, "Timeout"
	}


	if fault {
		task.SetLog("Не получилось разместить статью на сайте")
		task.SetError(wp.GetError().Error())
		go j.Cancel()
		return false, "Не получилось разместить статью на сайте"
	}

	task.SetLog("Статья размещена на сайте")

	task.SetFinished(1, "")
	fmt.Println(taskId)
	go j.Cancel()
	return true, "Задача #" + strconv.Itoa(taskId) + " была успешно выполнена"
}

func (j *JobHandler) ExtractSimilarWebData(link string) (*SimilarWebResp, error) {
	var jsonResp string

	//dswUrl := link
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
	jsonResp = strings.Replace(jsonResp, `"Description"`, `"Desc"`, 1)
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