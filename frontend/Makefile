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

build: clean build-library copy clean-dist-after ## Builds Nerd JS Library for distribution

build-debug: clean rollup ## Builds a non-minified version of Nerd JS Library

clean: ## Removes all files from the dist folder
	@rm -rf dist/*

build-library: build-css
	@parallel --colsep ' ' -a ./roster-js.txt ./minifyjs.sh {1} {2}

build-css: build-type-files
	@parallel --colsep ',' -a ./roster-css.txt ./minifycss.sh {1} {2}

build-type-files: setup-directories
	npx -p typescript tsc src/**/*.js --declaration --allowJs --emitDeclarationOnly --outDir dist

setup-directories:
	@mkdir -p dist/datetime
	@mkdir -p dist/color-picker
	@mkdir -p dist/dialogs
	@mkdir -p dist/menus
	@mkdir -p dist/message-bar
	@mkdir -p dist/shim
	@mkdir -p dist/spinner

copy:
	cp dist/*.css ../admin-static/css/
	cp dist/frame.min.js ../admin-static/js/
	cp dist/*.css ../cmd/frame/templates
	cp dist/frame.min.js ../cmd/frame/templates

clean-dist-after:
	npx -p typescript tsc dist/frame.js --declaration --allowJs --emitDeclarationOnly --outDir dist
	cd dist && find . -type f -name "*.js" ! -name "*.min.js" -exec rm -f {} +
