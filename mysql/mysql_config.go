package mysql

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type ConfigExtra struct {
	AdjacentKeys bool `json:"adjacent_keys"`
	CountStreams int `json:"count_streams"`
	LimitStreams int `json:"limit_streams"`
	CmdStreams string `json:"cmd_streams"`
}

func (m *Instance) GetConfig() Config {
	var result Config
	sqlQuery := "SELECT * FROM `config` LIMIT 1"

	err := m.db.Get(&result, sqlQuery)
	if err != nil {
		log.Println("MysqlConfig.GetConfig.HasError", err)
	}

	return result
}

func (m *Instance) SetExtra(extra ConfigExtra) error {
	extraJson, err := json.Marshal(extra)
	if err != nil {
		log.Println("MysqlConfig.SetExtra.HasError", err)
		return err
	}
	sqlQuery := "UPDATE `config` SET `extra` = :extra"
	data := map[string]interface{}{
		"extra": extraJson,
	}

	res, err := m.db.NamedExec(sqlQuery, data)
	fmt.Println(res)

	if err != nil {
		log.Println("MysqlConfig.SetExtra.HasError.2", err)
		return err
	}

	return nil
}

func (c *Config) GetVariants() []string {
	var result []string
	if c.Variants.Valid {
		result = strings.Split(c.Variants.String, ";")
	}
	return result
}

func (c *Config) GetExtra() ConfigExtra {
	Extra := ConfigExtra{}

	var extra map[string]interface{}
	_ = json.Unmarshal([]byte(c.Extra.String), &extra)
	if v, ok := extra["adjacent_keys"] ; ok {
		Extra.AdjacentKeys = v.(bool)
	}
	if v, ok := extra["count_streams"] ; ok {
		Extra.CountStreams = int(v.(float64))
	}
	if v, ok := extra["limit_streams"] ; ok {
		Extra.LimitStreams = int(v.(float64))
	}
	if v, ok := extra["cmd_streams"] ; ok {
		Extra.CmdStreams = v.(string)
	}

	return Extra
}