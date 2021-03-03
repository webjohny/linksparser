package mysql

import (
	"database/sql"
	"encoding/json"
	"time"
)

func (m *Instance) GetResultByQAndA(q string, a string) Result {
	var result Result
	sqlQuery := "SELECT * FROM `results` WHERE `q` = ? AND `a` = ? LIMIT 1"

	err := m.db.Get(&result, sqlQuery, q, a)
	if err != nil {
		//log.Println("MysqlDb.GetResultByQAndA.HasError", err)
	}

	return result
}

func (m *Instance) AddResult(item map[string]interface{}) (sql.Result, error) {
	t := time.Now()
	now := t.Format("2006-01-02 15:04:05")

	if _, ok := item["created_at"]; !ok {
		item["created_at"] = now
	}

	item["links"], _ = json.Marshal(item["links"])

	sqlQuery := "INSERT INTO `results` SET "
	sqlQuery += "`keyword` = :keyword, " +
		"`links` = :links, " +
		"`site_id` = :site_id, " +
		"`cat_id` = :cat_id, " +
		"`domain` = :domain, " +
		"`text` = :text, " +
		"`content` = :content, " +
		"`author` = :author, " +
		"`created_at` = :created_at"

	res, err := m.db.NamedExec(sqlQuery, item)

	return res, err
}
