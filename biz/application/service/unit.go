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
	uuMapper "github.com/xh-polaris/psych-user/biz/infrastructure/mapper/unit_user"
	usrMapper "github.com/xh-polaris/psych-user/biz/infrastructure/mapper/user"
	"github.com/xh-polaris/psych-user/biz/infrastructure/util/convert"
	"github.com/xh-polaris/psych-user/biz/infrastructure/util/encrypt"
	"github.com/xh-polaris/psych-user/biz/infrastructure/util/reg"
	"github.com/xh-polaris/psych-user/biz/infrastructure/util/result"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type IUnitService interface {
	UnitSignUp(ctx context.Context, req *u.UnitSignUpReq) (res *u.UnitSignUpResp, err error)
	UnitGetInfo(ctx context.Context, req *u.UnitGetInfoReq) (res *u.UnitGetInfoResp, err error)
	UnitUpdateInfo(ctx context.Context, req *u.UnitUpdateInfoReq) (res *basic.Response, err error)
	UnitUpdatePassword(ctx context.Context, req *u.UnitUpdatePasswordReq) (res *basic.Response, err error)
	UnitCreateAndLinkUser(ctx context.Context, req *u.UnitCreateAndLinkUserReq) (res *basic.Response, err error)
	UnitCreateAndLinkView(ctx context.Context, req *u.UnitCreateAndLinkViewReq) (res *basic.Response, err error)
	UnitSignIn(ctx context.Context, req *u.UnitSignInReq) (res *u.UnitSignInResp, err error)
	UnitCreateVerify(ctx context.Context, req *u.UnitCreateVerifyReq) (res *u.UnitCreateVerifyResp, err error)
	UnitUpdateVerifyPassword(ctx context.Context, req *u.UnitUpdateVerifyReq) (res *basic.Response, err error)
	UnitLinkUser(ctx context.Context, req *u.UnitLinkUserReq) (res *basic.Response, err error)
	UnitLinkView(ctx context.Context, req *u.UnitLinkViewReq) (res *basic.Response, err error)
	UnitPageQueryUser(ctx context.Context, req *u.UnitPageQueryUserReq) (res *u.UnitPageQueryUserResp, err error)
	UnitPageQueryView(ctx context.Context, req *u.UnitPageQueryViewReq) (res *u.UnitPageQueryViewResp, err error)
	UnitGetAppInfo(ctx context.Context, req *u.UnitGetAppInfoReq) (res *u.UnitGetAppInfoResp, err error)
	UnitModelGetInfo(ctx context.Context, req *u.UnitModelGetInfoReq) (res *u.UnitModelGetInfoResp, err error)
	UnitModelUpdateInfo(ctx context.Context, req *u.UnitModelUpdateInfoReq) (res *basic.Response, err error)
}

type UnitService struct {
	UnitMapper *untmapper.MongoMapper
	UserMapper *usrMapper.MongoMapper
	UUMapper   *uuMapper.MongoMapper
}

var UnitServiceSet = wire.NewSet(
	wire.Struct(new(UnitService), "*"),
	wire.Bind(new(IUnitService), new(*UnitService)),
)

// UnitSignUp 单位账号注册
func (s *UnitService) UnitSignUp(ctx context.Context, req *u.UnitSignUpReq) (res *u.UnitSignUpResp, err error) {
	// 参数校验
	if req.Unit == nil || !reg.CheckMobile(req.Unit.Phone) || req.Unit.Name == "" || req.Unit.Password == "" {
		logx.Error("UnitSignUp fail")
		return nil, consts.ErrUnitSignUp
	}

	// 检查手机号是否已注册
	if _, err = s.UnitMapper.FindOneByPhone(ctx, req.Unit.Phone); err != nil {
		return nil, consts.ErrUnitPhoneExist
	}

	// 密码加密
	hashedPwd, err := encrypt.BcryptEncrypt(req.Unit.Password)
	if err != nil {
		return nil, consts.ErrUnitSignUp
	}

	// 创建单位对象，要使用本地的unit进行解耦
	unit := &untmapper.Unit{
		Phone:    req.Unit.Phone,
		Password: hashedPwd,
		Name:     req.Unit.Name,
		Address:  req.Unit.Address,
		Contact:  req.Unit.Contact,
		Level:    req.Unit.Level,
		Status:   consts.Active,
	}

	// 保存到数据库
	unitId, err := s.UnitMapper.InsertWithEcho(ctx, unit)
	if err != nil {
		return nil, consts.ErrUnitSignUp
	}

	// 返回成功响应
	return &u.UnitSignUpResp{
		Unit: &u.Unit{
			Id:      unitId,
			Phone:   unit.Phone,
			Name:    unit.Name,
			Address: unit.Address,
			Contact: unit.Contact,
			Level:   unit.Level,
			Status:  unit.Status,
		},
	}, nil
}

