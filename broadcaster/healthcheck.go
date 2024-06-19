package broadcast

import (
	"net"
	logger "sidecarauth/utility"
	"time"
)

func CheckServer(ip string) bool {
	logger.LogF("Healthchecker File", ip)

	conn, err := net.DialTimeout("tcp", ip, 2*time.Second)
	if err != nil {
		logger.LogF("Healthchecker Error", err)

		return true
	}
	_ = conn.Close()
	return true
}
