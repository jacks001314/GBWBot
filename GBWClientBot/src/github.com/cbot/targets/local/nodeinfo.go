package local

import (
	"os"
	"os/user"
	"runtime"
)

type NodeInfo struct {
	IP string

	OutIP string

	Mac string

	OS string

	Arch string

	User string

	Hostname string
}

func GetNodeInfo() *NodeInfo {

	ip := GetWorkingIP()
	outIP := GetOutIP()
	mac := GetMacByIP(ip.String())

	u, _ := user.Current()
	hostname, _ := os.Hostname()

	return &NodeInfo{
		IP:       ip.String(),
		OutIP:    outIP,
		Mac:      mac,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		User:     u.Username,
		Hostname: hostname,
	}

}
