package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	waf_openapi "github.com/aliyun/alibaba-cloud-sdk-go/services/waf-openapi"
	"log"
)

func Waf_blacklist(Rule, Domain,wafRegion, accessKeyId, accessSecret  string) () {
	client, err := waf_openapi.NewClientWithAccessKey(wafRegion, accessKeyId, accessSecret)

	// get InstanceId
	instanceidRequest := waf_openapi.CreateDescribeInstanceInfoRequest()
	instanceidRequest.Scheme = "https"
	instanceidResponse, err := client.DescribeInstanceInfo(instanceidRequest)
	if err != nil {
		log.Println(err.Error())
	}
	InstanceId := instanceidResponse.InstanceInfo.InstanceId

	log.Println("waf_InstanceId:", instanceidResponse.InstanceInfo.InstanceId)
	log.Println("Domain: ", Domain)

	// get ruleid
	ruleidReq := waf_openapi.CreateDescribeProtectionModuleRulesRequest()
	ruleidReq.Scheme = "https"
	ruleidReq.InstanceId = InstanceId
	ruleidReq.Domain = Domain
	ruleidReq.DefenseType = "ac_blacklist"
	ruleidResponse, err := client.DescribeProtectionModuleRules(ruleidReq)
	log.Println("waf_blacklist_info:", ruleidResponse.GetHttpContentString())
	if err != nil {
		log.Println(err)
	}

	// update waf blacklist
	request := waf_openapi.CreateModifyProtectionModuleRuleRequest()
	request.Scheme = "https"
	request.Domain = Domain
	request.DefenseType = "ac_blacklist"
	request.Rule = Rule
	request.LockVersion = requests.NewInteger(1)
	request.InstanceId = InstanceId
	request.RuleId  = requests.NewInteger(int(ruleidResponse.Rules[0].RuleId))
	response, err := client.ModifyProtectionModuleRule(request)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(response)
}