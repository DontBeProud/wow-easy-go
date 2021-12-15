package aliyun_sms

import (
	"errors"
	aliOpenApi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

// 创建用户对象
func createAliYunSMSClient(credential AliYunSMSClientCredential, opt *AliYunSMSOptionalConfig) (*AliYunSMSClient, error) {

	// base config
	cfg := &aliOpenApi.Config{
		AccessKeyId:     tea.String(credential.AccessKeyId),
		AccessKeySecret: tea.String(credential.AccessKeySecret),
		Endpoint:        tea.String(AliYunSMSEndPointAddr),
	}

	// optional config
	if opt != nil {
		cfg.SecurityToken = opt.SecurityToken
		cfg.Protocol = opt.Protocol
		cfg.RegionId = opt.RegionId
		cfg.ReadTimeout = opt.ReadTimeout
		cfg.ConnectTimeout = opt.ConnectTimeout
		cfg.HttpProxy = opt.HttpProxy
		cfg.HttpsProxy = opt.HttpsProxy
		cfg.NoProxy = opt.NoProxy
		cfg.MaxIdleConns = opt.MaxIdleConns
		cfg.Network = opt.Network
		cfg.UserAgent = opt.UserAgent
		cfg.Suffix = opt.Suffix
		cfg.Socks5Proxy = opt.Socks5Proxy
		cfg.Socks5NetWork = opt.Socks5NetWork
		cfg.EndpointType = opt.EndpointType
		cfg.OpenPlatformEndpoint = opt.OpenPlatformEndpoint
		cfg.SignatureAlgorithm = opt.SignatureAlgorithm
	}

	c, err := dysmsapi20170525.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &AliYunSMSClient{
		config: cfg,
		client: c,
	}, nil
}

// 创建可发送短信的用户对象
func createAliYunSMSClientSender(credential AliYunSMSClientCredential, signName string, template SmsTemplateInterface, opt *AliYunSMSOptionalConfig) (*AliYunSMSClientSender, error) {
	client, err := createAliYunSMSClient(credential, opt)
	if err != nil {
		return nil, err
	}

	return client.createAliYunSMSClientSender(signName, template)
}

// 创建可发送短信的用户对象
func (c AliYunSMSClient) createAliYunSMSClientSender(signName string, template SmsTemplateInterface) (*AliYunSMSClientSender, error) {
	if len(signName) == 0 || len(template.GetTemplateId()) == 0 {
		return nil, errors.New("CreateAliYunSMSClient failed. signName and templateName can not be ''")
	}

	return &AliYunSMSClientSender{
		signName:                 signName,
		template:                 &template,
		AliYunSMSClientInterface: &c,
	}, nil
}
