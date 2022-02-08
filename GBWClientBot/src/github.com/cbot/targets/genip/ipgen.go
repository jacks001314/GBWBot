package genip

import (
	"fmt"
	"github.com/cbot/utils/netutils"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
	"math/big"
)

type IPGen struct {

	TengoObj

	wb *WBList
	cycle *Cycle

	first uint64
	last  uint64
	factor uint64
	modulus uint64
	current uint64

	maxIndex uint64
}


func (ipg *IPGen)initIPGen() {

	var numElts uint64 = ipg.cycle.Order
	var expBegin uint64 = uint64(uint64(ipg.cycle.Offset) % numElts)
	var expEnd uint64 = uint64(uint64(ipg.cycle.Offset) % numElts)

	// Multiprecision variants of everything above
	genM := new(big.Int).SetUint64(ipg.cycle.Generator)
	expBeginM := new(big.Int).SetUint64(expBegin)
	expEndM := new(big.Int).SetUint64(expEnd)
	primeM := new(big.Int).SetUint64(ipg.cycle.Group.prime)

	startM := new(big.Int)
	stopM := new(big.Int)

	startM = startM.Exp(genM,expBeginM,primeM)
	stopM = stopM.Exp(genM,expEndM,primeM)

	ipg.first = startM.Uint64()
	ipg.last = stopM.Uint64()
	ipg.factor = ipg.cycle.Generator
	ipg.modulus = ipg.cycle.Group.prime

	ipg.current = ipg.first

	ipg.roll2Valid()

}

func NewIPGen(whiteListFName string,blackListFName string,whiteListEntries []string,
	blakListEntries []string,ignoreInvalidHosts bool) (*IPGen,error) {

	var ipg IPGen
	var err error
	var max32Int uint64 = 1<<32

	ipg.wb,err = NewWblist(whiteListFName,blackListFName,whiteListEntries,blakListEntries,ignoreInvalidHosts)

	if err!=nil {

		return nil,err
	}

	numAddr := ipg.wb.WBListCountAllowed()

	group := GetGroup(numAddr)

	if group == nil {
		return nil,fmt.Errorf("Cannot get valid group for numAddr:%d",numAddr)
	}

	if numAddr> max32Int {
		ipg.maxIndex= 0xFFFFFFFF
	} else {
		ipg.maxIndex = numAddr
	}

	ipg.cycle = group.MakeCycle()

	ipg.initIPGen()

	return &ipg,nil
}


func (ipg *IPGen) GetCurIP() uint32 {

	return uint32(ipg.wb.LookupIndex(ipg.current - 1))
}

func (ipg *IPGen) getNextElem() uint32 {

	var max32Int uint64 = 1 << 32

	for {

		ipg.current = ipg.current*ipg.factor
		ipg.current = ipg.current%ipg.modulus

		if ipg.current <max32Int {

			break
		}
	}

	return uint32(ipg.current)
}

func (ipg *IPGen) GetNextIP() uint32 {

	var candidate uint64

	if ipg.current == 0 {
		return 0
	}

	for {
		candidate = uint64(ipg.getNextElem())
		if candidate == ipg.last {

			ipg.current = 0
			return 0
		}

		if candidate -1 < ipg.maxIndex {
			return uint32(ipg.wb.LookupIndex(candidate - 1))
		}
	}
}

func (ipg *IPGen) roll2Valid() uint32 {

	if ipg.current - 1 <ipg.maxIndex {
		return uint32(ipg.current)
	}

	return ipg.GetNextIP()
}


func (p *IPGen) IndexGet(index objects.Object)(value objects.Object,err error){

	key,ok := objects.ToString(index)

	if !ok {
		return nil,tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	switch key {

	case "curIP":

		return &TengoCurIP{
			TengoObj: TengoObj{name: "curIP"},
			ipgen:  p,
		}, nil

	case "nextIP":

		return &TengoNextIP{
			TengoObj: TengoObj{name: "nextIP"},
			ipgen:   p,
		}, nil

	}

	return nil,fmt.Errorf("Unknown ipgen method:%s",key)
}

type TengoCurIP struct {

	TengoObj

	ipgen *IPGen
}

func (t *TengoCurIP) Call(args ... objects.Object) (objects.Object,error) {

	ip := t.ipgen.GetCurIP()

	if ip ==0 {

		return objects.FromInterface("")
	}


	return objects.FromInterface(netutils.IPv4StrBig(ip))

}

type TengoNextIP struct {

	TengoObj

	ipgen *IPGen
}

func (t *TengoNextIP) Call(args ... objects.Object) (objects.Object,error) {

	ip := t.ipgen.GetNextIP()

	if ip ==0 {

		return objects.FromInterface("")
	}

	return objects.FromInterface(netutils.IPv4StrBig(ip))
}


func newTengoIPGenFromFile(args ... objects.Object) (objects.Object,error) {

	if len(args) != 2 {

		return nil, tengo.ErrWrongNumArguments
	}

	wlistFname, ok := objects.ToString(args[0])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "wlistFname",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	blistFname,ok := objects.ToString(args[1])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "blistFname",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	return NewIPGen(wlistFname,blistFname,[]string{},[]string{},true)
}


func newTengoIPGenFromArray(args ... objects.Object) (objects.Object,error) {

	if len(args) != 2 {

		return nil, tengo.ErrWrongNumArguments
	}


	wlistArr,ok:= objects.ToInterface(args[0]).([]interface{})

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:    "wlistArray",
			Expected: "[]string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	blistArr,ok:= objects.ToInterface(args[1]).([]interface{})

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "blistArray",
			Expected: "[]string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	wlists := make([]string,0)
	blists := make([]string,0)

	for _,w := range wlistArr {

		wlists = append(wlists,w.(string))
	}

	for _,b := range blistArr {

		blists = append(blists,b.(string))
	}

	return NewIPGen("","",wlists,blists,true)
}


var moduleMap objects.Object = &objects.ImmutableMap{

	Value: map[string]objects.Object{

		"newIPGenFromFile": &objects.UserFunction{
			Name:  "new_ipgen_fromfile",
			Value: newTengoIPGenFromFile,
		},

		"newIPGenFromArray": &objects.UserFunction{
			Name:  "new_ipgen_from_array",
			Value: newTengoIPGenFromArray,
		},

	},
}

func (IPGen) Import(moduleName string) (interface{}, error) {

	switch moduleName {
	case "ipgen":
		return moduleMap, nil
	default:
		return nil, fmt.Errorf("undefined module %q", moduleName)
	}
}


