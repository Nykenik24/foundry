# Foundry
A CLI-based workspace manager for much better development flow.

## Installation
Install with

```bash
go install github.com/Nykenik24/foundry@latest
```

### Usage
Create a new project with `project init`

```bash
foundry project init <name>
```

Add tasks with `task new`

```bash
foundry task new greet
```

To add a command to the task, you can either use the `--cmd` flag

```bash
foundry task new greet --cmd "echo 'Hello\!'"
```

Or edit `.foundry/tasks.toml`

```toml
[tasks.greet]
name = "greet"
cmd = "echo 'Hello!'"
```
