## 大数据量排序探索

背景：给超过内存存储空间的数据排序

- 测试数据量：1200000（个）

- 测试内存允许存放的最多数字：100（个）

(最优结果是`Test4`那一组)

#### Test1
2017/03/27 00:04:28 Config:  100 50 2
2017/03/27 00:04:28 Start Divide file into small files
2017/03/27 00:04:33 Done Divide file into small files,  12000
2017/03/27 00:04:33 Start File Merge sort round:  0
2017/03/27 00:04:46 Done File Merge sort round:  0 , files:  240
2017/03/27 00:04:46 Start File Merge sort round:  1
2017/03/27 00:04:55 Done File Merge sort round:  1 , files:  5
2017/03/27 00:04:55 Start File Merge sort round:  2
2017/03/27 00:04:59 Done File Merge sort round:  2 , files:  1
2017/03/27 00:04:59 Sort finished [test_data.out.2.0]
2017/03/27 00:04:59 Total Cost:  31 s

#### Test2
2017/03/27 00:03:02 Config:  100 34 3
2017/03/27 00:03:02 Start Divide file into small files
2017/03/27 00:03:07 Done Divide file into small files,  12000
2017/03/27 00:03:07 Start File Merge sort round:  0
2017/03/27 00:03:19 Done File Merge sort round:  0 , files:  353
2017/03/27 00:03:19 Start File Merge sort round:  1
2017/03/27 00:03:26 Done File Merge sort round:  1 , files:  11
2017/03/27 00:03:26 Start File Merge sort round:  2
2017/03/27 00:03:31 Done File Merge sort round:  2 , files:  1
2017/03/27 00:03:31 Sort finished [test_data.out.2.0]
2017/03/27 00:03:31 Total Cost:  29 s

#### Test3
2017/03/27 00:00:47 Config:  100 33 3
2017/03/27 00:00:47 Start Divide file into small files
2017/03/27 00:00:52 Done Divide file into small files,  12000
2017/03/27 00:00:52 Start File Merge sort round:  0
2017/03/27 00:01:05 Done File Merge sort round:  0 , files:  364
2017/03/27 00:01:05 Start File Merge sort round:  1
2017/03/27 00:01:12 Done File Merge sort round:  1 , files:  12
2017/03/27 00:01:12 Start File Merge sort round:  2
2017/03/27 00:01:17 Done File Merge sort round:  2 , files:  1
2017/03/27 00:01:17 Sort finished [test_data.out.2.0]
2017/03/27 00:01:17 Total Cost:  30 s

#### Test4
2017/03/26 23:58:53 Config:  100 25 4
2017/03/26 23:58:53 Start Divide file into small files
2017/03/26 23:58:58 Done Divide file into small files,  12000
2017/03/26 23:58:58 Start File Merge sort round:  0
2017/03/26 23:59:08 Done File Merge sort round:  0 , files:  480
2017/03/26 23:59:08 Start File Merge sort round:  1
2017/03/26 23:59:15 Done File Merge sort round:  1 , files:  20
2017/03/26 23:59:15 Start File Merge sort round:  2
2017/03/26 23:59:20 Done File Merge sort round:  2 , files:  1
2017/03/26 23:59:20 Sort finished [test_data.out.2.0]
2017/03/26 23:59:20 Total Cost:  27 s

#### Test5
2017/03/26 23:57:07 Config:  100 24 4
2017/03/26 23:57:07 Start Divide file into small files
2017/03/26 23:57:11 Done Divide file into small files,  12000
2017/03/26 23:57:11 Start File Merge sort round:  0
2017/03/26 23:57:23 Done File Merge sort round:  0 , files:  500
2017/03/26 23:57:23 Start File Merge sort round:  1
2017/03/26 23:57:29 Done File Merge sort round:  1 , files:  21
2017/03/26 23:57:29 Start File Merge sort round:  2
2017/03/26 23:57:35 Done File Merge sort round:  2 , files:  1
2017/03/26 23:57:35 Sort finished [test_data.out.2.0]
2017/03/26 23:57:35 Total Cost:  28 s

#### Test6
2017/03/26 23:55:33 Config:  100 20 5
2017/03/26 23:55:33 Start Divide file into small files
2017/03/26 23:55:38 Done Divide file into small files,  12000
2017/03/26 23:55:38 Start File Merge sort round:  0
2017/03/26 23:55:47 Done File Merge sort round:  0 , files:  600
2017/03/26 23:55:47 Start File Merge sort round:  1
2017/03/26 23:55:53 Done File Merge sort round:  1 , files:  30
2017/03/26 23:55:53 Start File Merge sort round:  2
2017/03/26 23:55:58 Done File Merge sort round:  2 , files:  2
2017/03/26 23:55:58 Start File Merge sort round:  3
2017/03/26 23:56:02 Done File Merge sort round:  3 , files:  1
2017/03/26 23:56:02 Sort finished [test_data.out.3.0]
2017/03/26 23:56:02 Total Cost:  29 s

#### Test7
2017/03/26 23:53:22 Config:  100 10 10
2017/03/26 23:53:23 Start Divide file into small files
2017/03/26 23:53:27 Done Divide file into small files,  12000
2017/03/26 23:53:27 Start File Merge sort round:  0
2017/03/26 23:53:35 Done File Merge sort round:  0 , files:  1200
2017/03/26 23:53:35 Start File Merge sort round:  1
2017/03/26 23:53:40 Done File Merge sort round:  1 , files:  120
2017/03/26 23:53:40 Start File Merge sort round:  2
2017/03/26 23:53:44 Done File Merge sort round:  2 , files:  12
2017/03/26 23:53:44 Start File Merge sort round:  3
2017/03/26 23:53:49 Done File Merge sort round:  3 , files:  2
2017/03/26 23:53:49 Start File Merge sort round:  4
2017/03/26 23:53:52 Done File Merge sort round:  4 , files:  1
2017/03/26 23:53:52 Sort finished [test_data.out.4.0]
2017/03/26 23:53:52 Total Cost:  30 s

#### Test8
2017/03/26 23:54:33 Config:  100 5 20
2017/03/26 23:54:33 Start Divide file into small files
2017/03/26 23:54:38 Done Divide file into small files,  12000
2017/03/26 23:54:38 Start File Merge sort round:  0
2017/03/26 23:54:45 Done File Merge sort round:  0 , files:  2400
2017/03/26 23:54:45 Start File Merge sort round:  1
2017/03/26 23:54:50 Done File Merge sort round:  1 , files:  480
2017/03/26 23:54:50 Start File Merge sort round:  2
2017/03/26 23:54:54 Done File Merge sort round:  2 , files:  96
2017/03/26 23:54:54 Start File Merge sort round:  3
2017/03/26 23:54:59 Done File Merge sort round:  3 , files:  20
2017/03/26 23:54:59 Start File Merge sort round:  4
2017/03/26 23:55:03 Done File Merge sort round:  4 , files:  4
2017/03/26 23:55:03 Start File Merge sort round:  5
2017/03/26 23:55:07 Done File Merge sort round:  5 , files:  1
2017/03/26 23:55:07 Sort finished [test_data.out.5.0]
2017/03/26 23:55:07 Total Cost:  34 s