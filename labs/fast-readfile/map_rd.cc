// mmap可以比read减少额外的拷贝
#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/mman.h>


const int MAXN = 10000000;
int numbers[MAXN];

void analyse(char *buf, int len)
{
    int i = 0;
    numbers[0] = 0;
    for (char *p = buf; p < buf + len; p++)
        if (*p == '\n') {
            numbers[++i] = 0;
        } else {
            numbers[i] = numbers[i] * 10 + *p - '0';
        }
}

void fread_analyse()
{
    int fd = open("/tmp/data.txt", O_RDONLY);
    int len = lseek(fd, 0, SEEK_END);
    char *mbuf = (char *) mmap(NULL, len, PROT_READ, MAP_PRIVATE, fd, 0);
    analyse(mbuf, len);
}

int main(void)
{
    fread_analyse();
    printf("%d %d\n", numbers[0], numbers[MAXN - 1]);
    return 0;
}



