package wow_regexp

import (
	"testing"
)

func TestDefaultExprEmail(t *testing.T) {
	exp, err := DefaultExprEmail.CreateRegExp()
	if err != nil {
		t.Error(err.Error())
	}
	if !exp.MatchString("12@163.com") {
		t.Error("Find bug in DefaultExprEmail")
	}
	if exp.MatchString("12#163.com") {
		t.Error("Find bug in DefaultExprEmail")
	}
}

func TestDefaultExprIpv4Addr(t *testing.T) {
	exp, err := DefaultExprIpv4Addr.CreateRegExp()
	if err != nil {
		t.Error(err.Error())
	}
	if !exp.MatchString("127.2.2.2") {
		t.Error("Find bug in DefaultExprIpv4Addr")
	}
	if exp.MatchString("225.225.256.256") {
		t.Error("Find bug in DefaultExprIpv4Addr")
	}
}

func TestDefaultExprChineseMobilePhone(t *testing.T) {
	exp, err := DefaultExprChineseMobilePhone.CreateRegExp()
	if err != nil {
		t.Error(err.Error())
	}
	if !exp.MatchString("15855533311") {
		t.Error("Find bug in DefaultExprChineseMobilePhone")
	}
	if exp.MatchString("666") {
		t.Error("Find bug in DefaultExprChineseMobilePhone")
	}
}
