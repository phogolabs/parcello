#!/bin/bash

echo > coverage.txt

find ./*.coverprofile -maxdepth 10 -type f -exec cat {} >> coverage.txt \;
