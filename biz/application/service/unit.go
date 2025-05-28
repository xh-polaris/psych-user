package service

import (
	"context"
	"errors"
	"github.com/xh-polaris/psych-pkg/wirex"

	"github.com/xh-polaris/psych-idl/kitex_gen/basic"
	u "github.com/xh-polaris/psych-idl/kitex_gen/user"
	"github.com/xh-polaris/psych-user/biz/infrastructure/consts"
	umapper "github.com/xh-polaris/psych-user/biz/infrastructure/mapper/unit"
	usrMapper "github.com/xh-polaris/psych-user/biz/infrastructure/mapper/user"
	"github.com/xh-polaris/psych-user/biz/infrastructure/util/encrypt"
)

type IUnitService interface {
	UnitSignUp(ctx context.Context, req *u.UnitSignUpReq) (res *basic.Response, err error)
	UnitGetInfo(ctx context.Context, req *u.UnitGetInfoReq) (res *u.UnitGetInfoResp, err error)
	UnitUpdateInfo(ctx context.Context, req *u.UnitUpdateInfoReq) (res *basic.Response, err error)
	UnitUpdatePassword(ctx context.Context, req *u.UnitUpdatePasswordReq) (res *basic.Response, err error)
	UnitStrongVerify(ctx context.Context, req *u.UnitStrongVerifyReq) (res *basic.Response, err error)
	UnitWeakVerify(ctx context.Context, req *u.UnitWeakVerifyReq) (res *basic.Response, err error)
	UnitCreateVerify(ctx context.Context, req *u.UnitCreateVerifyReq) (res *u.UnitCreateVerifyResp, err error)
	UnitUpdateVerifyPassword(ctx context.Context, req *u.UnitUpdateVerifyPasswordReq) (res *basic.Response, err error)
	UnitLinkUser(ctx context.Context, req *u.UnitLinkUserReq) (res *basic.Response, err error)
	UnitLinkView(ctx context.Context, req *u.UnitLinkViewReq) (res *basic.Response, err error)
	UnitPageQueryUser(ctx context.Context, req *u.UnitPageQueryUserReq) (res *u.UnitPageQueryUserResp, err error)
	UnitPageQueryView(ctx context.Context, req *u.UnitPageQueryViewReq) (res *u.UnitPageQueryViewResp, err error)
	UnitCreateAndLinkUser(ctx context.Context, req *u.UnitCreateAndLinkUserReq) (res *basic.Response, err error)
	UnitCreateAndLinkView(ctx context.Context, req *u.UnitCreateAndLinkViewReq) (res *basic.Response, err error)
	UnitGetAppInfo(ctx context.Context, req *u.UnitGetAppInfoReq) (res *u.UnitGetAppInfoResp, err error)
	UnitModelGetInfo(ctx context.Context, req *u.UnitModelGetInfoReq) (res *u.UnitModelGetInfoResp, err error)
	UnitModelUpdateInfo(ctx context.Context, req *u.UnitModelUpdateInfoReq) (res *basic.Response, err error)
}

type UnitService struct {
	UnitMapper *umapper.MongoMapper
	UserMapper *usrMapper.MongoMapper
}

var UnitServiceSet = wirex.NewWireSet[UnitService, IUnitService]()

// UnitSignUp 单位账号注册
func (s *UnitService) UnitSignUp(ctx context.Context, req *u.UnitSignUpReq) (res *basic.Response, err error) {
	res = &basic.Response{}

	// 参数校验
	if req.Unit == nil || req.Unit.Phone == "" || req.Unit.Password == "" || req.Unit.Name == "" {
		return nil, consts.ErrInvalidParams
	}

	// 检查手机号是否已注册
	existUnit, err := s.UnitMapper.FindOneByPhone(ctx, req.Unit.Phone)
	if err == nil && existUnit != nil {
		return nil, consts.ErrUnitPhoneExist
	} else if err != nil && !errors.Is(err, consts.ErrNotFound) {
		return nil, err
	}

	// 密码加密
	hashedPwd, err := encrypt.BcryptEncrypt(req.Unit.Password)
	if err != nil {
		return nil, consts.ErrUnitSignUp
	}

	// 创建单位对象
	unit := &umapper.Unit{
		Phone:    req.Unit.Phone,
		Password: hashedPwd,
		Name:     req.Unit.Name,
		Address:  req.Unit.Address,
		Contact:  req.Unit.Contact,
		Level:    req.Unit.Level,
		Status:   "active",
	}

	// 保存到数据库
	err = s.UnitMapper.Insert(ctx, unit)
	if err != nil {
		return nil, consts.ErrUnitSignUp
	}

	// 返回成功响应
	res.Code = 0
	res.Msg = "success"
	return res, nil
}

