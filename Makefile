.PHONY: build clean deploy

build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/create functions/lessons/create.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

test: build
	cd functions/lessons && go test
	cd internal/database && go test

run: 
	sudo docker-compose up -d
	until curl -s http://localhost:8000 > /dev/null; do echo "Waiting for DynamoDB to start..."; sleep 2; done 

	aws dynamodb create-table \
  		--table-name lessons \
  		--attribute-definitions AttributeName=ID,AttributeType=S \
  		--key-schema AttributeName=ID,KeyType=HASH \
  		--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
		--region us-west-1 \
		--endpoint-url http://localhost:8000

	aws dynamodb list-tables --endpoint-url http://localhost:8000 --region us-west-1

	sudo docker-compose logs -f

stop: 
	sudo docker-compose down