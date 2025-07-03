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
	"github.com/xh-polaris/psych-user/biz/infrastructure/util/reg"
)

type IUserService interface {
	UserSignUp(ctx context.Context, req *u.UserSignUpReq) (res *u.UserSignUpResp, err error)
	UserGetInfo(ctx context.Context, req *u.UserGetInfoReq) (res *u.UserGetInfoResp, err error)
	UserUpdateInfo(ctx context.Context, req *u.UserUpdateInfoReq) (res *basic.Response, err error)
	UserUpdatePassword(ctx context.Context, req *u.UserUpdatePasswordReq) (res *basic.Response, err error)
	UserBelongUnit(ctx context.Context, req *u.UserBelongUnitReq) (res *u.UserBelongUnitResp, err error)
	UserSignIn(ctx context.Context, req *u.UserSignInReq) (res *u.UserSignInResp, err error)
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

func (s *UserService) UserSignUp(ctx context.Context, req *u.UserSignUpReq) (res *u.UserSignUpResp, err error) {
	// 默认用户通过注册接口，使用手机号注册
	pwd, err := encrypt.BcryptEncrypt(req.User.Password)
	if err != nil {
		return nil, err
	}
	user := &usrmapper.User{
		Phone:    req.User.Phone,
		Password: pwd,
		Name:     req.User.Name,
		Birth:    req.User.Birth,
		Gender:   req.User.Gender,
		Status:   consts.Active,
	}
	userId, err := s.UserMapper.InsertWithEcho(ctx, user)
	if err != nil {
		return nil, err
	}
	return &u.UserSignUpResp{User: &u.User{
		Id:         userId,
		Phone:      user.Phone,
		Name:       user.Name,
		Birth:      user.Birth,
		Gender:     user.Gender,
		Status:     user.Status,
		CreateTime: user.CreateTime.Unix(),
		UpdateTime: user.UpdateTime.Unix(),
	}}, err
}

func (s *UserService) UserSignIn(ctx context.Context, req *u.UserSignInReq) (res *u.UserSignInResp, err error) {
	// 根据 authType 选择登录方式
	authType := req.AuthType
	switch authType {
	case consts.AuthPhoneAndPwd:
		{
			userId, err := s.signInWithPhoneAndPwd(ctx, req)
			if err != nil {
				return nil, err
			}
			return &u.UserSignInResp{UserId: userId}, nil
		}
	case consts.AuthPhoneAndCode:
		{
			userId, err := s.signInWithPhoneAndCode(ctx, req)
			if err != nil {
				return nil, err
			}
			return &u.UserSignInResp{UserId: userId}, nil
		}
	case consts.AuthStudentIdAndPwd:
		{
			userId, unitId, studentId, err := s.signInWithStudentIdAndPwd(ctx, req)
			if err != nil {
				return nil, err
			}
			return &u.UserSignInResp{
				UnitId:    *unitId,
				UserId:    *userId,
				StudentId: studentId,
			}, nil
		}
	}

	return nil, nil
}

func (s *UserService) signInWithPhoneAndPwd(ctx context.Context, req *u.UserSignInReq) (string, error) {
	// 通过手机号+密码登录

	// 获取手机号
	phone := req.GetAuthId()
	if !reg.CheckMobile(phone) {
		return "", consts.ErrInvalidParams
	}

	// 查询用户密码
	password := req.GetPassword()
	user, err := s.UserMapper.FindOneByPhone(ctx, phone)
	if err != nil {
		return "", err
	}

	// 校验密码
	if !encrypt.BcryptCheck(password, user.Password) {
		return "", consts.ErrUserPasswordMismatch
	}

	// 登录成功
	return user.ID.Hex(), nil
}

func (s *UserService) signInWithPhoneAndCode(ctx context.Context, req *u.UserSignInReq) (string, error) {
	// TODO: 手机号+验证码
	return "", nil
}

func (s *UserService) signInWithStudentIdAndPwd(ctx context.Context, req *u.UserSignInReq) (userId, unitId, studentId *string, err error) {
	// 根据 studentId + unitId 获取 userId
	uu, err := s.UUMapper.FindOneByUnitAndStu(ctx, req.UnitId, req.AuthId)
	if err != nil {
		return nil, nil, nil, err
	}

	// 根据 userId 获取密码
	user, err := s.UserMapper.FindOne(ctx, uu.UserId)
	if err != nil {
		return nil, nil, nil, err
	}

	// 校验密码
	if !encrypt.BcryptCheck(*req.Password, user.Password) {
		return nil, nil, nil, consts.ErrUserPasswordMismatch
	}

	hex := user.ID.Hex()
	return &hex, &req.UnitId, &req.AuthId, nil
}

func (s *UserService) UserGetInfo(ctx context.Context, req *u.UserGetInfoReq) (res *u.UserGetInfoResp, err error) {
	userId := req.GetUserId()
	logx.Info("正在查找用户信息: ", userId)
	user, err := s.UserMapper.FindOne(ctx, userId)
	if err != nil {
		return nil, consts.ErrUserNotExist
	}

	res = &u.UserGetInfoResp{User: &u.User{
		Id:         userId,
		Phone:      user.Phone,
		Name:       user.Name,
		Birth:      user.Birth,
		Gender:     user.Gender,
		Status:     user.Status,
		CreateTime: user.CreateTime.Unix(),
		UpdateTime: user.UpdateTime.Unix(),
	}}

	// 如果传入了 unitId，查询option
	if unitId := req.GetUnitId(); unitId != "" {
		// 查询关系表
		uu, err := s.UUMapper.FindOneByUU(ctx, userId, unitId)
		if err != nil {
			return nil, err
		}
		res.UnitId = &unitId
		res.StudentId = &uu.StudentId
		res.Form = uu.Options
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
