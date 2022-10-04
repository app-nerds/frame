.DEFAULT_GOAL := help
.PHONY: help setup

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

VERSION := $(shell cat ../VERSION)

setup: ## Sets up dependencies
	yarn global add rollup uglify-js uglifycss

build: clean build-library copy ## Builds Nerd JS Library for distribution

build-debug: clean rollup ## Builds a non-minified version of Nerd JS Library

clean: ## Removes all files from the dist folder
	@rm -rf dist/*

build-library: rollup
	uglifyjs dist/frame.js -c -m -o dist/frame.min.js --comments '/v\d.\d.\d/' --source-map
	uglifycss --output dist/base.min.css src/css/base.css 
	uglifycss --output dist/admin-left-side-nav.min.css src/css/base.css src/css/admin-left-side-nav.css
	uglifycss --output dist/components.min.css src/dialogs/alert.css \
		src/dialogs/confirm.css src/menus/popup-menu.css src/shim/shim.css \
		src/spinner/spinner.css src/members/member-login-bar.css \
		src/members/google-login-form.css src/message-bar/message-bar.css

rollup:
	rollup src/frame.js --name dist/frame.js --format es -o dist/frame.js --banner "/* Copyright © 2022 App Nerds LLC $(VERSION) */"

copy:
	cp dist/*.css examples/
	cp dist/*.js examples/
