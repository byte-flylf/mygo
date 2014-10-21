#!/bin/bash

~/bin/ninja all

echo -e '\ntime scanf_rd'
time ./scanf_rd.exe

echo -e '\ntime cin_rd'
time ./cin_rd.exe

echo -e '\ntime cin_nosync'
time ./cin_nosync.exe

echo -e '\ntime analyse_rd'
time ./analyse_rd.exe

echo -e '\ntime analyse2_rd'
time ./analyse2_rd.exe

echo -e '\ntime map_rd'
time ./map_rd.exe
