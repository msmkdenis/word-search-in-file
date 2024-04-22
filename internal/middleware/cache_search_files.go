package middleware

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CacheSearcher interface {
	GetFiles(path, word string) ([]string, bool)
}

type CacheSearchMiddleware struct {
	cache CacheSearcher
}

func NewCacheSearchMiddleware(cache CacheSearcher) *CacheSearchMiddleware {
	return &CacheSearchMiddleware{cache: cache}
}

func (m *CacheSearchMiddleware) GetFromCache() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			word := c.QueryParam("word")
			dir := c.QueryParam("dir")
			if word == "" || dir == "" {
				slog.Info("Bad request: word or dir is empty")
				return c.JSON(http.StatusBadRequest, nil)
			}

			// Проверяем есть ли в кэше индекс по полученной директории
			// Если есть - то ищем по нему
			if files, ok := m.cache.GetFiles(dir, word); ok {
				// Проставим заголовок, чтобы понять, что данные из кэша
				c.Response().Header().Set("X-Cache", "Cached")
				if len(files) == 0 {
					return c.JSON(http.StatusOK, nil)
				}
				return c.JSON(http.StatusOK, files)
			}

			// В кэше нет индекса по этой директории
			c.Response().Header().Set("X-Cache", "None")
			return next(c)
		}
	}
}
