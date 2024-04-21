.PHONY: all

# Generate mock search service
gen-mock-searcher-service:
	@mockgen --build_flags=--mod=mod \
			 -destination=internal/mocks/mock_searcher_service.go \
			 -package=mocks github.com/msmkdenis/word-search-in-files/internal/handler Searcher