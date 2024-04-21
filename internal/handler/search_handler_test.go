package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/msmkdenis/word-search-in-file/internal/mocks"
)

type SearchHandlerTestSuite struct {
	suite.Suite
	h               *SearchHandler
	searcherService *mocks.MockSearcher
	echo            *echo.Echo
	ctrl            *gomock.Controller
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(SearchHandlerTestSuite))
}

func (s *SearchHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.echo = echo.New()
	s.searcherService = mocks.NewMockSearcher(s.ctrl)
	s.h = NewSearchHandler(s.echo, s.searcherService)
}

func (s *SearchHandlerTestSuite) TestSearchWords_Handler() {
	testCases := []struct {
		name         string
		path         string
		method       string
		prepare      func()
		word         string
		dir          string
		expectedBody string
		expectedCode int
	}{
		{
			name:         "Bad request with empty word",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
			path:         "http://localhost:8080/files/search",
			word:         "",
			dir:          "./example",
			expectedBody: "null\n",
			prepare: func() {
				s.searcherService.EXPECT().Search(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:         "Bad request with empty dir",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
			path:         "http://localhost:8080/files/search",
			word:         "world",
			dir:          "",
			expectedBody: "null\n",
			prepare: func() {
				s.searcherService.EXPECT().Search(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:         "Bad request with empty word and dir",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
			path:         "http://localhost:8080/files/search",
			word:         "",
			dir:          "",
			expectedBody: "null\n",
			prepare: func() {
				s.searcherService.EXPECT().Search(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:         "Empty slice returned",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			path:         "http://localhost:8080/files/search",
			word:         "world",
			dir:          "./example",
			expectedBody: "null\n",
			prepare: func() {
				s.searcherService.EXPECT().Search(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil, nil)
			},
		},
		{
			name:         "Error returned with empty slice",
			method:       http.MethodGet,
			expectedCode: http.StatusInternalServerError,
			path:         "http://localhost:8080/files/search",
			word:         "world",
			dir:          "./example",
			expectedBody: "null\n",
			prepare: func() {
				s.searcherService.EXPECT().Search(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.New("error"))
			},
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			if test.prepare != nil {
				test.prepare()
			}

			request := httptest.NewRequest(test.method, test.path, strings.NewReader(""))
			q := request.URL.Query()
			q.Add("word", test.word)
			q.Add("dir", test.dir)
			request.URL.RawQuery = q.Encode()

			w := httptest.NewRecorder()
			l := s.echo.NewContext(request, w)

			err := s.h.SearchWords(l)
			assert.NoError(t, err)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}
