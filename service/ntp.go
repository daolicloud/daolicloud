package service

import (
    "github.com/beevik/ntp"
    "go.dedis.ch/onet/v3/log"
)

var NTP_SERVERS = [...]string{
    "ntp1.aliyun.com",
}

func getRemoteTime() string {
	for _, server := range NTP_SERVERS {
		time, err := ntp.Time(server)
		if err != nil {
			log.Error(err)
			continue
		}
		return time.String()
	}

	return ""
}
