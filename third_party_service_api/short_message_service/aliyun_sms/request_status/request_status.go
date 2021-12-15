package request_status

import (
	. "github.com/DontBeProud/wow-easy-go/third_party_service_api/short_message_service/aliyun_sms/utils"
)

// AliYunSmsApiRequestStatus 阿里云短信业务API请求返回的状态信息
// 详细内容参照 https://help.aliyun.com/document_detail/101346.htm?spm=a2c4g.11186623.0.0.450f15dbEpHIIi
type AliYunSmsApiRequestStatus struct {
	Code      string // 请求状态码, 返回OK代表请求成功.
	Message   string // 状态码的描述
	RequestId string // 请求ID
}

const (
	AliYunSmsApiRequestStatusSuccess = "OK"
)

// AliYunSmsApiRequestStatusInterface AliYunSmsApiRequestStatus interface
type AliYunSmsApiRequestStatusInterface interface {
	IsRequestSuccess() bool               // 请求调用阿里云API是否成功
	QueryRequestStatusErrorInfo() *string // 查询调用阿里云API的请求错误信息, 若无错误则返回nil
	GetRequestStatusCode() string         // 查询调用阿里云API的请求状态码
	GetRequestStatusMessage() string      // 查询调用阿里云API的请求状态码描述
	GetRequestStatusRequestId() string    // 查询调用阿里云API的请求Id
}

// CreateAliYunSmsApiRequestStatus create
func CreateAliYunSmsApiRequestStatus(Code *string, Message *string, RequestId *string) *AliYunSmsApiRequestStatus {
	return &AliYunSmsApiRequestStatus{
		Code:      ParseStrPointerIntoString(Code),
		Message:   ParseStrPointerIntoString(Message),
		RequestId: ParseStrPointerIntoString(RequestId),
	}
}

func (s AliYunSmsApiRequestStatus) IsRequestSuccess() bool {
	return s.Code == AliYunSmsApiRequestStatusSuccess
}

func (s AliYunSmsApiRequestStatus) GetRequestStatusCode() string {
	return s.Code
}

func (s AliYunSmsApiRequestStatus) GetRequestStatusMessage() string {
	return s.Message
}

func (s AliYunSmsApiRequestStatus) QueryRequestStatusErrorInfo() *string {
	if s.IsRequestSuccess() {
		return nil
	}
	return &s.Message
}

func (s AliYunSmsApiRequestStatus) GetRequestStatusRequestId() string {
	return s.RequestId
}
