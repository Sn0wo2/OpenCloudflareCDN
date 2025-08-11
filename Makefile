.PHONY: release

release:
	@echo Running release tool...
	@go run ./tool/release/main.go
