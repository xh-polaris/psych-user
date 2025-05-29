package service

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/psych-idl/kitex_gen/basic"
	u "github.com/xh-polaris/psych-idl/kitex_gen/user"
	"github.com/xh-polaris/psych-pkg/util/logx"
	"github.com/xh-polaris/psych-user/biz/infrastructure/consts"
	untmapper "github.com/xh-polaris/psych-user/biz/infrastructure/mapper/unit"
	uumapper "github.com/xh-polaris/psych-user/biz/infrastructure/mapper/unit_user"
	usrmapper "github.com/xh-polaris/psych-user/biz/infrastructure/mapper/user"
	"github.com/xh-polaris/psych-user/biz/infrastructure/util/encrypt"
	"github.com/xh-polaris/psych-user/biz/infrastructure/util/result"
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
	UserMapper *usrmapper.MongoMapper
	UnitMapper *untmapper.MongoMapper
	UUMapper   *uumapper.MongoMapper
}

var UserServiceSet = wire.NewSet(
	wire.Struct(new(UserService), "*"),
	wire.Bind(new(IUserService), new(*UserService)),
)

func (s *UserService) UserSignUp(ctx context.Context, req *u.UserSignUpReq) (res *basic.Response, err error) {
	// 1. 判断根据手机号还是学号注册
	return nil, err
}

func (s *UserService) UserSignIn(ctx context.Context, req *u.UserSignInReq) (res *basic.Response, err error) {
	// 根据 authType 选择登录方式
	authType := req.AuthType
	switch authType {
	case consts.AuthPhone:
		{
			// 通过手机号登录
			// 获取手机号
			phone := req.GetAuthId()
			password := req.GetPassword()
			verifyCode := req.GetVerifyCode()
			if phone == "" {
				return nil, consts.ErrUserNotExist
			}

			// 判断通过密码还是验证码登录
			if verifyCode != "" {
				// TODO: 验证码逻辑
				if *req.VerifyCode == "xh-polaris" {
					return result.ResponseOk(), nil
				} else {
					return nil, consts.ErrUserVerify
				}
			}

			// 查询用户密码
			user, err := s.UserMapper.FindOneByPhone(ctx, phone)
			if err != nil {
				return nil, consts.ErrUserNotExist
			}
			// 校验密码
			if !encrypt.BcryptCheck(password, user.Password) {
				return nil, consts.ErrUserPasswordMismatch
			}
			return result.ResponseOk(), nil
		}
	case consts.AuthStudentId:
		{
			// 通过学号登录
			studentId := req.GetAuthId()
			password := req.GetPassword()
			// 通过studentId + unitId查询
			user, err := s.UUMapper.FindOneByUnitAndStu(ctx, req.UnitId, studentId)
			if err != nil {
				return nil, consts.ErrUserNotExist
			}
			usr, err := s.UserMapper.FindOne(ctx, user.UserId)
			if err != nil {
				return nil, consts.ErrUserNotExist
			}
			if !encrypt.BcryptCheck(password, usr.Password) {
				return nil, consts.ErrUserPasswordMismatch
			}
			return result.ResponseOk(), nil
		}
	}
	return nil, consts.ErrUserSignIn
}
func (s *UserService) UserGetInfo(ctx context.Context, req *u.UserGetInfoReq) (res *u.UserGetInfoResp, err error) {
	userId := req.GetUserId()
	logx.Info("正在查找用户信息: userId = %d\n", userId)
	user, err := s.UserMapper.FindOne(ctx, userId)
	if err != nil {
		return nil, consts.ErrUserNotExist
	}

	res = &u.UserGetInfoResp{User: &u.User{
		Id:         user.Id,
		Phone:      user.Phone,
		Password:   "",
		Name:       user.Name,
		Birth:      user.Birth,
		Gender:     user.Gender,
		Status:     user.Status,
		CreateTime: user.CreateTime.Unix(),
		UpdateTime: user.DeleteTime.Unix(),
	}}
	unitId := req.GetUnitId()
	if unitId != "" {
		logx.Info("正在查找用户关联信息: userId = %s, unitId = %s\n", userId, unitId)
		uu, err := s.UUMapper.FindOneByUU(ctx, userId, unitId)
		if err != nil {
			return nil, err
		}
		arr := []*u.Option{uu.Options}
		arr = append(arr, uu.Options)
		res.Options = arr
		// res.Options = uu.Options
	}
	return res, nil
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
