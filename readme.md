# aliwaf auto set blacklist

## Dependences
- go 1.14
- docker
- docker-compose
- redis > 3.2
- [openresty waf](https://github.com/winewei/lua-waf)

## How to use
```shell script
git clone https://github.com/winewei/aliwaf-autoblacklist.git
cd aliwaf-autoblacklist
docker-compose up
```
## setup env in sys env or docker-compose
```shell script
wafRegion: cn-hangzhou
accessKeyId: xxx
accessSecret: xxx
KeyPrefix: super_blacklist:*
redisURL: redis://redis:6379/0
Domain: www.baidu.com
```


