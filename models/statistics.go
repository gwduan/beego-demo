package models

import (
	"beego-demo/models/myredis"
	"fmt"
	"github.com/gwduan/beego"
	"time"
)

func IncTotalUserCount(t time.Time) error {
	y, m, d := t.Date()
	y_key := fmt.Sprintf("sys.number.%04d", y)
	m_key := fmt.Sprintf("sys.number.%04d%02d", y, m)
	d_key := fmt.Sprintf("sys.number.%04d%02d%02d", y, m, d)

	conn := myredis.Conn()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HINCRBY", "sys.number.sum", "users", 1)
	conn.Send("HINCRBY", y_key, "new_users", 1)
	conn.Send("HINCRBY", m_key, "new_users", 1)
	conn.Send("HINCRBY", d_key, "new_users", 1)
	r, err := conn.Do("EXEC")
	if err != nil {
		beego.Error("MULTI HINCRBY for new user registeration:", err)
		return err
	}
	beego.Debug("MULTI HINCRBY for new user registeration reply:", r)

	return nil
}
