BINDIR ?= bin
APP_NAME := maze
BUILD_DIR := build
DIST_DIR := dist
DOCS_DIR := docs

.PHONY: all install uninstall clean dvi dist tests

all: $(BUILD_DIR)/$(APP_NAME)

$(BUILD_DIR)/$(APP_NAME):
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/maze_gui

install: all
	install -d $(BINDIR)
	install -m 0755 $(BUILD_DIR)/$(APP_NAME) $(BINDIR)/$(APP_NAME)

uninstall:
	rm -f $(BINDIR)/$(APP_NAME)

clean:
	rm -rf $(BUILD_DIR) $(DIST_DIR)
	go clean -testcache

rebuild: clean all

dvi:
	./dvi.sh $(DOCS_DIR) ./internal/...

dist: clean
	mkdir -p $(DIST_DIR)
	tar -czf $(DIST_DIR)/$(APP_NAME).tar.gz --exclude="$(DIST_DIR)/$(APP_NAME).tar.gz" *

tests:
	go test ./... -cover
