package threadsvc

import (
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"html/template"
	"time"

	"regexp"

	"go.uber.org/zap"
)

type ThreadService struct {
	threadRepo            ports.ThreadRepo
	memberRepo            ports.MemberRepo
	logger                *zap.SugaredLogger
	defaultMaxThreadLimit int
}

func NewThreadService(postRepo ports.ThreadRepo, memberRepo ports.MemberRepo, logger *zap.SugaredLogger, defaultMaxThreadLimit int) ThreadService {
	return ThreadService{
		threadRepo:            postRepo,
		logger:                logger,
		memberRepo:            memberRepo,
		defaultMaxThreadLimit: defaultMaxThreadLimit,
	}
}

func (s ThreadService) NewPost(body, ip, memberName string, threadID int) (int, error) {
	memberID, err := s.memberRepo.GetMemberIDByUsername(memberName)
	if err != nil {
		s.logger.Errorf("error getting member id by username: %v", err)
		return 0, err
	}

	id, err := s.threadRepo.SavePost(domain.ThreadPost{
		Body:     body,
		MemberIP: ip,
		ParentID: threadID,
		MemberID: memberID,
	})
	if err != nil {
		s.logger.Errorf("error saving post: %v", err)
		return 0, err
	}

	return id, nil
}

func (s ThreadService) GetPostByID(id int) (*domain.ThreadPost, error) {
	return s.threadRepo.GetPostByID(id)
}

func (s ThreadService) GetThreadByID(limit, offset, id int) (*domain.Thread, error) {
	posts, err := s.threadRepo.ListPostsForThread(limit, offset, id)
	if err != nil {
		s.logger.Errorf("error getting posts by thread id: %v", err)
		return nil, err
	}

	thread, err := s.threadRepo.GetThreadByID(id)
	if err != nil {
		s.logger.Errorf("error getting thread by id: %v", err)
		return nil, err
	}

	for idx, post := range posts {
		postPtr := &post
		postPtr.Position = idx + 1
		thread.Posts = append(thread.Posts, *postPtr)
	}

	return thread, nil
}

func (s ThreadService) GetThreadsWithCursorForward(limit int, firstPage bool, cursor *time.Time) (*domain.SiteContext, error) {
	if firstPage {
		start := time.Date(2999, 1, 1, 0, 0, 0, 0, time.UTC)
		threads, err := s.threadRepo.ListThreadsByCursorForward(limit, &start)
		if err != nil {
			s.logger.Errorf("error getting first page of threads by cursor: %v", err)
			return nil, err
		}

		site := &domain.SiteContext{
			ThreadPage: domain.ThreadPage{
				Threads: threads,
			},
		}

		if len(threads) > s.defaultMaxThreadLimit {
			site.ThreadPage.HasNextPage = true
			threads = threads[:s.defaultMaxThreadLimit]
		} else {
			site.ThreadPage.HasNextPage = false
		}

		site.ThreadPage.Threads = threads

		if len(site.ThreadPage.Threads) != 0 {
			site.PageCursor = site.ThreadPage.Threads[len(site.ThreadPage.Threads)-1].DateLastPosted
			site.PrevPageCursor = nil
			prevExists, err := s.threadRepo.PeekPrevious(threads[0].DateLastPosted)
			if err != nil {
				return nil, err
			}

			site.ThreadPage.HasPrevPage = prevExists
		}

		return site, nil
	}

	threads, err := s.threadRepo.ListThreadsByCursorForward(limit, cursor)
	if err != nil {
		s.logger.Errorf("error getting page of threads by cursor: %v", err)
		return nil, err
	}

	site := &domain.SiteContext{}

	if len(threads) > s.defaultMaxThreadLimit {
		site.ThreadPage.HasNextPage = true
		threads = threads[:s.defaultMaxThreadLimit]
	} else {
		site.ThreadPage.HasNextPage = false
	}

	site.ThreadPage.Threads = threads

	if len(site.ThreadPage.Threads) != 0 {
		site.PageCursor = site.ThreadPage.Threads[len(site.ThreadPage.Threads)-1].DateLastPosted
		site.PrevPageCursor = site.ThreadPage.Threads[0].DateLastPosted

		prevExists, err := s.threadRepo.PeekPrevious(threads[0].DateLastPosted)
		if err != nil {
			return nil, err
		}

		site.ThreadPage.HasPrevPage = prevExists
	}

	return site, nil
}

func (s ThreadService) GetThreadsWithCursorReverse(limit int, cursor *time.Time) (*domain.SiteContext, error) {
	threads, err := s.threadRepo.ListThreadsByCursorReverse(limit, cursor)
	if err != nil {
		s.logger.Errorf("error getting page of threads by cursor: %v", err)
		return nil, err
	}

	site := &domain.SiteContext{}

	site.ThreadPage.Threads = threads[:len(threads)-1]
	site.PrevPageCursor = threads[0].DateLastPosted

	prevExists, err := s.threadRepo.PeekPrevious(threads[0].DateLastPosted)
	if err != nil {
		return nil, err
	}

	site.ThreadPage.HasPrevPage = prevExists

	if len(threads) > s.defaultMaxThreadLimit {
		site.ThreadPage.HasNextPage = true
		site.ThreadPage.Threads = site.ThreadPage.Threads[:s.defaultMaxThreadLimit]
	} else {
		site.ThreadPage.HasNextPage = false
	}

	if len(site.ThreadPage.Threads) != 0 {
		site.PageCursor = site.ThreadPage.Threads[len(site.ThreadPage.Threads)-1].DateLastPosted
	}

	return site, nil
}

func (s ThreadService) ListThreads(limit, offset int) (*domain.SiteContext, error) {
	return s.threadRepo.ListThreads(limit, offset)
}

func (s ThreadService) NewThread(memberName, memberIP, body, subject string) (int, error) {
	id, err := s.memberRepo.GetMemberIDByUsername(memberName)
	if err != nil {
		s.logger.Errorf("error getting member id by username: %v", err)
		return 0, err
	}

	thread := domain.Thread{
		Subject:       subject,
		FirstPostText: body,
		MemberID:      id,
		LastPosterID:  id,
		MemberIP:      memberIP,
	}

	threadID, err := s.threadRepo.SaveThread(thread)
	if err != nil {
		s.logger.Errorf("error saving thread: %v", err)
		return 0, err
	}

	return threadID, nil
}

func (s ThreadService) ConvertPostBodyBbcodeToHtml(postBody string) (*template.HTML, error) {
	hrefRegexp := regexp.MustCompile(`((https?|ftp)://[^\s/$.?#].[^\s]*)`)
	convertedPostBody := hrefRegexp.ReplaceAllString(postBody, `<a href="$1" class="link" onclick="window.open(this.href); return false;">$1</a>`)
	htmlPostBody := template.HTML(convertedPostBody)
	return &htmlPostBody, nil
}
