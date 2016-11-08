#!/bin/bash -xe
cd /build/lightning
make clean
make
make -j3 check
make check-source
