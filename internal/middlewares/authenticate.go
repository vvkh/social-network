package middlewares

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/cookies"
	"github.com/vvkh/social-network/internal/domain/profiles"
	profilesEntity "github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/domain/users"
)

type ctxKey int

const (
	CtxKeyProfile = ctxKey(1)
)

func AuthenticateUser(log *zap.SugaredLogger, usersUseCase users.UseCase, profilesUseCase profiles.UseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			encodedToken, err := cookies.ReadAuthCookie(r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()
			token, err := usersUseCase.DecodeToken(ctx, encodedToken.Value)
			if err != nil {
				log.Warn("got error while decoding token, resetting it", "err", err)
				http.SetCookie(w, cookies.EmptyAuthCookie)
				next.ServeHTTP(w, r)
				return
			}

			profiles, err := profilesUseCase.GetByUserID(ctx, token.UserID)
			if err != nil {
				log.Error("got error while fetching user by id in auth MW", "err", err)
				// we are not sure that profile doesn't exist, maybe it's just some temporary problem
				next.ServeHTTP(w, r)
				return
			}

			for _, profile := range profiles {
				if profile.UserID == token.UserID && profile.ID == token.ProfileID {
					ctx = AddProfileToCtx(ctx, profile)
					r = r.WithContext(ctx)
					next.ServeHTTP(w, r)
					return
				}
			}
			http.SetCookie(w, cookies.EmptyAuthCookie)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func ProfileFromCtx(ctx context.Context) (profilesEntity.Profile, bool) {
	profile, ok := ctx.Value(CtxKeyProfile).(profilesEntity.Profile)
	return profile, ok
}

func AddProfileToCtx(ctx context.Context, profile profilesEntity.Profile) context.Context {
	return context.WithValue(ctx, CtxKeyProfile, profile)
}
