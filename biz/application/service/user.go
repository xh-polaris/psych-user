package service

import (
	"context"
	"github.com/xh-polaris/psych-idl/kitex_gen/basic"
	"github.com/xh-polaris/psych-idl/kitex_gen/user"
	"github.com/xh-polaris/psych-idl/kitex_gen/user/user"
	"github.com/xh-polaris/psych-user/biz/infrastructure/mapper/user"
)

type UserServiceImpl interface {
	UserSignUp(ctx context.Context, req *UserSignUpReq) (res *basic.Response, err error)
	UserGetInfo(ctx context.Context, req *UserGetInfoReq) (res *UserGetInfoResp, err error)
	UserUpdateInfo(ctx context.Context, req *UserUpdateInfoReq) (res *basic.Response, err error)
	UserUpdatePassword(ctx context.Context, req *UserUpdatePasswordReq) (res *basic.Response, err error)
	UserBelongUnit(ctx context.Context, req *UserBelongUnitReq) (res *UserBelongUnitResp, err error)
	UserSignIn(ctx context.Context, req *UserSignInReq) (res *basic.Response, err error)
}

type UserService struct {
	UserMapper *user.MongoMapper
	// UnitMapper
}

func (s *UserService) UserSignUp(ctx context.Context, req *UserSignUpReq) (res *basic.Response, err error) {
	return nil, err
}

func (s *UserService) UserGetInfo(ctx context.Context, req *UserGetInfoReq) (res *UserGetInfoResp, err error) {
	// meta?
	u, err := s.UserMapper.FindOne(ctx)
	return nil, err
}
func (s *UserService) UserUpdateInfo(ctx context.Context, req *UserUpdateInfoReq) (res *basic.Response, err error) {
	return nil, err
}
func (s *UserService) UserUpdatePassword(ctx context.Context, req *UserUpdatePasswordReq) (res *basic.Response, err error) {
	return nil, err
}
func (s *UserService) UserBelongUnit(ctx context.Context, req *UserBelongUnitReq) (res *UserBelongUnitResp, err error) {
	return nil, err
}
func (s *UserService) UserSignIn(ctx context.Context, req *UserSignInReq) (res *basic.Response, err error) {
	return nil, err
}
