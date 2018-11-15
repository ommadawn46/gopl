#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

go test -bench .


# ❯ bash bench.bash
# goos: darwin
# goarch: amd64
# pkg: github.com/ommadawn46/the_go_programming_language-training/ch11/ex06
# BenchmarkPopCountTable1-4               	 5000000	       304 ns/op
# BenchmarkPopCountTable10-4              	 5000000	       319 ns/op
# BenchmarkPopCountTable100-4             	 2000000	       642 ns/op
# BenchmarkPopCountTable1000-4            	  300000	      3629 ns/op
# BenchmarkPopCountTable10000-4           	   50000	     34010 ns/op
# BenchmarkPopCountBitShift1-4            	30000000	        45.5 ns/op
# BenchmarkPopCountBitShift10-4           	 3000000	       443 ns/op
# BenchmarkPopCountBitShift100-4          	  300000	      4670 ns/op
# BenchmarkPopCountBitShift1000-4         	   30000	     43916 ns/op
# BenchmarkPopCountBitShift10000-4        	    3000	    429865 ns/op
# BenchmarkPopCountLowerBitClear1-4       	100000000	        10.1 ns/op
# BenchmarkPopCountLowerBitClear10-4      	20000000	        93.9 ns/op
# BenchmarkPopCountLowerBitClear100-4     	 2000000	       716 ns/op
# BenchmarkPopCountLowerBitClear1000-4    	  300000	      6042 ns/op
# BenchmarkPopCountLowerBitClear10000-4   	   30000	     55208 ns/op
# PASS
# ok  	github.com/ommadawn46/the_go_programming_language-training/ch11/ex06	25.874s
