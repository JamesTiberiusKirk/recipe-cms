run: run_server run_styles run_proxy

run_styles:
	npm run tw:watch

run_server:
	air
	#templ generate --watch --proxy="http://127.0.0.1:5000" --cmd="go run main.go"

run_proxy:
	hrp -ignoreSuffix ".templ" -includeSuffix ".go,.css" -ignore "node_modules" -debug ./

gen: 
	templ generate
	npm run tw

kill_server:
	kill -9 $(lsof -ti:5000)

