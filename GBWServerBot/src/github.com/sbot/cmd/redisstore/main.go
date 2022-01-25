package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Entry struct {

	Name string  `json:"name"`
	Value string `json:"value"`
}

type Test struct {

	Time  uint64  `json:"time"`
	Entry []*Entry   `json:"entry"`

}

func makeTest(n int) (uint64, string) {


	ts := uint64(time.Now().UnixNano()/(1000*1000))

	result := &Test{
		Time:  ts,
		Entry: make([]*Entry,0),
	}

	for n>0 {

		name := fmt.Sprintf("name%d",n)
		value := fmt.Sprintf("value%d",n)

		result.Entry =append(result.Entry,&Entry{
			Name:  name,
			Value: value,
		})
		n--
	}

	data,_:= json.Marshal(result)

	return ts,string(data)
}


func main(){


	client := redis.NewClient(
		  &redis.Options{
			  Addr:               "192.168.198.128:6379",
		  })

	//r,_:=client.ConfigGet(context.Background(),"databases").Result()

/*
	t := `local key = KEYS[1]
local timeKey = KEYS[2]
local valueKey = KEYS[3]
local timeStamp = tonumber(ARGV[1])
local value = ARGV[2]

local zaddResult = redis.call("zadd",timeKey,timeStamp,key)
local hsetResult = redis.call("hset",valueKey,key,value)

return {zaddResult,hsetResult}`*/

	/*
	d := `local key = KEYS[1]
local timeKey = KEYS[2]
local valueKey = KEYS[3]

local zremResult = redis.call("zrem",timeKey,key)
local hdelResult = redis.call("hdel",valueKey,key)

return {zremResult,hdelResult}`*/


	//client.HSetNX(context.Background(),"redis.databases","node",1)
	//client.HSetNX(context.Background(),"redis.databases","attack",2)
	//client.HSetNX(context.Background(),"redis.databases","attack2",100)

	q := `
local timeKey = KEYS[1]
local valueKey = KEYS[2]

local timeStart = tonumber(ARGV[1])
local  timeEnd = tonumber(ARGV[2])

local key = ""
local value = ""
local tm = 0
local JsonValue

local cur = 0
local keys = {}
local tnum = 0
local ckeys = {}

repeat

    keys = redis.call("hscan",valueKey,cur)
    cur = tonumber(keys[1])

    for i = 2,#keys[2],2 do

        key = keys[2][1]
        value = keys[2][2]
        JsonValue = cjson.decode(value)
        tm = tonumber(JsonValue["time"])
        if tm>=timeStart and tm<=timeEnd then
            if %s then
                table.insert(ckeys,key)
                tnum = tnum+1
            end
        end
    end
until cur==0

for _,k in ipairs(ckeys) do

    if redis.call("zrem",timeKey,k) == 0 then
        repeat
            redis.call("zrem",timeKey,k)
        until redis.call("zscore",timeKey,k) == nil
    end

    if redis.call("hdel",valueKey,k) == 0 then
        repeat
            redis.call("hdel",valueKey,k)
        until redis.call("hget",valueKey,k) == nil
    end
end


return tnum
`

	qstr := fmt.Sprintf(q,`#JsonValue["entry"]==2`)

	//fmt.Println(qstr)
	//qstr,_:= tmp.Parse("1==1")

	script := redis.NewScript(qstr)
	r,err:=script.Run(context.Background(),client,[]string{"redis.node.info.times","redis.node.info.values"},
		0,time.Now().UnixNano()/(1000*1000),3,10,0).Result()

	fmt.Println(r,err)

	//fmt.Println(r,len(r.([]interface{})))

	//dscript := redis.NewScript(d)

	/*
	tm,value := makeTest(1)

	n := 100000

	for n >0 {

		key := fmt.Sprintf("192.168.1.%d",n)
		tm,value = makeTest(n%5+1)

		script.Run(context.Background(),client,[]string{key,"redis.node.info.times","redis.node.info.values"},tm+1000,value).Result()

		n--
	}*/

	//r,_:=client.HVals(context.Background(),"redis.databases").Result()
	//k,_:= client.HKeys(context.Background(),"redis.databases").Result()

	//client.ZPopMax()


//	fmt.Println(strconv.ParseInt(r[1].(string),10,32))
}