// UnitSignIn 单位账号登录
func (s *UnitService) UnitSignIn(ctx context.Context, req *u.UnitSignInReq) (res *u.UnitSignInResp, err error) {
	// 手机号校验
	if !reg.CheckMobile(req.Phone) {
		return nil, consts.ErrInvalidParams
	}
	// 根据authType选择登录类型
	unit := &untmapper.Unit{}
	switch req.AuthType {
	case consts.AuthPhoneAndPwd:
		unit, err = s.signInWithPhoneAndPwd(ctx, req)
	case consts.AuthPhoneAndCode:
		unit, err = s.signInWithPhoneAndCode(ctx, req)
	default:
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &u.UnitSignInResp{
		Unit: &u.Unit{
			Id:      unit.ID.Hex(),
			Phone:   unit.Phone,
			Name:    unit.Name,
			Address: unit.Address,
			Contact: unit.Contact,
			Level:   unit.Level,
			Status:  unit.Status,
			// TODO
			//		Verify:     convert.VerifyLoc2Gen(unit.Verify),
			CreateTime: unit.CreateTime.Unix(),
			UpdateTime: unit.UpdateTime.Unix(),
		},
	}, nil
}

func (s *UnitService) signInWithPhoneAndPwd(ctx context.Context, req *u.UnitSignInReq) (*untmapper.Unit, error) {
	// 手机号+密码登录
	password := req.GetPassword()
	// 获取密码
	unit, err := s.UnitMapper.FindOneByPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}

	// 校验密码
	if encrypt.BcryptCheck(password, unit.Password) {
		unit.Password = ""
		return unit, nil
	}
	return nil, nil
}

func (s *UnitService) signInWithPhoneAndCode(ctx context.Context, req *u.UnitSignInReq) (*untmapper.Unit, error) {
	// TODO: 手机号+验证码登录
	return nil, nil
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
		return nil, err
	}

	// 构建响应
	res = &u.UnitGetInfoResp{
		Unit: &u.Unit{
			Id:      req.Id,
			Phone:   unit.Phone,
			Name:    unit.Name,
			Address: unit.Address,
			Contact: unit.Contact,
			Level:   unit.Level,
			Status:  unit.Status,
			// TODO		Verify:     convert.VerifyLoc2Gen(unit.Verify),
			CreateTime: unit.CreateTime.Unix(),
			UpdateTime: unit.UpdateTime.Unix(),
		},
	}

	return res, nil
}

// UnitLinkUser 关联用户账号
func (s *UnitService) UnitLinkUser(ctx context.Context, req *u.UnitLinkUserReq) (res *basic.Response, err error) {
	// 参数校验
	if req.UnitId == "" || req.UserId == "" {
		return nil, consts.ErrInvalidParams
	}

	// 验证单位是否存在
	if _, err = s.UnitMapper.FindOne(ctx, req.UnitId); err != nil {
		return nil, consts.ErrUnitNotExist
	}

	// 验证用户是否存在
	if _, err = s.UserMapper.FindOne(ctx, req.UserId); err != nil {
		return nil, consts.ErrUserNotExist
	}

	// 检查关联是否已存在
	if _, err := s.UnitMapper.CheckLinkExists(ctx, req.UnitId, req.UserId); err != nil {
		return nil, consts.ErrUnitLinkUser
	}

	// 如果关联不存在，则创建关联
	if s.UnitMapper.LinkUser(ctx, req.UnitId, req.UserId) != nil {
		return nil, consts.ErrUnitLinkUser
	}

	// 返回成功响应
	return result.ResponseOk(), nil
}

