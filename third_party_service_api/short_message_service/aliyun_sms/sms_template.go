package aliyun_sms

// SmsTemplateInterface 短信模板接口
type SmsTemplateInterface interface {
	GetTemplateId() string                            // 获取模板ID
	GenerateTemplateParam(params interface{}) *string // 根据传入的结构体, 生成短信模板所需的参数
}