// UnitStrongVerify 单位账号认证
func (s *UnitService) UnitStrongVerify(ctx context.Context, req *u.UnitStrongVerifyReq) (res *basic.Response, err error) {
	res = &basic.Response{}

	// 参数校验
	if req.Phone == "" || req.Password == nil || *req.Password == "" {
		return nil, consts.ErrInvalidParams
	}

	// TODO: 验证码功能
	if req.VerifyCode != nil && *req.VerifyCode != "" {
		// 验证码校验逻辑
	}

	// 查询单位账号
	unit, err := s.UnitMapper.FindOneByPhone(ctx, req.Phone)
	if err != nil {
		if errors.Is(err, consts.ErrNotFound) {
			return nil, consts.ErrUnitNotExist
		}
		return nil, err
	}

	// 密码校验
	if !encrypt.BcryptCheck(*req.Password, unit.Password) {
		return nil, consts.ErrUnitPasswordMismatch
	}

	// 返回成功响应
	res.Code = 0
	res.Msg = "success"
	return res, nil
}

// UnitGetInfo 获取单位信息
func (s *UnitService) UnitGetInfo(ctx context.Context, req *u.UnitGetInfoReq) (res *u.UnitGetInfoResp, err error) {
	// 参数校验
	if req.Id == "" {
		return nil, consts.ErrInvalidParams
	}

	// 查询单位信息
	unit, err := s.UnitMapper.FindOne(ctx, req.Id)
	if err != nil {
		if errors.Is(err, consts.ErrNotFound) {
			return nil, consts.ErrUnitNotExist
		}
		return nil, consts.ErrUnitGetInfo
	}

	// 构建响应
	res = &u.UnitGetInfoResp{
		Unit: &u.Unit{
			Id:         unit.Id.Hex(),
			Phone:      unit.Phone,
			Password:   "", // 密码字段为空
			Name:       unit.Name,
			Address:    unit.Address,
			Contact:    unit.Contact,
			Level:      unit.Level,
			Status:     unit.Status,
			CreateTime: unit.CreateTime.Unix(),
			UpdateTime: unit.UpdateTime.Unix(),
		},
	}

	return res, nil
}

