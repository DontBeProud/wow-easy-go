package aliyun_sms

import (
	"github.com/DontBeProud/wow-easy-go/third_party_service_api/short_message_service/aliyun_sms/send_sms"
)

// AliYunSMSClientSender 可发送短信的阿里云短信客户对象
type AliYunSMSClientSender struct {
	signName string
	template *SmsTemplateInterface
	AliYunSMSClientInterface
	AliYunSMSClientSenderInterface
}

// CreateAliYunSMSClientSender  基于认证信息创建可发送短信的用户对象
func (c AliYunSMSClientCredential) CreateAliYunSMSClientSender(signName string, template SmsTemplateInterface, optCfg *AliYunSMSOptionalConfig) (
	*AliYunSMSClientSender, error) {
	return createAliYunSMSClientSender(c, signName, template, optCfg)
}

// CreateAliYunSMSClientSender 基于基础用户对象创建可发送短信的用户对象
func (c AliYunSMSClient) CreateAliYunSMSClientSender(signName string, template SmsTemplateInterface) (*AliYunSMSClientSender, error) {
	return c.createAliYunSMSClientSender(signName, template)
}

// AliYunSMSClientSenderInterface AliYunSMSClientSender interface
type AliYunSMSClientSenderInterface interface {
	AliYunSMSClientInterface
	// QueryAliYunSmsSendApiPhoneNumberLimit 查询阿里云短信发送功能单次请求的手机号码数量上限
	QueryAliYunSmsSendApiPhoneNumberLimit() int
	// SendSms 发送短信
	SendSms(PhoneNumberList []string, UpExtendCode *string, OutId *string, params interface{}) (*send_sms.AliYunSmsSendResponse, error)
}

// SendSms 发送短信
// PhoneNumberList: 接收短信的手机号列表. 0 < len(PhoneNumberList) <= QueryAliYunSmsSendApiPhoneNumberLimit()
// UpExtendCode: 可选,未使用该字段则传nil. 上行短信扩展码, 上行短信, 指发送给通信服务提供商的短信, 用于定制某种服务、完成查询，或是办理某种业务等, 需要收费的, 按运营商普通短信资费进行扣费
// OutId: 可选,未使用该字段则传nil. 外部流水扩展字段
// rawTemplateParam: 原始模板参数. 函数内部调用AliYunSMSClientSender.template.GenerateTemplateParam(params interface{})将该参数格式化为Json格式的参数
func (s AliYunSMSClientSender) SendSms(PhoneNumberList []string, UpExtendCode *string, OutId *string, rawTemplateParam interface{}) (*send_sms.AliYunSmsSendResponse, error) {
	return s.sendSms(PhoneNumberList, UpExtendCode, OutId, rawTemplateParam)
}

// QueryAliYunSmsSendApiPhoneNumberLimit 查询阿里云短信发送功能单次请求的手机号码数量上限
func (s AliYunSMSClientSender) QueryAliYunSmsSendApiPhoneNumberLimit() int {
	return send_sms.QueryAliYunSmsSendApiPhoneNumberLimit()
}

// 发送短信
func (s AliYunSMSClientSender) sendSms(PhoneNumberList []string, UpExtendCode *string, OutId *string, rawTemplateParam interface{}) (*send_sms.AliYunSmsSendResponse, error) {
	signName := s.signName
	templateId := (*s.template).GetTemplateId()
	templateParams := (*s.template).GenerateTemplateParam(rawTemplateParam)
	apiParams := send_sms.CreateAliYunSmsSendRequestParams(
		PhoneNumberList,
		signName,
		templateId,
		templateParams,
		UpExtendCode,
		OutId,
	)

	client := s.AliYunSMSClientInterface.GetClient()
	return send_sms.AliYunSmsSend(client, apiParams)
}
