STEAMPIPE_INSTALL_DIR ?= ~/.steampipe
BUILD_TAGS = netgo
install:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/hub.steampipe.io/plugins/turbot/kubernetes@latest/steampipe-plugin-kubernetes.plugin -tags "${BUILD_TAGS}" *.go
install-local:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/local/kubernetes/kubernetes.plugin -tags "${BUILD_TAGS}" *.go
