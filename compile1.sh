#! /bin/sh -e

race=
# race=-race
name=base50

cd cmd/$name

if ! go build $race; then
    exit 1
fi

bin=$name

if [ "x$(go env GOOS)" = "xwindows" ]; then
bin=$name.exe
fi

mv $bin $name.bin.$(go env GOOS)-$(go env GOARCH)

