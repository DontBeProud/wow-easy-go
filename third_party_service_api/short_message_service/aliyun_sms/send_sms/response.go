package send_sms

import (
	"errors"
	. "github.com/DontBeProud/wow-easy-go/third_party_service_api/short_message_service/aliyun_sms/request_status"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
)

type AliYunSmsSendResponse struct {
	AliYunSmsApiRequestStatus         // 基础的请求状态信息
	BizId                     *string // 回执ID. 可根据发送回执ID查询具体的发送状态. 若为nil或长度为0则说明该次请求未返回合法的回执ID
}

type AliYunSmsSendResponseInterface interface {
	AliYunSmsApiRequestStatusInterface
	GetBizId() string   // 获取回执ID
	IsBizIdValid() bool // 返回的回执ID是否合法
}

// GetBizId 获取回执ID
func (r AliYunSmsSendResponse) GetBizId() *string {
	return r.BizId
}

// IsBizIdValid 返回的回执ID是否合法
func (r AliYunSmsSendResponse) IsBizIdValid() bool {
	return r.BizId != nil && len(*r.BizId) == 0
}

func parseSendSmsRawResponse(raw *dysmsapi20170525.SendSmsResponse) (*AliYunSmsSendResponse, error) {
	if raw == nil {
		return nil, errors.New("SendSmsResponse == nil")
	}

	if raw.Body == nil {
		return nil, errors.New("SendSmsResponse.Body == nil")
	}

	return &AliYunSmsSendResponse{
		BizId:                     raw.Body.BizId,
		AliYunSmsApiRequestStatus: *CreateAliYunSmsApiRequestStatus(raw.Body.Code, raw.Body.Message, raw.Body.RequestId),
	}, nil
}
