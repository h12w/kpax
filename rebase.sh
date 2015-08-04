#!/bin/sh

set -e

rm -r -f proto
wget https://github.com/stealthly/siesta/archive/master.zip
unzip master.zip
mv siesta-master proto
rm master.zip

for f in proto/*.go; do
	sed -i.bak 's/package siesta/package proto/' $f
done

rm proto/*.bak
