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
	"github.com/xh-polaris/psych-user/biz/infrastructure/util/encrypt"
	"github.com/xh-polaris/psych-user/biz/infrastructure/util/reg"
	"github.com/xh-polaris/psych-user/biz/infrastructure/util/result"
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
	UnitMapper *untmapper.MongoMapper
	UserMapper *usrMapper.MongoMapper
	UUMapper   *uuMapper.MongoMapper
}

var UnitServiceSet = wire.NewSet(
	wire.Struct(new(UnitService), "*"),
	wire.Bind(new(IUnitService), new(*UnitService)),
)

// UnitSignUp 单位账号注册
func (s *UnitService) UnitSignUp(ctx context.Context, req *u.UnitSignUpReq) (res *basic.Response, err error) {
	res = &basic.Response{}

	// 参数校验
	if req.Unit == nil || req.Unit.Phone == "" || req.Unit.Password == "" || req.Unit.Name == "" || !reg.CheckMobile(req.Unit.Phone) {
		logx.Error("UnitSignUp fail")
		return nil, consts.ErrUnitSignUp
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
	err = s.UnitMapper.Insert(ctx, unit)
	if err != nil {
		return nil, consts.ErrUnitSignUp
	}

	// 返回成功响应
	return result.ResponseOk(), nil
}

// UnitStrongVerify 单位账号认证
func (s *UnitService) UnitStrongVerify(ctx context.Context, req *u.UnitStrongVerifyReq) (res *basic.Response, err error) {
	res = &basic.Response{}

	// 参数校验
	if req.Phone == "" || req.GetPhone() == "" || !reg.CheckMobile(req.Phone) {
		return nil, consts.ErrInvalidParams
	}

	// TODO: 验证码功能
	if req.GetVerifyCode() != "" {
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
	res.Code = 200
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
			Id:         unit.ID.Hex(),
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
	return result.ResponseOk(), nil
}

// UnitCreateAndLinkUser 批量创建并关联用户
func (s *UnitService) UnitCreateAndLinkUser(ctx context.Context, req *u.UnitCreateAndLinkUserReq) (res *basic.Response, err error) {
	authType := req.AuthType
	unitId := req.UnitId
	password, _ := encrypt.BcryptEncrypt(consts.DefaultPassword)
	for i := 0; i < len(req.GetAuthId()); i++ {
		authId := req.AuthId[i]
		name := req.UserName[i]
		gender := req.Gender[i]
		options := req.GetOptions()
		var option *u.Option = nil
		if options != nil {
			option = options[i]
		}

		// 判断 authType
		switch authType {
		case consts.AuthPhone:
			{
				// 关联手机号
				// TODO: 开启事务?
				// 先查找该手机号是否已经注册
				logx.Info("正在查询用户: phone = %s, unitId = %s\n", authId, unitId)
				user, err := s.UserMapper.FindOneByPhone(ctx, authId)
				if err != nil && !errors.Is(err, consts.ErrNotFound) {
					logx.Error("查询用户失败。phone = %s, unitId = %s\n", authId, unitId)
					continue
				}

				if errors.Is(err, consts.ErrNotFound) {
					// 未被注册，先在user表中添加数据
					user = &usrMapper.User{
						Phone:    authId,
						Name:     name,
						Gender:   gender,
						Password: password,
					}
					userId, err := s.UserMapper.InsertWithEcho(ctx, user)
					if err != nil {
						logx.Error("创建用户失败。phone = %s, unitId = %s\n", authId, unitId)
						continue
					}
					// 再在unit_user表中关联
					err = s.UUMapper.Insert(ctx, &uuMapper.UnitUser{
						UnitId:  unitId,
						UserId:  *userId,
						Options: option,
					})
					if err != nil {
						logx.Error("创建用户关联失败。userId = %s, unitId = %s, phone = %s\n", userId, unitId, authId)
						continue
					}
					logx.Info("创建用户关联成功！userId = %s, unitId = %s, phone = %s\n", userId, unitId, authId)
				} else {
					// 已经被注册，则判断是否已经被关联，若无关联则关联，已关联则跳过（不报错）
					userId := user.ID.Hex()
					r, err := s.UUMapper.FindOneByUU(ctx, userId, unitId)
					if err != nil {
						logx.Error("查询用户失败。phone = %s, unitId = %s\n", authId, unitId)
						continue
					}
					// 已经有关联，跳过
					if r != nil {
						logx.Info("用户和该单位已经绑定。userId = %s, unitId = %s, phone = %s\n", userId, unitId, authId)
						continue
					}
					// 无关联，则插入
					err = s.UUMapper.Insert(ctx, &uuMapper.UnitUser{
						UnitId:  unitId,
						UserId:  userId,
						Options: option,
					})
					if err != nil {
						logx.Error("创建用户关联失败。userId = %s, unitId = %s, phone = %s\n", userId, unitId, authId)
						continue
					}
					logx.Info("创建用户关联成功！userId = %s, unitId = %s, phone = %s\n", userId, unitId, authId)
				}
			}
		case consts.AuthStudentId:
			{
				// 关联学号
				// 先根据学号(authId -> studentId)和unitId查找是否已经存在
				// 此时如果记录存在，则表示已经有user账号且关联完成
				unitUser, err := s.UUMapper.FindOneByUnitAndStu(ctx, authId, unitId)
				if err != nil && !errors.Is(err, consts.ErrNotFound) {
					logx.Error("查询用户失败。studentId = %s, unitId = %s\n", authId, unitId)
					continue
				}

				if unitUser != nil {
					// 如果已经存在，则直接跳过
					logx.Info("该用户已经创建且关联。userId = %s, unitId = %s, studentId = %s,\n", unitUser.UserId, unitId, authId)
					continue
				}

				// 如果不存在，则先进行user创建
				userId, err := s.UserMapper.InsertWithEcho(ctx, &usrMapper.User{
					Password: password,
					Name:     name,
					Gender:   gender,
				})
				if err != nil {
					logx.Error("创建用户失败。studentId = %s, unitId = %s\n", authId, unitId)
					continue
				}

				// 创建后进行关联
				err = s.UUMapper.Insert(ctx, &uuMapper.UnitUser{
					UnitId:    unitId,
					UserId:    *userId,
					StudentId: authId,
					Options:   option,
				})
				if err != nil {
					logx.Error("创建用户关联失败。userId = %s, unitId = %s, studentId = %s\n", userId, unitId, authId)
					continue
				}
				logx.Info("创建用户关联成功！userId = %s, unitId = %s, studentId = %s\n", userId, unitId, authId)
			}
		}
	}
	return result.ResponseOk(), nil
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
