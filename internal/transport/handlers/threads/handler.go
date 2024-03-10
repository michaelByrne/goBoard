package threads

import (
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"goBoard/internal/transport/handlers/threads/views"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	threadService ports.ThreadService

	logger *zap.SugaredLogger
}

func NewHandler(threadService ports.ThreadService, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		threadService: threadService,
		logger:        logger,
	}
}

func (h *Handler) Register(r chi.Router) {
	r.Get("/", h.ThreadsHome)
	r.Get("/threads", h.Threads)
}

func (h *Handler) ThreadsHome(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	threads, cursorsOut, err := h.threadService.ListThreads(ctx, domain.Cursors{}, 50)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templ.Handler(views.Home(views.Threads(threads, cursorsOut))).Component.Render(ctx, w)
}

func (h *Handler) Threads(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	nextCursor := r.URL.Query().Get("next")
	prevCursor := r.URL.Query().Get("prev")
	dir := r.URL.Query().Get("dir")
	if dir != "next" && dir != "prev" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if dir == "next" {
		prevCursor = ""
	} else {
		nextCursor = ""
	}

	cursors := domain.Cursors{
		Next: nextCursor,
		Prev: prevCursor,
	}

	threads, cursors, err := h.threadService.ListThreads(ctx, cursors, 50)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isHx(r) {
		templ.Handler(views.Threads(threads, cursors)).Component.Render(ctx, w)
		return
	}

	templ.Handler((views.Home(views.Threads(threads, cursors)))).Component.Render(ctx, w)
}

func isHx(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
