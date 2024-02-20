.PHONY: build

build:
	cd lessons;
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go

clean: 
	rm -R tf-resources/.aws-sam \
		tf-resources/.aws-sam-iacs \
		tf-resources/.terraform

build-terraform: clean
	cd tf-resources; sam build --hook-name terraform --terraform-project-root-path ../

invoke:
	cd tf-resources; sam local invoke --hook-name terraform

start:
	cd tf-resources; sam local start-api