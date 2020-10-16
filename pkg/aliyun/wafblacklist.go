package aliyun

import (
	waf_openapi "github.com/aliyun/alibaba-cloud-sdk-go/services/waf-openapi"
	"log"
)

func Waf_blacklist(Rule, Domain,wafRegion, accessKeyId, accessSecret  string) () {
	client, err := waf_openapi.NewClientWithAccessKey(wafRegion, accessKeyId, accessSecret)

	// get InstanceId
	instanceid_request := waf_openapi.CreateDescribeInstanceInfoRequest()
	instanceid_request.Scheme = "https"
	instanceid_response, err := client.DescribeInstanceInfo(instanceid_request)
	InstanceId := instanceid_response.InstanceInfo.InstanceId

	log.Println("waf_InstanceId:", instanceid_response.InstanceInfo.InstanceId)
	log.Println("Domain: ", Domain)

	// update waf blacklist
	request := waf_openapi.CreateModifyProtectionModuleRuleRequest()
	request.Scheme = "https"

	request.Domain = Domain
	request.DefenseType = "ac_blacklist"
	request.Rule = Rule
	request.LockVersion = "1"
	request.InstanceId = InstanceId

	response, err := client.ModifyProtectionModuleRule(request)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(response)
}