package result

import "github.com/xh-polaris/psych-idl/kitex_gen/basic"

func ResponseOk() *basic.Response {
	return &basic.Response{
		Code: 200,
		Msg:  "ok",
	}
}
