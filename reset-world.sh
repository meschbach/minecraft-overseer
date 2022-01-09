#!/bin/bash

cd test/data
find . -d 1 |grep -v minecraft |grep -v '^.$'|xargs rm -fR