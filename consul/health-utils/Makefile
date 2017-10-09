build:
	GOOS=linux go build -o bin/mem-check_linux_amd64 mem-check/mem-check.go;
	GOOS=linux go build -o bin/disk-check_linux_amd64 disk-check/disk-check.go;
	GOOS=linux go build -o bin/cpu-check_linux_amd64 cpu-check/cpu-check.go;
	md5 bin/cpu-check_linux_amd64 bin/disk-check_linux_amd64 bin/mem-check_linux_amd64;