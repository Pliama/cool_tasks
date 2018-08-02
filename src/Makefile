LINTERCOMMAND=gometalinter --vendor --tests --skip=mock --deadline=1500s

packages = \
	./config \
	./database \
	./models \
    ./services/auth \
    ./services/users \
    ./services/tasks \
    ./services/events \
    ./services/flights \
    ./services/museums \
    ./services/trains \
    ./services/events \
    ./services/hotels \

all: test format migrate-up go-build dc-build dc-up

setup:
	go get -u gopkg.in/alecthomas/gometalinter.v2
	go get bitbucket.org/liamstask/goose/cmd/goose

test:
	@$(foreach package,$(packages), set -e; \
    		go test -coverprofile $(package)/cover.out -covermode=count $(package);)

perf-test:
	vegeta attack -targets=src/performance/test-login.txt -duration=5s -rate=150 | tee src/performance/results.bin | vegeta report
	cat src/performance/results.bin | vegeta report -reporter=plot > src/performance/plot.html

code-quality:
	$(LINTERCOMMAND) --checkstyle ./... > static-analysis.xml

code-quality-print:
	$(LINTERCOMMAND) ./...

format:
	go fmt $(packages)
	
migrate-down:
	goose -dir ./migrations postgres "user=ycomdeblqmlkxy password=16f2952f81f5b8fbde7aa903469618103d98615243a2ebcbc7c43c34888fc280 host=ec2-54-235-177-183.compute-1.amazonaws.com dbname=db1ut68n8ri2r3 sslmode=require" down-to 0

migrate-up:
	goose -dir ./migrations postgres "user=ycomdeblqmlkxy password=16f2952f81f5b8fbde7aa903469618103d98615243a2ebcbc7c43c34888fc280 host=ec2-54-235-177-183.compute-1.amazonaws.com dbname=db1ut68n8ri2r3 sslmode=require" up

doc:
	godoc -http=:6060
	http://localhost:6060/pkg/github.com/Nastya-Kruglikova/cool_tasks/

go-build:
	GOOS=linux GOARCH=amd64 go build -o cool_tasks

dc-build:
	docker-compose build

dc-up:
	docker-compose up
