#!/bin/bash

uin=0

while [[ $uin -le 99 ]]; do
    sql="DROP DATABASE petLog_${uin}"
    mysql -upetlog -pcGV0bG9nCg -h127.0.0.1 -e "$sql"
    uin=$((uin + 1))
done
echo "Finished"
