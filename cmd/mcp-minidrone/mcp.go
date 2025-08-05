package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var s *server.MCPServer
var httpSrv *server.StreamableHTTPServer

var mu sync.Mutex

var (
	errDroneNotAvailable = errors.New("Minidrone not available")
)

func startMCP(port string) {
	s = server.NewMCPServer(
		"TinyGo Minidrone",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	addToolTakeoff()
	addToolLand()
	addToolHover()
	addToolUp()
	addToolDown()
	addToolForward()
	addToolBackward()
	addToolRight()
	addToolLeft()
	addToolClockwise()
	addToolCounterClockwise()
	addToolFrontFlip()
	addToolBackFlip()

	addToolIsFlying()
	addResourceFlying()

	httpServer := server.NewStreamableHTTPServer(s)
	log.Printf("MCP server listening on http %s", port)
	if err := httpServer.Start(port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func addToolTakeoff() {
	tool := mcp.NewTool("takeoff",
		mcp.WithDescription("Causes the minidrone to takeoff"),
	)

	s.AddTool(tool, takeoffToolHandler)
}

func takeoffToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := "takeoff"
	if drone == nil {
		return mcpError(name, errDroneNotAvailable), nil
	}

	mu.Lock()
	defer mu.Unlock()

	err := drone.TakeOff()
	if err != nil {
		return mcpError(name, err), nil
	}

	return mcpSuccess(name, "minidrone taking off"), nil
}

func addToolLand() {
	tool := mcp.NewTool("land",
		mcp.WithDescription("Causes the minidrone to land"),
	)

	s.AddTool(tool, landToolHandler)
}

func landToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := "land"
	if drone == nil {
		return mcpError(name, errDroneNotAvailable), nil
	}

	mu.Lock()
	defer mu.Unlock()

	err := drone.Land()
	if err != nil {
		return mcpError(name, err), nil
	}

	return mcpSuccess(name, "minidrone landing"), nil
}

func addToolHover() {
	tool := mcp.NewTool("hover",
		mcp.WithDescription("Causes the minidrone to hover"),
	)

	s.AddTool(tool, hoverToolHandler)
}

func hoverToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := "hover"
	if drone == nil {
		return mcpError(name, errDroneNotAvailable), nil
	}

	mu.Lock()
	defer mu.Unlock()

	err := drone.Hover()
	if err != nil {
		return mcpError(name, err), nil
	}

	return mcpSuccess(name, "minidrone hovering"), nil
}

func addToolUp() {
	tool := mcp.NewTool("up",
		mcp.WithDescription("Causes the minidrone to move up"),
		mcp.WithNumber("speed",
			mcp.Description("speed from 0-100"),
			mcp.Required(),
		),
		mcp.WithNumber("duration",
			mcp.Description("for how long to move (should default to 500 milliseconds)"),
			mcp.Required(),
		),
	)

	s.AddTool(tool, upToolHandler)
}

func upToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := "up"
	if drone == nil {
		return mcpError(name, errDroneNotAvailable), nil
	}

	speed, err := request.RequireInt("speed")
	if err != nil {
		return mcpError(name, err), nil
	}

	duration, err := request.RequireInt("duration")
	if err != nil {
		return mcpError(name, err), nil
	}

	mu.Lock()
	defer mu.Unlock()

	err = drone.Up(speed)
	if err != nil {
		return mcpError(name, err), nil
	}

	time.AfterFunc(time.Duration(duration)*time.Millisecond, func() {
		drone.Hover()
	})

	return mcpSuccess(name, fmt.Sprintf("minidrone moving up at speed %d", speed)), nil
}

func addToolDown() {
	tool := mcp.NewTool("down",
		mcp.WithDescription("Causes the minidrone to move down"),
		mcp.WithNumber("speed",
			mcp.Description("speed from 0-100"),
			mcp.Required(),
		),
		mcp.WithNumber("duration",
			mcp.Description("for how long to move (should default to 500 milliseconds)"),
			mcp.Required(),
		),
	)

	s.AddTool(tool, downToolHandler)
}

func downToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := "down"
	if drone == nil {
		return mcpError(name, errDroneNotAvailable), nil
	}

	speed, err := request.RequireInt("speed")
	if err != nil {
		return mcpError(name, err), nil
	}

	duration, err := request.RequireInt("duration")
	if err != nil {
		return mcpError(name, err), nil
	}

	mu.Lock()
	defer mu.Unlock()

	err = drone.Down(speed)
	if err != nil {
		return mcpError(name, err), nil
	}

	time.AfterFunc(time.Duration(duration)*time.Millisecond, func() {
		drone.Hover()
	})

	return mcpSuccess(name, fmt.Sprintf("minidrone moving down at speed %d", speed)), nil
}

