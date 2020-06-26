BUILD_DIR=build
PACKAGE_DIR=package

gobuild:
	@mkdir -p $(BUILD_DIR)
	env GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/trigger

targz:
	@mkdir -p $(PACKAGE_DIR)
	tar zcvf $(PACKAGE_DIR)/trigger_linux_amd64.tar.gz $(BUILD_DIR)/trigger
	
