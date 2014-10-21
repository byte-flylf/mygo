#!/usr/bin/env python
# encoding: utf-8

""" load sql file  after logcons cored"""

import  sys
import glob
import os

path = sys.argv[1]

files = glob.glob(path + "/*.sql")
for f in files:
    cmdstr = "/data/qqpet/logcons/bin/logadmin load --filename %s" %(f)
    os.system(cmdstr)
    os.remove(f)
