package service

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/psych-idl/kitex_gen/basic"
	u "github.com/xh-polaris/psych-idl/kitex_gen/user"
	"github.com/xh-polaris/psych-user/biz/infrastructure/consts"
	nmapper "github.com/xh-polaris/psych-user/biz/infrastructure/mapper/unit"
	umapper "github.com/xh-polaris/psych-user/biz/infrastructure/mapper/user"
	"github.com/xh-polaris/psych-user/biz/infrastructure/util/encrypt"
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
	UnitMapper *nmapper.MongoMapper
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
	// 1. 判断使用学号还是手机号登录。学号->密码；手机号->密码或验证码
	if req.GetPhone() != "" {
		// 手机号校验
		// 2. 判断手机号是否存在
		user, err := s.UserMapper.FindOneByPhone(ctx, *req.Phone)
		if err != nil {
			return nil, consts.ErrUserNotExist
		}

		// 3. 如果验证码不为空，则走验证码校验流程
		if req.GetVerifyCode() != "" {
			// TODO: 验证码
			if *req.VerifyCode == "xh-polaris" {
				return &basic.Response{
					Code: 200,
					Msg:  "",
				}, nil
			} else {
				return res, consts.ErrUserVerify
			}
		}

		// 4. 如果验证码为空，走密码验证
		if req.Password == nil || *req.Password == "" {
			return nil, consts.ErrUserPasswordMismatch
		}

		// 5. 校验密码
		if user.Password != *req.Password {
			return nil, consts.ErrUserPasswordMismatch
		}

	} else if len(req.StudentId) > 0 {
		// 学号校验
		// 2. 遍历学号，检查是否对应唯一user
		userIdSet := make(map[string]struct{})
		for _, sid := range req.StudentId {
			record, err := s.UserMapper.FindUSULinkBySid(ctx, sid)
			if err != nil || record == nil {
				return nil, consts.ErrUserStudentIdNotExist
			}
			userIdSet[record.UserId.Hex()] = struct{}{}
		}

		if len(userIdSet) != 1 {
			return nil, consts.ErrUserStudentIdDuplicate
		}

		// 2. 查询用户信息
		var userId string
		for uid := range userIdSet {
			userId = uid
		}
		user, err := s.UserMapper.FindOne(ctx, userId)
		if err != nil || user == nil {
			return nil, consts.ErrUserNotExist
		}

		if req.Password == nil || !encrypt.BcryptCheck(req.GetPassword(), user.Password) {
			return nil, consts.ErrUserPasswordMismatch
		}

		return &basic.Response{
			Code: 200,
			Msg:  "",
		}, nil
	}
	return nil, consts.ErrUserSignIn
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
