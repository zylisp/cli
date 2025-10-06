# zylisp/cli

Command-line REPL for Zylisp with support for local, server, and client modes.

## Features

- **Local REPL**: Interactive Zylisp interpreter
- **Server Mode**: Start a REPL server for remote connections
- **Client Mode**: Connect to a remote REPL server
- **Multiple Transports**: In-process, Unix domain sockets, and TCP
- **Context-aware**: Graceful shutdown with signal handling
- **Stateful Sessions**: Persistent environment across evaluations

## Usage

### Local Mode (Default)

Start an interactive REPL:

```bash
zylisp
```

This runs a local REPL session with an in-process evaluator.

### Server Mode

Start a REPL server that accepts remote connections:

```bash
# TCP server on port 5555
zylisp --mode=server --transport=tcp --addr=:5555

# Unix domain socket server
zylisp --mode=server --transport=unix --addr=/tmp/zylisp.sock
```

The server will run until interrupted (Ctrl-C).

### Client Mode

Connect to a remote REPL server:

```bash
# Connect via TCP
zylisp --mode=client --addr=localhost:5555

# Connect via Unix domain socket
zylisp --mode=client --addr=/tmp/zylisp.sock
```

## Command-Line Options

- `--mode`: Operation mode - `local` (default), `server`, or `client`
- `--transport`: Transport type - `in-process` (default), `unix`, or `tcp`
- `--addr`: Server address (required for server/client modes)
- `--codec`: Message codec - `json` (default) or `msgpack` (when supported)

## REPL Commands

- `exit`, `quit` - Exit the REPL
- `:reset` - Reset the environment (clear all definitions)
- `:help` - Show help message

## Examples

### Basic Arithmetic

```lisp
> (+ 1 2)
3

> (* 5 (+ 2 3))
25
```

### Defining Functions

```lisp
> (define square (lambda (x) (* x x)))
<function>

> (square 5)
25
```

### Stateful Operations

```lisp
> (define x 10)
10

> (define y 20)
20

> (+ x y)
30
```

For more examples, see [EXAMPLES.md](EXAMPLES.md).

## Remote REPL Example

Terminal 1 (Server):
```bash
$ zylisp --mode=server --transport=tcp --addr=:5555
Starting Zylisp REPL server on :5555 (tcp)
```

Terminal 2 (Client):
```bash
$ zylisp --mode=client --addr=localhost:5555
Connecting to Zylisp REPL at localhost:5555...
Connected!

> (+ 1 2)
3
```

## Status

✅ Production-ready implementation with full REPL protocol support
