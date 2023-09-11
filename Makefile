start_sp:
	cd service_provider/wgp_backend && npm run start

build_frontend:
	cd service_provider/wgp_frontend && npm run build

build_interceptor: ## must do from a bash terminal
	cd layer8/interceptor/ && GOARCH=wasm GOOS=js go build -o ../../service_provider/wgp_frontend/src/assets/interceptor.wasm