# 夜 (Yoru)

Yoru is a fast, elegant terminal-based SSH and Telnet client with credential storage and tabbed sessions.

## What does it do?

- **Multiple Connections** - Open and manage several SSH/Telnet sessions in tabs
- **Credential Management** - Store and use SSH keys and passwords locally
- **Connection History** - Keep track of all your past connections with logs
- **Known Hosts Management** - View and manage SSH fingerprints for security
- **Offline First** - All data is stored locally, no cloud dependencies

## Building and Running

To build and run Yoru locally, you need to have [Go](https://golang.org/dl/) installed on your system.

The project uses [Makefile](https://www.gnu.org/software/make/) for build automation. Start by cloning the repository:

```bash
git clone https://github.com/luciferreeves/yoru.git
cd yoru
```

Then, you can build the project using:

```bash
make build
```

This will compile the source code and create an executable in the `./bin` directory. You can run Yoru with:

```bash
./bin/yoru
```

To start a development version:

```bash
make dev
```
