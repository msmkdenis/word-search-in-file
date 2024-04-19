package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/msmkdenis/word-search-in-file/pkg/searcher"
)

type SearchHandler struct {
	e        *echo.Echo
	searcher *searcher.Searcher
}

func NewSearchHandler(e *echo.Echo, searcher *searcher.Searcher) *SearchHandler {
	handler := &SearchHandler{
		e:        e,
		searcher: searcher,
	}

	e.GET("/files/search", handler.SearchWords)

	return handler
}

func (s *SearchHandler) SearchWords(c echo.Context) error {
	answer, _ := s.searcher.Search("World")
	return c.JSON(http.StatusOK, answer)
}
