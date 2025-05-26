package adaptor

import (
	"context"
	"errors"

	"encoding/json"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v4"
	"github.com/xh-polaris/gopkg/util"
	"github.com/xh-polaris/gopkg/util/log"
	"github.com/xh-polaris/psych-idl-gen/kitex_gen/basic"
	"github.com/xh-polaris/psych-user/biz/infrastructure/config"
)

const hertzContext = "hertz_context"

func InjectContext(ctx context.Context, c *app.RequestContext) context.Context {
	return context.WithValue(ctx, hertzContext, c)
}

func ExtractContext(ctx context.Context) (*app.RequestContext, error) {
	c, ok := ctx.Value(hertzContext).(*app.RequestContext)
	if !ok {
		return nil, errors.New("hertz context not found")
	}
	return c, nil
}

func ExtractMeta(ctx context.Context) (*basic.UserMeta, *basic.Extra) {
	return ExtractUserMeta(ctx), ExtractExtra(ctx)
}

func ExtractUserMeta(ctx context.Context) (user *basic.UserMeta) {
	user = new(basic.UserMeta)
	var err error
	defer func() {
		if err != nil {
			log.CtxInfo(ctx, "extract user meta fail, err=%v", err)
		}
	}()
	c, err := ExtractContext(ctx)
	if err != nil {
		return
	}
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.Parse(string(tokenString), func(_ *jwt.Token) (interface{}, error) {
		return jwt.ParseECPublicKeyFromPEM([]byte(config.GetConfig().Auth.PublicKey))
	})
	if err != nil {
		return
	}
	if !token.Valid {
		err = errors.New("token is not valid")
		return
	}
	data, err := json.Marshal(token.Claims)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, user)
	if err != nil {
		return
	}
	if user.SessionUserId == "" {
		user.SessionUserId = user.UserId
	}
	if user.SessionAppId == 0 {
		user.SessionAppId = user.AppId
	}
	if user.SessionDeviceId == "" {
		user.SessionDeviceId = user.DeviceId
	}
	log.CtxInfo(ctx, "userMeta=%s", util.JSONF(user))
	return
}

func ExtractExtra(ctx context.Context) (extra *basic.Extra) {
	extra = new(basic.Extra)
	var err error
	defer func() {
		if err != nil {
			log.CtxInfo(ctx, "extract extra fail, err=%v", err)
		}
	}()
	c, err := ExtractContext(ctx)
	if err != nil {
		return
	}
	extra.ClientIP = c.ClientIP()
	err = c.Bind(extra)
	if err != nil {
		return
	}
	log.CtxInfo(ctx, "extra=%s", util.JSONF(extra))
	return
}
