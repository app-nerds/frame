#!/bin/bash

input="./$2"
output="./dist/$1"

uglifycss --output $output $input
