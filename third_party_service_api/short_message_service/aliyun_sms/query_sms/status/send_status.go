package status

// AliYunSMSSendStatus 阿里云短信发送状态码
type AliYunSMSSendStatus int64

const (
	// AliYunSMSSendStatusUnknown 			阿里云短信发送状态_未知状态(官方仅定义了1/2/3三种状态, 其他status均为未知状态)
	AliYunSMSSendStatusUnknown AliYunSMSSendStatus = 0
	// AliYunSMSSendStatusWaitingForReceipt 阿里云短信发送状态_等待回执
	AliYunSMSSendStatusWaitingForReceipt AliYunSMSSendStatus = 1
	// AliYunSMSSendStatusSendFailed 		阿里云短信发送状态_发送失败
	AliYunSMSSendStatusSendFailed AliYunSMSSendStatus = 2
	// AliYunSMSSendStatusSendSuccess 		阿里云短信发送状态_发送成功
	AliYunSMSSendStatusSendSuccess AliYunSMSSendStatus = 3
)

// AliYunSMSSendStatusInterface AliYunSMSSendStatus interface
type AliYunSMSSendStatusInterface interface {
	QuerySendStatus() AliYunSMSSendStatus // 查询短信的状态
	IsSendSuccess() bool                  // 判断短信是否发送成功
	IsWaitingForReceipt() bool            // 判断短信是否处于等待回执的状态
	QuerySendStatusRawCode() int64        // 查询int64类型的原生状态码
}

// IsSendSuccess 判断短信是否发送成功
func (s AliYunSMSSendStatus) IsSendSuccess() bool {
	return s == AliYunSMSSendStatusSendSuccess
}

// IsWaitingForReceipt 判断短信是否处于等待回执的状态
func (s AliYunSMSSendStatus) IsWaitingForReceipt() bool {
	return s == AliYunSMSSendStatusWaitingForReceipt
}

// QuerySendStatusRawCode 查询int64类型的原生状态码
func (s AliYunSMSSendStatus) QuerySendStatusRawCode() int64 {
	return int64(s)
}

// QuerySendStatus 查询短信的状态
func (s AliYunSMSSendStatus) QuerySendStatus() AliYunSMSSendStatus {
	switch s {
	case AliYunSMSSendStatusWaitingForReceipt:
		fallthrough
	case AliYunSMSSendStatusSendSuccess:
		fallthrough
	case AliYunSMSSendStatusSendFailed:
		return s
	default:
		return AliYunSMSSendStatusUnknown
	}
}