func addToolForward() {
	tool := mcp.NewTool("forward",
		mcp.WithDescription("Causes the minidrone to move forward"),
		mcp.WithNumber("speed",
			mcp.Description("speed from 0-100"),
			mcp.Required(),
		),
		mcp.WithNumber("duration",
			mcp.Description("for how long to move (should default to 500 milliseconds)"),
			mcp.Required(),
		),
	)

	s.AddTool(tool, forwardToolHandler)
}

func forwardToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := "forward"
	if drone == nil {
		return mcpError(name, errDroneNotAvailable), nil
	}

	speed, err := request.RequireInt("speed")
	if err != nil {
		return mcpError(name, err), nil
	}

	duration, err := request.RequireInt("duration")
	if err != nil {
		return mcpError(name, err), nil
	}

	mu.Lock()
	defer mu.Unlock()

	err = drone.Forward(speed)
	if err != nil {
		return mcpError(name, err), nil
	}

	time.AfterFunc(time.Duration(duration)*time.Millisecond, func() {
		drone.Hover()
	})

	return mcpSuccess(name, fmt.Sprintf("minidrone moving forward at speed %d", speed)), nil
}

func addToolBackward() {
	tool := mcp.NewTool("backward",
		mcp.WithDescription("Causes the minidrone to move backward"),
		mcp.WithNumber("speed",
			mcp.Description("speed from 0-100"),
			mcp.Required(),
		),
		mcp.WithNumber("duration",
			mcp.Description("for how long to move (should default to 500 milliseconds)"),
			mcp.Required(),
		),
	)

	s.AddTool(tool, backwardToolHandler)
}

func backwardToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := "backward"
	if drone == nil {
		return mcpError(name, errDroneNotAvailable), nil
	}

	speed, err := request.RequireInt("speed")
	if err != nil {
		return mcpError(name, err), nil
	}

	duration, err := request.RequireInt("duration")
	if err != nil {
		return mcpError(name, err), nil
	}

	mu.Lock()
	defer mu.Unlock()

	err = drone.Backward(speed)
	if err != nil {
		return mcpError(name, err), nil
	}

	time.AfterFunc(time.Duration(duration)*time.Millisecond, func() {
		drone.Hover()
	})

	return mcpSuccess(name, fmt.Sprintf("minidrone moving backward at speed %d", speed)), nil
}

func addToolRight() {
	tool := mcp.NewTool("right",
		mcp.WithDescription("Causes the minidrone to move right"),
		mcp.WithNumber("speed",
			mcp.Description("speed from 0-100"),
			mcp.Required(),
		),
		mcp.WithNumber("duration",
			mcp.Description("for how long to move (should default to 500 milliseconds)"),
			mcp.Required(),
		),
	)

	s.AddTool(tool, rightToolHandler)
}

func rightToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := "right"
	if drone == nil {
		return mcpError(name, errDroneNotAvailable), nil
	}

	speed, err := request.RequireInt("speed")
	if err != nil {
		return mcpError(name, err), nil
	}

	duration, err := request.RequireInt("duration")
	if err != nil {
		return mcpError(name, err), nil
	}

	mu.Lock()
	defer mu.Unlock()

	err = drone.Right(speed)
	if err != nil {
		return mcpError(name, err), nil
	}

	time.AfterFunc(time.Duration(duration)*time.Millisecond, func() {
		drone.Hover()
	})

	return mcpSuccess(name, fmt.Sprintf("minidrone moving right at speed %d", speed)), nil
}

func addToolLeft() {
	tool := mcp.NewTool("left",
		mcp.WithDescription("Causes the minidrone to move left"),
		mcp.WithNumber("speed",
			mcp.Description("speed from 0-100"),
			mcp.Required(),
		),
		mcp.WithNumber("duration",
			mcp.Description("for how long to move (should default to 500 milliseconds)"),
			mcp.Required(),
		),
	)

	s.AddTool(tool, leftToolHandler)
}

func leftToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := "left"
	if drone == nil {
		return mcpError(name, errDroneNotAvailable), nil
	}

	speed, err := request.RequireInt("speed")
	if err != nil {
		return mcpError(name, err), nil
	}

	duration, err := request.RequireInt("duration")
	if err != nil {
		return mcpError(name, err), nil
	}

	mu.Lock()
	defer mu.Unlock()

	err = drone.Left(speed)
	if err != nil {
		return mcpError(name, err), nil
	}

	time.AfterFunc(time.Duration(duration)*time.Millisecond, func() {
		drone.Hover()
	})

	return mcpSuccess(name, fmt.Sprintf("minidrone moving left at speed %d", speed)), nil
}

