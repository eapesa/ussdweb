###
GO=go
PWD=$(CURDIR)
MKDIR_P=mkdir -p
BUILD_DIR=$(PWD)/builds
APPNAME=ussdweb

.PHONY: install clean run

echo:
	@echo $(APPNAME)

install:
	@$(MKDIR_P) $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/$(APPNAME)
	@echo "Binary installed at $(BUILD_DIR)/$(APPNAME)"

clean:
	rm $(BUILD_DIR)/*

run:
	@$(BUILD_DIR)/$(APPNAME)
