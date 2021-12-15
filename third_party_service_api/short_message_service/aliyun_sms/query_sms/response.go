package query_sms

import (
	"errors"
	. "github.com/DontBeProud/wow-easy-go/third_party_service_api/short_message_service/aliyun_sms/request_status"
	. "github.com/DontBeProud/wow-easy-go/third_party_service_api/short_message_service/aliyun_sms/utils"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
)

// AliYunSmsStatusQueryResponse 调用阿里云短信状态查询接口返回的结果(对阿里云api原生返回值进行了可用性的优化)
type AliYunSmsStatusQueryResponse struct {
	AliYunSmsStatusQueryResponseInterface
	AliYunSmsApiRequestStatus                         // 基础的请求状态信息
	totalCount                int                     // 短信发送总条数
	detailCount               int                     // 返回的详细信息数量, 即SmsSendDetailDTOs数组成员的数量
	smsSendDetailDTOs         []AliYunSmsStatusDetail // 每条短信对应的详细数据
}

type AliYunSmsStatusQueryResponseInterface interface {
	AliYunSmsApiRequestStatusInterface
	GetTotalCount() int                           // 查询返回的短信总条数
	GetDetailCount() int                          // 查询返回的详细信息总条数
	GetDetailByIndex(uint) *AliYunSmsStatusDetail // 根据序号获取详细信息指针. 0 <= index < detailCount
	GetAllDetails() *[]AliYunSmsStatusDetail      // 获取全部详细信息
}

// GetTotalCount 查询返回的短信总条数
func (r AliYunSmsStatusQueryResponse) GetTotalCount() int {
	return r.totalCount
}

// GetDetailCount 查询返回的详细信息总条数
func (r AliYunSmsStatusQueryResponse) GetDetailCount() int {
	return r.detailCount
}

// GetDetailByIndex 根据序号获取详细信息指针. 0 <= index < detailCount
func (r AliYunSmsStatusQueryResponse) GetDetailByIndex(index uint) *AliYunSmsStatusDetail {
	if index < 0 || int(index) >= r.detailCount {
		return nil
	}
	return &r.smsSendDetailDTOs[index]
}

// GetAllDetails 获取全部详细信息
func (r AliYunSmsStatusQueryResponse) GetAllDetails() *[]AliYunSmsStatusDetail {
	return &r.smsSendDetailDTOs
}

// 解析阿里云的短信状态查询返回结果的MsgBody
func parseQuerySendDetailsResponseBody(raw *dysmsapi20170525.QuerySendDetailsResponseBody) (*AliYunSmsStatusQueryResponse, error) {
	if raw == nil {
		return nil, errors.New("QuerySendDetailsResponseBody == nil")
	}

	cnt, err := ParseStrPointerIntoInt64(raw.TotalCount)
	if err != nil {
		return nil, err
	}

	res := AliYunSmsStatusQueryResponse{
		totalCount:                cnt,
		AliYunSmsApiRequestStatus: *CreateAliYunSmsApiRequestStatus(raw.Code, raw.Message, raw.RequestId),
	}

	res.smsSendDetailDTOs, res.detailCount = refineDetailListFromAliYunRawResponse(raw.SmsSendDetailDTOs)
	return &res, err
}
