package mysql

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type FreeTask struct {
	Id int `db:"id" json:"id"`
	SiteId int `db:"site_id" json:"site_id"`
	CatId int `db:"cat_id" json:"cat_id"`
	Keyword string `db:"keyword" json:"keyword"`
	Cat string `db:"cat" json:"cat"`
	TryCount int `db:"try_count" json:"try_count"`
	Log []string `db:"log" json:"log"`

	Language string `db:"language" json:"language"`
	Theme string `db:"theme" json:"theme"`
	Domain string `db:"domain" json:"domain"`
	Login string `db:"login" json:"login"`
	Password string `db:"password" json:"password"`
	From int `db:"from" json:"from"`
	To int `db:"to" json:"to"`
	QstsLimit int `db:"qsts_limit" json:"qsts_limit"`
	Linking int `db:"linking" json:"linking"`
	Header int `db:"header" json:"header"`
	SubHeaders int `db:"subheaders" json:"subheaders"`
	ParseDates int `db:"parse_dates" json:"parse_dates"`
	ParseDoubles int `db:"parse_doubles" json:"parse_doubles"`
	PubImage int `db:"pub_image" json:"pub_image"`
	VideoStep int `db:"video_step" json:"video_step"`
	QaCountFrom int `db:"qa_count_from" json:"qa_count_from"`
	QaCountTo int `db:"qa_count_to" json:"qa_count_to"`
	ParseFast int `db:"parse_fast" json:"parse_fast"`
	ParseSearch4     int              `db:"parse_search4" json:"parse_search4"`
	ImageKey         int              `db:"image_key" json:"image_key"`
	H1               int              `db:"h1" json:"h1"`
	ShOrder          int              `db:"sh_order" json:"sh_order"`
	ShFormat         int              `db:"sh_format" json:"sh_format"`
	ImageSource      int              `db:"image_source" json:"image_source"`
	MoreTags         string           `db:"more_tags" json:"more_tags"`
	SymbMicroMarking string           `db:"symb_micro_marking" json:"symb_micro_marking"`
	CountRows        int              `db:"count_rows" json:"count_rows"`
	SavingAvailable  bool             `db:"saving_available" json:"saving_available"`
	Extra            ConfigExtra `db:"extra" json:"extra"`

	db *Instance
}

func ArrayRand(arr []string) string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int() % len(arr)
	return strings.Trim(arr[n], " ")
}

func (t *FreeTask) MergeTask(task Task) {
	fmt.Println(task.Id)
	t.Id = int(task.Id.Int64)
	t.Keyword = task.Keyword.String
	t.CatId = int(task.CatId.Int64)
	t.Cat = task.Cat.String
	t.TryCount = int(task.TryCount.Int32)
}

func (t *FreeTask) MergeSite(site Site){
	t.SiteId = int(site.Id.Int64)
	t.Language = site.Language.String
	t.Theme = site.Theme.String
	t.Domain = site.Domain.String
	t.Login = site.Login.String
	t.Password = site.Password.String
	t.From = int(site.From.Int64)
	t.To = int(site.To.Int64)
	t.QstsLimit = int(site.QstsLimit.Int64)
	t.Linking = int(site.Linking.Int64)
	t.Header = int(site.Header.Int64)
	t.SubHeaders = int(site.SubHeaders.Int64)
	t.ParseDates = int(site.ParseDates.Int64)
	t.ParseDoubles = int(site.ParseDoubles.Int64)
	t.PubImage = int(site.PubImage.Int64)
	t.VideoStep = int(site.VideoStep.Int64)
	t.QaCountFrom = int(site.QaCountFrom.Int32)
	t.QaCountTo = int(site.QaCountTo.Int32)
	t.ParseFast = int(site.ParseFast.Int32)
	t.ParseSearch4 = int(site.ParseSearch4.Int32)
	t.ImageKey = int(site.ImageKey.Int64)
	t.H1 = int(site.H1.Int32)
	t.ShOrder = int(site.ShOrder.Int32)
	t.ShFormat = int(site.ShFormat.Int32)
	t.ImageSource = int(site.ImageSource.Int64)
	t.CountRows = int(site.CountRows.Int64)
	t.MoreTags = site.MoreTags.String
	t.SymbMicroMarking = site.SymbMicroMarking.String
	t.Extra = ConfigExtra{}

	var extra map[string]interface{}
	_ = json.Unmarshal([]byte(site.Extra.String), &extra)
	if v, ok := extra["adjacent_keys"] ; ok {
		t.Extra.AdjacentKeys = v.(bool)
	}
	if v, ok := extra["count_streams"] ; ok {
		t.Extra.CountStreams = v.(int)
	}
	if v, ok := extra["limit_streams"] ; ok {
		t.Extra.LimitStreams = v.(int)
	}
	if v, ok := extra["cmd_streams"] ; ok {
		t.Extra.CmdStreams = v.(string)
	}
}

