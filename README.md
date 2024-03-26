This is an attempt to solve 1 billion row challenge in fastest way possible using golang.

Run the below command to clear cache before running the program to get accurate results.
```go clean -cache```

Refer to this github repo [text](https://github.com/gunnarmorling/1brc) to create measurements file with 1 BILLION rows!

Run the code using
```
go build .
./1brc.exe -filepath="../measurements.txt" -cpuprofile="./result/cpuprofile.prof" -revision="rev1"
```