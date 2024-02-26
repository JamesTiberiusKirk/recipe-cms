run: run_server run_styles run_proxy

run_styles:
	npm run tw:watch

run_server:
	air

run_proxy:
	hrp -ignoreSuffix ".templ" -includeSuffix ".go,.css" -ignore "node_modules" -debug ./

gen: 
	templ generate
	npm run tw

setup:
	cp example.env .env
	migrator shema-up

install_deps: 
	go get ./...
	go install github.com/cosmtrek/air@latest
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/JamesTiberiusKirk/migrator/cmd/migrator@latest
	go install github.com/JamesTiberiusKirk/hot-reloader-proxy/cmd/hrp@latest

