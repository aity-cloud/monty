package main

import (
	"github.com/aity-cloud/monty/pkg/monty"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	monty.Execute()
}
