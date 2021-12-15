package aliyun_sms

import (
	"github.com/DontBeProud/wow-easy-go/third_party_service_api/short_message_service/aliyun_sms/query_sms"
	aliOpenApi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
)

// AliYunSMSClientCredential 阿里云短信服务认证信息
type AliYunSMSClientCredential struct {
	AccessKeyId     string
	AccessKeySecret string
}

// AliYunSMSClient 阿里云短信服务用户
type AliYunSMSClient struct {
	AliYunSMSClientInterface
	config *aliOpenApi.Config
	client *dysmsapi20170525.Client
}

type AliYunSMSClientInterface interface {
	// GetClient 获取阿里云短信平台用户对象
	GetClient() *dysmsapi20170525.Client
	// QuerySmsStatus 查询短信的状态
	QuerySmsStatus(phoneNumber string, date string, page uint, size uint, bizId *string) (*query_sms.AliYunSmsStatusQueryResponse, error)
	// QuerySingleSmsStatus 查询单条短信的状态
	QuerySingleSmsStatus(phoneNumber string, date string, bizId *string) (*query_sms.AliYunSmsStatusQueryResponse, error)
}

// CreateAliYunSMSClient 基于认证信息创建用户对象
func (c AliYunSMSClientCredential) CreateAliYunSMSClient(optCfg *AliYunSMSOptionalConfig) (*AliYunSMSClient, error) {
	return createAliYunSMSClient(c, optCfg)
}

// QuerySmsStatus 查询短信状态
// 	phoneNumber:	国内短信：11位手机号码，例如1590000****. 国际/港澳台消息：国际区号+号码，例如8520000****
// 	date:			短信发送日期，支持查询最近30天的记录,格式为yyyyMMdd, 例如20181225。
// 	page:			分页查看发送记录，指定发送记录的当前页码。
// 	size:			分页查看发送记录，指定每页显示的短信记录数量。取值范围为1~50。
// 	bizId:			发送回执ID，即发送流水号。
func (c AliYunSMSClient) QuerySmsStatus(phoneNumber string, date string, page uint, size uint, bizId *string) (*query_sms.AliYunSmsStatusQueryResponse, error) {
	return query_sms.AliYunQuerySmsStatus(c.client, phoneNumber, date, page, size, bizId)
}

// QuerySingleSmsStatus 查询单条短信的发送状态
// 	phoneNumber:	国内短信：11位手机号码，例如1590000****. 国际/港澳台消息：国际区号+号码，例如8520000****
// 	date:			短信发送日期，支持查询最近30天的记录,格式为yyyyMMdd, 例如20181225
// 	bizId:			发送回执ID，即发送流水号
func (c AliYunSMSClient) QuerySingleSmsStatus(phoneNumber string, date string, bizId *string) (*query_sms.AliYunSmsStatusQueryResponse, error) {
	return query_sms.AlliYunQuerySingleSmsStatus(c.client, phoneNumber, date, bizId)
}

// GetClient 获取阿里云短信平台用户对象
func (c AliYunSMSClient) GetClient() *dysmsapi20170525.Client {
	return c.client
}
