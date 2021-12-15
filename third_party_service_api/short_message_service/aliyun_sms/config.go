package aliyun_sms

const (
	// AliYunSMSEndPointAddr 固定值. 短信服务统一使用以下公网服务地址作为Endpoint. 参见阿里云文档(https://help.aliyun.com/document_detail/101511.html)
	AliYunSMSEndPointAddr = "dysmsapi.aliyuncs.com"
)

// AliYunSMSOptionalConfig 阿里云短信服务可选配置项, 具体含义参照aliOpenApi.Config的定义或在线文档
type AliYunSMSOptionalConfig struct {
	SecurityToken        *string // security token
	Protocol             *string // http protocol
	RegionId             *string // region id
	HttpProxy            *string // http proxy
	HttpsProxy           *string // https proxy
	Socks5Proxy          *string // socks5 proxy
	Socks5NetWork        *string // socks5 network
	NoProxy              *string // proxy white list
	Network              *string // network for endpoint
	UserAgent            *string // user agent
	Suffix               *string // suffix for endpoint
	EndpointType         *string // endpoint type
	OpenPlatformEndpoint *string // OpenPlatform endpoint
	SignatureAlgorithm   *string // Signature Algorithm
	ReadTimeout          *int    // read timeout
	ConnectTimeout       *int    // connect timeout
	MaxIdleConns         *int    // max idle conns
}
