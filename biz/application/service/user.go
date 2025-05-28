package service

import (
	"context"
	"github.com/xh-polaris/psych-idl/kitex_gen/basic"
	u "github.com/xh-polaris/psych-idl/kitex_gen/user"
	"github.com/xh-polaris/psych-user/biz/infrastructure/consts"
	umapper "github.com/xh-polaris/psych-user/biz/infrastructure/mapper/user"
)

type IUserService interface {
	UserGetInfo(ctx context.Context, req *u.UserGetInfoReq) (res *u.UserGetInfoResp, err error)
	UserSignIn(ctx context.Context, req *u.UserSignInReq) (res *basic.Response, err error)

	UserSignUp(ctx context.Context, req *u.UserSignUpReq) (res *basic.Response, err error)
	UserUpdateInfo(ctx context.Context, req *u.UserUpdateInfoReq) (res *basic.Response, err error)
	UserUpdatePassword(ctx context.Context, req *u.UserUpdatePasswordReq) (res *basic.Response, err error)
	UserBelongUnit(ctx context.Context, req *u.UserBelongUnitReq) (res *u.UserBelongUnitResp, err error)
}

type UserService struct {
	UserMapper *umapper.MongoMapper
	// UnitMapper
}

func (s *UserService) UserSignUp(ctx context.Context, req *u.UserSignUpReq) (res *basic.Response, err error) {
	res = &basic.Response{}

	// 1. 判断根据手机号还是学号注册

	return nil, err
}

func (s *UserService) UserSignIn(ctx context.Context, req *u.UserSignInReq) (res *basic.Response, err error) {
	res = &basic.Response{}
	// 1. 判断使用学号还是手机号登录。学号-密码；手机号-密码或验证码

	// 2. 如果验证码不为空，则走验证码校验流程
	if req.VerifyCode != nil && *req.VerifyCode != "" {
		// TODO: 验证码
		if *req.VerifyCode == "?" {
			res.Code = 0
			res.Msg = "success"
			return res, nil
		} else {
			res.Code = 1
			res.Msg = "fail"
			return res, nil
		}
	}

	// 3. 如果验证码为空，走密码验证
	if req.Password == nil || *req.Password == "" {
		res.Code = 1
		res.Msg = "password cannot be empty"
		return res, nil
	}

	// 查询用户信息
	user, err := s.UserMapper.FindOneByPhone(ctx, req.Phone)
	if err != nil {
		res.Code = 1
		res.Msg = "user not found"
		return res, nil
	}

	if user.Password != *req.Password {
		res.Code = 1
		res.Msg = "password mismatch"
		return res, nil
	}

	res.Code = 0
	res.Msg = "success"
	return res, nil
}
func (s *UserService) UserGetInfo(ctx context.Context, req *u.UserGetInfoReq) (res *u.UserGetInfoResp, err error) {
	id := req.Id
	user, err := s.UserMapper.FindOne(ctx, id)
	if err != nil {
		return nil, consts.ErrNotFound
	}

	return &u.UserGetInfoResp{User: &u.User{
		Id:         user.Id.Hex(),
		Phone:      user.Phone,
		Password:   "",
		Name:       user.Name,
		Birth:      user.Birth,
		Gender:     user.Gender,
		Status:     user.Status,
		CreateTime: user.CreateTime.Unix(),
		UpdateTime: user.DeleteTime.Unix(),
	}}, err
}
func (s *UserService) UserUpdateInfo(ctx context.Context, req *u.UserUpdateInfoReq) (res *basic.Response, err error) {
	return nil, err
}
func (s *UserService) UserUpdatePassword(ctx context.Context, req *u.UserUpdatePasswordReq) (res *basic.Response, err error) {
	return nil, err
}
func (s *UserService) UserBelongUnit(ctx context.Context, req *u.UserBelongUnitReq) (res *u.UserBelongUnitResp, err error) {
	return nil, err
}
