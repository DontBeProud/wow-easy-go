package verification_code_rdb

import (
	"context"
	"errors"
	"github.com/DontBeProud/wow-easy-go/redis_support/base"
	"github.com/DontBeProud/wow-easy-go/utils/wow_time"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

// 添加并记录验证码(添加该用户的验证码缓存，并且向该用户未核销的验证码集合中添加该验证码)
func (r VerificationCodeRdb) setAndRegisterVerificationCode(objName string, verCode string, expireNanoDuration time.Duration) error {
	if err := r.setVerificationCode(objName, verCode, expireNanoDuration); err != nil {
		return err
	}
	r.addUnusedVerificationCode(objName, verCode)
	return nil
}

// 核销验证码
func (r VerificationCodeRdb) verifyAndUseVerificationCode(objName string, verCode string) (exist bool, success bool, err error) {
	exist, code, err := r.getVerificationCode(objName)
	if err != nil || !exist {
		// 执行出错或验证码不存在
		return exist, false, err
	}

	success = code == verCode
	if success {
		r.verifySuccess(objName, verCode)
	} else {
		r.verifyFail(objName)
	}

	return exist, success, nil
}

// 发送验证码前的校验(组合校验用户当前状态是否合法)
// 校验请求是否过于频繁、验证错误次数是否过多、未核销的验证码是否过多(是否频繁请求验证码但不进行验证)
func (r VerificationCodeRdb) preCheckBeforeSendVerificationCode(objName string) (it InvalidType, err error) {
	return r.combineCheckIsUserValid(objName, map[InvalidType]func(string) (bool, error){
		InvalidTypeRequestTooFrequently:    r.CheckIsRequestTooFrequently,
		InvalidTypeVerifyFailTooFrequently: r.CheckIsVerifyFailTooFrequently,
		InvalidTypeUnusedCodeTooMany:       r.CheckIsUnusedCodeTooMany,
	})
}

// 核销验证码前的校验(组合校验用户当前状态是否合法)
// 校验验证错误次数是否过多、未核销的验证码是否过多(是否频繁请求验证码但不进行验证)
func (r VerificationCodeRdb) preCheckBeforeVerifyAndUseVerificationCode(objName string) (it InvalidType, err error) {
	return r.combineCheckIsUserValid(objName, map[InvalidType]func(string) (bool, error){
		InvalidTypeVerifyFailTooFrequently: r.CheckIsVerifyFailTooFrequently,
		InvalidTypeUnusedCodeTooMany:       r.CheckIsUnusedCodeTooMany,
	})
}

// 组合校验用户当前状态是否合法
// 支持传入	CheckIsRequestTooFrequently/CheckIsVerifyFailTooFrequently/CheckIsUnusedCodeTooMany
func (r VerificationCodeRdb) combineCheckIsUserValid(objName string, fnList map[InvalidType]func(string) (bool, error)) (it InvalidType, err error) {

	type fnRes struct {
		invalid bool
		err     error
		it      InvalidType
	}
	resChan := make(chan fnRes)

	fnGo := func(it InvalidType, fn func(string) (bool, error)) {
		iv, er := fn(objName)
		resChan <- fnRes{
			invalid: iv,
			err:     er,
			it:      it,
		}
	}

	for it, fn := range fnList {
		go fnGo(it, fn)
	}

	for i, l := 0, len(fnList); i < l; i++ {
		res := <-resChan
		if res.err != nil || res.invalid {
			return res.it, err
		}
	}
	return UserIsValid, nil
}

// 设置验证码
func (r VerificationCodeRdb) setVerificationCode(objName string, verCode string, expireNanoDuration time.Duration) error {
	_, err := r.rDb.Set(context.TODO(), r.getRedisFieldNameVerificationCode(objName), verCode, expireNanoDuration).Result()
	return err
}

// 查询用户的验证码 exist: 验证码是否存在 code: 验证码内容
func (r VerificationCodeRdb) getVerificationCode(objName string) (exist bool, code string, err error) {
	result, err := r.rDb.Get(context.TODO(), r.getRedisFieldNameVerificationCode(objName)).Result()
	if err == redis.Nil {
		return false, result, nil
	}
	return err == nil, result, err
}

// 核销验证码失败(失败次数+1，更新最后一次失败的时间)
func (r VerificationCodeRdb) verifyFail(objName string) {
	r.increaseErrorCount(objName)
	r.updateLastErrorTime(objName)
}

// 核销验证码成功(删除该用户的验证码缓存，并且从该用户未核销的验证码集合中删除该验证码)
func (r VerificationCodeRdb) verifySuccess(objName string, verCode string) {
	r.rDb.Del(context.TODO(), r.getRedisFieldNameVerificationCode(objName))
	r.rDb.SRem(context.TODO(), r.getRedisFieldNameVerificationCodeSet(objName), verCode)
}

// 将验证码加入到该用户当日待核销的验证码集合中
func (r VerificationCodeRdb) addUnusedVerificationCode(objName string, verCode string) {
	f := r.getRedisFieldNameVerificationCodeSet(objName)
	r.rDb.SAdd(context.TODO(), f, verCode)
	r.rDb.ExpireAt(context.TODO(), f, wow_time.GetTomorrowZeroTime()) // 设置有效期到第二天的零时
}

// 更新该用户当日最后一次验证错误的时间
func (r VerificationCodeRdb) updateLastErrorTime(objName string) {
	fl := r.getRedisFieldNameVerificationCodeLastFailedTime(objName)
	r.rDb.Set(context.TODO(), fl, time.Now().Unix(), time.Duration(wow_time.GetTomorrowZeroTime().UnixNano()-time.Now().UnixNano())) // 设置有效期到第二天的零时
}

// 该用户当日累计错误次数 +1
func (r VerificationCodeRdb) increaseErrorCount(objName string) {
	fc := r.getRedisFieldNameVerificationCodeErrorCount(objName)
	r.rDb.Incr(context.TODO(), fc)
	r.rDb.ExpireAt(context.TODO(), fc, wow_time.GetTomorrowZeroTime()) // 设置有效期到第二天的零时
}

// 判断申请验证码是否过于频繁, 若上一次请求的验证码尚未被核销，且生命周期尚未结束，则返回true. threshold: 阈值(单位为秒)
func (r VerificationCodeRdb) checkIsRequestTooFrequently(objName string, threshold int64) (bool, error) {
	ttl, err := r.queryVerificationCodeTTL(objName)
	return err == nil && ttl > 0 && (r.strategy.ValidityDuration-ttl) <= threshold, err
}

// 判断当日未使用的验证码是否过多(用于防止恶意刷接口) threshold: 阈值
func (r VerificationCodeRdb) checkIsUnusedCodeTooMany(objName string, threshold int) (bool, error) {
	cnt, err := r.queryCountOfUnusedVerificationCode(objName)
	if err != nil {
		return false, err
	}

	// 低于阈值
	if cnt < threshold {
		return false, nil
	}

	// 高于阈值
	if cnt > threshold && threshold > 0 {
		return true, nil
	}

	// cnt == threshold 判断是否仍有未使用的验证码
	exist, _, err := r.getVerificationCode(objName)
	return exist, err
}

// 获取验证码的剩余有效时长
func (r VerificationCodeRdb) queryVerificationCodeTTL(objName string) (int64, error) {
	result, err := r.rDb.TTL(context.TODO(), r.getRedisFieldNameVerificationCode(objName)).Result()
	if err == redis.Nil {
		return 0, nil
	}
	return int64(result.Seconds()), err
}

// 获取验证码已等待核销的时长
func (r VerificationCodeRdb) queryVerificationCodeRegisteredPeriod(objName string) (bool, int64, error) {
	ttl, err := r.queryVerificationCodeTTL(objName)
	return err != nil && ttl == 0, r.strategy.ValidityDuration - ttl, err
}

// 获取该用户最后一次验证错误的时间
func (r VerificationCodeRdb) queryLastErrorTime(objName string) (exist bool, lastTime time.Time, err error) {
	tm, err := r.rDb.Get(context.TODO(), r.getRedisFieldNameVerificationCodeLastFailedTime(objName)).Int64()
	if err == redis.Nil {
		return false, time.Time{}, nil
	}
	return err == nil, time.Unix(tm, 0), err
}

// 判断用户是否验证错误过于频繁
func (r VerificationCodeRdb) checkIsVerifyFailTooFrequently(objName string, threshold int, temporarilyBanStrategy *sync.Map) (bool, error) {
	// 获取当日失败次数
	cnt, err := r.queryErrorsCountToday(objName)
	if err != nil || cnt == 0 {
		// 失败次数为0，则说明当日无失败记录，也就无需根据失败次数和最后一次失败时间来判断失败频率
		return false, err
	}

	// 高于失败次数阈值
	if threshold > 0 && cnt >= threshold {
		return true, nil
	}

	// 无暂时封禁策略
	if temporarilyBanStrategy == nil {
		return false, nil
	}

	exist, lastErrTime, err := r.queryLastErrorTime(objName)
	if err != nil || !exist {
		return false, err
	}

	invalid := false
	temporarilyBanStrategy.Range(func(t, banDuration interface{}) bool {
		// 错误次数高于判定阈值，同时当前时间距离最后一次验证错误的时间差小于设定的时间范围，则判定当前时刻仍处于封禁状态
		invalid = cnt >= t.(int) && (time.Now().Unix()-lastErrTime.Unix()) <= banDuration.(int64)
		return invalid == false
	})
	return invalid, nil
}

// 查询该用户当日未核销成功的验证码数量
func (r VerificationCodeRdb) queryCountOfUnusedVerificationCode(objName string) (int, error) {
	f := r.getRedisFieldNameVerificationCodeSet(objName)
	cnt, err := r.rDb.SCard(context.TODO(), f).Uint64()
	// 无待核销的验证码
	if err == redis.Nil {
		return 0, nil
	}
	return int(cnt), err
}

// 统计该用户当日验证错误的次数
func (r VerificationCodeRdb) queryErrorsCountToday(objName string) (int, error) {
	fc := r.getRedisFieldNameVerificationCodeErrorCount(objName)
	cnt, err := r.rDb.Get(context.TODO(), fc).Int()
	// 无错误记录
	if err == redis.Nil {
		return 0, nil
	}
	return cnt, err
}

// 根据对象名称生成存储验证码的字段名称
func (r VerificationCodeRdb) getRedisFieldNameVerificationCode(objName string) string {
	return r.ModuleName + "VerificationCode" + objName
}

// 根据对象名称生成存储该手机当日待核销的验证码的集合的字段名称
func (r VerificationCodeRdb) getRedisFieldNameVerificationCodeSet(objName string) string {
	return r.ModuleName + "VerificationCodeSet" + objName + time.Now().Format("20060102")
}

// 根据对象名称生成存储该手机当日验证错误的次数
func (r VerificationCodeRdb) getRedisFieldNameVerificationCodeErrorCount(objName string) string {
	return r.ModuleName + "VerificationCodeErrorCount" + objName + time.Now().Format("20060102")
}

// 根据对象名称生成存储该手机当日最后一次验证错误的时间
func (r VerificationCodeRdb) getRedisFieldNameVerificationCodeLastFailedTime(objName string) string {
	return r.ModuleName + "VerificationCodeLastErrorTime" + objName + time.Now().Format("20060102")
}

func createVerificationCodeRdb(rdb *redis.Client, moduleName string, strategy VerificationCodeServiceStrategy) (*VerificationCodeRdb, error) {
	if rdb == nil {
		return nil, errors.New("rdb == nil")
	}

	if moduleName == "" {
		return nil, errors.New("ModuleName == \"\"")
	}

	// 测试redis是否可用
	if _, err := base.VerifyConnection(rdb); err != nil {
		return nil, err
	}

	return &VerificationCodeRdb{
		ModuleName: moduleName,
		rDb:        rdb,
		strategy:   strategy,
	}, nil
}