// UnitLinkUser 关联用户账号
func (s *UnitService) UnitLinkUser(ctx context.Context, req *u.UnitLinkUserReq) (res *basic.Response, err error) {
	res = &basic.Response{}

	// 参数校验
	if req.UnitId == "" || req.UserId == "" {
		return nil, consts.ErrInvalidParams
	}

	// 验证单位是否存在
	_, err = s.UnitMapper.FindOne(ctx, req.UnitId)
	if err != nil {
		if err == consts.ErrNotFound {
			return nil, consts.ErrUnitNotExist
		}
		return nil, err
	}

	// 验证用户是否存在
	_, err = s.UserMapper.FindOne(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	// 检查关联是否已存在
	exists, err := s.UnitMapper.CheckLinkExists(ctx, req.UnitId, req.UserId)
	if err != nil {
		return nil, err
	}

	// 如果关联不存在，则创建关联
	if !exists {
		err = s.UnitMapper.LinkUser(ctx, req.UnitId, req.UserId)
		if err != nil {
			return nil, consts.ErrUnitLinkUser
		}
	}

	// 返回成功响应
	res.Code = 0
	res.Msg = "success"
	return res, nil
}

// UnitCreateAndLinkUser 批量创建并关联用户
func (s *UnitService) UnitCreateAndLinkUser(ctx context.Context, req *u.UnitCreateAndLinkUserReq) (res *basic.Response, err error) {
	res = &basic.Response{}

	// 参数校验
	if req.UnitId == "" || len(req.UserPhone) == 0 {
		return nil, consts.ErrInvalidParams
	}

	// 验证单位是否存在
	_, err = s.UnitMapper.FindOne(ctx, req.UnitId)
	if err != nil {
		if errors.Is(err, consts.ErrNotFound) {
			return nil, consts.ErrUnitNotExist
		}
		return nil, err
	}

	// 遍历手机号列表
	for _, phone := range req.UserPhone {
		// 查询用户是否已存在
		user, err := s.UserMapper.FindOneByPhone(ctx, phone)

		var userId string
		if errors.Is(err, consts.ErrNotFound) {
			// 用户不存在，创建新用户
			// 生成统一默认密码
			defaultPassword, err := encrypt.BcryptEncrypt("123456") // 默认密码示例
			if err != nil {
				continue
			}

			newUser := &usrMapper.User{
				Phone:    phone,
				Password: defaultPassword,
				Name:     phone,           // 默认使用手机号作为名称
				Status:   "need_complete", // 需要完善信息
			}

			err = s.UserMapper.Insert(ctx, newUser)
			if err != nil {
				continue
			}

			userId = newUser.Id.Hex()
		} else if err != nil {
			// 查询出错，跳过当前手机号
			continue
		} else {
			// 用户已存在
			userId = user.Id.Hex()
		}

		// 检查关联是否已存在
		exists, err := s.UnitMapper.CheckLinkExists(ctx, req.UnitId, userId)
		if err != nil {
			continue
		}

		// 如果关联不存在，则创建关联
		if !exists {
			err = s.UnitMapper.LinkUser(ctx, req.UnitId, userId)
			if err != nil {
				continue
			}
		}
	}

	// 返回成功响应
	res.Code = 0
	res.Msg = "success"
	return res, nil
}

func (s *UnitService) UnitUpdateInfo(ctx context.Context, req *u.UnitUpdateInfoReq) (res *basic.Response, err error) {
	return nil, err
}

func (s *UnitService) UnitUpdatePassword(ctx context.Context, req *u.UnitUpdatePasswordReq) (res *basic.Response, err error) {
	return nil, err
}

func (s *UnitService) UnitWeakVerify(ctx context.Context, req *u.UnitWeakVerifyReq) (res *basic.Response, err error) {
	return nil, err
}

func (s *UnitService) UnitCreateVerify(ctx context.Context, req *u.UnitCreateVerifyReq) (res *u.UnitCreateVerifyResp, err error) {
	return nil, err
}

func (s *UnitService) UnitUpdateVerifyPassword(ctx context.Context, req *u.UnitUpdateVerifyPasswordReq) (res *basic.Response, err error) {
	return nil, err
}

func (s *UnitService) UnitLinkView(ctx context.Context, req *u.UnitLinkViewReq) (res *basic.Response, err error) {
	return nil, err
}

func (s *UnitService) UnitPageQueryUser(ctx context.Context, req *u.UnitPageQueryUserReq) (res *u.UnitPageQueryUserResp, err error) {
	return nil, err
}

func (s *UnitService) UnitPageQueryView(ctx context.Context, req *u.UnitPageQueryViewReq) (res *u.UnitPageQueryViewResp, err error) {
	return nil, err
}

func (s *UnitService) UnitCreateAndLinkView(ctx context.Context, req *u.UnitCreateAndLinkViewReq) (res *basic.Response, err error) {
	return nil, err
}

func (s *UnitService) UnitGetAppInfo(ctx context.Context, req *u.UnitGetAppInfoReq) (res *u.UnitGetAppInfoResp, err error) {
	return nil, err
}

func (s *UnitService) UnitModelGetInfo(ctx context.Context, req *u.UnitModelGetInfoReq) (res *u.UnitModelGetInfoResp, err error) {
	return nil, err
}

func (s *UnitService) UnitModelUpdateInfo(ctx context.Context, req *u.UnitModelUpdateInfoReq) (res *basic.Response, err error) {
	return nil, err
}
