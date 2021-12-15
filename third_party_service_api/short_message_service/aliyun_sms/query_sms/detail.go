package query_sms

import (
	. "github.com/DontBeProud/wow-easy-go/third_party_service_api/short_message_service/aliyun_sms/query_sms/status"
	. "github.com/DontBeProud/wow-easy-go/third_party_service_api/short_message_service/aliyun_sms/utils"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"time"
)

const (
	// AliYunSMSStatusTimeFormat 阿里云短信状态查询返回的时间格式
	AliYunSMSStatusTimeFormat = "2006-01-02 15:04:05"
)

// AliYunSmsStatusDetail 调用阿里云短信状态查询接口返回的单条短信的详细信息
type AliYunSmsStatusDetail struct {
	AliYunSmsStatusDetailInterface
	sendStatus   AliYunSMSSendStatus     // 短信发送状态. 1: 等待回执 2: 发送失败 3: 发送成功 其他: 未知状态
	errCode      AliYunSMSOperatorStatus // 运营商短信状态码. 短信发送成功：DELIVERED. https://help.aliyun.com/document_detail/102352.html
	phoneNum     string                  // 接收短信的手机号码
	content      string                  // 短信内容
	outId        string                  // 外部流水扩展字段
	templateCode string                  // 短信模板ID
	sendDate     time.Time               // 短信发送时间
	receiveDate  time.Time               // 短信接收时间
}

type AliYunSmsStatusDetailInterface interface {
	AliYunSMSSendStatusInterface     // IsSendSuccess() bool	/ QuerySendStatusCode() int64
	AliYunSMSOperatorStatusInterface // QueryOperatorStatusErrorInfo() *string / QueryOperatorStatusCode() string
	GetPhoneNumber() string          // 获取接收短信的手机号码
	GetSmsContent() string           // 获取短信内容
	GetOutId() string                // 获取外部流水扩展字段
	GetTemplateCode() string         // 获取短信模板ID
	GetSendTime() time.Time          // 获取短信发送时间
	GetReceiveTime() time.Time       // 获取短信接收时间
}

// IsSendSuccess 短信是否发送成功
func (d AliYunSmsStatusDetail) IsSendSuccess() bool {
	return d.sendStatus.IsSendSuccess() && d.QueryOperatorStatusErrorInfo() == nil
}

// GetPhoneNumber 获取接收短信的手机号码
func (d AliYunSmsStatusDetail) GetPhoneNumber() string {
	return d.phoneNum
}

// GetSmsContent 获取短信内容
func (d AliYunSmsStatusDetail) GetSmsContent() string {
	return d.content
}

// GetOutId 获取外部流水扩展字段
func (d AliYunSmsStatusDetail) GetOutId() string {
	return d.outId
}

// GetTemplateCode 获取短信模板ID
func (d AliYunSmsStatusDetail) GetTemplateCode() string {
	return d.templateCode
}

// GetSendTime 获取短信发送时间
func (d AliYunSmsStatusDetail) GetSendTime() time.Time {
	return d.sendDate
}

// GetReceiveTime 获取短信接收时间
func (d AliYunSmsStatusDetail) GetReceiveTime() time.Time {
	return d.receiveDate
}

// 从原始返回数据中提炼短信的详细信息
func refineDetailListFromAliYunRawResponse(raw *dysmsapi20170525.QuerySendDetailsResponseBodySmsSendDetailDTOs) ([]AliYunSmsStatusDetail, int) {
	size := 0
	if raw != nil {
		size = len(raw.SmsSendDetailDTO)
	}
	res := make([]AliYunSmsStatusDetail, size)
	for i := 0; i < size; i++ {
		tmp := raw.SmsSendDetailDTO[i]
		if tmp != nil {
			res[i] = AliYunSmsStatusDetail{
				sendStatus:   AliYunSMSSendStatus(ParseIntPointerIntoInt(tmp.SendStatus)),
				errCode:      AliYunSMSOperatorStatus(ParseStrPointerIntoString(tmp.ErrCode)),
				phoneNum:     ParseStrPointerIntoString(tmp.PhoneNum),
				content:      ParseStrPointerIntoString(tmp.Content),
				outId:        ParseStrPointerIntoString(tmp.OutId),
				templateCode: ParseStrPointerIntoString(tmp.TemplateCode),
				receiveDate:  ParseStrPointerIntoTime(tmp.ReceiveDate, AliYunSMSStatusTimeFormat),
				sendDate:     ParseStrPointerIntoTime(tmp.SendDate, AliYunSMSStatusTimeFormat),
			}
		}
	}
	return res, size
}
