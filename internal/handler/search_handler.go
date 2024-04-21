package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Searcher interface {
	Search(ctx context.Context, word string, dirPath string) (files []string, err error)
}

type SearchHandler struct {
	e        *echo.Echo
	searcher Searcher
}

func NewSearchHandler(e *echo.Echo, searcher Searcher) *SearchHandler {
	handler := &SearchHandler{
		e:        e,
		searcher: searcher,
	}

	e.GET("/files/search", handler.SearchWords)

	return handler
}

func (s *SearchHandler) SearchWords(c echo.Context) error {
	word := c.QueryParam("word")
	dir := c.QueryParam("dir")
	if word == "" || dir == "" {
		slog.Info("Bad request: word or dir is empty")
		return c.JSON(http.StatusBadRequest, nil)
	}

	files, err := s.searcher.Search(c.Request().Context(), word, dir)
	if err != nil {
		slog.Error("Error while searching in files", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, files)
}
