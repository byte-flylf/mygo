#include <stdio.h>
#include <iostream>

using namespace std;

const int MAXN = 10000000;
int numbers[MAXN];

void cin_read()
{
    freopen("/tmp/data.txt", "r", stdin);
    for (int i = 0; i < MAXN; i++) {
        std::cin >> numbers[i];
    }
}

int main(void)
{
    cin_read();
    printf("%d %d\n", numbers[0], numbers[MAXN - 1]);
    return 0;
}
