Golang key/value db-bench (pogreb-bench)
========================================

pogreb-bench is a key-value store benchmarking tool. Currently it supports pogreb, goleveldb, bolt,  badgerdb, slowpoke and pudge.


Some tests, MacBook Pro (Retina, 13-inch, Early 2015)
=====================================================


### Round 1
Number of keys: 1000000
Minimum key size: 16, maximum key size: 64
Minimum value size: 128, maximum value size: 512
Concurrency: 2


|                       | pogreb  | goleveldb | bolt   | badgerdb | pudge  | slowpoke | pudge(mem) |
|-----------------------|---------|-----------|--------|----------|--------|----------|------------|
| 1M (Put+Get), seconds | 187     | 38        | 126    | 34       | 23     | 23       | 2          |
| 1M Put, ops/sec       | 5336    | 34743     | 8054   | 33539    | 47298  | 46789    | 439581     |
| 1M Get, ops/sec       | 1782423 | 98406     | 499871 | 220597   | 499172 | 445783   | 1652069    |
| FileSize,Mb           | 568     | 357       | 552    | 487      | 358    | 358      | 358        |



### Round 2
Number of keys: 2000000
Key size: 16
Value size: 128
Concurrency: 1


|                       | pogreb  | goleveldb | bolt   | badgerdb | pudge  | slowpoke | pudge(mem) |
|-----------------------|---------|-----------|--------|----------|--------|----------|------------|
| 2M (Put+Get), seconds | 512     | 59        | 199    | 89       | 62     | 56       | 5          |
| 2M Put, ops/sec       | 3922    | 69029     | 10344  | 27368    | 58135  | 59590    | 553112     |
| 2M Get, ops/sec       | 947348  | 64561     | 329248 | 125174   | 70613  | 86120    | 1014628    |
| FileSize,Mb           | 1010    | 296       | 456    | 516      | 305    | 305      | 305        |


### Round 3
Number of keys: 10000000
Key size: 8
Value size: 16
Concurrency: 10


|                       | goleveldb | badgerdb | pudge  |
|-----------------------|-----------|----------|--------|
| 10M (Put+Get), seconds| 216       | 190      | 253    |
| 10M Put, ops/sec      | 95497     | 70840    | 42116  |
| 10M Get, ops/sec      | 89390     | 202284   | 617683 |
| FileSize,Mb           | 608       | 1870     | 686    |