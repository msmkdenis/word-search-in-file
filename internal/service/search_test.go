package service

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/msmkdenis/word-search-in-file/internal/mocks"
	"github.com/msmkdenis/word-search-in-file/internal/model"
)

type SearchHandlerTestSuite struct {
	suite.Suite
	searcherService *SearcherService
	mockCache       *mocks.MockIndexCache
	ctrl            *gomock.Controller
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(SearchHandlerTestSuite))
}

func (s *SearchHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockCache = mocks.NewMockIndexCache(s.ctrl)
	s.searcherService = NewSearcher(s.mockCache, 5)
}

func (s *SearchHandlerTestSuite) TestSearch() {
	fsTest := fstest.MapFS{
		"file1.txt": {Data: []byte("World")},
		"file2.txt": {Data: []byte("World1")},
		"file3.txt": {Data: []byte("Hello World")},
	}

	fs := model.NewFileSystem("./examples", fsTest)

	testCases := []struct {
		name         string
		word         string
		fs           model.FileSystem
		prepare      func()
		expectedBody []string
		expectedErr  error
	}{
		{
			name: "Success",
			word: "World",
			fs:   fs,
			prepare: func() {
				s.mockCache.EXPECT().SetIndex(gomock.Any(), gomock.Any()).Times(1)
			},
			expectedBody: []string{"file1", "file3"},
			expectedErr:  nil,
		},
	}
	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			if test.prepare != nil {
				test.prepare()
			}

			files, err := s.searcherService.Search(context.Background(), test.word, test.fs)
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedBody, files)
		})
	}
}
