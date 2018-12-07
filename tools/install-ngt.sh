#!/bin/sh

curl -LO https://github.com/yahoojapan/NGT/archive/v$NGT_VERSION.tar.gz
tar zxf v$NGT_VERSION.tar.gz
cd NGT-$NGT_VERSION
cmake .
make -j
make install
