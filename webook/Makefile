.PHONY: docker
docker:
	@rm webook || true
	@go mod tidy
	@go build -o webook .
	@docker rmi -f zzlemi/webook:v0.0.1
	@docker build -t zzlemi/webook:v0.0.1 .