#!/bin/bash

go build -o ./test-binary main.go
trap "rm -rf ./test-binary" EXIT

for y in {2016..2021}; do 
	for m in {1..12}; do
		for d in {1..28}; do
			./test-binary \
				--driver=local tests/file*.txt \
				--date $(printf "%02d-%02d-%02d_%02d-%02d-%02d" ${y} ${m} ${d} $((1 + $RANDOM % 10)) $((1 + $RANDOM % 10)) $((1 + $RANDOM % 10))) \
				"$@"
		done
	done
done
