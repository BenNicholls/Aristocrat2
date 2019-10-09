#!/bin/bash

echo "Cleaning up!"

go fmt 

echo "Building Aristocrat 2"
echo "----------"

go build && ./Aristocrat2.exe
echo "----------"
echo "Yay, we had fun!"