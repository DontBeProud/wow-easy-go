package wow_random

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// GenerateVerificationCode 生成随机的验证码
// length: 位数
func GenerateVerificationCode(length uint) string {
	var s strings.Builder
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < int(length); i++ {
		_, _ = fmt.Fprintf(&s, "%d", rand.Intn(10))
	}
	return s.String()
}