// UnitCreateAndLinkUser 批量创建并关联用户
func (s *UnitService) UnitCreateAndLinkUser(ctx context.Context, req *u.UnitCreateAndLinkUserReq) (res *basic.Response, err error) {
	authType := req.AuthType
	unitId := req.UnitId
	if unitId == "" {
		return nil, consts.ErrInvalidParams
	}

	for _, user := range req.Users {
		switch authType {
		case consts.CreateByPhone:
			s.createUserByPhone(ctx, unitId, user)
		case consts.CreateByStudentId:
			s.createUserByStudentId(ctx, unitId, user)
		}
	}
	return result.ResponseOk(), nil
}

func (s *UnitService) createUserByPhone(ctx context.Context, unitId string, user *u.UnitCreateAndLinkUserReq_U) {
	// 关联手机号
	// TODO: 开启事务?
	phone := user.AuthId
	defaultPwd, err := encrypt.BcryptEncrypt(consts.DefaultPassword)
	if err != nil {
		return
	}
	// 获取 option，判断是否为空
	option := convert.OptionGen2Loc(user.GetOptions())

	// 先查找该手机号是否已经注册
	logx.Info("正在查询用户: phone = %s, unitId = %s\n", phone, unitId)
	us, err := s.UserMapper.FindOneByPhone(ctx, phone)
	if errors.Is(err, consts.ErrNotFound) {
		// 未被注册，先在user表中添加数据
		newUser := &usrMapper.User{
			Phone:    phone,
			Password: defaultPwd,
			Name:     user.Name,
			Gender:   user.Gender,
			Status:   consts.Active,
		}
		userId, err := s.UserMapper.InsertWithEcho(ctx, newUser)
		if err != nil || userId == "" {
			logx.Error("创建用户失败。phone = %s, unitId = %s\n", phone, unitId)
			return
		}

		// 在unit_user表中关联
		if s.UUMapper.Insert(ctx, &uuMapper.UnitUser{
			UnitId:  unitId,
			UserId:  userId,
			Options: option,
		}) != nil {
			logx.Error("创建用户关联失败。userId = %s, unitId = %s, phone = %s\n", userId, unitId, phone)
			return
		}
		logx.Info("创建用户关联成功！userId = %s, unitId = %s, phone = %s\n", userId, unitId, phone)
	} else if err == nil && us != nil {
		// 已经被注册，则判断是否已经被关联，若无关联则关联，已关联则跳过（不报错）
		userId := us.ID.Hex()

		// 查询关联
		r, err := s.UUMapper.FindOneByUU(ctx, userId, unitId)
		if err != nil {
			logx.Error("查询用户失败。phone = %s, unitId = %s\n", phone, unitId)
			return
		}

		// 已经有关联，跳过
		if r != nil {
			logx.Info("用户和该单位已经绑定。userId = %s, unitId = %s, phone = %s\n", userId, unitId, phone)
			return
		}

		// 无关联，则插入
		err = s.UUMapper.Insert(ctx, &uuMapper.UnitUser{
			UnitId:  unitId,
			UserId:  userId,
			Options: option,
		})

		if err != nil {
			logx.Error("创建用户关联失败。userId = %s, unitId = %s, phone = %s\n", userId, unitId, phone)
			return
		}
		logx.Info("创建用户关联成功！userId = %s, unitId = %s, phone = %s\n", userId, unitId, phone)
	} else {
		logx.Error("查询用户失败。phone = %s, unitId = %s\n", phone, unitId)
	}
}

