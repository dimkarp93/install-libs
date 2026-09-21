MODULE  := github.com/dimkarp93/install-libs
NAME    := install-libs
VERSION := $(shell tr -d '[:space:]' < versions.txt)
TAG     := v$(VERSION)
DIST    := dist
PROXY   := $(DIST)/proxy/$(MODULE)/@v
STAGE   := $(DIST)/stage/$(MODULE)@$(TAG)

.PHONY: build test test-v test-run cover vet fmt check pack clean bump-patch bump-minor bump-major

build:
	go build ./...

test:
	go test ./...

test-v:
	go test -v ./...

test-run:
	@test -n "$(T)" || { echo "specify a test: make test-run T=TestNormalizeURL"; exit 1; }
	go test -v -run '$(T)' ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1
	@echo "detailed report: go tool cover -html=coverage.out"

vet:
	go vet ./...

fmt:
	gofmt -l -w .

check: vet test

pack: build test
	rm -rf $(DIST)
	mkdir -p $(PROXY) $(STAGE)
	tar -c --exclude=./$(DIST) --exclude=./.git --exclude=./go.work --exclude=./go.work.sum . | tar -x -C $(STAGE)
	cd $(DIST)/stage && zip -qr ../proxy/$(MODULE)/@v/$(TAG).zip $(MODULE)@$(TAG)
	cp go.mod $(PROXY)/$(TAG).mod
	printf '{"Version":"%s","Time":"%s"}\n' "$(TAG)" "$$(date -u +%Y-%m-%dT%H:%M:%SZ)" > $(PROXY)/$(TAG).info
	printf '%s\n' "$(TAG)" > $(PROXY)/list
	rm -rf $(DIST)/stage
	tar -czf $(DIST)/$(NAME)-$(VERSION)-proxy.tar.gz -C $(DIST) proxy
	@echo
	@echo "Package: $(DIST)/$(NAME)-$(VERSION)-proxy.tar.gz  ($(MODULE) $(TAG))"
	@echo
	@echo "How to use it before the module is published:"
	@echo "  tar -xzf $(NAME)-$(VERSION)-proxy.tar.gz -C /opt"
	@echo "  GOFLAGS=-mod=mod GOPROXY=file:///opt/proxy,https://proxy.golang.org,direct \\"
	@echo "  GONOSUMDB='github.com/dimkarp93/*' GONOSUMCHECK=1 GOSUMDB=off go build ./..."

clean:
	rm -rf $(DIST) coverage.out

bump-patch:
	@v=$$(tr -d '[:space:]' < versions.txt); \
	MAJ=$${v%%.*}; rest=$${v#*.}; MIN=$${rest%%.*}; PAT=$${rest##*.}; \
	printf '%s.%s.%s\n' "$$MAJ" "$$MIN" "$$((PAT + 1))" > versions.txt; \
	cat versions.txt

bump-minor:
	@v=$$(tr -d '[:space:]' < versions.txt); \
	MAJ=$${v%%.*}; rest=$${v#*.}; MIN=$${rest%%.*}; \
	printf '%s.%s.0\n' "$$MAJ" "$$((MIN + 1))" > versions.txt; \
	cat versions.txt

bump-major:
	@v=$$(tr -d '[:space:]' < versions.txt); \
	MAJ=$${v%%.*}; \
	printf '%s.0.0\n' "$$((MAJ + 1))" > versions.txt; \
	cat versions.txt
