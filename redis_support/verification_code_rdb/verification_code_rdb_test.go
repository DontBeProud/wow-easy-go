package verification_code_rdb

import (
	"context"
	"github.com/go-redis/redis/v8"
	"testing"
)

const (
	redisCon     = "localhost:6379"
	redisPsw     = ""
	redisDb      = 1
	testPhoneNum = "TestPhoneNumber001"
	testVerCode  = "TestTest"
)

var (
	r = redis.NewClient(&redis.Options{
		Addr:     redisCon,
		Password: redisPsw,
		DB:       redisDb,
	})

	strategy, _ = CreateVerificationCodeServiceStrategy(300,
		60,
		5,
		10,
		&map[int]int64{
			3: 40,
			5: 120,
		})

	rdb, _ = CreateVerificationCodeRdb(r, "SMS", *strategy)
)

func TestCreateRdbSmsVerificationCode(t *testing.T) {
	_, err := CreateVerificationCodeRdb(r, "SMS", *strategy)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestPreCheckBeforeSendVerificationCode(t *testing.T) {
	code, _ := rdb.PreCheckBeforeSendVerificationCode(testPhoneNum)
	println(code)
	_ = rdb.SetAndRegisterVerificationCode(testPhoneNum, testVerCode)
	_ = rdb.SetAndRegisterVerificationCode(testPhoneNum, testVerCode+"1")
	_ = rdb.SetAndRegisterVerificationCode(testPhoneNum, testVerCode+"2")
	_ = rdb.SetAndRegisterVerificationCode(testPhoneNum, testVerCode+"3")
	code, _ = rdb.PreCheckBeforeSendVerificationCode(testPhoneNum)
	println(code)
	_ = rdb.SetAndRegisterVerificationCode(testPhoneNum, testVerCode+"4")
	code, _ = rdb.PreCheckBeforeSendVerificationCode(testPhoneNum)
	println(code)
	clear(rdb)
}

func TestPreCheckBeforeVerifyAndUseVerificationCode(t *testing.T) {
	code, _ := rdb.PreCheckBeforeVerifyAndUseVerificationCode(testPhoneNum)
	println(code)
	_ = rdb.SetAndRegisterVerificationCode(testPhoneNum, testVerCode)
	_, _, _ = rdb.VerifyAndUseVerificationCode(testPhoneNum, testVerCode+"1")
	code, _ = rdb.PreCheckBeforeVerifyAndUseVerificationCode(testPhoneNum)
	println(code)
	_ = rdb.SetAndRegisterVerificationCode(testPhoneNum, testVerCode)
	_, _, _ = rdb.VerifyAndUseVerificationCode(testPhoneNum, testVerCode+"2")
	code, _ = rdb.PreCheckBeforeVerifyAndUseVerificationCode(testPhoneNum)
	println(code)
	_ = rdb.SetAndRegisterVerificationCode(testPhoneNum, testVerCode)
	_, _, _ = rdb.VerifyAndUseVerificationCode(testPhoneNum, testVerCode+"3")
	code, _ = rdb.PreCheckBeforeVerifyAndUseVerificationCode(testPhoneNum)
	println(code)
	_ = rdb.SetAndRegisterVerificationCode(testPhoneNum, testVerCode)
	_, _, _ = rdb.VerifyAndUseVerificationCode(testPhoneNum, testVerCode+"4")
	code, _ = rdb.PreCheckBeforeVerifyAndUseVerificationCode(testPhoneNum)
	println(code)
	_ = rdb.SetAndRegisterVerificationCode(testPhoneNum, testVerCode)
	_, _, _ = rdb.VerifyAndUseVerificationCode(testPhoneNum, testVerCode+"5")
	code, _ = rdb.PreCheckBeforeVerifyAndUseVerificationCode(testPhoneNum)
	println(code)
	clear(rdb)
}

func TestUseVerificationCode(t *testing.T) {
	err := rdb.SetAndRegisterVerificationCode(testPhoneNum, testVerCode)
	if err != nil {
		t.Error(err.Error())
	}

	_, success, err := rdb.VerifyAndUseVerificationCode(testPhoneNum, testVerCode)
	if err != nil {
		t.Error(err.Error())
	}
	if !success {
		t.Error("核销失败")
	}

	_, success, err = rdb.VerifyAndUseVerificationCode(testPhoneNum, testVerCode+"1")
	if err != nil {
		println(err.Error())
	}
	clear(rdb)
}

func TestVerifyError(t *testing.T) {
	err := rdb.SetAndRegisterVerificationCode(testPhoneNum, testVerCode)
	if err != nil {
		t.Error(err.Error())
	}
	b, _ := rdb.CheckIsRequestTooFrequently(testPhoneNum)
	if !b {
		t.Error("频率判定模块有bug")
	}

	_, _, _ = rdb.VerifyAndUseVerificationCode(testPhoneNum, testVerCode+"fake")
	_, _, _ = rdb.VerifyAndUseVerificationCode(testPhoneNum, testVerCode+"fake")

	b, _ = rdb.CheckIsVerifyFailTooFrequently(testPhoneNum)
	println(b)

	_, _, _ = rdb.VerifyAndUseVerificationCode(testPhoneNum, testVerCode+"fake")
	b, _ = rdb.CheckIsVerifyFailTooFrequently(testPhoneNum)
	println(b)

	b, tm, err := rdb.queryLastErrorTime(testPhoneNum)
	if err != nil {
		println(err.Error())
	}
	if !b {
		t.Error("错误查询模块有bug")
	}
	println(tm.Unix())
	clear(rdb)
}

func TestUncheckedVerCode(t *testing.T) {

	_ = rdb.SetAndRegisterVerificationCode(testPhoneNum, testVerCode)
	_ = rdb.SetAndRegisterVerificationCode(testPhoneNum, testVerCode+"1")
	_ = rdb.SetAndRegisterVerificationCode(testPhoneNum, testVerCode+"2")
	_ = rdb.SetAndRegisterVerificationCode(testPhoneNum, testVerCode+"3")
	b, _ := rdb.CheckIsUnusedCodeTooMany(testPhoneNum)
	if b {
		t.Error("未核销验证码数量统计模块有bug")
	}

	_ = rdb.SetAndRegisterVerificationCode(testPhoneNum, testVerCode+"4")
	b, _ = rdb.CheckIsUnusedCodeTooMany(testPhoneNum)
	if !b {
		t.Error("未核销验证码数量统计模块有bug")
	}
	clear(rdb)
}

func clear(r *VerificationCodeRdb) {
	r.rDb.Del(context.TODO(), r.getRedisFieldNameVerificationCodeErrorCount(testPhoneNum))
	r.rDb.Del(context.TODO(), r.getRedisFieldNameVerificationCodeLastFailedTime(testPhoneNum))
	r.rDb.Del(context.TODO(), r.getRedisFieldNameVerificationCodeSet(testPhoneNum))
	r.rDb.Del(context.TODO(), r.getRedisFieldNameVerificationCode(testPhoneNum))
}
