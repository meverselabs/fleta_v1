#! /bin/sh
set -eu

go get github.com/dvyukov/go-fuzz/go-fuzz
go get github.com/dvyukov/go-fuzz/go-fuzz-build

if [ ! -e toml-fuzz.zip ]; then
    go-fuzz-build github.com/fletaio/fleta_v1/cmd/config/go-toml
fi

rm -fr fuzz
mkdir -p fuzz/corpus
cp *.toml fuzz/corpus

go-fuzz -bin=toml-fuzz.zip -workdir=fuzz
