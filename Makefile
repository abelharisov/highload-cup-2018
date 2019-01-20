run:
	docker build -t highloadcup .
	docker run -p 8000:80 -p 27017:27017 -v $(CURDIR)/test_accounts_291218/data:/tmp/data --rm -t -i highloadcup

test:
	docker build -t highloadcup .
	docker run -v $(CURDIR)/test_accounts_291218/data:/tmp/data --rm -t highloadcup "test app -v"

push:
	docker tag highloadcup stor.highloadcup.ru/accounts/tall_beaver
	docker push stor.highloadcup.ru/accounts/tall_beaver
