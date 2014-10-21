#include <stdio.h>
#include <fcntl.h>
#include <unistd.h>

const int MAXN = 10000000;
int numbers[MAXN];

const int MAXS = 100 * 1024 * 1024;
char buf[MAXS];

void analyse(char *buf)
{
    int i = 0;
    numbers[0] = 0;
    for (char *p = buf; *p; p++)
        if (*p == '\n') {
            i++;
            numbers[i] = 0;
        }else {
            numbers[i] = numbers[i] * 10 + *p - '0';
        }
}

void fread_analyse()
{
    int fd = open("/tmp/data.txt", O_RDONLY);
    int len = read(fd, buf, MAXS);
    buf[len] = '\0';
    analyse(buf);
}

int main(void)
{
    fread_analyse();
    printf("%d %d\n", numbers[0], numbers[MAXN - 1]);
    return 0;
}
