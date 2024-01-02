APP=secret-scanner

.PHONY: build-and-execute
build-and-execute:
	go build -o ${APP} cmd/secret-scanner/main.go && chmod +x ./${APP} && ./${APP}

.PHONY: build
build:
	go build -o ${APP} cmd/secret-scanner/main.go

.PHONY: run
run:
	./${APP}

.PHONY: debug
debug: 
	export DEBUG=True && make build-and-execute
	
.PHONY: prod
prod: 
	export PROD=True && make build-and-execute
