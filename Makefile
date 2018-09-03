.PHONY: install test build serve clean pack deploy ship

TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)

export TAG

install:
	go get .

test: install
	go test ./...

build: install
	go build -ldflags "-X main.version=$(TAG)" -o news .

serve: build
	./news

clean:
	rm ./news

pack:
	GOOS=linux make build
	docker build -t eu.gcr.io/myproject/news-service:$(TAG) .

upload:
	gcloud docker -- push eu.gcr.io/myproject/news-service:$(TAG)

deploy:
	envsubst < k8s/deployment.yml | kubectl apply -f -

ship: test pack upload deploy clean
