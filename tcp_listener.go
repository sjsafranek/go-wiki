package main

import (
	"encoding/json"
	"net"
	"runtime"
	"time"

	"github.com/sjsafranek/socket2em"
)

var (
	TCP_SERVER socket2em.Server
	tcp_port   int = 9622
	startTime  time.Time
)

func init() {
	startTime = time.Now()
}

func RunTcpServer() {

	TCP_SERVER = socket2em.Server{
		LoggingHandler: func(message string) { logger.Info(message) },
		Port:           tcp_port,
	}

	// Simple ping method
	TCP_SERVER.RegisterMethod("ping", func(message socket2em.Message, conn net.Conn) {
		// {"method": "ping"}
		TCP_SERVER.HandleSuccess(`{"message": "pong"}`, conn)
	})

	// Returns runtime and system information
	TCP_SERVER.RegisterMethod("get_runtime_stats", func(message socket2em.Message, conn net.Conn) {
		// {"method": "get_runtime_stats"}
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		results := make(map[string]interface{})
		results["NumGoroutine"] = runtime.NumGoroutine()
		results["Alloc"] = ms.Alloc / 1024
		results["TotalAlloc"] = ms.TotalAlloc / 1024
		results["Sys"] = ms.Sys / 1024
		results["NumGC"] = ms.NumGC
		results["Registered"] = startTime.UTC()
		results["Uptime"] = time.Since(startTime).Seconds()
		results["NumCPU"] = runtime.NumCPU()
		results["GOOS"] = runtime.GOOS
		TCP_SERVER.SendResponseFromStruct(results, conn)
	})

	// Create new user
	TCP_SERVER.RegisterMethod("create_user", func(message socket2em.Message, conn net.Conn) {
		// {"method": "create_user", "data":{"username":"stefan","password":"test"}}
		// logger.Info(message)
		var data map[string]string
		json.Unmarshal(message.Data, &data)
		logger.Info(data)

		// set user
		user := User{Username: data["username"]}
		user.SetPassword(data["password"])
		USERS.Add(&user)
		USERS.Save(users_file)

		results := make(map[string]interface{})
		TCP_SERVER.SendResponseFromStruct(results, conn)
	})

	go TCP_SERVER.Start()

}
