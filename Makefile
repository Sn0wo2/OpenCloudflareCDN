.PHONY: release
.PHONY: mvstatic

release:
	@echo Running release tool...
	@go run ./tool/release/main.go


# 'C:\Users\user\AppData\Local\go-build\Frontend\dist‘
mvstatic:
	@echo Running mvstatic tool...
	@go run -ldflags "-X main.WorkDir=$(CURDIR)" ./tool/mvstatic/main.go