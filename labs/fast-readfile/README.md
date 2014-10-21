## 探寻C++最快的读取文件的方案

验证下面文章是否正确
https://www.byvoid.com/blog/fast-readfile

数据输入: urandom.go生成随机10000000个整数

环境:
vm虚拟机， 1G内存


结果如下：

    $ time  ./scanf_rd
    real    0m1.619s
    user    0m1.524s
    sys 0m0.093s


    $ time ./cin_rd
    real    0m6.121s
    user    0m5.199s
    sys 0m0.922s


    $ time ./cin_nosync
    real    0m2.176s
    user    0m1.637s
    sys 0m0.536s


