package broadcast

import (
	"net"
	"time"
)

func CheckServer(ip string) bool {
	conn, err := net.DialTimeout("tcp", ip, 2*time.Second)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
