package mysql

import (
	"database/sql"
	"log"
	"time"
)

var userAgents []Uagent

func (m *Instance) GetAgents() []Uagent {
	var agents []Uagent

	t := time.Now()
	now := t.Format("2006-01-02 15:04:05")
	sqlQuery := "SELECT * FROM `user_agents` WHERE (status is NULL OR status = 0) AND (timeout is NULL OR timeout < '" + now + "') ORDER BY RAND() LIMIT 100"

	err := m.db.Select(&agents, sqlQuery)
	if err != nil {
		log.Println("MysqlUAgent.GetAgents.HasError", err)
	}

	return agents
}

func (m *Instance) GetAgent() *Uagent {
	if len(userAgents) < 1 {
		userAgents = m.GetAgents()
	}
	var agent Uagent
	agent, userAgents = userAgents[0], userAgents[1:]

	return &agent
}

func (m *Instance) AddUAgent(sign string) (sql.Result, error) {
	sqlQuery := "INSERT INTO `user_agents` SET "
	sqlQuery += "`sign` = :sign, " +
		"`status` = 0, " +
		"`timeout` = NULL"

	res, err := m.db.NamedExec(sqlQuery, map[string]interface{}{"sign":sign})

	return res, err
}