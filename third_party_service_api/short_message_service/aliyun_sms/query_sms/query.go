package query_sms

import (
	"errors"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

// AliYunQuerySmsStatus 查询短信状态
// 	phoneNumber:	国内短信：11位手机号码，例如1590000****. 国际/港澳台消息：国际区号+号码，例如8520000****
// 	date:			短信发送日期，支持查询最近30天的记录,格式为yyyyMMdd, 例如20181225。
// 	page:			分页查看发送记录，指定发送记录的当前页码。
// 	size:			分页查看发送记录，指定每页显示的短信记录数量。取值范围为1~50。
// 	bizId:			发送回执ID，即发送流水号, 不需要使用该参数则传入nil
func AliYunQuerySmsStatus(c *dysmsapi20170525.Client, phoneNumber string, date string, page uint, size uint, bizId *string) (*AliYunSmsStatusQueryResponse, error) {
	return querySmsStatus(c, phoneNumber, date, page, size, bizId)
}

// AlliYunQuerySingleSmsStatus 查询单条短信的状态
// 	phoneNumber:	国内短信：11位手机号码，例如1590000****. 国际/港澳台消息：国际区号+号码，例如8520000****
// 	size:			分页查看发送记录，指定每页显示的短信记录数量。取值范围为1~50。
// 	bizId:			发送回执ID，即发送流水号, 不需要使用该参数则传入nil
func AlliYunQuerySingleSmsStatus(c *dysmsapi20170525.Client, phoneNumber string, date string, bizId *string) (*AliYunSmsStatusQueryResponse, error) {
	return AliYunQuerySmsStatus(c, phoneNumber, date, 1, 1, bizId)
}

// 查询短信信息
func querySmsStatus(c *dysmsapi20170525.Client, phoneNumber string, date string, page uint, size uint, bizId *string) (*AliYunSmsStatusQueryResponse, error) {
	if c == nil {
		return nil, errors.New("QuerySendDetails failed. dysmsapi20170525.Client == nil")
	}

	querySendDetailsRequest := &dysmsapi20170525.QuerySendDetailsRequest{
		PhoneNumber: tea.String(phoneNumber),
		SendDate:    tea.String(date),
		PageSize:    tea.Int64(int64(page)),
		CurrentPage: tea.Int64(int64(size)),
		BizId:       bizId,
	}

	raw, err := c.QuerySendDetails(querySendDetailsRequest)
	if err != nil {
		return nil, err
	}

	if raw == nil {
		return nil, errors.New("QuerySendDetails return nil")
	}

	return parseQuerySendDetailsResponseBody(raw.Body)
}
