#include <stdio.h>

const int MAXN = 10000000;
int numbers[MAXN];

void scanf_read()
{
    freopen("/tmp/data.txt", "r", stdin);
    for (int i = 0; i < MAXN; i++) {
        scanf("%d", &numbers[i]);
    }
}


int main(void)
{
    scanf_read();
    printf("%d %d\n", numbers[0], numbers[MAXN - 1]);
    return 0;
}

