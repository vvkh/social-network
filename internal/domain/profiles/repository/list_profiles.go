package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"

	"github.com/vvkh/social-network/internal/domain/profiles"
	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/domain/profiles/repository/dto"
)

func (r *repo) ListProfiles(ctx context.Context, firstNamePrefix string, lastNamePrefix string, limit int) ([]entity.Profile, bool, error) {
	var query = goqu.Dialect("mysql").From(`profiles`).Select("*")

	if lastNamePrefix != "" {
		query = query.Where(goqu.Ex{
			"last_name": goqu.Op{
				"ILIKE": lastNamePrefix + "%",
			},
		})
	}

	if firstNamePrefix != "" {
		query = query.Where(goqu.Ex{
			"first_name": goqu.Op{
				"ILIKE": firstNamePrefix + "%",
			},
		})
	}

	appliedLimit := false
	if limit != profiles.ShowAllProfiles {
		appliedLimit = true
		// we request "limit" + 1 so that later we would know if there are more profiles to show
		// if number greater than "limit" returned then we have at least one more profile
		query = query.Limit(uint(limit) + 1)
	}

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, false, err
	}

	var profileDtos []dto.Profile
	if err = r.db.SelectContext(ctx, &profileDtos, sql, args...); err != nil {
		return nil, false, err
	}

	var hasMore bool
	if appliedLimit && len(profileDtos) > limit {
		profileDtos = profileDtos[:limit]
		hasMore = true
	} else {
		hasMore = false
	}
	return dto.ToProfiles(profileDtos), hasMore, nil
}
