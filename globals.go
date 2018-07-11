package main

import (
	"github.com/sjsafranek/ligneous"
)

const (
	DEFAULT_PORT = 1337
)

var (
	logger = ligneous.New().Log
	USERS  Users
	PORT   int = DEFAULT_PORT
)
