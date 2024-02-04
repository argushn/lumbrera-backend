.PHONY: build

build:
	cd lessons; GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go

clean: 
	rm -R tf-resources/.aws-sam \
		tf-resources/.aws-sam-iacs \
		tf-resources/.terraform

build-terraform:
	cd tf-resources; sam build --hook-name terraform --terraform-project-root-path ../

run:
	cd tf-resources; sam local invoke --hook-name terraform