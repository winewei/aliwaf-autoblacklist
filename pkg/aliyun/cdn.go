package aliyun

import (
	cdn20180510 "github.com/alibabacloud-go/cdn-20180510/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
	"log"
)

func createClient (accessKeyId *string, accessKeySecret *string) (_result *cdn20180510.Client, _err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("cdn.aliyuncs.com")
	_result = &cdn20180510.Client{}
	_result, _err = cdn20180510.NewClient(config)
	return _result, _err
}

func CdnBlackList (domainname, blacklist, accessKeyId, accessKeySecret string) {
	client, err := createClient(tea.String(accessKeyId), tea.String(accessKeySecret))
	if err != nil {
		log.Println("newAliCdnClentErr:", err)
	}

	setIpBlackListConfigRequest := &cdn20180510.SetIpBlackListConfigRequest{
		DomainName: tea.String(domainname),
		BlockIps: tea.String(blacklist),
	}
	// 复制代码运行请自行打印 API 的返回值
	rep, err := client.SetIpBlackListConfig(setIpBlackListConfigRequest)
	if err != nil {
		log.Println(err)
	}
	log.Println(domainname, rep.Body)
}