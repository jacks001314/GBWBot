package redisstore

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sbot/store"
	"strconv"
	"time"
)

const (
	REDIS_CONFIG_DATABASES = "databases"

	//store databasename----index map in hset
	REDIS_DATABASE_KEY = "Redis.DataBaseStore"


)

var putSCRIPT = `
local key = KEYS[1]
local timeKey = KEYS[2]
local valueKey = KEYS[3]
local timeStamp = tonumber(ARGV[1])
local value = ARGV[2]

local zaddResult = redis.call("zadd",timeKey,timeStamp,key)
local hsetResult = redis.call("hset",valueKey,key,value)

return {zaddResult,hsetResult}
`
var delScript = `
local key = KEYS[1]
local timeKey = KEYS[2]
local valueKey = KEYS[3]

local zremResult = redis.call("zrem",timeKey,key)
local hdelResult = redis.call("hdel",valueKey,key)

return {zremResult,hdelResult}
`

type RedisStore struct {

	cfg *store.Config
	redisClient *redis.Client

	ctx context.Context

	db int

	storeTimeKey string

	storeValueKey string

	putScript *redis.Script

	delScript *redis.Script

}

func getMaxDatabases(ctx context.Context,redisClient *redis.Client) int   {

	r,err:=redisClient.ConfigGet(ctx,REDIS_CONFIG_DATABASES).Result()

	if err!=nil {

		return 16 // default
	}

	v,err:= strconv.ParseInt(r[1].(string),10,32)

	if err!=nil {
		return 16
	}

	return int(v)
}

/*redis hset store the database name correspond to db index values*/

func getDatabaseIndex(ctx context.Context,redisClient *redis.Client,name string) (int,error) {

	scriptStr := `
local hname = KEYS[1]
local dbname = ARGV[1]

local maxDatabaseIndex =  tonumber(ARGV[2])
local value = 0
local maxV = 0

local valueS = redis.call("hget",hname,dbname)

if valueS then

    return tonumber(valueS)
end

local values = redis.call("hvals",hname)

for _,v in ipairs(values) do

    value = tonumber(v)

    if value>maxV then

        maxV = value
    end
end

if maxV+1 >= maxDatabaseIndex then
    return maxDatabaseIndex-1
end

redis.call("hset",hname,dbname,maxV+1)

return maxV+1
`

	script := redis.NewScript(scriptStr)

	result,err:=script.Run(ctx,redisClient,[]string{REDIS_DATABASE_KEY},name,getMaxDatabases(ctx,redisClient)).Result()

	if err!=nil {

		return -1,nil
	}

	return int(result.(int64)),nil
}


//open a database to store
func (s *RedisStore)Open(cfg *store.Config) (store.Store,error) {

	ctx := context.Background()
	addr :=  fmt.Sprintf("%s:%d",cfg.Host,cfg.Port)
	timeOut := time.Duration(cfg.Timeout)*time.Millisecond

	//create a temp redis client to get database index
	redisClient := redis.NewClient(&redis.Options{

		Addr:              addr,
		Password:           cfg.Pass,
		DialTimeout:       timeOut ,
		DB: 0,
	})

	defer redisClient.Close()

	//get database index
	db,err := getDatabaseIndex(ctx,redisClient,cfg.DB)

	if err!=nil {

		return nil,err
	}

	//create a redis client by database index
	client := redis.NewClient(&redis.Options{

		Addr:               addr,
		Password:           cfg.Pass,
		DialTimeout:        timeOut,
		DB: db,
	})

	return &RedisStore{
		cfg:         cfg,
		redisClient: client,
		ctx:         ctx,
		db:          db,
		storeTimeKey: fmt.Sprintf("redis.%s.%s.times",cfg.DB,cfg.Table),
		storeValueKey: fmt.Sprintf("redis.%s.%s.values",cfg.DB,cfg.Table),
		putScript: redis.NewScript(putSCRIPT),
		delScript: redis.NewScript(delScript),
	},nil

}


func (s *RedisStore) Close() error {

	return s.redisClient.Close()
}


func (s *RedisStore) Put(key string,time uint64,value interface{}) error {

	script := s.putScript

	data,err := json.Marshal(value)

	if err !=nil {

		return err
	}

	script.Run(s.ctx,s.redisClient,[]string{key,s.storeTimeKey,s.storeValueKey},time,string(data)).Result()

	return nil
}

func (s *RedisStore) Get(key string,value interface{}) (bool,error) {

	r,err := s.redisClient.HGet(s.ctx,s.storeValueKey,key).Result()

	if err!=nil {

		return false,err
	}

	err = json.Unmarshal([]byte(r),value)

	if err != nil {

		return false,err
	}

	return true,nil
}

func (s *RedisStore) Del(key string) error {

	script := s.delScript

	script.Run(s.ctx,s.redisClient,[]string{key,s.storeTimeKey,s.storeValueKey}).Result()

	return nil
}

