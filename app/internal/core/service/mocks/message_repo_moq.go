// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"sync"
	"time"
)

// Ensure, that MessageRepoMock does implement ports.MessageRepo.
// If this is not the case, regenerate this file with moq.
var _ ports.MessageRepo = &MessageRepoMock{}

// MessageRepoMock is a mock implementation of ports.MessageRepo.
//
//	func TestSomethingThatUsesMessageRepo(t *testing.T) {
//
//		// make and configure a mocked ports.MessageRepo
//		mockedMessageRepo := &MessageRepoMock{
//			DeleteMessageFunc: func(ctx context.Context, memberID int, messageID int) error {
//				panic("mock out the DeleteMessage method")
//			},
//			GetMessageByIDFunc: func(ctx context.Context, messageID int, memberID int) (*domain.Message, error) {
//				panic("mock out the GetMessageByID method")
//			},
//			GetMessageParticipantsFunc: func(ctx context.Context, messageID int) ([]string, error) {
//				panic("mock out the GetMessageParticipants method")
//			},
//			GetMessagePostByIDFunc: func(id int) (*domain.MessagePost, error) {
//				panic("mock out the GetMessagePostByID method")
//			},
//			GetMessagePostsByIDFunc: func(memberID int, messageID int, limit int) ([]domain.MessagePost, error) {
//				panic("mock out the GetMessagePostsByID method")
//			},
//			GetMessagePostsCollapsibleFunc: func(ctx context.Context, viewable int, messageID int, memberID int) ([]domain.MessagePost, int, error) {
//				panic("mock out the GetMessagePostsCollapsible method")
//			},
//			GetMessagesWithCursorForwardFunc: func(memberID int, limit int, cursor *time.Time) ([]domain.Message, error) {
//				panic("mock out the GetMessagesWithCursorForward method")
//			},
//			GetMessagesWithCursorReverseFunc: func(memberID int, limit int, cursor *time.Time) ([]domain.Message, error) {
//				panic("mock out the GetMessagesWithCursorReverse method")
//			},
//			GetNewMessageCountsFunc: func(ctx context.Context, memberID int) (*domain.MessageCounts, error) {
//				panic("mock out the GetNewMessageCounts method")
//			},
//			ListMessagesFunc: func(ctx context.Context, cursors domain.Cursors, limit int, memberID int) ([]domain.Message, domain.Cursors, error) {
//				panic("mock out the ListMessages method")
//			},
//			PeekPreviousFunc: func(timestamp *time.Time) (bool, error) {
//				panic("mock out the PeekPrevious method")
//			},
//			SaveMessageFunc: func(message domain.Message) (int, error) {
//				panic("mock out the SaveMessage method")
//			},
//			SavePostFunc: func(post domain.MessagePost) (int, error) {
//				panic("mock out the SavePost method")
//			},
//			ViewMessageFunc: func(ctx context.Context, memberID int, messageID int) (int, error) {
//				panic("mock out the ViewMessage method")
//			},
//		}
//
//		// use mockedMessageRepo in code that requires ports.MessageRepo
//		// and then make assertions.
//
//	}
type MessageRepoMock struct {
	// DeleteMessageFunc mocks the DeleteMessage method.
	DeleteMessageFunc func(ctx context.Context, memberID int, messageID int) error

	// GetMessageByIDFunc mocks the GetMessageByID method.
	GetMessageByIDFunc func(ctx context.Context, messageID int, memberID int) (*domain.Message, error)

	// GetMessageParticipantsFunc mocks the GetMessageParticipants method.
	GetMessageParticipantsFunc func(ctx context.Context, messageID int) ([]string, error)

	// GetMessagePostByIDFunc mocks the GetMessagePostByID method.
	GetMessagePostByIDFunc func(id int) (*domain.MessagePost, error)

	// GetMessagePostsByIDFunc mocks the GetMessagePostsByID method.
	GetMessagePostsByIDFunc func(memberID int, messageID int, limit int) ([]domain.MessagePost, error)

	// GetMessagePostsCollapsibleFunc mocks the GetMessagePostsCollapsible method.
	GetMessagePostsCollapsibleFunc func(ctx context.Context, viewable int, messageID int, memberID int) ([]domain.MessagePost, int, error)

	// GetMessagesWithCursorForwardFunc mocks the GetMessagesWithCursorForward method.
	GetMessagesWithCursorForwardFunc func(memberID int, limit int, cursor *time.Time) ([]domain.Message, error)

	// GetMessagesWithCursorReverseFunc mocks the GetMessagesWithCursorReverse method.
	GetMessagesWithCursorReverseFunc func(memberID int, limit int, cursor *time.Time) ([]domain.Message, error)

	// GetNewMessageCountsFunc mocks the GetNewMessageCounts method.
	GetNewMessageCountsFunc func(ctx context.Context, memberID int) (*domain.MessageCounts, error)

	// ListMessagesFunc mocks the ListMessages method.
	ListMessagesFunc func(ctx context.Context, cursors domain.Cursors, limit int, memberID int) ([]domain.Message, domain.Cursors, error)

	// PeekPreviousFunc mocks the PeekPrevious method.
	PeekPreviousFunc func(timestamp *time.Time) (bool, error)

	// SaveMessageFunc mocks the SaveMessage method.
	SaveMessageFunc func(message domain.Message) (int, error)

	// SavePostFunc mocks the SavePost method.
	SavePostFunc func(post domain.MessagePost) (int, error)

	// ViewMessageFunc mocks the ViewMessage method.
	ViewMessageFunc func(ctx context.Context, memberID int, messageID int) (int, error)

	// calls tracks calls to the methods.
	calls struct {
		// DeleteMessage holds details about calls to the DeleteMessage method.
		DeleteMessage []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// MemberID is the memberID argument value.
			MemberID int
			// MessageID is the messageID argument value.
			MessageID int
		}
		// GetMessageByID holds details about calls to the GetMessageByID method.
		GetMessageByID []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// MessageID is the messageID argument value.
			MessageID int
			// MemberID is the memberID argument value.
			MemberID int
		}
		// GetMessageParticipants holds details about calls to the GetMessageParticipants method.
		GetMessageParticipants []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// MessageID is the messageID argument value.
			MessageID int
		}
		// GetMessagePostByID holds details about calls to the GetMessagePostByID method.
		GetMessagePostByID []struct {
			// ID is the id argument value.
			ID int
		}
		// GetMessagePostsByID holds details about calls to the GetMessagePostsByID method.
		GetMessagePostsByID []struct {
			// MemberID is the memberID argument value.
			MemberID int
			// MessageID is the messageID argument value.
			MessageID int
			// Limit is the limit argument value.
			Limit int
		}
		// GetMessagePostsCollapsible holds details about calls to the GetMessagePostsCollapsible method.
		GetMessagePostsCollapsible []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Viewable is the viewable argument value.
			Viewable int
			// MessageID is the messageID argument value.
			MessageID int
			// MemberID is the memberID argument value.
			MemberID int
		}
		// GetMessagesWithCursorForward holds details about calls to the GetMessagesWithCursorForward method.
		GetMessagesWithCursorForward []struct {
			// MemberID is the memberID argument value.
			MemberID int
			// Limit is the limit argument value.
			Limit int
			// Cursor is the cursor argument value.
			Cursor *time.Time
		}
		// GetMessagesWithCursorReverse holds details about calls to the GetMessagesWithCursorReverse method.
		GetMessagesWithCursorReverse []struct {
			// MemberID is the memberID argument value.
			MemberID int
			// Limit is the limit argument value.
			Limit int
			// Cursor is the cursor argument value.
			Cursor *time.Time
		}
		// GetNewMessageCounts holds details about calls to the GetNewMessageCounts method.
		GetNewMessageCounts []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// MemberID is the memberID argument value.
			MemberID int
		}
		// ListMessages holds details about calls to the ListMessages method.
		ListMessages []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Cursors is the cursors argument value.
			Cursors domain.Cursors
			// Limit is the limit argument value.
			Limit int
			// MemberID is the memberID argument value.
			MemberID int
		}
		// PeekPrevious holds details about calls to the PeekPrevious method.
		PeekPrevious []struct {
			// Timestamp is the timestamp argument value.
			Timestamp *time.Time
		}
		// SaveMessage holds details about calls to the SaveMessage method.
		SaveMessage []struct {
			// Message is the message argument value.
			Message domain.Message
		}
		// SavePost holds details about calls to the SavePost method.
		SavePost []struct {
			// Post is the post argument value.
			Post domain.MessagePost
		}
		// ViewMessage holds details about calls to the ViewMessage method.
		ViewMessage []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// MemberID is the memberID argument value.
			MemberID int
			// MessageID is the messageID argument value.
			MessageID int
		}
	}
	lockDeleteMessage                sync.RWMutex
	lockGetMessageByID               sync.RWMutex
	lockGetMessageParticipants       sync.RWMutex
	lockGetMessagePostByID           sync.RWMutex
	lockGetMessagePostsByID          sync.RWMutex
	lockGetMessagePostsCollapsible   sync.RWMutex
	lockGetMessagesWithCursorForward sync.RWMutex
	lockGetMessagesWithCursorReverse sync.RWMutex
	lockGetNewMessageCounts          sync.RWMutex
	lockListMessages                 sync.RWMutex
	lockPeekPrevious                 sync.RWMutex
	lockSaveMessage                  sync.RWMutex
	lockSavePost                     sync.RWMutex
	lockViewMessage                  sync.RWMutex
}

