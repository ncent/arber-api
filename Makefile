HANDLER=handler
PACKAGE=package
ifdef DOTENV
	DOTENV_TARGET=dotenv
else
	DOTENV_TARGET=./.env
endif

.PHONY: build clean

build: clean

	env GOOS=linux go build -ldflags="-s -w -d" -a -tags netgo -installsuffix netgo -o bin/kinesis/archiver handlers/aws/kinesis/archive/main.go
	env GOOS=linux go build -ldflags="-s -w -d" -a -tags netgo -installsuffix netgo -o bin/kinesis/publisher handlers/aws/kinesis/publish/main.go
	env GOOS=linux go build -ldflags="-s -w -d" -a -tags netgo -installsuffix netgo -o bin/kinesis/consumer handlers/aws/kinesis/consume/main.go
	env GOOS=linux go build -ldflags="-s -w" -a -tags netgo -installsuffix netgo -o bin/user/google/new handlers/google/oauth/new/main.go
	env GOOS=linux go build -ldflags="-s -w" -a -tags netgo -installsuffix netgo -o bin/emailer/send handlers/aws/ses/send/main.go
	env GOOS=linux go build -ldflags="-s -w" -a -tags netgo -installsuffix netgo -o bin/google/gmail/send handlers/google/gmail/send/main.go
	env GOOS=linux go build -ldflags="-s -w" -a -tags netgo -installsuffix netgo -o bin/emailer/receive handlers/aws/ses/receive/main.go
	env GOOS=linux go build -ldflags="-s -w" -a -tags netgo -installsuffix netgo -o bin/user/google/contacts/new handlers/google/oauth/contacts/new/main.go
	env GOOS=linux go build -ldflags="-s -w" -a -tags netgo -installsuffix netgo -o bin/mail/reshare handlers/mail/reshare/main.go
	chmod +x bin/kinesis/archiver
	chmod +x bin/kinesis/publisher
	chmod +x bin/kinesis/consumer
	chmod +x bin/user/google/new
	chmod +x bin/emailer/send
	chmod +x bin/google/gmail/send
	chmod +x bin/emailer/receive
	chmod +x bin/user/google/contacts/new
	chmod +x bin/mail/reshare
	zip -j bin/user/google/contacts/new.zip bin/user/google/contacts/new
	zip -j bin/user/google/new.zip bin/user/google/new
	zip -j bin/emailer/send.zip bin/emailer/send
	zip -j bin/google/gmail/send.zip bin/google/gmail/send
	zip -j bin/emailer/receive.zip bin/emailer/receive
	zip -j bin/mail/reshare.zip bin/mail/reshare


clean:
	-rm -rf ./bin

test: build
	go test -race $$(go list ./... | grep -v /vendor/) -v -coverprofile=coverage.out
	go tool cover -func=coverage.out