func (t *FreeTask) SetFinished(status int, errorMsg string) {
	now := time.Now()
	formattedDate := now.Format("2006-01-02 15:04:05")

	lastLog := ""
	if len(t.Log) > 0 {
		lastLog = t.Log[len(t.Log)-1]
	}

	data := map[string]interface{}{}
	data["status"] = strconv.Itoa(status)
	data["log"] = strings.Join(t.Log, "\n")
	data["log_last"] = lastLog
	data["error"] = errorMsg
	data["stream"] = "NULL"
	data["timeout"] = "NULL"
	data["parse_date"] = formattedDate

	_, err := t.db.UpdateTask(data, t.Id)
	if err != nil {
		log.Println("MysqlFreeTask.SetFinished.HasError", err)
	}
}

func (t *FreeTask) FreeTask() {
	lastLog := ""
	if len(t.Log) > 0 {
		lastLog = t.Log[len(t.Log)-1]
	}

	if t.TryCount > 0 {
		t.TryCount -= 1
	}

	data := map[string]interface{}{}
	data["log"] = strings.Join(t.Log, "\n")
	data["log_last"] = lastLog
	data["stream"] = "NULL"
	data["timeout"] = "NULL"
	data["try_count"] = t.TryCount

	_, err := t.db.UpdateTask(data, t.Id)
	if err != nil {
		log.Println("MysqlFreeTask.FreeTask.HasError", err)
	}
}

func (t *FreeTask) SetTimeout(parser int) {
	now := time.Now().Local().Add(time.Minute * time.Duration(5))
	formattedDate := now.Format("2006-01-02 15:04:05")

	lastLog := ""
	if len(t.Log) > 0 {
		lastLog = t.Log[len(t.Log)-1]
	}

	data := map[string]interface{}{}
	data["log"] = strings.Join(t.Log, "\n")
	data["log_last"] = lastLog
	data["stream"] = strconv.Itoa(parser)
	data["timeout"] = formattedDate

	_, err := t.db.UpdateTask(data, t.Id)
	if err != nil {
		log.Println("MysqlFreeTask.SetTimeout.HasError", err)
	}
}

func (t *FreeTask) SetError(error string) {
	if error == "" {
		return
	}
	now := time.Now().Local().Add(time.Minute * time.Duration(5))
	formattedDate := now.Format("2006-01-02 15:04:05")
	t.SetLog(error)

	data := map[string]interface{}{}
	data["log"] = strings.Join(t.Log, "\n")
	data["log_last"] = error
	data["error"] = error
	data["status"] = 2
	data["stream"] = ""
	data["timeout"] = "NULL"
	data["parse_date"] = formattedDate

	_, err := t.db.UpdateTask(data, t.Id)
	if err != nil {
		log.Println("MysqlFreeTask.SetError.HasError", err)
	}
}

func (t *FreeTask) SetLog(text string) {
	if text == "" {
		return
	}

	timePoint := time.Now()
	text = timePoint.Format("2006-01-02 15:04:05") + " #" + strconv.Itoa(t.Id) + ": " + text
	fmt.Println(text)
	t.Log = append(t.Log, text)
	t.SaveLog()
}

func (t *FreeTask) SaveLog() {
	data := map[string]interface{}{}
	data["log"] = strings.Join(t.Log, "\n")
	data["log_last"] = t.Log[len(t.Log) - 1]

	_, err := t.db.UpdateTask(data, t.Id)
	if err != nil {
		log.Println("MysqlFreeTask.SaveLog.HasError", err)
	}
}

