.PHONY: all

# Generate mock search service
gen-mock-searcher-service:
	@mockgen --build_flags=--mod=mod \
			 -destination=internal/mocks/mock_searcher_service.go \
			 -package=mocks github.com/msmkdenis/word-search-in-file/internal/handler Searcher

gen-mock-index-cache:
	@mockgen --build_flags=--mod=mod \
			 -destination=internal/mocks/mock_index_cache.go \
			 -package=mocks github.com/msmkdenis/word-search-in-file/internal/service IndexCache