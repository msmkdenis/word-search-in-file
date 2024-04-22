package handler

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	"github.com/msmkdenis/word-search-in-file/internal/middleware"
	"github.com/msmkdenis/word-search-in-file/internal/model"
)

type Searcher interface {
	Search(ctx context.Context, word string, fs model.FileSystem) (files []string, err error)
}

type SearchHandler struct {
	e               *echo.Echo
	searcher        Searcher
	cacheMiddleware *middleware.CacheSearchMiddleware
}

func NewSearchHandler(e *echo.Echo, searcher Searcher, cache *middleware.CacheSearchMiddleware) *SearchHandler {
	handler := &SearchHandler{
		e:               e,
		searcher:        searcher,
		cacheMiddleware: cache,
	}

	e.GET("/files/search", handler.SearchWords, handler.cacheMiddleware.GetFromCache())

	return handler
}

func (s *SearchHandler) SearchWords(c echo.Context) error {
	word := c.QueryParam("word")
	dir := c.QueryParam("dir")
	if word == "" || dir == "" {
		slog.Info("Bad request: word or dir is empty")
		return c.JSON(http.StatusBadRequest, nil)
	}

	files, err := s.searcher.Search(c.Request().Context(), word, model.NewFileSystem(dir, os.DirFS(dir)))
	if err != nil {
		slog.Error("Error while searching in files", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, files)
}