func (t *FreeTask) GetRandDomain() string {
	domains := t.Domain
	if domains != "" && domains != "[]" {
		var arr []string
		err := json.Unmarshal([]byte(domains), &arr)
		if err != nil {
			log.Println("MysqlFreeTask.GetRandDomain.HasError", err)
		}else {
			return ArrayRand(arr)
		}
	}
	return ""
}

func (t *FreeTask) GetRandSymb() string {
	symbs := t.SymbMicroMarking
	if symbs != "" && symbs != "[]" {
		var arr []string
		err := json.Unmarshal([]byte(symbs), &arr)
		if err != nil {
			log.Println("MysqlFreeTask.GetRandSymb.HasError", err)
		}else {
			return ArrayRand(arr)
		}
	}
	return ""
}

func (t *FreeTask) GetRandTag() string {
	moreTags := t.MoreTags
	if moreTags != "" && moreTags != "[]" {
		var arr []string
		err := json.Unmarshal([]byte(moreTags), &arr)
		if err != nil {
			log.Println("MysqlFreeTask.GetRandTag.HasError", err)
		}else {
			return ArrayRand(arr)
		}
	}
	return ""
}

func (m *Instance) GetFreeTask(id int) FreeTask {
	var freeTask FreeTask
	var sites []Site

	sqlCount := "SELECT COUNT(*) FROM `tasks` WHERE `site_id` = s.id"
	sqlSelectSite := "s.id, s.extra, s.qsts_limit, s.more_tags, s.symb_micro_marking, s.language, s.theme, s.from, s.to, s.qa_count_from, s.qa_count_to, s.login, s.password, s.domain, s.h1, s.sh_format, s.sh_order, s.video_step, s.linking, s.parse_dates, s.parse_doubles, s.parse_fast, s.parse_search4, s.image_source, s.image_key, s.pub_image, (" + sqlCount + ") as count_rows"
	sqlSite := "SELECT " + sqlSelectSite + " FROM sites s"

	err := m.db.Select(&sites, sqlSite)
	if err != nil{
		log.Println("MysqlDb.GetFreeTask.HasError", err)
	}
	sites = ShuffleSites(sites)

	var site Site
	var siteId int64
	var siteCountTasks int64
	for _, item := range sites {
		if item.CountRows.Int64 > 0 {
			site = item
			siteId = item.Id.Int64
			siteCountTasks = item.CountRows.Int64
			break
		}
	}


	if siteId > 0 {
		freeTask.MergeSite(site)

		t := time.Now()
		now := t.Format("2006-01-02 15:04:05")

		randomOffset := int(siteCountTasks)
		if randomOffset < 1 {
			return freeTask
		}
		randomOffset = rand.Intn(randomOffset)

		var sqlQuery string
		if id > 0 {
			sqlQuery = "SELECT t.id, t.keyword, t.try_count, c.title AS cat, t.site_id, t.cat_id FROM tasks t"
			sqlQuery += " LEFT JOIN cats c ON (c.id = t.cat_id)"
			sqlQuery += " AND t.id = " + strconv.Itoa(id)
		}else{
			sqlQuery = "SELECT t.id, t.keyword, t.try_count, c.title AS cat, t.site_id, t.cat_id FROM tasks t"
			sqlQuery += " LEFT JOIN cats c ON (c.id = t.cat_id)"
			sqlQuery += " WHERE t.site_id = "
			sqlQuery += strconv.Itoa(int(siteId))
			sqlQuery += " AND (t.try_count IS NULL OR t.try_count <= 5)"
			sqlQuery += " AND (t.status IS NULL OR t.status = 0) AND (t.timeout is NULL OR t.timeout < '"
			sqlQuery += now
			sqlQuery += "') ORDER BY RAND() LIMIT 1"
		}

		var task Task
		err := m.db.Get(&task, sqlQuery)
		if err != nil{
			log.Println("MysqlDb.GetFreeTask.2.HasError", err)
		}
		freeTask.MergeTask(task)
		freeTask.SavingAvailable = freeTask.QstsLimit > freeTask.CountRows
	}
	freeTask.db = m
	return freeTask
}