# Minidrone MCP Server

This program creates an [MCP server](https://modelcontextprotocol.io/overview) that provides tools for controlling the Minidrone from any model that has tool calling support.

## Building

```shell
go build -o mcp-minidrone .
```

## Running

```shell
mcp-minidrone [MAC address or Bluetooth ID]
```

You can also use the `-port` flag to set a specific port for the MCP server.

```shell

```

Once it is running, you can call it from whatever MCP host/client that you wish.
