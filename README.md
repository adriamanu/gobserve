# gobserve
Golang file watcher

```bash
go run main.go -c "go build main.go"
```

```bash
go run main.go -c "go build main.go" -files "**/**/*.txt"
```

```bash
echo "trigger watcher" > test.txt
go run main.go -c "go build main.go" -files test.txt
```