.phony: deploy bootstrap build

deploy:
	npx cdk deploy

bootstrap:
	npx cdk bootstrap

build:
	cd src && GOOS=linux CGO_ENABLED=0 go build main.go

