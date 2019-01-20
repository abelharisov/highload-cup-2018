run:
	docker build -t highloadcup .
	docker run --rm -t highloadcup

test:
	docker build -t highloadcup .
	docker run --rm -t highloadcup "test app -v"
