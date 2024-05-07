# Lumbrera

## Functions

```
sudo docker run --rm -ti \
  -v $(pwd)/bin:/var/task \
  -v ${HOME}/.aws/:/root/.aws \
  --name lambda-env \
  -p 8080:8080 \
  public.ecr.aws/lambda/go:latest \
  create
```

```
curl -XPOST \
  "http://localhost:8080/2015-03-31/functions/function/invocations"  \
  -d '{"payload":"hello world!"}'
```