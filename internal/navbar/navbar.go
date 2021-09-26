package navbar

import (
	"context"

	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/domain/chats"
	"github.com/vvkh/social-network/internal/domain/friendship"
	"github.com/vvkh/social-network/internal/middlewares"
)

type Navbar struct {
	friendshipUC friendship.UseCase
	chatsUC      chats.UseCase
	log          *zap.SugaredLogger
}

func New(log *zap.SugaredLogger, friendshipUC friendship.UseCase, chatsUC chats.UseCase) *Navbar {
	return &Navbar{
		friendshipUC: friendshipUC,
		chatsUC:      chatsUC,
		log:          log,
	}
}

func (n *Navbar) GetContext(ctx context.Context) *Context {
	navbarContext := Context{}

	profile, ok := middlewares.ProfileFromCtx(ctx)
	if !ok {
		return nil
	}
	navbarContext.SelfProfileID = &profile.ID

	pendingRequests, err := n.friendshipUC.ListPendingRequests(ctx, profile.ID)
	if err == nil && len(pendingRequests) > 0 {
		requestCount := len(pendingRequests)
		navbarContext.PendingFriendshipRequestsCount = &requestCount
	}

	if n.chatsUC != nil {
		messagesCount, err := n.chatsUC.GetUnreadMessagesCount(ctx, profile.ID)
		if err == nil && messagesCount > 0 {
			navbarContext.UnreadMessagesCount = &messagesCount
		}
	}

	return &navbarContext
}
