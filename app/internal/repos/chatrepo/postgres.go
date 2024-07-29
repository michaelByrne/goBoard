package chatrepo

import (
	"context"
	"go.uber.org/zap"
	"goBoard/internal/core/domain"
	"goBoard/internal/repos/chatrepo/db"
)

type ChatRepo struct {
	queries *db.Queries

	logger *zap.Logger
}

func NewChatRepo(queries *db.Queries, logger *zap.Logger) *ChatRepo {
	return &ChatRepo{
		queries: queries,
		logger:  logger,
	}
}

func (r *ChatRepo) GetChatGroupsForMember(ctx context.Context, memberID int) ([]domain.ChatGroup, error) {
	chatGroups, err := r.queries.GetChatGroupsForMember(ctx, int32(memberID))
	if err != nil {
		r.logger.Error("failed to get chat groups for member", zap.Error(err))
		return nil, err
	}

	return fromDBToDomainChatGroups(chatGroups), nil
}

func (r *ChatRepo) InsertChatGroup(ctx context.Context, chatGroup domain.ChatGroup) (int, error) {
	dbChatGroup := fromDomainToDBChatGroup(chatGroup)
	chatGroupOut, err := r.queries.InsertChatGroup(ctx, dbChatGroup.Topic)
	if err != nil {
		r.logger.Error("failed to insert chat group", zap.Error(err))
		return 0, err
	}

	return int(chatGroupOut.ID), nil
}
