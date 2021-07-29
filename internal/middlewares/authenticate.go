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
				log.Warnw("got error while decoding token, resetting it", "err", err)
				http.SetCookie(w, cookies.EmptyAuthCookie)
				next.ServeHTTP(w, r)
				return
			}

			profilesByID, err := profilesUseCase.GetByID(ctx, token.ProfileID)
			if err != nil || len(profilesByID) == 0 {
				log.Errorw("got error while fetching profilesByID by id in auth MW", "err", err)
				http.SetCookie(w, cookies.EmptyAuthCookie)
				next.ServeHTTP(w, r)
				return
			}
			profile := profilesByID[0]
			if profile.UserID != token.UserID {
				log.Warnw("profile belongs to user different from one in auth token, resetting token",
					"profile", profile.ID,
					"expected_user", profile.UserID,
					"user", token.UserID)
				http.SetCookie(w, cookies.EmptyAuthCookie)
				next.ServeHTTP(w, r)
				return
			}

			ctx = AddProfileToCtx(ctx, profile)
			r = r.WithContext(ctx)
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
