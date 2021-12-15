package verification_code_rdb

import (
	"errors"
	"sync"
)

// VerificationCodeServiceStrategy 验证码服务策略
type VerificationCodeServiceStrategy struct {
	vcsStrategyInterface
	ValidityDuration             int64     // 验证码有效期时长(秒), 必须大于0
	RequestTimeIntervalThreshold int64     // 验证码请求间隔时长限制(秒)，即两次获取验证码的时间差值下限. 不需要该项限制则填0
	DenyThresholdOfUnusedCode    int       // 单日未核销的验证码数量阈值，超过该值后禁止手机号使用短信验证码业务. 不需要该项限制则填0
	DenyThresholdOfFailedCount   int       // 单日验证错误次数阈值，高于该值后禁止手机号使用短信验证码业务. 不需要该项限制则填0
	TemporarilyBanStrategy       *sync.Map // 短暂禁止手机号使用短信验证码业务的策略，key:失败次数的阈值;value:禁止时长(秒), 详见CheckIsVerifyFailTooFrequently
}

func CreateVerificationCodeServiceStrategy(duration int64, intervalThreshold int64, unusedThreshold int,
	failThreshold int, tempBanStrategy *map[int]int64) (*VerificationCodeServiceStrategy, error) {
	if duration <= 0 {
		return nil, errors.New("CreateVerificationCodeServiceStrategy ValidityDuration == 0")
	}

	res := VerificationCodeServiceStrategy{
		ValidityDuration:             duration,
		RequestTimeIntervalThreshold: intervalThreshold,
		DenyThresholdOfUnusedCode:    unusedThreshold,
		DenyThresholdOfFailedCount:   failThreshold,
		TemporarilyBanStrategy:       &sync.Map{},
	}

	if tempBanStrategy != nil {
		for threshold, duration := range *tempBanStrategy {
			res.TemporarilyBanStrategy.Store(threshold, duration)
		}
	}

	return &res, nil
}

type vcsStrategyInterface interface {
	QueryValidityDuration() int64
	QueryRequestTimeIntervalThreshold() int64
	QueryDenyThresholdOfFailedCount() int
	QueryDenyThresholdOfUnusedCode() int
	QueryTemporarilyBanStrategy() *map[int]int64
	AddTemporarilyBanStrategy(threshold int, duration int64)
	DelTemporarilyBanStrategy(threshold int)
	ModifyValidityDuration(duration int64) error
	ModifyRequestTimeIntervalThreshold(intervalThreshold int64)
	ModifyDenyThresholdOfFailedCount(threshold int)
	ModifyDenyThresholdOfUnusedCode(threshold int)
}

func (s VerificationCodeServiceStrategy) QueryValidityDuration() int64 {
	return s.ValidityDuration
}

func (s VerificationCodeServiceStrategy) QueryRequestTimeIntervalThreshold() int64 {
	return s.RequestTimeIntervalThreshold
}

func (s VerificationCodeServiceStrategy) QueryDenyThresholdOfFailedCount() int {
	return s.DenyThresholdOfFailedCount
}

func (s VerificationCodeServiceStrategy) QueryDenyThresholdOfUnusedCode() int {
	return s.DenyThresholdOfUnusedCode
}

func (s VerificationCodeServiceStrategy) QueryTemporarilyBanStrategy() *map[int]int64 {
	result := make(map[int]int64)

	s.TemporarilyBanStrategy.Range(func(threshold, banDuration interface{}) bool {
		result[threshold.(int)] = banDuration.(int64)
		return true
	})

	return &result
}

// AddTemporarilyBanStrategy 添加短期禁止策略
func (s *VerificationCodeServiceStrategy) AddTemporarilyBanStrategy(threshold int, duration int64) {
	s.TemporarilyBanStrategy.Store(threshold, duration)
}

// DelTemporarilyBanStrategy 删除短期禁止策略
func (s *VerificationCodeServiceStrategy) DelTemporarilyBanStrategy(threshold int) {
	s.TemporarilyBanStrategy.Delete(threshold)
}

// ModifyDenyThresholdOfFailedCount 修改单日验证错误次数阈值
func (s *VerificationCodeServiceStrategy) ModifyDenyThresholdOfFailedCount(threshold int) {
	s.DenyThresholdOfFailedCount = threshold
}

// ModifyDenyThresholdOfUnusedCode 修改单日未核销的验证码数量阈值
func (s *VerificationCodeServiceStrategy) ModifyDenyThresholdOfUnusedCode(threshold int) {
	s.DenyThresholdOfUnusedCode = threshold
}

// ModifyValidityDuration 修改验证码有效期时长
func (s *VerificationCodeServiceStrategy) ModifyValidityDuration(duration int64) error {
	if duration <= 0 {
		return errors.New("CreateVerificationCodeServiceStrategy ValidityDuration == 0")
	}
	s.ValidityDuration = duration
	return nil
}

// ModifyRequestTimeIntervalThreshold 修改验证码请求间隔时长限制
func (s *VerificationCodeServiceStrategy) ModifyRequestTimeIntervalThreshold(intervalThreshold int64) {
	s.RequestTimeIntervalThreshold = intervalThreshold
}
