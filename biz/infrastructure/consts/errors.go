package consts

import (
	"errors"
	"github.com/xh-polaris/psych-pkg/errorx"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Errno struct {
	err  error
	code codes.Code
}

// GRPCStatus 实现 GRPCStatus 方法
func (en *Errno) GRPCStatus() *status.Status {
	return status.New(en.code, en.err.Error())
}

// 实现 Error 方法
func (en *Errno) Error() string {
	return en.err.Error()
}

// NewErrno 创建自定义错误
func NewErrno(code codes.Code, err error) *Errno {
	return &Errno{
		err:  err,
		code: code,
	}
}

// 定义常量错误
// TODO: 定义错误常量
var (
	// User模块相关错误码
	// User模块相关错误码
	ErrUserSignUp             = errorx.New(10001, "用户注册失败，请重试")
	ErrUserPhoneExist         = errorx.New((10002), "该手机号已被注册")
	ErrUserVerify             = errorx.New((10003), "验证码错误")
	ErrUserPasswordMismatch   = errorx.New((10004), "密码不匹配")
	ErrUserNotExist           = errorx.New((10005), "用户账号不存在")
	ErrUserGetInfo            = errorx.New((10006), "获取用户信息失败")
	ErrUserSignIn             = errorx.New((10007), "用户登录失败")
	ErrUserStudentIdNotExist  = errorx.New((10008), "学号不存在")
	ErrUserStudentIdDuplicate = errorx.New((10009), "学号对应多个用户")

	// Unit模块相关错误码
	ErrUnitSignUp           = errorx.New(12001, "单位注册失败，请重试")
	ErrUnitPhoneExist       = errorx.New(12002, "该手机号已被注册为单位账号")
	ErrUnitVerify           = errorx.New(12003, "单位账号验证失败")
	ErrUnitPasswordMismatch = errorx.New(12004, "密码不匹配")
	ErrUnitNotExist         = errorx.New(12005, "单位账号不存在")
	ErrUnitLinkUser         = errorx.New(12006, "关联用户失败")
	ErrUnitGetInfo          = errorx.New(12007, "获取单位信息失败")
	ErrUnitCreateUser       = errorx.New(12008, "创建用户失败")
)

// ErrInvalidParams 调用时错误
var (
	ErrInvalidParams = NewErrno(codes.InvalidArgument, errors.New("参数错误"))
	ErrCall          = NewErrno(codes.Unknown, errors.New("调用接口失败，请重试"))
)

// 数据库相关错误
var (
	ErrNotFound        = NewErrno(codes.NotFound, errors.New("not found"))
	ErrInvalidObjectId = NewErrno(codes.InvalidArgument, errors.New("无效的id "))
	ErrUpdate          = NewErrno(codes.Code(2001), errors.New("更新失败"))
)
