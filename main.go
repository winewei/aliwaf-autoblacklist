package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	. "github.com/winewei/aliwaf-autoblacklist/pkg/utils"
	. "github.com/winewei/aliwaf-autoblacklist/pkg/aliyun"
)

var ctx = context.Background()

func init()  {
	os.Mkdir("./logs", 0755)
	file := "./" + "logs/aliwaf-autoblacklist.log"
	logFile, err := os.OpenFile(file, os.O_RDWR | os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
}

func main() {
	// take credentials from sys environment
	//wafRegion := GetEnvDefault("wafRegion", "cn-hangzhou")
	accessKeyId := GetEnvDefault("accessKeyId", "xxx")
	accessSecret := GetEnvDefault("accessSecret", "sss")
	//Domain := GetEnvDefault("Domain", "localhost")
	KeyPrefix := GetEnvDefault("KeyPrefix", "super_blacklist:*")
	redisURL := GetEnvDefault("redisURL", "redis://localhost:6379/0")
	Interval, _ := strconv.Atoi(GetEnvDefault("Interval", "5"))

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	log.Println("sys_info:", "accessKeyId:", accessKeyId, "match:", KeyPrefix, "redis:", redisURL, "Interval:", Interval)

	rdb := redis.NewClient(opt)
	defer rdb.Close()
	steps := 0
	for {
		// collect all domains from black keys
		tmpPre := strings.Split(KeyPrefix, ":")
		BasePrefix := tmpPre[0]
		go func() {
			var DL []string
			for {
				keys, cursor, err  := rdb.Scan(ctx,0, KeyPrefix, 200).Result()
				if err != nil {
					log.Println(err)
				}
				if len(keys) > 0 {
					for _, i := range keys {
						v := strings.Split(i, ":")
						if len(v) == 3 {
							DL = append(DL, v[1])
						} else {
							log.Println("DomainList got illegality redis key:", i)
						}
					}
				}
				if cursor == 0 {
					break
				}
			}

			DomainList := RemoveRep(DL)
			NewDomainListJson, _ := json.Marshal(DomainList)

			// store domain list
			rdb.Set(ctx,"protect:domains", NewDomainListJson, 0)

			//log.Println("domainList:", DomainList)
			// collect domain's attack ip address
			for _, domain := range DomainList {
				// key name: super_blacklist:www.google.com:*
				rKey := BasePrefix + ":" + domain + ":*"
				//log.Println("banList: scan redis key:", rKey)

				var IP []string
				for {
					keys, cursor, err  := rdb.Scan(ctx,0, rKey, 200).Result()
					if err != nil {
						log.Println(err)
					}
					if len(keys) > 0 {
						for _, i := range keys {
							v := strings.Split(i, ":")
							if len(v) == 3 {
								IP = append(IP, v[2])
							} else {
								log.Println("banList: got illegality redis key:", i)
							}
						}
					}
					if cursor == 0 {
						break
					}
				}

				// store domain's black ip list to redis
				t := RemoveRep(IP)
				BlackIpList := strings.Join(t, ",")
				DomainBanKeyName := "ban:" + domain
				domainIpList, _ := json.Marshal(t)
				//log.Println("DomainBanKeyName:", DomainBanKeyName, "iplist:", IP)
				rdb.Set(ctx, DomainBanKeyName, domainIpList,0)

				// write ip blacklist to cdn
				// check set to cdn?
				checkIpListKey := DomainBanKeyName + ":check"
				oldIpList, _ := rdb.Get(ctx, checkIpListKey).Result()

				if oldIpList != string(domainIpList) {
					rdb.Set(ctx, checkIpListKey, domainIpList, 0)
					log.Println("update cdn ip black list:", domain, BlackIpList)
					go CdnBlackList(domain, BlackIpList, accessKeyId, accessSecret)
				}
			}
			// get domains from redis key
			//var M []string
			//domains, _ := rdb.Get(ctx, "domains").Result()
			//if err := json.Unmarshal([]byte(domains), &M); err == nil {
			//}
		}()

		time.Sleep(time.Second * time.Duration(Interval))
		if steps == 20 {
			log.Println("Next 20 cycles ...")
			steps = 0
		}
		steps += 1
	}
}