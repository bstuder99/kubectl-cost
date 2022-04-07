.PHONY: darwin-amd64
darwin-amd64:
	cd cmd/kubectl-cost && GOOS=darwin GOARCH=amd64 go build -o kubectl-cost-darwin-amd64
	tar --transform 's|.*kubectl-cost-darwin-amd64|kubectl-cost|' \
		-czf cmd/kubectl-cost/kubectl-cost-darwin-amd64.tar.gz \
		cmd/kubectl-cost/kubectl-cost-darwin-amd64 \
		LICENSE

.PHONY: darwin-arm64
darwin-arm64:
	cd cmd/kubectl-cost && GOOS=darwin GOARCH=arm64 go build -o kubectl-cost-darwin-arm64
	tar --transform 's|.*kubectl-cost-darwin-arm64|kubectl-cost|' \
		-czf cmd/kubectl-cost/kubectl-cost-darwin-arm64.tar.gz \
		cmd/kubectl-cost/kubectl-cost-darwin-arm64 \
		LICENSE

.PHONY: linux-amd64
linux-amd64:
	cd cmd/kubectl-cost && GOOS=linux GOARCH=amd64 go build -o kubectl-cost-linux-amd64
	tar --transform 's|.*kubectl-cost-linux-amd64|kubectl-cost|' \
		-czf cmd/kubectl-cost/kubectl-cost-linux-amd64.tar.gz \
		cmd/kubectl-cost/kubectl-cost-linux-amd64 \
		LICENSE

.PHONY: windows-amd64
windows-amd64:
	cd cmd/kubectl-cost && GOOS=windows GOARCH=amd64 go build -o kubectl-cost-windows-amd64
	tar --transform 's|.*kubectl-cost-windows-amd64|kubectl-cost|' \
		-czf cmd/kubectl-cost/kubectl-cost-windows-amd64.tar.gz \
		cmd/kubectl-cost/kubectl-cost-windows-amd64 \
		LICENSE

.PHONY: release
release: darwin-amd64 darwin-arm64 linux-amd64 windows-amd64
	sha256sum cmd/kubectl-cost/kubectl-cost-darwin-amd64.tar.gz
	sha256sum cmd/kubectl-cost/kubectl-cost-darwin-arm64.tar.gz
	sha256sum cmd/kubectl-cost/kubectl-cost-linux-amd64.tar.gz
	sha256sum cmd/kubectl-cost/kubectl-cost-windows-amd64.tar.gz

.PHONY: build
build:
	cd cmd/kubectl-cost && go build

.PHONY: install
install: build
	chmod +x ./cmd/kubectl-cost/kubectl-cost
	cp ./cmd/kubectl-cost/kubectl-cost ~/go/bin/kubectl-cost
