# mcprouter

OpenRouter for MCP Servers

![mcp-infa](./mcp-infa.png)

## Start Proxy Server

1. edit config file

```shell
cp .env.example.toml .env.toml
```

edit `.env.toml` as needed.

2. start proxy server

```shell
go run main.go proxy
```

3. add Proxy URL to MCP Client like Cursor

`http://localhost:8025/sse/fetch`

make sure you have set `mcp_server_commands.fetch` in `.env.toml`

## Start API Server

1. edit config file

```shell
cp .env.example.toml .env.toml
```

edit `.env.toml` as needed.

2. start api server

```shell
go run main.go api
```

3. request api with curl

```shell
curl -X POST http://127.0.0.1:8027/v1/list-tools \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer fetch'
```

make sure you have set `mcp_server_commands.fetch` in `.env.toml`

## Manage MCP Servers

### List Running MCP Servers

You can list currently running MCP servers:

```shell
curl -X GET http://127.0.0.1:8027/v1/list-running-servers \
-H 'Authorization: Bearer fetch' \
-H 'X-Client-Info: {"name": "admin", "version": "1.0.0"}'
```

### Stop MCP Server

You can stop a specific MCP server using the stop-server endpoint:

```shell
# 停止指定的MCP服务器
curl -X POST http://127.0.0.1:8027/v1/stop-server \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer fetch' \
-H 'X-Client-Info: {"name": "admin", "version": "1.0.0"}' \
-d '{"server_key": "fetch", "force": false}'

# 强制停止指定的MCP服务器
curl -X POST http://127.0.0.1:8027/v1/stop-server \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer fetch' \
-H 'X-Client-Info: {"name": "admin", "version": "1.0.0"}' \
-d '{"server_key": "fetch", "force": true}'
```

The `server_key` parameter specifies which MCP server to stop. You can get available server keys from the `list-running-servers` endpoint.

### Direct Proxy Management

You can also manage MCP servers directly through the proxy server:

```shell
# 列出运行中的客户端
curl -X GET http://127.0.0.1:8025/admin/list-clients

# 直接停止MCP服务器
curl -X POST http://127.0.0.1:8025/admin/stop-server \
-H 'Content-Type: application/json' \
-d '{"server_key": "fetch", "force": false}'
```
