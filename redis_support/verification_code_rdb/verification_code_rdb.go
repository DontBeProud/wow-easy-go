package verification_code_rdb

import (
	"github.com/DontBeProud/wow-easy-go/redis_support/base"
	"github.com/go-redis/redis/v8"
	"time"
)

// InvalidType 违规类型
type InvalidType int

const (
	UserIsValid                        = iota + 1
	InvalidTypeUnusedCodeTooMany       // 未核销的验证码过多(频繁请求验证码但不进行验证)
	InvalidTypeRequestTooFrequently    // 请求验证码过于频繁(短时间内连续多次请求验证码)
	InvalidTypeVerifyFailTooFrequently // 验证码核销失败过于频繁
)

// VerificationCodeRdb 用于验证码相关服务的通用Rdb结构
type VerificationCodeRdb struct {
	ModuleName string                          // 业务模块名称, 不同业务对应不同的名称，防止发生不同业务的数据碰撞(部分redis-key与该字段关联)
	rDb        *redis.Client                   // redis对象
	strategy   VerificationCodeServiceStrategy // 策略
	VerificationCodeRdbInterface
}

// CreateVerificationCodeRdb 创建用于验证码服务的Rdb
func CreateVerificationCodeRdb(rdb *redis.Client, moduleName string, strategy VerificationCodeServiceStrategy) (*VerificationCodeRdb, error) {
	return createVerificationCodeRdb(rdb, moduleName, strategy)
}

type VerificationCodeRdbInterface interface {
	vcsStrategyInterface
	base.RdbBaseInterface
	PreCheckBeforeSendVerificationCode(objName string) (it InvalidType, err error)
	SetAndRegisterVerificationCode(objName string, verCode string) error
	PreCheckBeforeVerifyAndUseVerificationCode(objName string) (it InvalidType, err error)
	VerifyAndUseVerificationCode(objName string, verCode string) (exist bool, success bool, err error)
	CheckIsUnusedCodeTooMany(objName string) (bool, error)
	CheckIsRequestTooFrequently(objName string) (bool, error)
	CheckIsVerifyFailTooFrequently(objName string) (bool, error)
	QueryErrorsCountToday(objName string) (int, error)
	QueryLastErrorTime(objName string) (exist bool, lastTime time.Time, err error)
	QueryCountOfUnusedVerificationCode(objName string) (int, error)
	QueryVerificationCodeTTL(objName string) (int64, error)
	QueryVerificationCodeRegisteredPeriod(objName string) (invalid bool, period int64, err error)
}

// VerifyConnection 判断redis是否成功连接并可用(在执行关键步骤前应先调用本函数验证redis是否可用，避免无谓的资源消耗，包括但不限于验证码发送费用、服务端资源等)
func (r VerificationCodeRdb) VerifyConnection() (bool, error) {
	return base.VerifyConnection(r.rDb)
}

// PreCheckBeforeSendVerificationCode 发送验证码前的校验(组合校验用户当前状态是否合法)
// 校验请求是否过于频繁、验证错误次数是否过多、未核销的验证码是否过多(是否频繁请求验证码但不进行验证)
func (r VerificationCodeRdb) PreCheckBeforeSendVerificationCode(objName string) (it InvalidType, err error) {
	return r.preCheckBeforeSendVerificationCode(objName)
}

// SetAndRegisterVerificationCode 添加并记录验证码(添加该用户的验证码缓存，并且向该用户未核销的验证码集合中添加该验证码)
func (r VerificationCodeRdb) SetAndRegisterVerificationCode(objName string, verCode string) error {
	return r.setAndRegisterVerificationCode(objName, verCode, time.Duration(r.strategy.ValidityDuration)*time.Second)
}

// PreCheckBeforeVerifyAndUseVerificationCode 核销验证码前的校验(组合校验用户当前状态是否合法)
// 校验验证错误次数是否过多、未核销的验证码是否过多(是否频繁请求验证码但不进行验证)
func (r VerificationCodeRdb) PreCheckBeforeVerifyAndUseVerificationCode(objName string) (it InvalidType, err error) {
	return r.preCheckBeforeVerifyAndUseVerificationCode(objName)
}

// VerifyAndUseVerificationCode 核销验证码
func (r VerificationCodeRdb) VerifyAndUseVerificationCode(objName string, verCode string) (exist bool, success bool, err error) {
	return r.verifyAndUseVerificationCode(objName, verCode)
}

// CheckIsUnusedCodeTooMany 判断当日未使用的验证码是否过多(用于防止恶意刷接口) threshold: 阈值
// 一般在请求验证码和核销验证码前调用判断
func (r VerificationCodeRdb) CheckIsUnusedCodeTooMany(objName string) (bool, error) {
	return r.checkIsUnusedCodeTooMany(objName, r.strategy.DenyThresholdOfUnusedCode)
}

