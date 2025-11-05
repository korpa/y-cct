# Y - Composite Command Tree

The idea is borrowed from Robert S. Muhlestein's https://github.com/rwxrob/bonzai

## Run

### All tools

```go
go run cmd/main.go
```

### Specific tools only

Every tool can run as own tool

```go
go run commands/retention/cmd/main.go
```

## Tools

### Retention

Delete old backups. Keeps one backup per day of last 4 weeks. Delete all backups older than 4 weeks, but keep Monday backups and 1st of month.

```
backup_2025-09-01.tar.gz
backup_2025-09-02.tar.gz
.
.
backup_2025-09-29.tar.gz
backup_2025-09-30.tar.gz
backup_2025-10-01.tar.gz
backup_2025-10-02.tar.gz
.
.
backup_2025-10-30.tar.gz
backup_2025-10-31.tar.gz
backup_2025-11-01.tar.gz
backup_2025-11-02.tar.gz
```

```shell
ls | y-retention -f 2006-01-02 | xargs rm -r
```
