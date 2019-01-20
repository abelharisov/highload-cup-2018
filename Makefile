run:
	docker build -t highloadcup .
	docker run -p 8000:80 -p 27017:27017 --rm -t -i highloadcup

test:
	docker build -t highloadcup .
	docker run --rm -t highloadcup "test app -v"
