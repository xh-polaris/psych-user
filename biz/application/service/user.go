package service

import (
	"context"
	"errors"
	"github.com/google/wire"
	"github.com/xh-polaris/psych-idl/kitex_gen/basic"
	u "github.com/xh-polaris/psych-idl/kitex_gen/user"
	"github.com/xh-polaris/psych-pkg/util/logx"
	"github.com/xh-polaris/psych-user/biz/infrastructure/consts"
	untmapper "github.com/xh-polaris/psych-user/biz/infrastructure/mapper/unit"
	uumapper "github.com/xh-polaris/psych-user/biz/infrastructure/mapper/unit_user"
	usrmapper "github.com/xh-polaris/psych-user/biz/infrastructure/mapper/user"
	"github.com/xh-polaris/psych-user/biz/infrastructure/util/convert"
	"github.com/xh-polaris/psych-user/biz/infrastructure/util/encrypt"
	"github.com/xh-polaris/psych-user/biz/infrastructure/util/reg"
	"github.com/xh-polaris/psych-user/biz/infrastructure/util/result"
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
	// 参数校验
	if req.User == nil || !reg.CheckMobile(req.User.Phone) || req.User.Name == "" || req.User.Password == "" {
		logx.Error("UserSignUp fail")
		return nil, consts.ErrUnitSignUp
	}

	// 检查手机号是否已注册
	if _, err = s.UserMapper.FindOneByPhone(ctx, req.User.Phone); !errors.Is(err, consts.ErrNotFound) {
		return nil, err
	}
	if _, err = s.UnitMapper.FindOneByPhone(ctx, req.User.Phone); !errors.Is(err, consts.ErrNotFound) {
		return nil, err
	}

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

	// 如果传入了 unitId，需要进行关联
	if unitId := req.GetUnitId(); unitId != "" {
		uu := &uumapper.UnitUser{
			UserId:    userId,
			StudentId: req.GetStudentId(),
			UnitId:    unitId,
		}

		if err := s.UUMapper.Insert(ctx, uu); err != nil {
			return nil, err
		}
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
			// 手机号 + 密码
			userId, err := s.signInWithPhoneAnwPwd(ctx, req.AuthId, req.VerifyCode)
			if err != nil {
				return nil, err
			}
			return &u.UserSignInResp{
				UserId: userId,
				Strong: true,
			}, nil
		}
	case consts.AuthPhoneAndCode:
		{
			// 手机号 + 验证码
			userId, err := s.signInWithPhoneAndCode(ctx, req.AuthId, req.VerifyCode)
			if err != nil {
				return nil, err
			}
			return &u.UserSignInResp{
				UserId: userId,
				Strong: true,
			}, nil
		}
	case consts.AuthStudentIdAndPwd:
		{
			// 学号 + 密码
			userId, err := s.signInWithStuIdAndPwd(ctx, req.UnitId, req.AuthId, req.VerifyCode)
			if err != nil {
				return nil, err
			}
			return &u.UserSignInResp{
				UnitId:    req.UnitId,
				UserId:    userId,
				StudentId: &req.AuthId,
				Strong:    true,
			}, nil
		}
	case consts.AuthWeakAccountAndPwd:
		{
			// 弱验证账号 + 密码
			unitId, err := s.signInWithWeakAccountAndPwd(ctx, req.AuthId, req.VerifyCode)
			if err != nil {
				return nil, err
			}
			return &u.UserSignInResp{
				UnitId: unitId,
				Strong: false,
			}, nil
		}
	}
	return nil, consts.ErrUserSignIn
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
		form, err := convert.FormDB2Gen(uu.Options)
		if err != nil {
			return nil, err
		}
		res.UnitId = &unitId
		res.StudentId = &uu.StudentId
		res.Form = form
	}

	return res, nil
}
func (s *UserService) UserUpdateInfo(ctx context.Context, req *u.UserUpdateInfoReq) (res *basic.Response, err error) {
	// 修改user私有信息
	err = s.UserMapper.UpdateBasicInfo(ctx, req.User)
	if err != nil {
		return nil, err
	}

	// 修改关联信息
	if unitId := req.GetUnitId(); unitId != "" {
		form, err := convert.FormGen2DB(req.Form)
		if err != nil {
			return nil, err
		}
		err = s.UUMapper.UpdateBasicInfo(ctx, req.User.Id, unitId, form)
	}

	return result.ResponseOk(), err
}
func (s *UserService) UserUpdatePassword(ctx context.Context, req *u.UserUpdatePasswordReq) (res *basic.Response, err error) {
	authType := req.AuthType
	flag := false

	// 校验
	switch authType {
	case consts.UpdateByOldPwd:
		// 旧密码
		flag, err = s.updatePwdByOldPwd(ctx, req.Id, req.AuthValue, req.NewPassword)
	case consts.UpdateByCode:
		// TODO: 验证码
		flag, err = s.updatePwdByCode(ctx, req.Id, req.AuthValue, req.NewPassword)
	}

	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}

	if flag {
		// 修改密码
		newPwd, err := encrypt.BcryptEncrypt(req.NewPassword)
		if err != nil {
			return nil, err
		}
		if err := s.UserMapper.UpdatePassword(ctx, req.Id, newPwd); err != nil {
			return nil, err
		}
		return result.ResponseOk(), nil
	}
	return nil, nil
}

