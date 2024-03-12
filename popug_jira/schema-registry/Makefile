GO_PKG_DST := go/
GO_PKG_PB_DST := pkg/pb/
GO_PKG_PB_PATH_DST = $(GO_PKG_DST)/$(GO_PKG_PB_DST)

.PHONY: all

all: buf-lint buf-generate-go

buf-lint:
	buf lint

buf-generate-go:
	rm -rf $(GO_PKG_PB_PATH_DST)
	buf generate --template=./buf.gen.go.yaml -o $(GO_PKG_PB_PATH_DST)
	if [ ! -f $(GO_PKG_DST)/go.mod ]; then \
		echo "File does not exist, create init go package"; \
		(cd $(GO_PKG_DST) && go mod init github.com/lyckety/async_arch/popug_jira/schema-registry); \
	fi
	(cd $(GO_PKG_DST) && go mod tidy)
