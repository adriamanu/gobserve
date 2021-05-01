# Goverwatch
## Golang file watcher

Specify a command to execute and a list of files to watch on.<br>
If one of the file is modified, your command will be re-executed.
It is useful in a development environment to have hot reloading.

```bash
echo "trigger watcher" > test.txt
go run main.go -c "go build main.go" -files test.txt
```

```bash
go run main.go -c "go run test/server.go" -files "test/*.go"
```