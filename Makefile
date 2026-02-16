.PHONY: all help deps test build-desktop build-android-lib new-skill clean

PROJECT_NAME := upcraft-agent
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod

all: test build-desktop

## help: Show available targets
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'

## deps: Download and tidy dependencies
deps:
	$(GOMOD) tidy

## test: Run core and backend tests (legacy code excluded by default tags)
test:
	$(GOTEST) -v ./core/... ./backend/...

## build-desktop: Build desktop CLI harness
build-desktop:
	@mkdir -p bin
	$(GOBUILD) -o bin/upcraft-cli ./core/cmd/desktop
	@echo "Build complete: bin/upcraft-cli"

## build-android-lib: Build Android AAR from Go mobile bridge
build-android-lib:
	@mkdir -p app/shared/libs
	gomobile bind -target=android -o app/shared/libs/upcraft_core.aar ./core/mobile
	@echo "Android library compiled: app/shared/libs/upcraft_core.aar"

## new-skill: Scaffold a new skill (usage: make new-skill name=weather)
new-skill:
	@test "$(name)" || (echo "Error: name argument required (e.g. make new-skill name=weather)" && exit 1)
	$(GOCMD) run ./scripts/scaffold_skill.go $(name)

## clean: Remove local build outputs
clean:
	@rm -rf bin
