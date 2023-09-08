#!/bin/bash

year=$(date +%Y)
version=$(cat ../VERSION)
banner="/* Copyright © $year App Nerds LLC $version */"

input="./src/$1"
output="./dist/$1"
minoutput="./dist/$2"

rollup $input --format es -o $output --banner "$banner"
uglifyjs $output -c -m --keep-fargs -o $minoutput --source-map

