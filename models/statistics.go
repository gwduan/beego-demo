package models

import (
	"fmt"
	"time"

	"beego-demo/models/myredis"

	"github.com/astaxie/beego"
)

// IncTotalUserCount increase user count in redis.
func IncTotalUserCount(t time.Time) error {
	y, m, d := t.Date()
	yKey := fmt.Sprintf("sys.number.%04d", y)
	mKey := fmt.Sprintf("sys.number.%04d%02d", y, m)
	dKey := fmt.Sprintf("sys.number.%04d%02d%02d", y, m, d)

	conn := myredis.Conn()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HINCRBY", "sys.number.sum", "users", 1)
	conn.Send("HINCRBY", yKey, "new_users", 1)
	conn.Send("HINCRBY", mKey, "new_users", 1)
	conn.Send("HINCRBY", dKey, "new_users", 1)
	r, err := conn.Do("EXEC")
	if err != nil {
		beego.Error("MULTI HINCRBY for new user registeration:", err)
		return err
	}
	beego.Debug("MULTI HINCRBY for new user registeration reply:", r)

	return nil
}
