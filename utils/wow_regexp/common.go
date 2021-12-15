package wow_regexp

type DefaultRegExpr string

func (exp DefaultRegExpr) CreateRegExp() (*RegExp, error) {
	return CreateRegExp(string(exp))
}

const (
	DefaultExprEmail              DefaultRegExpr = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	DefaultExprChineseMobilePhone DefaultRegExpr = "^(13[0-9]|14[5|7]|15[0|1|2|3|4|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9])\\d{8}$"
	DefaultExprIpv4Addr           DefaultRegExpr = "((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}"
)
