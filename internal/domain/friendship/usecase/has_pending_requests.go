package usecase

import "context"

func (u *usecase) HasPendingRequest(ctx context.Context, profileFromID uint64, profileToID uint64) (bool, error) {
	requests, err := u.repo.ListPendingRequests(ctx, profileToID)
	if err != nil {
		return false, err
	}
	for _, requestedProfile := range requests {
		if requestedProfile == profileFromID {
			return true, nil
		}
	}
	return false, nil
}
