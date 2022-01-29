package server

import (
	"github.com/sbot/utils/netutils"
	"golang.org/x/net/dns/dnsmessage"
	"net"
	"strconv"
	"strings"
	"sync"
)

type DNSServer struct {

	lock sync.Mutex

	cfg *DNSServerConfig


	subs []*DNSRequestSub
}

type DNSServerConfig struct {


	SubDomain  string  `json:"subDomain"`

	DefaultIP string `json:defIP`

}

type DNSRequest struct {

	Domain string

	AttackType string
	AttackIP string
	TargetIP string
	TargetPort int

	TargetOutIP string

}

type DNSRequestSub struct {

	requests chan *DNSRequest
}


func (d *DNSServer) pub(req *DNSRequest) {

	d.lock.Lock()
	defer d.lock.Unlock()

	for _,sub := range d.subs {

		sub.requests<-req
	}

}

func (d *DNSServer) NewDNSRequestSub() *DNSRequestSub{

	d.lock.Lock()
	defer d.lock.Unlock()

	sub := &DNSRequestSub{requests:make(chan *DNSRequest)}

	d.subs = append(d.subs,sub)

	return sub
}

func (s *DNSRequestSub) Sub() chan *DNSRequest {

	return s.requests
}

func (d *DNSServer) makeDnsRequest(rip string,domain string ,dcr *netutils.DNSDomainCrypt) *DNSRequest {



	return &DNSRequest{
		Domain:      domain,
		AttackType:  dcr.AttackType,
		AttackIP:    dcr.AttackIP,
		TargetIP:    dcr.TargetIP,
		TargetPort:  dcr.TargetPort,
		TargetOutIP: rip,
	}

}

func (d *DNSServer) getResponseIP() [4]byte {

	arr := strings.Split(d.cfg.DefaultIP,".")

	aa,_:= strconv.ParseInt(arr[0],10,32)
	bb,_:= strconv.ParseInt(arr[1],10,32)
	cc,_:= strconv.ParseInt(arr[2],10,32)
	dd,_:= strconv.ParseInt(arr[3],10,32)

	return [4]byte{byte(aa),byte(bb),byte(cc),byte(dd)}
}

func (d *DNSServer) serverDNS(addr *net.UDPAddr,conn *net.UDPConn,msg *dnsmessage.Message) {


	domain := msg.Questions[0].Name.String()

	if domain[len(domain)-1] == '.' {

		domain = domain[:len(domain)-1]

	}

	log.Infof("Reveive a Dns Request from:%v,domain:%s",addr,domain)

	if !strings.HasSuffix(domain,d.cfg.SubDomain){

		log.Errorf("Receive a Dns Request from:%v,but this domain:%s is not a subdomain of dnslog domain:%s ",
			addr,domain,d.cfg.SubDomain)
		return
	}

	indx := strings.Index(domain,".")

	if indx <=0 {

		log.Errorf("Receive a Dns Request from:%v,but this domain:%s format is error ",
			addr,domain)

		return
	}

	dcr,err:= netutils.DeCryptToDNSDomain(domain[:indx])

	if err!=nil  {
		log.Errorf("Receive a Dns Request from:%v,but this domain:%s format is error:%v ",
			addr,domain,err)
		return
	}

	dnsReq := d.makeDnsRequest(addr.IP.String(),domain,dcr)

	if dnsReq == nil {
		log.Errorf("Receive a Dns Request from:%v,but this domain:%s format is error ",
			addr,domain)
		return
	}

	d.pub(dnsReq)

	arecord := dnsmessage.Resource{
		Header: dnsmessage.ResourceHeader{
			Name:  msg.Questions[0].Name,
			Class: dnsmessage.ClassINET,
			TTL:   600,
		},
		Body: &dnsmessage.AResource{
			A: d.getResponseIP(),
		},
	}
	msg.Response = true
	msg.Answers = append(msg.Answers,arecord)

	packed, err := msg.Pack()
	if err != nil {
		log.Errorf("Response DNS pack failed:%v ",
			err)

		return
	}

	if _, err := conn.WriteToUDP(packed, addr); err != nil {

		log.Errorf("Response write failed:%v",err)
	}

}


func (d *DNSServer) Start() error {

	conn,err:= net.ListenUDP("udp",&net.UDPAddr{
		Port: 53,
	})

	if err!=nil {

		log.Errorf("UDP Cannot Listen on port:53 for dnslog server")
		return err
	}

	defer conn.Close()

	for {

		buf := make([]byte,512)

		n,addr,err := conn.ReadFromUDP(buf)

		if err!=nil ||n<=0 {

			log.Errorf("Read DNS Packet error:%v,addr:%s",err,addr)
			continue
		}

		var msg dnsmessage.Message
		if err = msg.Unpack(buf);err!=nil {

			log.Errorf("Parse DNS Packet error:%v,addr:%s",err,addr)
			continue
		}

		if len(msg.Questions)<=0 {

			log.Errorf("Parse DNS Packet error:%v,addr:%s",err,addr)
			continue
		}

		go d.serverDNS(addr,conn,&msg)

	}
}

func NewDNSServer(cfg *DNSServerConfig) *DNSServer {


	return &DNSServer{
		lock: sync.Mutex{},
		cfg:  cfg,
		subs: make([]*DNSRequestSub,0),
	}
}