package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

var checkLoopCollect = false

func ShuffleSites(sites []Site) []Site {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(sites), func(i, j int) { sites[i], sites[j] = sites[j], sites[i] })
	return sites
}

func (m *Instance) CountWorkingTasks() int {
	rows, _ := m.db.Query("SELECT COUNT(*) as count FROM `tasks` WHERE `timeout` IS NOT NULL AND `parser` IS NOT NULL")
	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Println("MysqlTasks.CountWorkingTasks.HasError", err)
		}
	}
	return count
}

func (m *Instance) GetTaskByKeyword(k string) Task {
	var result Task
	sqlQuery := "SELECT * FROM `tasks` WHERE `keyword` = ? LIMIT 1"

	err := m.db.Get(&result, sqlQuery, k)
	if err != nil {
		//log.Println("MysqlDb.GetTaskByKeyword.HasError", err)
	}

	return result
}

func (m *Instance) GetCountTasks(params map[string]interface{}) int {
	var count int
	rows, err := m.db.Query("SELECT COUNT(*) as `count` FROM `tasks`")
	fmt.Println(rows)
	if err == nil {
		for rows.Next() {
			err := rows.Scan(&count)
			if err != nil {
				log.Println("MysqlDb.GetCountTasks.HasError", err)
			}
		}
	} else {
		fmt.Println("ERR.MysqlTasks.GetCountTasks", err)
	}
	return count
}

func (m *Instance) GetTasks(params map[string]interface{}) []Task {
	var results []Task

	//fmt.Println(task.ParseDate.String)
	sqlQuery := "SELECT * FROM `tasks`"

	if len(params) > 0{
		if params["isStat"] != 0 {
			sqlQuery = "SELECT id, site_id, cat_id, status FROM `tasks`"
		}
	}
	sqlQuery = sqlQuery + " ORDER BY `id`"

	if len(params) > 0{
		if params["limit"] != 0 {
			if params["offset"] != 0 {
				sqlQuery = sqlQuery + "LIMIT " + strconv.Itoa(params["offset"].(int)) + ", " + strconv.Itoa(params["limit"].(int))
			}else{
				sqlQuery = sqlQuery + " LIMIT " + strconv.Itoa(params["limit"].(int))
			}
		}
	}

	err := m.db.Select(&results, sqlQuery)
	if err != nil {
		log.Println("MysqlDb.GetTasks.HasError", err)
	}

	return results
}

func (m *Instance) UpdateTask(data map[string]interface{}, id int) (sql.Result, error) {
	sqlQuery := "UPDATE `tasks` SET "

	if len(data) > 0 {
		updateQuery := ""
		i := 0
		for k, v := range data {
			if i > 0 {
				updateQuery += ", "
			}
			updateQuery += "`" + k + "` = "
			if v == "NULL" {
				updateQuery += "NULL"
			}else{
				updateQuery += ":" + k
			}
			//data[k] = UTILS.MysqlRealEscapeString(v.(string))
			i++
		}
		sqlQuery += updateQuery
	}

	sqlQuery += " WHERE `id` = " + strconv.Itoa(id)

	res, err := m.db.NamedExec(sqlQuery, data)

	return res, err
}

func (m *Instance) AddTask(item map[string]interface{}) (sql.Result, error) {
	sqlQuery := "INSERT INTO `tasks` SET "
	sqlQuery += "`site_id` = :site_id, " +
		"`cat_id` = :cat_id, " +
		"`keyword` = :keyword, " +
		"`parent_id` = :parent_id, " +
		"`parser` = NULL, " +
		"`error` = NULL"

	res, err := m.db.NamedExec(sqlQuery, item)

	return res, err
}

func (m *Instance) LoopCollectStats() {
	if ! checkLoopCollect {
		checkLoopCollect = true
		for {
			count := m.GetCountTasks(map[string]interface{}{})
			fmt.Println(count)
			if count < 20000 {
				break
			}
			m.CollectStats()
			time.Sleep(2 * time.Minute)
		}
	}
}


func (m *Instance) CollectStats() map[int64]map[string]interface{} {
	params := make(map[string]interface{})

	limit := 20000
	offset := 0
	stat := map[int64]map[string]interface{}{}

	count := m.GetCountTasks(map[string]interface{}{})
	var parts int

	if limit > count {
		parts = 1
	} else {
		parts = int(math.Ceil(float64(count) / float64(limit)))
	}
	//fmt.Println(parts)
	//parts = 1

	cats := m.GetCats(map[string]interface{}{}, map[string]interface{}{})
	sites := m.GetSites(map[string]interface{}{}, map[string]interface{}{})

	//notCorrectData := make([]interface{}, 0)
	if true {
		for i := 0; i < parts; i++ {
			params["limit"] = limit
			params["offset"] = offset
			params["isStat"] = true
			//tasks := []MysqlTask{}
			tasks := m.GetTasks(params)
			if true {
				for _, task := range tasks {
					SiteId := task.SiteId.Int64
					CatId := task.CatId.Int64
					Status := task.Status.Int32

					var Site Site
					var Cat Cat

					if SiteId == 0 {
						//notCorrectData = append(notCorrectData, question)
						continue
					}

					for _, cat := range cats {
						if CatId == cat.Id.Int64 {
							Cat = cat
						}
					}

					for _, site := range sites {
						if SiteId == site.Id.Int64 {
							Site = site
						}
					}

					if CatId != 0 {
						site := map[string]interface{}{}

						if item, ok := stat[SiteId]; ok {
							site = item
						}

						if Site.Domain.Valid {
							site["domain"] = Site.Domain.String
						}

						if _, ok := site["ready"]; ! ok {
							site["ready"] = 0
						}

						if _, ok := site["error"]; ! ok {
							site["error"] = 0
						}

						if _, ok := site["total"]; ! ok {
							site["total"] = 0
						}

						cetegors := map[int64]interface{}{}
						cat := map[string]interface{}{}

						_, ok := site["cats"]
						if ok && len(site["cats"].(map[int64]interface{})) > 0 {
							cetegors = site["cats"].(map[int64]interface{})

							_, ok := cetegors[CatId]
							if ok && len(cetegors[CatId].(map[string]interface{})) > 0 {
								cat = cetegors[CatId].(map[string]interface{})
							}
						}

						cat["title"] = Cat.Title.String

						if _, ok := cat["ready"]; ! ok {
							cat["ready"] = 0
						}

						if _, ok := cat["error"]; ! ok {
							cat["error"] = 0
						}

						if _, ok := cat["total"]; ! ok {
							cat["total"] = 0
						}

						if Status == 2 {
							site["error"] = site["error"].(int) + 1
							cat["error"] = cat["error"].(int) + 1
						} else if Status == 1 {
							site["ready"] = site["ready"].(int) + 1
							cat["ready"] = cat["ready"].(int) + 1
						}

						site["total"] = site["total"].(int) + 1
						cat["total"] = cat["total"].(int) + 1

						cetegors[CatId] = cat
						site["cats"] = cetegors

						stat[SiteId] = site
					}
				}

				if count < offset {
					offset = offset - (offset - count)
				} else {
					offset = offset + limit
				}
				time.Sleep(time.Second)
			}
		}

		if count > 20000 {
			for k, v := range stat {
				info, err := json.Marshal(v)
				if err != nil {
					log.Println("MysqlDb.CollectStats.HasError", err)
				}
				item := map[string]interface{}{
					"info": info,
				}
				res, err := m.UpdateSite(item, int(k))
				fmt.Println(res)
				if err != nil {
					log.Println("MysqlDb.CollectStats.2.HasError", err)
				}
			}
		}
	}
	return stat
}