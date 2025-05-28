package service

import (
	"context"
	"github.com/google/wire"

	"github.com/xh-polaris/psych-idl/kitex_gen/basic"
	"github.com/xh-polaris/psych-idl/kitex_gen/unit"
	"github.com/xh-polaris/psych-user/biz/infrastructure/mapper/unit"
)

type IUnitService interface {
	UnitSignUp(ctx context.Context, req *unit.UnitSignUpReq) (res *basic.Response, err error)
	UnitGetInfo(ctx context.Context, req *unit.UnitGetInfoReq) (res *unit.UnitGetInfoResp, err error)
	UnitUpdateInfo(ctx context.Context, req *unit.UnitUpdateInfoReq) (res *basic.Response, err error)
	UnitUpdatePassword(ctx context.Context, req *unit.UnitUpdatePasswordReq) (res *basic.Response, err error)
	UnitStrongVerify(ctx context.Context, req *unit.UnitStrongVerifyReq) (res *basic.Response, err error)
	UnitWeakVerify(ctx context.Context, req *unit.UnitWeakVerifyReq) (res *basic.Response, err error)
	UnitCreateVerify(ctx context.Context, req *unit.UnitCreateVerifyReq) (res *unit.UnitCreateVerifyResp, err error)
	UnitUpdateVerifyPassword(ctx context.Context, req *unit.UnitUpdateVerifyPasswordReq) (res *basic.Response, err error)
	UnitLinkUser(ctx context.Context, req *unit.UnitLinkUserReq) (res *basic.Response, err error)
	UnitLinkView(ctx context.Context, req *unit.UnitLinkViewReq) (res *basic.Response, err error)
	UnitPageQueryUser(ctx context.Context, req *unit.UnitPageQueryUserReq) (res *unit.UnitPageQueryUserResp, err error)
	UnitPageQueryView(ctx context.Context, req *unit.UnitPageQueryViewReq) (res *unit.UnitPageQueryViewResp, err error)
	UnitCreateAndLinkUser(ctx context.Context, req *unit.UnitCreateAndLinkUserReq) (res *basic.Response, err error)
	UnitCreateAndLinkView(ctx context.Context, req *unit.UnitCreateAndLinkViewReq) (res *basic.Response, err error)
	UnitGetAppInfo(ctx context.Context, req *unit.UnitGetAppInfoReq) (res *unit.UnitGetAppInfoResp, err error)
	UnitModelGetInfo(ctx context.Context, req *unit.UnitModelGetInfoReq) (res *unit.UnitModelGetInfoResp, err error)
	UnitModelUpdateInfo(ctx context.Context, req *unit.UnitModelUpdateInfoReq) (res *basic.Response, err error)
}

type UnitService struct {
	UnitMapper *unit.MongoMapper
}

var PsychUnitSet = wire.NewSet()

func (s *UnitService) UnitSignUp(ctx context.Context, req *unit.UnitSignUpReq) (res *basic.Response, err error) {
	return nil, err
}

func (s *UnitService) UnitGetInfo(ctx context.Context, req *unit.UnitGetInfoReq) (res *unit.UnitGetInfoResp, err error) {
	return nil, err
}

func (s *UnitService) UnitUpdateInfo(ctx context.Context, req *unit.UnitUpdateInfoReq) (res *basic.Response, err error) {
	return nil, err
}

func (s *UnitService) UnitUpdatePassword(ctx context.Context, req *unit.UnitUpdatePasswordReq) (res *basic.Response, err error) {
	return nil, err
}

func (s *UnitService) UnitStrongVerify(ctx context.Context, req *unit.UnitStrongVerifyReq) (res *basic.Response, err error) {
	return nil, err
}

func (s *UnitService) UnitWeakVerify(ctx context.Context, req *unit.UnitWeakVerifyReq) (res *basic.Response, err error) {
	return nil, err
}

func (s *UnitService) UnitCreateVerify(ctx context.Context, req *unit.UnitCreateVerifyReq) (res *unit.UnitCreateVerifyResp, err error) {
	return nil, err
}

func (s *UnitService) UnitUpdateVerifyPassword(ctx context.Context, req *unit.UnitUpdateVerifyPasswordReq) (res *basic.Response, err error) {
	return nil, err
}

func (s *UnitService) UnitLinkUser(ctx context.Context, req *unit.UnitLinkUserReq) (res *basic.Response, err error) {
	return nil, err
}

func (s *UnitService) UnitLinkView(ctx context.Context, req *unit.UnitLinkViewReq) (res *basic.Response, err error) {
	return nil, err
}

func (s *UnitService) UnitPageQueryUser(ctx context.Context, req *unit.UnitPageQueryUserReq) (res *unit.UnitPageQueryUserResp, err error) {
	return nil, err
}

func (s *UnitService) UnitPageQueryView(ctx context.Context, req *unit.UnitPageQueryViewReq) (res *unit.UnitPageQueryViewResp, err error) {
	return nil, err
}

func (s *UnitService) UnitCreateAndLinkUser(ctx context.Context, req *unit.UnitCreateAndLinkUserReq) (res *basic.Response, err error) {
	return nil, err
}

func (s *UnitService) UnitCreateAndLinkView(ctx context.Context, req *unit.UnitCreateAndLinkViewReq) (res *basic.Response, err error) {
	return nil, err
}

func (s *UnitService) UnitGetAppInfo(ctx context.Context, req *unit.UnitGetAppInfoReq) (res *unit.UnitGetAppInfoResp, err error) {
	return nil, err
}

func (s *UnitService) UnitModelGetInfo(ctx context.Context, req *unit.UnitModelGetInfoReq) (res *unit.UnitModelGetInfoResp, err error) {
	return nil, err
}

func (s *UnitService) UnitModelUpdateInfo(ctx context.Context, req *unit.UnitModelUpdateInfoReq) (res *basic.Response, err error) {
	return nil, err
}
