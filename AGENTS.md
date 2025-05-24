# AGENTS.md

## Overview

This document explains how to set up, configure, and run agents in this repository using [OpenAI Codex](https://platform.openai.com/docs/guides/code).

---

## 1. Install Dependencies

Clone the repository and download all Go dependencies:

```sh
git clone https://github.com/DmytroShumeyko/Scanner.git
cd Scanner
go mod download
```

---

## 2. Configuration

Create a `.env` file in the repository root and provide your API keys and other
secrets. An example `.env` file is included in the repository. The application
uses these values to authenticate with supported exchanges and services.

---

## 3. Running the Scanner

Run the program using one of the available modes:

```sh
go run main.go -mode=spot          # spot markets only
go run main.go -mode=futures       # futures markets only
go run main.go -mode=spotfutures   # spot and futures together
go run main.go -mode=pump          # pump monitor
```

The application is written in Go (Go 1.24+) and requires the dependencies
downloaded above. The Telegram notifier requires valid credentials in the `.env`
file.

---

## 4. Running Checks

Before running checks, ensure all dependencies are downloaded:

```sh
go mod download
```

Then build the project and execute the tests:

```sh
go build ./...
go test ./...
```

`go test` needs network access to fetch modules, so some tests may fail if the environment is offline.
