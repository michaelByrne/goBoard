package chatrepo

import (
	"database/sql"
	"goBoard/internal/core/domain"
	"goBoard/internal/repos/chatrepo/db"
)

func fromDomainToDBChatGroup(chatGroup domain.ChatGroup) db.ChatGroup {
	return db.ChatGroup{
		ID:    int32(chatGroup.ID),
		Topic: sql.NullString{String: chatGroup.Topic, Valid: true},
	}
}

func fromDBToDomainChatGroup(chatGroup db.ChatGroup) domain.ChatGroup {
	return domain.ChatGroup{
		ID:    int(chatGroup.ID),
		Topic: chatGroup.Topic.String,
	}
}

func fromDomainToDBChatGroups(chatGroups []domain.ChatGroup) []db.ChatGroup {
	var dbChatGroups []db.ChatGroup
	for _, chatGroup := range chatGroups {
		dbChatGroups = append(dbChatGroups, fromDomainToDBChatGroup(chatGroup))
	}
	return dbChatGroups
}

func fromDBToDomainChatGroups(chatGroups []db.ChatGroup) []domain.ChatGroup {
	var domainChatGroups []domain.ChatGroup
	for _, chatGroup := range chatGroups {
		domainChatGroups = append(domainChatGroups, fromDBToDomainChatGroup(chatGroup))
	}
	return domainChatGroups
}
