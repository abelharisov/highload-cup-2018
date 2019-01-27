
DATA_SAMPLE = test_accounts_240119

build:
	docker build -t highloadcup .

run:
	docker build -t highloadcup .
	docker run -m=2G --memory-swap=2G -p 8000:80 -p 27017:27017 -v $(CURDIR)/${DATA_SAMPLE}/data:/tmp/data --rm -t -i highloadcup

test:
	docker build -t highloadcup .
	docker run -m=2G --memory-swap=2G -v $(CURDIR)/${DATA_SAMPLE}/data:/tmp/data --rm -t highloadcup "test app -v"

push:
	docker build -t highloadcup .
	docker tag highloadcup stor.highloadcup.ru/accounts/tall_beaver
	docker push stor.highloadcup.ru/accounts/tall_beaver

tester_install:
	go get -u github.com/atercattus/highloadcup_tester

tester_run: 
	docker build -t highloadcup .
	docker run -m=2G --memory-swap=2G -p 8000:80 -p 27017:27017 -v $(CURDIR)/${DATA_SAMPLE}/data:/tmp/data --rm --name highloadcup-tester-run -d highloadcup
	sleep 10
	# highloadcup_tester -addr http://127.0.0.1:8000 -hlcupdocs ./${DATA_SAMPLE}/ -test -phase 1 -uri \/accounts\/filter\/ -diff true -tank 100
	# highloadcup_tester -addr http://127.0.0.1:8000 -hlcupdocs ./${DATA_SAMPLE}/ -test -phase 1 -uri \/accounts\/group\/ -diff true -tank 100
	# highloadcup_tester -addr http://127.0.0.1:8000 -hlcupdocs ./${DATA_SAMPLE}/ -test -phase 1 -diff true -tank 400
	highloadcup_tester -addr http://127.0.0.1:8000 -hlcupdocs ./${DATA_SAMPLE}/ -test -phase 2 -diff true -tank 400
	docker stop highloadcup-tester-run