func addToolClockwise() {
	tool := mcp.NewTool("clockwise",
		mcp.WithDescription("Causes the minidrone to rotate in a clockwise direction"),
		mcp.WithNumber("speed",
			mcp.Description("speed from 0-100"),
			mcp.Required(),
		),
		mcp.WithNumber("duration",
			mcp.Description("for how long to move (should default to 500 milliseconds)"),
			mcp.Required(),
		),
	)

	s.AddTool(tool, clockwiseToolHandler)
}

func clockwiseToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := "clockwise"
	if drone == nil {
		return mcpError(name, errDroneNotAvailable), nil
	}

	speed, err := request.RequireInt("speed")
	if err != nil {
		return mcpError(name, err), nil
	}

	duration, err := request.RequireInt("duration")
	if err != nil {
		return mcpError(name, err), nil
	}

	mu.Lock()
	defer mu.Unlock()

	err = drone.Clockwise(speed)
	if err != nil {
		return mcpError(name, err), nil
	}

	time.AfterFunc(time.Duration(duration)*time.Millisecond, func() {
		drone.Hover()
	})

	return mcpSuccess(name, fmt.Sprintf("minidrone rotating clockwise at speed %d", speed)), nil
}

func addToolCounterClockwise() {
	tool := mcp.NewTool("counter_clockwise",
		mcp.WithDescription("Causes the minidrone to rotate in a counter-clockwise direction"),
		mcp.WithNumber("speed",
			mcp.Description("speed from 0-100"),
			mcp.Required(),
		),
		mcp.WithNumber("duration",
			mcp.Description("for how long to move (should default to 500 milliseconds)"),
			mcp.Required(),
		),
	)

	s.AddTool(tool, counterClockwiseToolHandler)
}

func counterClockwiseToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := "counter_clockwise"
	if drone == nil {
		return mcpError(name, errDroneNotAvailable), nil
	}

	speed, err := request.RequireInt("speed")
	if err != nil {
		return mcpError(name, err), nil
	}

	duration, err := request.RequireInt("duration")
	if err != nil {
		return mcpError(name, err), nil
	}

	mu.Lock()
	defer mu.Unlock()

	err = drone.CounterClockwise(speed)
	if err != nil {
		return mcpError(name, err), nil
	}

	time.AfterFunc(time.Duration(duration)*time.Millisecond, func() {
		drone.Hover()
	})

	return mcpSuccess(name, fmt.Sprintf("minidrone rotating counter-clockwies at speed %d", speed)), nil
}

func addToolFrontFlip() {
	tool := mcp.NewTool("front_flip",
		mcp.WithDescription("Causes the minidrone to perform a front flip"),
	)

	s.AddTool(tool, frontFlipToolHandler)
}

func frontFlipToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := "front_flip"
	if drone == nil {
		return mcpError(name, errDroneNotAvailable), nil
	}

	mu.Lock()
	defer mu.Unlock()

	err := drone.FrontFlip()
	if err != nil {
		return mcpError(name, err), nil
	}

	return mcpSuccess(name, "minidrone performing a front flip"), nil
}

func addToolBackFlip() {
	tool := mcp.NewTool("back_flip",
		mcp.WithDescription("Causes the minidrone to perform a back flip"),
	)

	s.AddTool(tool, backFlipToolHandler)
}

func backFlipToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := "back_flip"
	if drone == nil {
		return mcpError(name, errDroneNotAvailable), nil
	}

	mu.Lock()
	defer mu.Unlock()

	err := drone.BackFlip()
	if err != nil {
		return mcpError(name, err), nil
	}

	return mcpSuccess(name, "minidrone performing a back flip"), nil
}

func addToolIsFlying() {
	tool := mcp.NewTool("is_flying",
		mcp.WithDescription("Checks to see if the Minidrone is currently in flight"),
	)

	s.AddTool(tool, isFlyingToolHandler)
}

func isFlyingToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := "is_flying"
	if drone == nil {
		return mcpError(name, errDroneNotAvailable), nil
	}

	response := `{"flying":"false"}`
	if drone.Flying {
		response = `{"flying":"true"}`
	}

	return mcpSuccess(name, response), nil
}

func addResourceFlying() {
	resource := mcp.NewResource(
		"drone://flying",
		"drone flight state",
		mcp.WithResourceDescription("Returns true if the minidrone is currently flying, otherwise returns false."),
		mcp.WithMIMEType("application/json"),
	)

	s.AddResource(resource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		response := `{"flying":"false"}`
		if drone.Flying {
			response = `{"flying":"true"}`
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "drone://flying",
				MIMEType: "application/json",
				Text:     response,
			},
		}, nil
	})
}

func mcpSuccess(name, content string) *mcp.CallToolResult {
	return mcp.NewToolResultText(fmt.Sprintf("{\"tool_name\": \"%s\", \"content\": \"%s\"}", name, content))
}

func mcpError(name string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("{\"tool_name\": \"%s\", \"error\": \"%s\"}", name, err.Error()))
}