func getQueryScript(q string) string {

	script := `
local timeKey = KEYS[1]
local valueKey = KEYS[2]

local timeStart = tonumber(ARGV[1])
local  timeEnd = tonumber(ARGV[2])

local page = tonumber(ARGV[3])
local size = tonumber(ARGV[4])
local isDec = tonumber(ARGV[5])
local results = {}
local result = ""
local key = ""
local tm = 0

local cmd = "zrange"
if isDec == 1 then
    cmd = "zrevrange"
end

local startIndex = (page-1)*size+1
local endIndex = startIndex+size-1
local tnum = 0

local timeAndKeys = redis.call(cmd,timeKey,0,-1,"withscores")

for i = 1,#timeAndKeys,2 do
    tm = tonumber(timeAndKeys[i+1])
    if tm>=timeStart and tm<=timeEnd then
        key = timeAndKeys[i]
        result = redis.call("hget",valueKey,key)
        local JsonValue = cjson.decode(result)

        if %s then
            tnum = tnum+1

            if tnum>=startIndex and tnum<=endIndex then
                table.insert(results,key)
                table.insert(results,result)
            end
        end
    end
end

table.insert(results,tnum)

return results
`
	return fmt.Sprintf(script,q)

}

func (s *RedisStore) Query(queryString string,timeRange [2]uint64,pageable *store.Pageable) (*store.QueryResult,error) {

	var tnum uint64
	var tpage uint64

	queryScript := getQueryScript(queryString)

	script := redis.NewScript(queryScript)

	keys := []string{s.storeTimeKey,s.storeValueKey}

	isDec := 0
	if pageable.ISDec {
		isDec = 1
	}

	r,err := script.Run(s.ctx,s.redisClient,keys,timeRange[0],timeRange[1],pageable.Page,pageable.Size,isDec).Result()

	if err!=nil {

		return nil,err
	}

	qr := r.([]interface{})
	n := len(qr)

	if n <=1 {

		return &store.QueryResult{
			Page:    pageable.Page,
			Size:    pageable.Size,
			TPage:   0,
			TNum:    0,
			Results: []*store.ResultEntry{},
		},nil
	}

	tnum = uint64(qr[n-1].(int64))
	tpage = tnum/pageable.Size

	if tnum%pageable.Size !=0 {
		tpage++
	}

	results := &store.QueryResult{
		Page:    pageable.Page,
		Size:    pageable.Size,
		TPage:   tpage,
		TNum:    tnum,
		Results: make([]*store.ResultEntry,0),
	}
	for i:=0;i<n-1;i+=2 {

		entry := &store.ResultEntry{
			Key:   qr[i].(string),
			Value: qr[i+1].(string),
		}

		results.Results = append(results.Results,entry)
	}
	
	return results,nil
}

func (s *RedisStore) FlushDB() error {

	return s.redisClient.FlushDB(s.ctx).Err()

}

func (s *RedisStore) Count() uint64 {

	v,err:= s.redisClient.ZCard(s.ctx,s.storeTimeKey).Result()

	if err!=nil {
		return 0
	}

	return uint64(v)
}

func getFacetScript( query string,term string) string  {

	fscritp := `

local valueKey = KEYS[1]

local num = tonumber(ARGV[1])
local isDec = tonumber(ARGV[2])

local buckets = {}

local cur = 0
local values = {}
local sortValues = {}
local results = {}

repeat

    values = redis.call("hscan",valueKey,cur)
    cur = tonumber(values[1])
    local kvals = values[2]

    for i=1,#kvals,2 do

        local JsonValue = cjson.decode(kvals[i+1])
        if %s then
            local vv = %s

            if buckets[vv] then
                buckets[vv] = buckets[vv]+1
            else
                buckets[vv] = 1
            end
        end
    end

until cur == 0


for k,v in pairs(buckets) do

    table.insert(sortValues,{key=k,count=v})
end

local cmp = function(a,b)
    if a and b then
        if isDec == 1 then
            return a.count>b.count
        else
            return a.count<b.count
        end
    end
    return false
end

table.sort(sortValues,cmp)

if num >= #sortValues then
    return cjson.encode(sortValues)
end

for i=1,#sortValues,1 do

    table.insert(results,sortValues[i])
    if i>= num then
        break
    end
end

return cjson.encode(results)

`
	return fmt.Sprintf(fscritp,query,term)

}

func (s *RedisStore) Facet(query string ,term string,num uint64,isDec bool) ([]*store.TermFacet,error) {

	results := make([]*store.TermFacet,0)

	src := getFacetScript(query,term)
	script := redis.NewScript(src)

	keys := []string{s.storeValueKey}

	dec := 0
	if isDec {
		dec = 1
	}

	r,err := script.Run(s.ctx,s.redisClient,keys,num,dec).Result()

	if err!=nil {

		return nil,err
	}

	data := []byte(r.(string))

	err =json.Unmarshal(data,&results)

	return results,err
}
