package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"os/exec"
	"reflect"
	"time"
)

type Instance struct {
	db *sqlx.DB
}

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

func CreateConnection(host string, db string, login string, pass string) Instance {
	conn, err := sqlx.Connect("mysql", login + ":" + pass + "@tcp(" + host + ")/" + db)
	if err != nil {
		panic(err)
	}
	conn.SetMaxIdleConns(20)
	conn.SetConnMaxLifetime(time.Minute * 2)
	conn.SetMaxOpenConns(100)
	fmt.Println("Connected to MysqlDB!")

	return Instance{conn }
}

func (m *Instance) Disconnect() {
	err := m.db.Close()

	if err != nil {
		panic(err)
	}
	fmt.Println("Connection to MySQL closed.")
}

func (m *Instance) Restart() {
	cmd := exec.Command("service", "mysql restart")
	log.Printf("Mysql restarting and waiting for it to finish...")
	err := cmd.Run()
	log.Printf("Command finished with error: %v.HasError", err)
	time.Sleep(time.Second * 5)
}

