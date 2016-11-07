#!/bin/bash -xe
cd /build/lightning
make clean
make
make check
make check-source