// 关联学号
func (s *UnitService) createUserByStudentId(ctx context.Context, unitId string, user *u.UnitCreateAndLinkUserReq_U) {
	studentId := user.AuthId
	password := user.Password
	if password == "" {
		password = encrypt.GetDefaultPwd()
	} else {
		pwd, err := encrypt.BcryptEncrypt(password)
		if err != nil {
		}
		password = pwd
	}
	option := convert.OptionGen2Loc(user.GetOptions())

	// 先根据学号(authId -> studentId)和unitId查找是否已经存在
	// 此时如果记录存在，则表示已经有user账号且关联完成
	link, err := s.UUMapper.FindOneByUnitAndStu(ctx, studentId, unitId)
	if errors.Is(err, consts.ErrNotFound) {
		// 如果不存在，则先进行user创建
		userId, err := s.UserMapper.InsertWithEcho(ctx, &usrMapper.User{
			Password: password,
			Name:     user.Name,
			Gender:   user.Gender,
		})

		if err != nil {
			logx.Error("创建用户失败。studentId = %s, unitId = %s\n", studentId, unitId)
			return
		}

		// 创建后进行关联
		if s.UUMapper.Insert(ctx, &uuMapper.UnitUser{
			UnitId:    unitId,
			UserId:    userId,
			StudentId: studentId,
			Options:   option,
		}) != nil {
			logx.Error("创建用户关联失败。userId = %s, unitId = %s, studentId = %s\n", userId, unitId, studentId)
			return
		}
		logx.Info("创建用户关联成功！userId = %s, unitId = %s, studentId = %s\n", userId, unitId, studentId)
	} else if err == nil && link != nil {
		// 如果已经存在，则直接跳过
		logx.Info("该用户已经创建且关联。userId = %s, unitId = %s, studentId = %s,\n", link.UserId, unitId, studentId)
	} else {
		logx.Error("查询用户失败。studentId = %s, unitId = %s\n", studentId, unitId)
	}
}

func (s *UnitService) UnitUpdateInfo(ctx context.Context, req *u.UnitUpdateInfoReq) (res *basic.Response, err error) {
	// 不允许修改手机号、密码、验证方式、level
	// 密码、验证方式需要通过其他接口修改
	unitId, err := primitive.ObjectIDFromHex(req.Unit.Id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}
	unit := &untmapper.Unit{
		ID:         unitId,
		Name:       req.Unit.Name,
		Address:    req.Unit.Address,
		Contact:    req.Unit.Contact,
		Status:     req.Unit.Status,
		UpdateTime: time.Time{},
	}
	if err := s.UnitMapper.UpdateBasicInfo(ctx, unit); err != nil {
		return nil, err
	}

	return result.ResponseOk(), nil
}

func (s *UnitService) UnitUpdatePassword(ctx context.Context, req *u.UnitUpdatePasswordReq) (res *basic.Response, err error) {
	authType := req.AuthType
	unitId := req.Id
	flag := false

	// 校验
	switch authType {
	case consts.UpdateByOldPwd:
		// 旧密码
		unit, err := s.UnitMapper.FindOne(ctx, unitId)
		if err != nil {
			return nil, err
		}

		// 校验密码
		oldPwd := req.AuthId
		if !encrypt.BcryptCheck(oldPwd, unit.Password) {
			return nil, consts.ErrUnitPasswordMismatch
		}
		flag = true
	case consts.UpdateByCode:
		// TODO: 验证码
		unit, err := s.UnitMapper.FindOne(ctx, unitId)
		if err != nil {
			return nil, err
		}
		phone := unit.Phone
		logx.Info(phone)
		// flag == true
	}

	if flag {
		// 修改密码
		if err := s.UnitMapper.UpdatePassword(ctx, req.Id, req.NewPassword); err != nil {
			return nil, err
		}
		return result.ResponseOk(), nil
	}
	return nil, err
}

func (s *UnitService) UnitCreateVerify(ctx context.Context, req *u.UnitCreateVerifyReq) (res *u.UnitCreateVerifyResp, err error) {
	// 创建和修改二合一接口
	/*	unitId := req.UnitVerify.UserId
		verifyType := req.UnitVerify.VerifyType
		verify := &untmapper.UnitVerify{
			UnitId:     unitId,
			VerifyType: verifyType,
		}

		// 如果是弱验证，判断是否给出了 account 和 password，否则account随机生成，password默认
		if verifyType == consts.Weak {

		}*/

	return nil, err
}

// 生成不重复账号
func (s *UnitService) generateUniqueAccount(ctx context.Context, collection *mongo.Collection) (string, error) {
	/*	for i := 0; i < 10; i++ { // 最多尝试10次
			account, err := random.GenerateRandomAccount()
			if err != nil {
				return "", err
			}

			// 检查数据库中是否存在
			count, err := s.UnitMapper.
			if err != nil {
				return "", err
			}

			if count == 0 {
				return account, nil // 找到了唯一账号
			}
		}
		logx.Info("生成唯一账号失败，请重试")*/
	return "", consts.ErrUnitVerify
}

func (s *UnitService) UnitUpdateVerifyPassword(ctx context.Context, req *u.UnitUpdateVerifyReq) (res *basic.Response, err error) {
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