// DeleteMessage calls DeleteMessageFunc.
func (mock *MessageRepoMock) DeleteMessage(ctx context.Context, memberID int, messageID int) error {
	if mock.DeleteMessageFunc == nil {
		panic("MessageRepoMock.DeleteMessageFunc: method is nil but MessageRepo.DeleteMessage was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		MemberID  int
		MessageID int
	}{
		Ctx:       ctx,
		MemberID:  memberID,
		MessageID: messageID,
	}
	mock.lockDeleteMessage.Lock()
	mock.calls.DeleteMessage = append(mock.calls.DeleteMessage, callInfo)
	mock.lockDeleteMessage.Unlock()
	return mock.DeleteMessageFunc(ctx, memberID, messageID)
}

// DeleteMessageCalls gets all the calls that were made to DeleteMessage.
// Check the length with:
//
//	len(mockedMessageRepo.DeleteMessageCalls())
func (mock *MessageRepoMock) DeleteMessageCalls() []struct {
	Ctx       context.Context
	MemberID  int
	MessageID int
} {
	var calls []struct {
		Ctx       context.Context
		MemberID  int
		MessageID int
	}
	mock.lockDeleteMessage.RLock()
	calls = mock.calls.DeleteMessage
	mock.lockDeleteMessage.RUnlock()
	return calls
}

// GetMessageByID calls GetMessageByIDFunc.
func (mock *MessageRepoMock) GetMessageByID(ctx context.Context, messageID int, memberID int) (*domain.Message, error) {
	if mock.GetMessageByIDFunc == nil {
		panic("MessageRepoMock.GetMessageByIDFunc: method is nil but MessageRepo.GetMessageByID was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		MessageID int
		MemberID  int
	}{
		Ctx:       ctx,
		MessageID: messageID,
		MemberID:  memberID,
	}
	mock.lockGetMessageByID.Lock()
	mock.calls.GetMessageByID = append(mock.calls.GetMessageByID, callInfo)
	mock.lockGetMessageByID.Unlock()
	return mock.GetMessageByIDFunc(ctx, messageID, memberID)
}

// GetMessageByIDCalls gets all the calls that were made to GetMessageByID.
// Check the length with:
//
//	len(mockedMessageRepo.GetMessageByIDCalls())
func (mock *MessageRepoMock) GetMessageByIDCalls() []struct {
	Ctx       context.Context
	MessageID int
	MemberID  int
} {
	var calls []struct {
		Ctx       context.Context
		MessageID int
		MemberID  int
	}
	mock.lockGetMessageByID.RLock()
	calls = mock.calls.GetMessageByID
	mock.lockGetMessageByID.RUnlock()
	return calls
}

// GetMessageParticipants calls GetMessageParticipantsFunc.
func (mock *MessageRepoMock) GetMessageParticipants(ctx context.Context, messageID int) ([]string, error) {
	if mock.GetMessageParticipantsFunc == nil {
		panic("MessageRepoMock.GetMessageParticipantsFunc: method is nil but MessageRepo.GetMessageParticipants was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		MessageID int
	}{
		Ctx:       ctx,
		MessageID: messageID,
	}
	mock.lockGetMessageParticipants.Lock()
	mock.calls.GetMessageParticipants = append(mock.calls.GetMessageParticipants, callInfo)
	mock.lockGetMessageParticipants.Unlock()
	return mock.GetMessageParticipantsFunc(ctx, messageID)
}

// GetMessageParticipantsCalls gets all the calls that were made to GetMessageParticipants.
// Check the length with:
//
//	len(mockedMessageRepo.GetMessageParticipantsCalls())
func (mock *MessageRepoMock) GetMessageParticipantsCalls() []struct {
	Ctx       context.Context
	MessageID int
} {
	var calls []struct {
		Ctx       context.Context
		MessageID int
	}
	mock.lockGetMessageParticipants.RLock()
	calls = mock.calls.GetMessageParticipants
	mock.lockGetMessageParticipants.RUnlock()
	return calls
}

// GetMessagePostByID calls GetMessagePostByIDFunc.
func (mock *MessageRepoMock) GetMessagePostByID(id int) (*domain.MessagePost, error) {
	if mock.GetMessagePostByIDFunc == nil {
		panic("MessageRepoMock.GetMessagePostByIDFunc: method is nil but MessageRepo.GetMessagePostByID was just called")
	}
	callInfo := struct {
		ID int
	}{
		ID: id,
	}
	mock.lockGetMessagePostByID.Lock()
	mock.calls.GetMessagePostByID = append(mock.calls.GetMessagePostByID, callInfo)
	mock.lockGetMessagePostByID.Unlock()
	return mock.GetMessagePostByIDFunc(id)
}

// GetMessagePostByIDCalls gets all the calls that were made to GetMessagePostByID.
// Check the length with:
//
//	len(mockedMessageRepo.GetMessagePostByIDCalls())
func (mock *MessageRepoMock) GetMessagePostByIDCalls() []struct {
	ID int
} {
	var calls []struct {
		ID int
	}
	mock.lockGetMessagePostByID.RLock()
	calls = mock.calls.GetMessagePostByID
	mock.lockGetMessagePostByID.RUnlock()
	return calls
}

// GetMessagePostsByID calls GetMessagePostsByIDFunc.
func (mock *MessageRepoMock) GetMessagePostsByID(memberID int, messageID int, limit int) ([]domain.MessagePost, error) {
	if mock.GetMessagePostsByIDFunc == nil {
		panic("MessageRepoMock.GetMessagePostsByIDFunc: method is nil but MessageRepo.GetMessagePostsByID was just called")
	}
	callInfo := struct {
		MemberID  int
		MessageID int
		Limit     int
	}{
		MemberID:  memberID,
		MessageID: messageID,
		Limit:     limit,
	}
	mock.lockGetMessagePostsByID.Lock()
	mock.calls.GetMessagePostsByID = append(mock.calls.GetMessagePostsByID, callInfo)
	mock.lockGetMessagePostsByID.Unlock()
	return mock.GetMessagePostsByIDFunc(memberID, messageID, limit)
}

// GetMessagePostsByIDCalls gets all the calls that were made to GetMessagePostsByID.
// Check the length with:
//
//	len(mockedMessageRepo.GetMessagePostsByIDCalls())
func (mock *MessageRepoMock) GetMessagePostsByIDCalls() []struct {
	MemberID  int
	MessageID int
	Limit     int
} {
	var calls []struct {
		MemberID  int
		MessageID int
		Limit     int
	}
	mock.lockGetMessagePostsByID.RLock()
	calls = mock.calls.GetMessagePostsByID
	mock.lockGetMessagePostsByID.RUnlock()
	return calls
}

// GetMessagePostsCollapsible calls GetMessagePostsCollapsibleFunc.
func (mock *MessageRepoMock) GetMessagePostsCollapsible(ctx context.Context, viewable int, messageID int, memberID int) ([]domain.MessagePost, int, error) {
	if mock.GetMessagePostsCollapsibleFunc == nil {
		panic("MessageRepoMock.GetMessagePostsCollapsibleFunc: method is nil but MessageRepo.GetMessagePostsCollapsible was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		Viewable  int
		MessageID int
		MemberID  int
	}{
		Ctx:       ctx,
		Viewable:  viewable,
		MessageID: messageID,
		MemberID:  memberID,
	}
	mock.lockGetMessagePostsCollapsible.Lock()
	mock.calls.GetMessagePostsCollapsible = append(mock.calls.GetMessagePostsCollapsible, callInfo)
	mock.lockGetMessagePostsCollapsible.Unlock()
	return mock.GetMessagePostsCollapsibleFunc(ctx, viewable, messageID, memberID)
}

// GetMessagePostsCollapsibleCalls gets all the calls that were made to GetMessagePostsCollapsible.
// Check the length with:
//
//	len(mockedMessageRepo.GetMessagePostsCollapsibleCalls())
func (mock *MessageRepoMock) GetMessagePostsCollapsibleCalls() []struct {
	Ctx       context.Context
	Viewable  int
	MessageID int
	MemberID  int
} {
	var calls []struct {
		Ctx       context.Context
		Viewable  int
		MessageID int
		MemberID  int
	}
	mock.lockGetMessagePostsCollapsible.RLock()
	calls = mock.calls.GetMessagePostsCollapsible
	mock.lockGetMessagePostsCollapsible.RUnlock()
	return calls
}

// GetMessagesWithCursorForward calls GetMessagesWithCursorForwardFunc.
func (mock *MessageRepoMock) GetMessagesWithCursorForward(memberID int, limit int, cursor *time.Time) ([]domain.Message, error) {
	if mock.GetMessagesWithCursorForwardFunc == nil {
		panic("MessageRepoMock.GetMessagesWithCursorForwardFunc: method is nil but MessageRepo.GetMessagesWithCursorForward was just called")
	}
	callInfo := struct {
		MemberID int
		Limit    int
		Cursor   *time.Time
	}{
		MemberID: memberID,
		Limit:    limit,
		Cursor:   cursor,
	}
	mock.lockGetMessagesWithCursorForward.Lock()
	mock.calls.GetMessagesWithCursorForward = append(mock.calls.GetMessagesWithCursorForward, callInfo)
	mock.lockGetMessagesWithCursorForward.Unlock()
	return mock.GetMessagesWithCursorForwardFunc(memberID, limit, cursor)
}

// GetMessagesWithCursorForwardCalls gets all the calls that were made to GetMessagesWithCursorForward.
// Check the length with:
//
//	len(mockedMessageRepo.GetMessagesWithCursorForwardCalls())
func (mock *MessageRepoMock) GetMessagesWithCursorForwardCalls() []struct {
	MemberID int
	Limit    int
	Cursor   *time.Time
} {
	var calls []struct {
		MemberID int
		Limit    int
		Cursor   *time.Time
	}
	mock.lockGetMessagesWithCursorForward.RLock()
	calls = mock.calls.GetMessagesWithCursorForward
	mock.lockGetMessagesWithCursorForward.RUnlock()
	return calls
}

// GetMessagesWithCursorReverse calls GetMessagesWithCursorReverseFunc.
func (mock *MessageRepoMock) GetMessagesWithCursorReverse(memberID int, limit int, cursor *time.Time) ([]domain.Message, error) {
	if mock.GetMessagesWithCursorReverseFunc == nil {
		panic("MessageRepoMock.GetMessagesWithCursorReverseFunc: method is nil but MessageRepo.GetMessagesWithCursorReverse was just called")
	}
	callInfo := struct {
		MemberID int
		Limit    int
		Cursor   *time.Time
	}{
		MemberID: memberID,
		Limit:    limit,
		Cursor:   cursor,
	}
	mock.lockGetMessagesWithCursorReverse.Lock()
	mock.calls.GetMessagesWithCursorReverse = append(mock.calls.GetMessagesWithCursorReverse, callInfo)
	mock.lockGetMessagesWithCursorReverse.Unlock()
	return mock.GetMessagesWithCursorReverseFunc(memberID, limit, cursor)
}

// GetMessagesWithCursorReverseCalls gets all the calls that were made to GetMessagesWithCursorReverse.
// Check the length with:
//
//	len(mockedMessageRepo.GetMessagesWithCursorReverseCalls())
func (mock *MessageRepoMock) GetMessagesWithCursorReverseCalls() []struct {
	MemberID int
	Limit    int
	Cursor   *time.Time
} {
	var calls []struct {
		MemberID int
		Limit    int
		Cursor   *time.Time
	}
	mock.lockGetMessagesWithCursorReverse.RLock()
	calls = mock.calls.GetMessagesWithCursorReverse
	mock.lockGetMessagesWithCursorReverse.RUnlock()
	return calls
}

// GetNewMessageCounts calls GetNewMessageCountsFunc.
func (mock *MessageRepoMock) GetNewMessageCounts(ctx context.Context, memberID int) (*domain.MessageCounts, error) {
	if mock.GetNewMessageCountsFunc == nil {
		panic("MessageRepoMock.GetNewMessageCountsFunc: method is nil but MessageRepo.GetNewMessageCounts was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		MemberID int
	}{
		Ctx:      ctx,
		MemberID: memberID,
	}
	mock.lockGetNewMessageCounts.Lock()
	mock.calls.GetNewMessageCounts = append(mock.calls.GetNewMessageCounts, callInfo)
	mock.lockGetNewMessageCounts.Unlock()
	return mock.GetNewMessageCountsFunc(ctx, memberID)
}

// GetNewMessageCountsCalls gets all the calls that were made to GetNewMessageCounts.
// Check the length with:
//
//	len(mockedMessageRepo.GetNewMessageCountsCalls())
func (mock *MessageRepoMock) GetNewMessageCountsCalls() []struct {
	Ctx      context.Context
	MemberID int
} {
	var calls []struct {
		Ctx      context.Context
		MemberID int
	}
	mock.lockGetNewMessageCounts.RLock()
	calls = mock.calls.GetNewMessageCounts
	mock.lockGetNewMessageCounts.RUnlock()
	return calls
}

// ListMessages calls ListMessagesFunc.
func (mock *MessageRepoMock) ListMessages(ctx context.Context, cursors domain.Cursors, limit int, memberID int) ([]domain.Message, domain.Cursors, error) {
	if mock.ListMessagesFunc == nil {
		panic("MessageRepoMock.ListMessagesFunc: method is nil but MessageRepo.ListMessages was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		Cursors  domain.Cursors
		Limit    int
		MemberID int
	}{
		Ctx:      ctx,
		Cursors:  cursors,
		Limit:    limit,
		MemberID: memberID,
	}
	mock.lockListMessages.Lock()
	mock.calls.ListMessages = append(mock.calls.ListMessages, callInfo)
	mock.lockListMessages.Unlock()
	return mock.ListMessagesFunc(ctx, cursors, limit, memberID)
}

// ListMessagesCalls gets all the calls that were made to ListMessages.
// Check the length with:
//
//	len(mockedMessageRepo.ListMessagesCalls())
func (mock *MessageRepoMock) ListMessagesCalls() []struct {
	Ctx      context.Context
	Cursors  domain.Cursors
	Limit    int
	MemberID int
} {
	var calls []struct {
		Ctx      context.Context
		Cursors  domain.Cursors
		Limit    int
		MemberID int
	}
	mock.lockListMessages.RLock()
	calls = mock.calls.ListMessages
	mock.lockListMessages.RUnlock()
	return calls
}

// PeekPrevious calls PeekPreviousFunc.
func (mock *MessageRepoMock) PeekPrevious(timestamp *time.Time) (bool, error) {
	if mock.PeekPreviousFunc == nil {
		panic("MessageRepoMock.PeekPreviousFunc: method is nil but MessageRepo.PeekPrevious was just called")
	}
	callInfo := struct {
		Timestamp *time.Time
	}{
		Timestamp: timestamp,
	}
	mock.lockPeekPrevious.Lock()
	mock.calls.PeekPrevious = append(mock.calls.PeekPrevious, callInfo)
	mock.lockPeekPrevious.Unlock()
	return mock.PeekPreviousFunc(timestamp)
}

// PeekPreviousCalls gets all the calls that were made to PeekPrevious.
// Check the length with:
//
//	len(mockedMessageRepo.PeekPreviousCalls())
func (mock *MessageRepoMock) PeekPreviousCalls() []struct {
	Timestamp *time.Time
} {
	var calls []struct {
		Timestamp *time.Time
	}
	mock.lockPeekPrevious.RLock()
	calls = mock.calls.PeekPrevious
	mock.lockPeekPrevious.RUnlock()
	return calls
}

// SaveMessage calls SaveMessageFunc.
func (mock *MessageRepoMock) SaveMessage(message domain.Message) (int, error) {
	if mock.SaveMessageFunc == nil {
		panic("MessageRepoMock.SaveMessageFunc: method is nil but MessageRepo.SaveMessage was just called")
	}
	callInfo := struct {
		Message domain.Message
	}{
		Message: message,
	}
	mock.lockSaveMessage.Lock()
	mock.calls.SaveMessage = append(mock.calls.SaveMessage, callInfo)
	mock.lockSaveMessage.Unlock()
	return mock.SaveMessageFunc(message)
}

// SaveMessageCalls gets all the calls that were made to SaveMessage.
// Check the length with:
//
//	len(mockedMessageRepo.SaveMessageCalls())
func (mock *MessageRepoMock) SaveMessageCalls() []struct {
	Message domain.Message
} {
	var calls []struct {
		Message domain.Message
	}
	mock.lockSaveMessage.RLock()
	calls = mock.calls.SaveMessage
	mock.lockSaveMessage.RUnlock()
	return calls
}

// SavePost calls SavePostFunc.
func (mock *MessageRepoMock) SavePost(post domain.MessagePost) (int, error) {
	if mock.SavePostFunc == nil {
		panic("MessageRepoMock.SavePostFunc: method is nil but MessageRepo.SavePost was just called")
	}
	callInfo := struct {
		Post domain.MessagePost
	}{
		Post: post,
	}
	mock.lockSavePost.Lock()
	mock.calls.SavePost = append(mock.calls.SavePost, callInfo)
	mock.lockSavePost.Unlock()
	return mock.SavePostFunc(post)
}

// SavePostCalls gets all the calls that were made to SavePost.
// Check the length with:
//
//	len(mockedMessageRepo.SavePostCalls())
func (mock *MessageRepoMock) SavePostCalls() []struct {
	Post domain.MessagePost
} {
	var calls []struct {
		Post domain.MessagePost
	}
	mock.lockSavePost.RLock()
	calls = mock.calls.SavePost
	mock.lockSavePost.RUnlock()
	return calls
}

// ViewMessage calls ViewMessageFunc.
func (mock *MessageRepoMock) ViewMessage(ctx context.Context, memberID int, messageID int) (int, error) {
	if mock.ViewMessageFunc == nil {
		panic("MessageRepoMock.ViewMessageFunc: method is nil but MessageRepo.ViewMessage was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		MemberID  int
		MessageID int
	}{
		Ctx:       ctx,
		MemberID:  memberID,
		MessageID: messageID,
	}
	mock.lockViewMessage.Lock()
	mock.calls.ViewMessage = append(mock.calls.ViewMessage, callInfo)
	mock.lockViewMessage.Unlock()
	return mock.ViewMessageFunc(ctx, memberID, messageID)
}

// ViewMessageCalls gets all the calls that were made to ViewMessage.
// Check the length with:
//
//	len(mockedMessageRepo.ViewMessageCalls())
func (mock *MessageRepoMock) ViewMessageCalls() []struct {
	Ctx       context.Context
	MemberID  int
	MessageID int
} {
	var calls []struct {
		Ctx       context.Context
		MemberID  int
		MessageID int
	}
	mock.lockViewMessage.RLock()
	calls = mock.calls.ViewMessage
	mock.lockViewMessage.RUnlock()
	return calls
}
