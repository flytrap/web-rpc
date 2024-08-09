package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/flytrap/web-rpc/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var VERSION = "0.0.1"

// @BasePath /api/rpc/v1
// @title rpc service
// @in header
// @name Authorization
// @scheme bearer
// @schemes http https
func main() {
	ctx := context.Background()
	a := cli.NewApp()
	a.Name = "rbc-service"
	a.Version = VERSION
	a.Usage = "rbc service"
	a.Commands = []*cli.Command{
		newWebCmd(ctx),
	}
	err := a.Run(os.Args)
	if err != nil {
		logrus.Error(err.Error())
	}
}

func newWebCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "Run http server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "conf",
				Aliases: []string{"c"},
				Usage:   "App configuration file(.json,.yaml,.toml)",
			},
		},
		Action: func(c *cli.Context) error {
			config.MustLoad(c.String("conf"))
			config.PrintWithJSON()
			return runWeb(c.Context)
		},
	}
}

func runWeb(ctx context.Context) error {
	r := gin.Default()
	r.Any("/api/rpc/v1/run/:code", func(c *gin.Context) {
		if c.Request.Method != config.C.HTTP.Method {
			c.JSON(400, gin.H{"error": "method not allowed"})
			return
		}
		command := ""
		for _, cmd := range config.C.Commands {
			if cmd.Code == c.Param("code") {
				command = cmd.Exec
				break
			}
		}
		if command == "" {
			c.JSON(400, gin.H{"error": "command not found"})
			return
		}
		err := Exec(ctx, command)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "ok"})
	})
	r.Run(fmt.Sprintf("%s:%d", config.C.HTTP.Host, config.C.HTTP.Port))
	return nil
}

func Exec(ctx context.Context, cmd string) error {
	cs := strings.Split(cmd, " ")
	c := exec.Command(cs[0], cs[1:]...)
	out, err := c.CombinedOutput()
	logrus.Info(string(out))
	return err
}
