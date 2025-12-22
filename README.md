# Y - Composite Command Tree

The idea is borrowed from Robert S. Muhlestein's https://github.com/rwxrob/bonzai

## Run

### All tools

```go
go run cmd/main.go
```

### Specific tools only

Every tool can run or build as own tool

```go
go run commands/retention/cmd/y-<tool name>.go
```

## Tools

### Replacer

Replacer is a CLI tool for replacing or appending lines in a file based on regular expression patterns.

For each --replace argument, replacer searches for a matching line. If a match is found, the line is replaced.
If no match is found, the replacement line can optionally be appended at the end of the file (this is the default).

- Each --replace must be in the format REGEXP:LINE.
- You can specify --replace multiple times.
- You are responsible for adding ^ (line start) or $ (line end) anchors to your regular expressions if desired.
- For replacements where the pattern is not found, the line is appended only if --append-missing is true (default: true).

Examples:

#### Replace all lines starting with "Hello" with "Greetings!"<br><br>

```shell
y-replacer --file text.txt --replace "^Hello:Greetings!"
```

#### Replace lines starting with "But" or "Hello"

```shell
y-replacer --file text.txt --replace "^But:HOWEVER" --replace "^Hello:GOOD DAY"
```

#### Append the replacement line at the end if no match is found (default behavior)

```shell
y-replacer --file text.txt --replace "^NotFound:New line"
```

#### Do NOT append if no match is found (replace only)

```shell
y-replacer --file text.txt --replace "^NotFound:New line" --append-missing=false
```

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

### Rousego

Startes multiple processes in parallel and unifies output of these processes. A `rousego.toml` file has to be in the directory in which rousego starts. Normally your project main directory.

Example 1: rousego.toml

```toml
[[cmds]]
label = "Backend"
cmd = "make serve"

[[cmds]]
label = "Frontend"
cmd = "cd frontend ; npm run dev"
```

Example 2: rousego.toml

```toml
[[cmds]]
label = "Backend2"
cmd = "sleep 3 ; echo 'hello from backend via stdout'; sleep 4"


[[cmds]]
label = "Frontend"
cmd = "sleep 2 ; echo 'hello from frontend via stdout' ; sleep 1"

[[cmds]]
label = "Frontend2"
cmd = "sleep 1 ; echo 'hello from frontend2 via stderr' 1>&2 ; sleep 2"
```
