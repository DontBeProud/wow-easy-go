package status

// AliYunSMSOperatorStatus 阿里云短信服务运营商状态码
// 详细内容参照 https://help.aliyun.com/document_detail/102352.html
type AliYunSMSOperatorStatus string

const (
	// AliYunSMSOperatorStatusCodeSuccess   阿里云短信服务运营商状态码_发送成功
	AliYunSMSOperatorStatusCodeSuccess AliYunSMSOperatorStatus = "DELIVERED"
)

type AliYunSMSOperatorStatusInterface interface {
	QueryOperatorStatusErrorInfo() *string // 查询运营商状态码包含的错误信息, 若无错误则返回nil
	QueryOperatorStatusCode() string       // 查询运营商状态码
}

// QueryOperatorStatusCode 查询运营商状态码
func (s AliYunSMSOperatorStatus) QueryOperatorStatusCode() string {
	return string(s)
}

// QueryOperatorErrorInfo 查询运营商端上报的错误信息, 若无错误则返回nil
func (s AliYunSMSOperatorStatus) QueryOperatorErrorInfo() *string {
	if s == AliYunSMSOperatorStatusCodeSuccess {
		return nil
	}
	e := string(s)
	return &e
}
