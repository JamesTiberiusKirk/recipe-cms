install_dep:
	go install https://github.com/cespare/reflex@latest

run_dev: 
	reflex -d none -sr '.*\.(go|sql|templ)' -- templ generate && go run ./main.go
	# reflex -d none -sr '.*\.(go|sql|html|css)' -- go run ./main.go

get: 
	go get ./...
