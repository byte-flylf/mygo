测试几种链表遍历查找的性能差异。

实验结果：

    dada-imac:labs dada$ go test -test.bench="." labs07
    PASS
    Benchmark_Loop1    500000          4386 ns/op
    Benchmark_Loop2    500000          5206 ns/op
    Benchmark_Loop3    100000         14821 ns/op
    Benchmark_Loop4     10000        135384 ns/op
    Benchmark_Loop5      5000        607735 ns/op
    Benchmark_Loop6    500000          7067 ns/op
    ok      labs07  14.618s    

确认算法有效：

    dada-imac:labs dada$ go test labs07 -v
    === RUN Test_Loop4
    --- PASS: Test_Loop4 (0.00 seconds)
    === RUN Test_Loop5
    --- PASS: Test_Loop5 (0.00 seconds)
    === RUN Test_Loop6
    --- PASS: Test_Loop6 (0.00 seconds)
    PASS
    ok      labs07  0.016s

结论：有空做个查询表达式解析器吧。