// CheckIsRequestTooFrequently 判断申请验证码是否过于频繁, threshold: 阈值(单位为秒)
// 若上一次请求的验证码尚未被核销，且当前时间距离上次请求的时间差小于等于阈值，则返回true.
// 一般在请求验证码前调用判断
func (r VerificationCodeRdb) CheckIsRequestTooFrequently(objName string) (bool, error) {
	return r.checkIsRequestTooFrequently(objName, r.strategy.RequestTimeIntervalThreshold)
}

// CheckIsVerifyFailTooFrequently 判断用户是否验证错误过于频繁
func (r VerificationCodeRdb) CheckIsVerifyFailTooFrequently(objName string) (bool, error) {
	return r.checkIsVerifyFailTooFrequently(objName, r.strategy.DenyThresholdOfFailedCount, r.strategy.TemporarilyBanStrategy)
}

// QueryErrorsCountToday 查询用户当日失败的次数
func (r VerificationCodeRdb) QueryErrorsCountToday(objName string) (int, error) {
	return r.queryErrorsCountToday(objName)
}

// QueryLastErrorTime 查询用户最后一次验证失败的时间
func (r VerificationCodeRdb) QueryLastErrorTime(objName string) (exist bool, lastTime time.Time, err error) {
	return r.queryLastErrorTime(objName)
}

// QueryCountOfUnusedVerificationCode 查询用户当日未核销的验证码数量
func (r VerificationCodeRdb) QueryCountOfUnusedVerificationCode(objName string) (int, error) {
	return r.queryCountOfUnusedVerificationCode(objName)
}

// QueryVerificationCodeTTL 查询验证码剩余的有效时长(单位为秒)
func (r VerificationCodeRdb) QueryVerificationCodeTTL(objName string) (int64, error) {
	return r.queryVerificationCodeTTL(objName)
}

// QueryVerificationCodeRegisteredPeriod 获取验证码已等待核销的时长
// invalid: 验证码是否已失效. 若invalid为true，则说明验证码已失效, period的大小无意义
// period: 验证码已等待核销的时长, 单位为秒
func (r VerificationCodeRdb) QueryVerificationCodeRegisteredPeriod(objName string) (invalid bool, period int64, err error) {
	return r.queryVerificationCodeRegisteredPeriod(objName)
}

// QueryValidityDuration 查询验证码的默认有效期
func (r VerificationCodeRdb) QueryValidityDuration() int64 {
	return r.strategy.QueryValidityDuration()
}

// QueryRequestTimeIntervalThreshold 查询验证码请求间隔阈值
func (r VerificationCodeRdb) QueryRequestTimeIntervalThreshold() int64 {
	return r.strategy.QueryRequestTimeIntervalThreshold()
}

// QueryDenyThresholdOfFailedCount 查询因失败次数过多而拒绝请求的数量阈值
func (r VerificationCodeRdb) QueryDenyThresholdOfFailedCount() int {
	return r.strategy.QueryDenyThresholdOfFailedCount()
}

// QueryDenyThresholdOfUnusedCode 查询因未核销的验证码数量过多而拒绝请求的数量阈值
func (r VerificationCodeRdb) QueryDenyThresholdOfUnusedCode() int {
	return r.strategy.QueryDenyThresholdOfUnusedCode()
}

// QueryTemporarilyBanStrategy 查询临时封禁策略
func (r *VerificationCodeRdb) QueryTemporarilyBanStrategy() *map[int]int64 {
	return r.strategy.QueryTemporarilyBanStrategy()
}

// AddTemporarilyBanStrategy 添加临时封禁策略
func (r *VerificationCodeRdb) AddTemporarilyBanStrategy(threshold int, duration int64) {
	r.strategy.AddTemporarilyBanStrategy(threshold, duration)
}

// DelTemporarilyBanStrategy 删除临时封禁策略
func (r *VerificationCodeRdb) DelTemporarilyBanStrategy(threshold int) {
	r.strategy.DelTemporarilyBanStrategy(threshold)
}

// ModifyValidityDuration 修改临时封禁策略
func (r *VerificationCodeRdb) ModifyValidityDuration(duration int64) error {
	return r.strategy.ModifyValidityDuration(duration)
}

// ModifyRequestTimeIntervalThreshold 修改验证码请求间隔阈值
func (r *VerificationCodeRdb) ModifyRequestTimeIntervalThreshold(intervalThreshold int64) {
	r.strategy.ModifyRequestTimeIntervalThreshold(intervalThreshold)
}

// ModifyDenyThresholdOfFailedCount 修改因失败次数过多而拒绝请求的数量阈值
func (r *VerificationCodeRdb) ModifyDenyThresholdOfFailedCount(threshold int) {
	r.strategy.ModifyDenyThresholdOfFailedCount(threshold)
}

// ModifyDenyThresholdOfUnusedCode 修改因未核销的验证码数量过多而拒绝请求的数量阈值
func (r *VerificationCodeRdb) ModifyDenyThresholdOfUnusedCode(threshold int) {
	r.strategy.ModifyDenyThresholdOfUnusedCode(threshold)
}
