package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
	"time"
	. "aliwaf-blacklist/pkg/utils"
	. "aliwaf-blacklist/pkg/aliyun"
)

var ctx = context.Background()

func init()  {
	os.Mkdir("./logs", 0755)
	file := "./" + "logs/aliwaf-autoblacklist.txt"
	logFile, err := os.OpenFile(file, os.O_RDWR | os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
}

func main() {
	// take credentials from sys environment
	wafRegion := GetEnvDefault("wafRegion", "cn-hangzhou")
	accessKeyId := GetEnvDefault("accessKeyId", "xxx")
	accessSecret := GetEnvDefault("accessSecret", "sss")
	Domain := GetEnvDefault("Domain", "localhost")
	KeyPrefix := GetEnvDefault("KeyPrefix", "super_blacklist:*")
	redisURL := GetEnvDefault("redisURL", "redis://localhost:6379/0")
	Interval, _ := strconv.Atoi(GetEnvDefault("Interval", "5"))

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	log.Println("sys_info:", Domain, wafRegion, KeyPrefix, redisURL, "Interval:", Interval)

	rdb := redis.NewClient(opt)
	defer rdb.Close()
	steps := 0
	for {
		go func() {
			var Ipaddress = []string{}
			type M map[string]interface{}
			tempMap := make(M)
			for {
				keys, cursor, err  := rdb.Scan(ctx,0, KeyPrefix, 200).Result()
				if err != nil {
					log.Println(err.Error())
				}
				// if blacklist not 0, will enable waf blacklist
				if len(keys) > 0 {
					rdb.Set(ctx, "enable_waf_black", 1, 0)
					for _, i := range keys {
						val, err := rdb.Get(ctx, i).Result()
						if err != nil{
							log.Println(err.Error())
						}
						//log.Println("found_black_key:", i, "value:", val)
						// waf blacklist max lenth is 200
						if len(Ipaddress) <= 200 {
							Ipaddress = append(Ipaddress, val)
						} else {
							log.Println("out of max blacklist lenth 200:", val)
						}
					}
				}
				if cursor == 0 {
					break
				}
			}
			tempMap["remoteAddr"] = RemoveRep(Ipaddress)
			t_data, _ := json.Marshal(tempMap)

			rdb.Set(ctx,"new_waf_blacklist", t_data, 0)
			//log.Println("new_waf_blacklist:", string(t_data))
		}()

		// waf black
		is_waf_black, _ := rdb.Get(ctx, "enable_waf_black").Result()
		if is_waf_black == "1" {
			new_waf_blacklist, _ := rdb.Get(ctx, "new_waf_blacklist").Result()
			old_waf_blacklist, _ := rdb.Get(ctx, "old_waf_blacklist").Result()
			if old_waf_blacklist != new_waf_blacklist {
				log.Println("new_waf_blacklist:", new_waf_blacklist, "old_waf_blacklist:", old_waf_blacklist)
				// update old_waf_blacklist from new_waf_blacklist
				rdb.Set(ctx, "old_waf_blacklist", new_waf_blacklist, 0)
				go Waf_blacklist(new_waf_blacklist, Domain, wafRegion,accessKeyId, accessSecret)
				log.Println("waf_blacklist:", new_waf_blacklist)
			}
		}

		time.Sleep(time.Second * time.Duration(Interval))
		if steps == 20 {
			log.Println("Next 20 cycles ...")
			steps = 0
		}
		steps += 1
	}
}