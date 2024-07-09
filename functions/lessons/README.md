# Lessons

The Create Lesson Lambda function is designed to facilitate the creation of lesson records in a DynamoDB table. It takes lesson details as input, marshals these details into the appropriate format, and inserts the lesson as a new item into the DynamoDB table named lessons. This function ensures that each lesson is stored with a unique identifier and associated attributes, such as the lesson title.

## Local Environment

- Run the following to start a local environment in the root directory of this project
```shell
make run
```

- Create a new table in local dynamoDB
```shell
aws dynamodb create-table \
  --table-name lessons \
  --attribute-definitions AttributeName=ID,AttributeType=S \
  --key-schema AttributeName=ID,KeyType=HASH \
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
  --endpoint-url http://localhost:8000
```

- Create a lesson 
<!-- not working yet -->
```shell
curl -X POST http://localhost:8080/2015-03-31/functions/function/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "httpMethod": "POST",
    "headers": {
      "Content-Type": "application/json"
    },
    "body": "{\"Name\": \"Introduction to Go2\"}"
  }'
```

- List the lessons
```
aws dynamodb scan \
  --table-name lessons \
  --endpoint-url http://localhost:8000
```