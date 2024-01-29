run: run_server run_styles

run_server:
	npm run tw:watch

run_styles:
	templ generate --watch --proxy="http://127.0.0.1:5000" --cmd="go run main.go"

gen: 
	templ generate
	npm run tw


kill_server:
	kill -9 $(lsof -ti:5000)
