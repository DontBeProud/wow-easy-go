package send_sms

import (
	"errors"
	"fmt"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"strconv"
)

const (
	// AliYunSmsSendPhoneNumberLimit 批量发送短信功能单次请求的手机号码数量上限
	AliYunSmsSendPhoneNumberLimit = 1000
)

// QueryAliYunSmsSendApiPhoneNumberLimit 查询阿里云短信发送功能单次请求的手机号码数量上限
func QueryAliYunSmsSendApiPhoneNumberLimit() int {
	return AliYunSmsSendPhoneNumberLimit
}

// AliYunSmsSendRequestParams 阿里云平台请求发送短信API需要的参数
type AliYunSmsSendRequestParams struct {
	AliYunSmsSendRequestParamsInterface
	PhoneNumberList []string // 接收短信的手机号码. 支持同时向多个手机发送, 数量上限为 AliYunSmsSendPhoneNumberLimit
	SignName        string   // 短信签名名称
	TemplateCode    string   // 短信模板ID。
	TemplateParam   *string  // 可选. 短信模板变量对应的实际值，JSON格式。支持传入多个参数，示例：{"name":"张三","number":"15038****76"}
	UpExtendCode    *string  // 可选. 上行短信扩展码，上行短信，指发送给通信服务提供商的短信，用于定制某种服务、完成查询，或是办理某种业务等，需要收费的，按运营商普通短信资费进行扣费
	OutId           *string  // 可选. 外部流水扩展字段
}

type AliYunSmsSendRequestParamsInterface interface {
	aliYunSmsSendRequestParamsInterface
	CheckError() error // 判断各参数是否合法
}

// CreateAliYunSmsSendRequestParams create
func CreateAliYunSmsSendRequestParams(PhoneNumberList []string, SignName string, TemplateCode string, TemplateParam *string, UpExtendCode *string, OutId *string) *AliYunSmsSendRequestParams {
	return &AliYunSmsSendRequestParams{
		PhoneNumberList: PhoneNumberList,
		SignName:        SignName,
		TemplateCode:    TemplateCode,
		TemplateParam:   TemplateParam,
		UpExtendCode:    UpExtendCode,
		OutId:           OutId,
	}
}

// AliYunSmsSend 调用阿里云SMS-API发送短信
func AliYunSmsSend(c *dysmsapi20170525.Client, params *AliYunSmsSendRequestParams) (*AliYunSmsSendResponse, error) {
	return sendSms(c, params)
}

// CheckError 判断各参数是否非法
func (p AliYunSmsSendRequestParams) CheckError() error {
	return p.checkError()
}

// interface
type aliYunSmsSendRequestParamsInterface interface {
	converseToNativeApiParam() (*dysmsapi20170525.SendSmsRequest, error) // 转换成原生API的参数
	checkError() error                                                   // 判断各参数是否非法
	checkPhoneNumbers() error                                            // 判断手机号列表是否合法
	checkSignName() error                                                // 判断短信签名名称是否合法
	checkTemplatedCode() error                                           // 判断短信模板ID是否合法
	serializePhoneList() *string                                         // 将手机号码列表按API文档规范进行序列化
}

// 判断各参数是否非法
func (p AliYunSmsSendRequestParams) checkError() error {
	for _, fn := range []func() error{
		p.checkPhoneNumbers,
		p.checkSignName,
		p.checkTemplatedCode,
	} {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

// 判断手机号列表是否合法
func (p AliYunSmsSendRequestParams) checkPhoneNumbers() error {
	phoneNumberSize := len(p.PhoneNumberList)

	// 未传入手机号
	if phoneNumberSize == 0 {
		return errors.New("len(AliYunSmsSendRequestParams.phoneNumberSize) == 0")
	}

	// 是否超过单次请求的手机号码上限
	if phoneNumberSize > AliYunSmsSendPhoneNumberLimit {
		return errors.New("sendSms failed. 单次请求传入的手机号数量超过上限. 上限数量为: " + strconv.Itoa(AliYunSmsSendPhoneNumberLimit))
	}

	for index, number := range p.PhoneNumberList {
		if len(number) == 0 {
			return errors.New(fmt.Sprintf("AliYunSmsSendRequestParams.PhoneNumberList 第 %d 个手机号为空值", index+1))
		}
	}

	return nil
}

// 判断短信签名名称是否合法
func (p AliYunSmsSendRequestParams) checkSignName() error {
	if len(p.SignName) == 0 {
		return errors.New("len(AliYunSmsSendRequestParams.SignName) == 0")
	}
	return nil
}

// 判断短信模板ID是否合法
func (p AliYunSmsSendRequestParams) checkTemplatedCode() error {
	if len(p.TemplateCode) == 0 {
		return errors.New("len(AliYunSmsSendRequestParams.TemplateCode) == 0")
	}
	return nil
}

// 将手机号码列表按API文档规范进行序列化
func (p AliYunSmsSendRequestParams) serializePhoneList() *string {
	res := ""
	for index, number := range p.PhoneNumberList {
		if index != 0 {
			res += ","
		}
		res += number
	}
	return &res
}

// 转换成原生API的参数
func (p AliYunSmsSendRequestParams) converseToNativeApiParam() (*dysmsapi20170525.SendSmsRequest, error) {
	if err := p.CheckError(); err != nil {
		return nil, err
	}
	return &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:    p.serializePhoneList(),
		SignName:        &p.SignName,
		TemplateCode:    &p.TemplateCode,
		SmsUpExtendCode: p.UpExtendCode,
		TemplateParam:   p.TemplateParam,
		OutId:           p.OutId,
	}, nil
}

// send
func sendSms(c *dysmsapi20170525.Client, params *AliYunSmsSendRequestParams) (*AliYunSmsSendResponse, error) {

	if c == nil || params == nil {
		return nil, errors.New("QuerySendDetails failed. param is nil")
	}

	p, err := params.converseToNativeApiParam()
	if err != nil {
		return nil, err
	}

	res, err := c.SendSms(p)
	if err != nil {
		return nil, err
	}

	return parseSendSmsRawResponse(res)
}
