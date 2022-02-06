package exid

import (
	"fmt"
	"time"

	"github.com/ergoapi/util/exhash"
)

func GenUID(username string) string {
	return exhash.MD5(username + fmt.Sprint(time.Now().UnixNano()))
}
