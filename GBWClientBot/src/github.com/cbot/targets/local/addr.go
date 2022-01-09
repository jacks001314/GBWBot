package local

import (
	"fmt"
	"github.com/cbot/proto/http"
	"net"
	"strings"
)


type NetInterface struct {

	name string
	ip net.IP
	mask net.IPMask
	mac string
	gw string
}

func (n *NetInterface) Name() string {

	return n.name
}

func (n *NetInterface) IP4() string {

	return n.ip.To4().String()
}

func (n *NetInterface) IP6() string {

	return n.ip.To16().String()
}


func (n *NetInterface) Mask() string {

	return n.mask.String()
}

func (n *NetInterface) NetWorkRange() string {

	prefixLen,_ := n.mask.Size()

	return fmt.Sprintf("%s/%d",n.ip.Mask(n.mask).String(),prefixLen)
}

func (n *NetInterface) Mac() string {

	return n.mac
}

func (n *NetInterface) GW() string {

	return n.gw
}

func (n *NetInterface) String() string {

	return fmt.Sprintf("{name:%s,ip:%s,mask:%s,mac:%s,gw:%s}",
		n.name,n.ip,n.mask,n.mac,n.gw)
}

func Addrs(isv4 bool) []*NetInterface {


	interfaces := make([]*NetInterface,0)

	ifaces,err := net.Interfaces()

	if err!=nil {
		return interfaces
	}

	for _,iface := range ifaces {

		if iface.Flags&net.FlagUp == 0 ||iface.Flags&net.FlagLoopback!=0 {

			continue
		}

		addrs,err := iface.Addrs()
		if err!=nil {
			return interfaces
		}

		for _,addr := range addrs {

			ip,mask := getIpFromAddr(addr)

			if ip == nil {
				continue
			}
			prefixlen,_:=mask.Size()

			if isv4 && prefixlen>32 {
				continue
			}
			interfaces = append(interfaces,makeInterface(iface,ip,mask))
		}
	}

	return interfaces
}

func makeInterface(p net.Interface,ip net.IP,mask net.IPMask) (*NetInterface){

	return &NetInterface{
		name: p.Name,
		ip:   ip,
		mask: mask,
		mac:  p.HardwareAddr.String(),
		gw:   "",
	}
}

func getIpFromAddr(addr net.Addr) (net.IP,net.IPMask) {

	var ip net.IP
	var mask net.IPMask

	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
		mask = v.Mask

	case *net.IPAddr:

		ip = v.IP
		mask = ip.DefaultMask()
	}

	if ip == nil || ip.IsLoopback() {
		return nil,nil
	}

	return ip,mask
}


// Get preferred outbound ip of this machine
func GetWorkingIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func GetWorkingIPRange(isv4 bool) string {

	ip := GetWorkingIP()

	if ip == nil {

		return ""
	}

	addrs := Addrs(isv4)

	for _,addr := range addrs {

		if strings.EqualFold(addr.ip.String(),ip.String()){

			return addr.NetWorkRange()
		}
	}

	return ""
}

/*get out ip*/
func GetOutIP() string {

	client := http.NewHttpClient("ip.sb",80,false,10000)

	req := http.NewHttpRequest("GET","/").AddHeader("User-Agent"," curl/7.61.1")

	res,err:= client.Send(req)

	if err!=nil {
		fmt.Println(err)
		return ""
	}

	ip,_:= res.GetBodyAsString()
	return strings.TrimSpace(ip)
}

