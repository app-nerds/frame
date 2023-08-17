.DEFAULT_GOAL := help
.PHONY: help setup

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

VERSION := $(shell cat ../VERSION)
YEAR := $(shell date +%Y)
BANNER := "/* Copyright © $(YEAR) App Nerds LLC $(VERSION) */"

setup: ## Sets up dependencies
	yarn global add rollup uglify-js uglifycss

run: ## Runs the example
	@cd examples && python3 -m http.server 8000

watch: ## Watches for changes and rebuilds
	watchman-make -p 'src/**/*.css' 'src/**/*.js' -t build

build: clean build-library copy ## Builds Nerd JS Library for distribution

build-debug: clean rollup ## Builds a non-minified version of Nerd JS Library

clean: ## Removes all files from the dist folder
	@rm -rf dist/*

build-library: rollup
	uglifyjs dist/frame.js -c -m --keep-fargs -o dist/frame.min.js --comments all --source-map
	uglifyjs dist/color-picker/color-picker.js -c -m --keep-fargs -o dist/color-picker/color-picker.min.js --comments all --source-map
	uglifyjs dist/date-time/date-time-picker.js -c -m --keep-fargs -o dist/date-time/date-time-picker.min.js --comments all --source-map
	uglifyjs dist/date-time/date-time-service.js -c -m --keep-fargs -o dist/date-time/date-time-service.min.js --comments all --source-map
	uglifyjs dist/date-time/date-time.js -c -m --keep-fargs -o dist/date-time/date-time.min.js --comments all --source-map
	uglifyjs dist/dialogs/dialogs.js -c -m --keep-fargs -o dist/dialogs/dialogs.min.js --comments all --source-map
	uglifyjs dist/http/fetcher.js -c -m --keep-fargs -o dist/http/fetcher.min.js --comments all --source-map
	uglifyjs dist/http/graphql.js -c -m --keep-fargs -o dist/http/graphql.min.js --comments all --source-map
	uglifyjs dist/http/http.js -c -m --keep-fargs -o dist/http/http.min.js --comments all --source-map
	uglifyjs dist/menus/popup-menu.js -c -m --keep-fargs -o dist/menus/popup-menu.min.js --comments all --source-map
	uglifyjs dist/message-bar/message-bar.js -c -m --keep-fargs -o dist/message-bar/message-bar.min.js --comments all --source-map
	uglifyjs dist/sessions/session-service.js -c -m --keep-fargs -o dist/sessions/session-service.min.js --comments all --source-map
	uglifyjs dist/shim/shim.js -c -m --keep-fargs -o dist/shim/shim.min.js --comments all --source-map
	uglifyjs dist/spa/spa.js -c -m --keep-fargs -o dist/spa/spa.min.js --comments all --source-map
	uglifyjs dist/spinner/spinner.js -c -m --keep-fargs -o dist/spinner/spinner.min.js --comments all --source-map
	uglifyjs dist/utilities.js -c -m --keep-fargs -o dist/utilities.min.js --comments all --source-map

rollup: build-css
	rollup src/frame.js --name dist/frame.js --format es -o dist/frame.js --banner $(BANNER)
	rollup src/color-picker/color-picker.js --format es -o dist/color-picker/color-picker.js --banner $(BANNER)
	rollup src/datetime/date-time-picker.js --format es -o dist/date-time/date-time-picker.js --banner $(BANNER)
	rollup src/datetime/date-time-service.js --format es -o dist/date-time/date-time-service.js --banner $(BANNER)
	rollup src/datetime/index.js --format es -o dist/date-time/date-time.js --banner $(BANNER)
	rollup src/dialogs/index.js --format es -o dist/dialogs/dialogs.js --banner $(BANNER)
	rollup src/http/fetcher.js --format es -o dist/http/fetcher.js --banner $(BANNER)
	rollup src/http/graphql.js --format es -o dist/http/graphql.js --banner $(BANNER)
	rollup src/http/index.js --format es -o dist/http/http.js --banner $(BANNER)
	rollup src/menus/popup-menu.js --format es -o dist/menus/popup-menu.js --banner $(BANNER)
	rollup src/message-bar/message-bar.js --format es -o dist/message-bar/message-bar.js --banner $(BANNER)
	rollup src/sessions/session-service.js --format es -o dist/sessions/session-service.js --banner $(BANNER)
	rollup src/shim/shim.js --format es -o dist/shim/shim.js --banner $(BANNER)
	rollup src/spa/spa.js --format es -o dist/spa/spa.js --banner $(BANNER)
	rollup src/spinner/spinner.js --format es -o dist/spinner/spinner.js --banner $(BANNER)
	rollup src/utilities/index.js --format es -o dist/utilities.js --banner $(BANNER)

build-css:
	uglifycss --output dist/base.min.css src/css/base.css
	uglifycss --output dist/icons.min.css src/css/icons.css
	uglifycss --output dist/admin-left-side-nav.min.css src/css/base.css src/css/admin-left-side-nav.css
	uglifycss --output dist/admin-top-nav.min.css src/css/base.css src/css/admin-top-nav.css
	uglifycss --output dist/components.min.css src/dialogs/alert.css \
		src/datetime/date-time-picker.css \
		src/dialogs/confirm.css src/menus/popup-menu.css src/shim/shim.css \
		src/spinner/spinner.css src/members/member-login-bar.css \
		src/members/google-login-form.css src/message-bar/message-bar.css \
		src/color-picker/color-picker.css

copy:
	cp dist/*.css ../admin-static/css/
	cp dist/frame.min.js ../admin-static/js/
	cp dist/*.css ../cmd/frame/templates
	cp dist/frame.min.js ../cmd/frame/templates

