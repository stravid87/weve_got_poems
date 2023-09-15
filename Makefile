.PHONY: app

start_sp:
	cd service_provider/wgp_backend && npm run start

build_frontend:
	cd service_provider/wgp_frontend && npm run build

run_dev:
	cd service_provider/wgp_frontend && npm run dev

build_interceptor: ## must do from a bash terminal
	cd layer8/interceptor/ && GOARCH=wasm GOOS=js go build -o ../../service_provider/wgp_frontend/public/interceptor.wasm

### Load Balancer
build-load-balancer:
	cd layer8/go-load-balancer/cmd && go build -o ../../bin/load-balancer

load-balancer: # Port 8000
	make build-load-balancer && ./layer8/bin/load-balancer

### Proxy Master 
build-layer8-master:
	cd layer8/proxy_master/cmd && go build -o ../../bin/layer8-master

layer8-master: # Port 9001
	make build-layer8-master && ./layer8/bin/layer8-master
	
generate-layer8-master-proto:
	cd layer8/proxy_master && protoc --go_out=. --go-grpc_out=. proto/Layer8Master.proto

## Proxy Slave
build-layer8-slave-one:
	cd layer8/proxy_slave/layer8-slave-one/cmd && go build -o ../../../bin/layer8-slave-one

layer8-slave-one: # Port 8001
	make build-layer8-slave-one && ./layer8/bin/layer8-slave-one

generate-layer8-slave-proto:
	cd go-layer8-slaves && protoc --go_out=. --go-grpc_out=. proto/Layer8Slave.proto