build:
	docker build -t highloadcup .

run:
	docker build -t highloadcup .
	docker run -p 8000:80 -p 27017:27017 -v $(CURDIR)/test_accounts_220119/data:/tmp/data --rm -t -i highloadcup

test:
	docker build -t highloadcup .
	docker run -v $(CURDIR)/test_accounts_220119/data:/tmp/data --rm -t highloadcup "test app -v"

push:
	docker build -t highloadcup .
	docker tag highloadcup stor.highloadcup.ru/accounts/tall_beaver
	docker push stor.highloadcup.ru/accounts/tall_beaver

tester_install:
	go get -u github.com/atercattus/highloadcup_tester

tester_run: 
	docker build -t highloadcup .
	docker run -p 8000:80 -p 27017:27017 -v $(CURDIR)/test_accounts_220119/data:/tmp/data --rm --name highloadcup-tester-run -d highloadcup
	sleep 10
	highloadcup_tester -addr http://127.0.0.1:8000 -hlcupdocs ./test_accounts_220119/ -test -phase 1 -uri \/accounts\/filter\/ -diff true -tank 100
	# highloadcup_tester -addr http://127.0.0.1:8000 -hlcupdocs ./test_accounts_220119/ -test -phase 1 -uri \/accounts\/group\/ -diff true -tank 100
	docker stop highloadcup-tester-run
