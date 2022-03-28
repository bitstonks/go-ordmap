# Benchmarking generic implementation vs codegenerated implementation of ordmap

To run tests run
```
make static_generics
make static_codegen
```

Look at Makefile to see exact commands that are run.

```
Test name timed in ns/op                          Codegen   	Generics  	Regression [%]
===========================================================================================
BenchmarkNodeBuiltin_Insert/small_value-12        877.6     	1012.0    	15.31
BenchmarkNodeBuiltin_Insert/small_value-12        868.8     	925.7     	6.55
BenchmarkNodeBuiltin_Insert/small_value-12        886.7     	885.7     	-0.11
BenchmarkNodeBuiltin_Insert/big_value-12          1045.0    	1064.0    	1.82
BenchmarkNodeBuiltin_Insert/big_value-12          1035.0    	1045.0    	0.97
BenchmarkNodeBuiltin_Insert/big_value-12          1055.0    	1088.0    	3.13
BenchmarkNode_Insert/small_value-12               1311.0    	1394.0    	6.33
BenchmarkNode_Insert/small_value-12               1282.0    	1382.0    	7.80
BenchmarkNode_Insert/small_value-12               1289.0    	1327.0    	2.95
BenchmarkNode_Insert/big_value-12                 1429.0    	1510.0    	5.67
BenchmarkNode_Insert/big_value-12                 1440.0    	1529.0    	6.18
BenchmarkNode_Insert/big_value-12                 1420.0    	1450.0    	2.11
BenchmarkNodeBuiltin_Iterate/small_value-12       85989346.0	111130503.0	29.24
BenchmarkNodeBuiltin_Iterate/small_value-12       87352686.0	112079874.0	28.31
BenchmarkNodeBuiltin_Iterate/small_value-12       88038056.0	110768368.0	25.82
BenchmarkNodeBuiltin_Iterate/big_value-12         109166776.0	111814402.0	2.43
BenchmarkNodeBuiltin_Iterate/big_value-12         112630307.0	113435613.0	0.71
BenchmarkNodeBuiltin_Iterate/big_value-12         107427107.0	114676874.0	6.75
BenchmarkNode_Iterate/small_value-12              104613330.0	111448318.0	6.53
BenchmarkNode_Iterate/small_value-12              109555176.0	113058099.0	3.20
BenchmarkNode_Iterate/small_value-12              99827096.0	113208582.0	13.40
BenchmarkNode_Iterate/big_value-12                112429322.0	118176978.0	5.11
BenchmarkNode_Iterate/big_value-12                113255481.0	115493919.0	1.98
BenchmarkNode_Iterate/big_value-12                111984078.0	116434568.0	3.97
BenchmarkNodeBuiltin_Get/small_value-12           138.9     	197.9     	42.48
BenchmarkNodeBuiltin_Get/small_value-12           138.6     	197.6     	42.57
BenchmarkNodeBuiltin_Get/small_value-12           139.4     	196.8     	41.18
BenchmarkNodeBuiltin_Get/big_value-12             140.0     	193.0     	37.86
BenchmarkNodeBuiltin_Get/big_value-12             140.3     	194.5     	38.63
BenchmarkNodeBuiltin_Get/big_value-12             139.7     	195.4     	39.87
BenchmarkNode_Get/small_value-12                  178.2     	219.5     	23.18
BenchmarkNode_Get/small_value-12                  176.2     	217.4     	23.38
BenchmarkNode_Get/small_value-12                  176.2     	216.3     	22.76
BenchmarkNode_Get/big_value-12                    183.4     	217.6     	18.65
BenchmarkNode_Get/big_value-12                    170.5     	218.2     	27.98
BenchmarkNode_Get/big_value-12                    170.2     	217.0     	27.50
BenchmarkNodeBuiltin_Remove/small_value-12        959.8     	992.7     	3.43
BenchmarkNodeBuiltin_Remove/small_value-12        952.2     	960.2     	0.84
BenchmarkNodeBuiltin_Remove/small_value-12        934.0     	989.6     	5.95
BenchmarkNodeBuiltin_Remove/big_value-12          983.4     	978.4     	-0.51
BenchmarkNodeBuiltin_Remove/big_value-12          962.9     	987.9     	2.60
BenchmarkNodeBuiltin_Remove/big_value-12          974.4     	979.9     	0.56
BenchmarkNode_Remove/small_value-12               1150.0    	1201.0    	4.43
BenchmarkNode_Remove/small_value-12               1106.0    	1163.0    	5.15
BenchmarkNode_Remove/small_value-12               1067.0    	1133.0    	6.19
BenchmarkNode_Remove/big_value-12                 1054.0    	1135.0    	7.69
BenchmarkNode_Remove/big_value-12                 1098.0    	1130.0    	2.91
BenchmarkNode_Remove/big_value-12                 1098.0    	1106.0    	0.73
```
