package threadsvc

import (
	"context"
	"fmt"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"html/template"
	"math/rand"
	"strconv"
	"strings"
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

func (s ThreadService) GetCollapsibleThreadByID(ctx context.Context, viewable, threadID, memberID int) (*domain.Thread, error) {
	posts, count, err := s.threadRepo.ListPostsCollapsible(ctx, viewable, threadID, memberID)
	if err != nil {
		s.logger.Errorf("error getting collapsible posts by thread id: %v", err)
		return nil, err
	}

	thread, err := s.threadRepo.GetThreadByID(threadID, memberID)
	if err != nil {
		s.logger.Errorf("error getting thread by id: %v", err)
		return nil, err
	}

	thread.Posts = posts
	thread.NumCollapsed = count

	return thread, nil
}

func (s ThreadService) GetThreadByID(limit, offset, id, memberID int) (*domain.Thread, error) {
	posts, err := s.threadRepo.ListPostsForThread(limit, offset, id, memberID)
	if err != nil {
		s.logger.Errorf("error getting posts by thread id: %v", err)
		return nil, err
	}

	thread, err := s.threadRepo.GetThreadByID(id, memberID)
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

func (s ThreadService) ListThreads(ctx context.Context, cursors domain.Cursors, limit, memberID int) ([]domain.Thread, domain.Cursors, error) {
	threads, cursors, err := s.threadRepo.ListThreads(ctx, cursors, limit, memberID)
	if err != nil {
		s.logger.Errorf("error getting threads: %v", err)
		return nil, domain.Cursors{}, err
	}

	return threads, cursors, nil
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
	// set up tag lists & utility vars
	formattingTags := []string{"b", "i", "em", "u", "strong", "strike", "sub", "sup", "code"}
	supportedMediaAndFilterTags := []string{"img", "youtube", "vimeo", "soundcloud", "quote", "spoiler", "trigger"}
	mediaTagRegexes := map[string]*regexp.Regexp{}
	convertedPostBody := postBody

	// convert formatting tag
	for _, tag := range formattingTags {
		convertedPostBody = strings.Replace(convertedPostBody, "["+tag+"]", "<"+tag+">", -1)
		convertedPostBody = strings.Replace(convertedPostBody, "[/"+tag+"]", "</"+tag+">", -1)
	}

	// convert untagged HTML links
	untaggedHrefRegexp := regexp.MustCompile(`([^=\]](https?|ftp)://[^\s/$.?#].[^\s]*)`)
	convertedPostBody = untaggedHrefRegexp.ReplaceAllString(convertedPostBody, `<a href="$1" class="link" onclick="window.open(this.href); return false;">$1</a>`)

	// convert text link tags
	textLinkRegexp := regexp.MustCompile(`(\[url=(.[^\]]*)\](.[^\[]*)\[\/url\])`)
	convertedPostBody = textLinkRegexp.ReplaceAllString(convertedPostBody, `<a href="$2" class="link" onclick="window.open(this.href); return false;">$3</a>`)

	// generate media tag regex
	for _, tag := range supportedMediaAndFilterTags {
		mediaTagRegexes[tag] = regexp.MustCompile(`(\[` + tag + `\](.[^\[]*)\[\/` + tag + `\])`)
	}

	// convert img tags
	convertedPostBody = mediaTagRegexes["img"].ReplaceAllString(convertedPostBody, `<img src="$2" ondblclick="window.open(this.src);">`)

	// convert soundcloud tags
	soundcloudElmtHtml := `<object height="81" width="100%"><param name="wmode" value="opaque"><param name="movie" value="$2"><param name="allowscriptaccess" value="always"><embed allowscriptaccess="always" height="81" src="$2" type="video/mp4" width="100%"></object>`
	convertedPostBody = mediaTagRegexes["soundcloud"].ReplaceAllString(convertedPostBody, soundcloudElmtHtml)

	// convert youtube tags
	youtubeElmtHtml := `<object width="425" height="355"><param name="movie" value="$2"><param name="wmode" value="transparent"><embed src="$2" type="video/mp4" wmode="transparent" width="425" height="355"></object>`
	convertedPostBody = mediaTagRegexes["youtube"].ReplaceAllString(convertedPostBody, youtubeElmtHtml)

	// convert vimeo tags (this is just a duplicate of the youtube code right now, probably the same)
	vimeoElmtHtml := `<object width="425" height="355"><param name="movie" value="$2"><param name="wmode" value="transparent"><embed src="$2" type="video/mp4" wmode="transparent" width="425" height="355"></object>`
	convertedPostBody = mediaTagRegexes["vimeo"].ReplaceAllString(convertedPostBody, vimeoElmtHtml)

	// convert tweet tags
	seed := rand.NewSource(time.Now().UnixNano())
	rando := rand.New(seed)
	spanId := strconv.Itoa(rando.Intn(99999999-1000000+1) + 1000000)
	tweetTagRegexp := regexp.MustCompile(`\[tweet\].*\/status\/(\d+)\[\/tweet\]`)
	matches := tweetTagRegexp.FindStringSubmatch(convertedPostBody)
	if len(matches) > 0 {
		tweetId := matches[1]
		tweetScript := fmt.Sprintf("<script>twttr.widgets.createTweet(\"%s\",document.getElementById(\"tt-%s\"),{ dnt: true, theme: \"dark\" });</script>", tweetId, spanId)
		tweetElmtHtml := fmt.Sprintf("<span id=\"tt-%s\"></span>%s", spanId, tweetScript)
		convertedPostBody = tweetTagRegexp.ReplaceAllString(convertedPostBody, tweetElmtHtml)
	}

	// convert quote tags: this doesn't appear to make any actual formatting changes...
	// does it need some css to actually work?
	convertedPostBody = mediaTagRegexes["quote"].ReplaceAllString(convertedPostBody, `<blockquote>$2</blockquote>`)

	// convert spoiler & trigger tags
	convertedPostBody = mediaTagRegexes["spoiler"].ReplaceAllString(convertedPostBody, `<span class="spoiler" onclick="$(this).next().show();$(this).remove()">show spoiler</span><span style="display:none">$2</span>`)
	convertedPostBody = mediaTagRegexes["trigger"].ReplaceAllString(convertedPostBody, `<span class="trigger" onclick="$(this).next().show();$(this).remove()">show trigger</span><span style="display:none">$2</span>`)

	// recognize the prepared post string as HTML
	htmlPostBody := template.HTML(convertedPostBody)
	return &htmlPostBody, nil
}

func (s ThreadService) UndotThread(ctx context.Context, memberID, threadID int) error {
	return s.threadRepo.UndotThread(ctx, memberID, threadID)
}

func (s ThreadService) DotThread(ctx context.Context, memberID, threadID int) error {
	return s.threadRepo.DotThread(ctx, memberID, threadID)
}

func (s ThreadService) ToggleIgnore(ctx context.Context, memberID, threadID int, ignore bool) error {
	return s.threadRepo.ToggleIgnore(ctx, memberID, threadID, ignore)
}

func (s ThreadService) ToggleDot(ctx context.Context, memberID, threadID int) (bool, error) {
	return s.threadRepo.ToggleDot(ctx, memberID, threadID)
}