func (s *UserService) UserBelongUnit(ctx context.Context, req *u.UserBelongUnitReq) (res *u.UserBelongUnitResp, err error) {
	// 查询当前用户所属单位
	return nil, err
}

func (s *UserService) signInWithPhoneAnwPwd(ctx context.Context, phone, password string) (string, error) {
	// 通过手机号+密码登录

	// 获取手机号
	if !reg.CheckMobile(phone) {
		return "", consts.ErrInvalidParams
	}

	// 查询用户密码
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

func (s *UserService) signInWithPhoneAndCode(ctx context.Context, phone, code string) (string, error) {
	// TODO: 手机号+验证码
	return "", nil
}

func (s *UserService) signInWithStuIdAndPwd(ctx context.Context, unitId, studentId, password string) (userId string, err error) {
	// 根据 studentId + unitId 获取 userId
	uu, err := s.UUMapper.FindOneByUnitAndStu(ctx, unitId, studentId)
	if err != nil {
		return "", err
	}

	// 根据 userId 获取密码
	user, err := s.UserMapper.FindOne(ctx, uu.UserId)
	if err != nil {
		return "", err
	}

	// 校验密码
	if !encrypt.BcryptCheck(password, user.Password) {
		return "", consts.ErrUserPasswordMismatch
	}

	return user.ID.Hex(), nil
}

func (s *UserService) signInWithWeakAccountAndPwd(ctx context.Context, account, password string) (unitId string, err error) {
	// 查询弱验证账号所属 unit
	unit, err := s.UnitMapper.FindOneByAccount(ctx, account)
	if err != nil {
		return "", err
	}

	// 校验密码
	if !encrypt.BcryptCheck(password, unit.Password) {
		return "", consts.ErrUserPasswordMismatch
	}

	// 登录成功
	return unit.ID.Hex(), nil
}

func (s *UserService) updatePwdByOldPwd(ctx context.Context, userId string, oldPwd string, newPwd string) (bool, error) {
	// 查找用户
	user, err := s.UserMapper.FindOne(ctx, userId)
	if err != nil {
		return false, err
	}

	// 判断密码
	if !encrypt.BcryptCheck(oldPwd, user.Password) {
		return false, consts.ErrUserPasswordMismatch
	}

	return true, err
}

func (s *UserService) updatePwdByCode(ctx context.Context, userId string, code string, newPwd string) (bool, error) {
	// TODO: 根据验证码修改密码
	return false, nil
}
