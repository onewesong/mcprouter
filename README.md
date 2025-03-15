# mcprouter

SSE Proxy for MCP Servers

## Quick Start

1. edit config file

```shell
cp .env.toml.example .env.toml
```

edit `.env.toml` as needed.

2. start http server

```shell
go run main.go server
```

3. add Proxy URL to MCP Client like Cursor

`http://localhost:8025/sse/github`

make sure you have set `mcp_server_commands.github` in `.env.toml`
