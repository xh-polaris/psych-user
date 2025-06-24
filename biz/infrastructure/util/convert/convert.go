package convert

import (
	u "github.com/xh-polaris/psych-idl/kitex_gen/user"
	"github.com/xh-polaris/psych-user/biz/infrastructure/consts"
	"github.com/xh-polaris/psych-user/biz/infrastructure/mapper/common"
	"github.com/xh-polaris/psych-user/biz/infrastructure/mapper/unit"
)

func OptionGen2Loc(options *u.Form) *common.Form {
	option := &common.Form{}
	if options != nil {
		// 不空，则提取options中的kv对，插入option（要解耦数据定义）
		objs := options.GetOptions()
		var tmp []*common.Obj

		for _, obj := range objs {
			kv := &common.Obj{
				Key:   obj.Key,
				Value: obj.Value,
			}
			tmp = append(tmp, kv)
		}
		option.Options = tmp
	} else {
		option = nil
	}
	return option
}

func OptionLoc2Gen(options *common.Form) *u.Form {
	option := &u.Form{}
	if options != nil {
		// 不空，则提取options中的kv对，插入option（要解耦数据定义）
		objs := options.Options
		var tmp []*u.Form_Obj

		for _, obj := range objs {
			kv := &u.Form_Obj{
				Key:   obj.Key,
				Value: obj.Value,
			}
			tmp = append(tmp, kv)
		}
		option.Options = tmp
	} else {
		option = nil
	}
	return option
}

func VerifyGen2Loc(verify *u.UnitVerify) *unit.UnitVerify {
	res := &unit.UnitVerify{
		UnitId:     verify.UserId,
		VerifyType: verify.VerifyType,
	}

	if verify.VerifyType == consts.Weak {
		res.Account = &verify.Account
		res.Password = &verify.Password
		res.Form = OptionGen2Loc(verify.Form)
	}

	return res
}

func VerifyLoc2Gen(verify *unit.UnitVerify) *u.UnitVerify {
	res := &u.UnitVerify{
		UserId:     verify.UnitId,
		VerifyType: verify.VerifyType,
	}

	if verify.VerifyType == consts.Weak {
		res.Account = *verify.Account
		res.Password = *verify.Password
		res.Form = OptionLoc2Gen(verify.Form)
	}

	return res
}
