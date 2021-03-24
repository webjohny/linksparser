package main

import (
	"database/sql"
	"log"
	"strconv"
	"time"
)

type Proxy struct {
	Id int
	Host string
	Port string
	Login string
	Password string
	Agent string

	IsUsing bool

	Ip string
	LocalIp string
	FullIp string
	Log []string
}

func NewProxy() *Proxy {
	proxy := MYSQL.GetFreeProxy()
	if proxy.Id.Valid {
		agent := MYSQL.GetAgent()
		instance := &Proxy{}
		instance.Id = int(proxy.Id.Int64)
		instance.Host = proxy.Host.String
		instance.Port = proxy.Port.String
		instance.Login = proxy.Login.String
		instance.Password = proxy.Password.String
		if proxy.Agent.Valid {
			instance.Agent = proxy.Agent.String
		}else{
			instance.Agent = agent.Sign.String
		}
		instance.LocalIp = instance.Host + ":" + instance.Port
		instance.FullIp = "http://" + instance.Login + ":" + instance.Password + "@" + instance.Host + ":" + instance.Port

		return instance
	}
	return nil
}

func (p *Proxy) setTimeout(parser int, minutes int) sql.Result {
	now := time.Now().Local().Add(time.Minute * time.Duration(minutes))
	formattedDate := now.Format("2006-01-02 15:04:05")

	data := map[string]interface{}{}
	data["stream"] = strconv.Itoa(parser)
	data["timeout"] = formattedDate

	res, err := MYSQL.UpdateProxy(data, p.Id)
	if err != nil {
		log.Println("Proxy.SetTimeout.HasError", err)
	}

	return res
}

func (p *Proxy) freeProxy() {
	now := time.Now().Local().Add(time.Minute * time.Duration(2))
	formattedDate := now.Format("2006-01-02 15:04:05")

	data := map[string]interface{}{}
	data["stream"] = "NULL"
	data["timeout"] = formattedDate

	_, err := MYSQL.UpdateProxy(data, p.Id)
	if err != nil {
		log.Println("Proxy.freeProxy.HasError", err)
	}
	p.Id = 0
	p.Host = ""
}