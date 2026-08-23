package base

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"
)

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// copied from "github.com/golang/go/main.go"

// Execute execute the commands
func Execute() {
	flag.Parse() // flag参数是以-开头的，例如-c=config.json(和-c config.json等价)
	args := flag.Args()
	fmt.Printf("flag args num %d\n", len(args))
	fmt.Printf("flag args %v\n", flag.Args())
	fmt.Printf("os args num %d\n", len(os.Args))
	fmt.Printf("os args %v\n", os.Args)
	if len(args) < 1 {
		PrintUsage(os.Stderr, RootCommand)
		return
	}
	cmdName := args[0] // for error messages
	fmt.Printf("cmdName:%s\n", cmdName)
	if args[0] == "help" {
		Help(os.Stdout, args[1:])
		return
	}

BigCmdLoop:
	for bigCmd := RootCommand; ; {
		fmt.Println("---------------")
		fmt.Println("bigCmd LongName:", bigCmd.LongName())
		fmt.Println("bigCmd UsageLine:", bigCmd.UsageLine)
		fmt.Println("bigCmd long:", strings.TrimSpace(bigCmd.Long))
		fmt.Println("bigCmd short:", bigCmd.Short)
		fmt.Println("bigCmd cmd Long name:", bigCmd.LongName(), "short name:", bigCmd.Name())
		fmt.Println("---------------")
		fmt.Println("info CommandEnv.CommandsWidth:", CommandEnv.CommandsWidth) // Only the package "help" will set the value of this variable

		maxLongNameWidth := 0
		for _, cmd := range bigCmd.Commands {
			if len(cmd.LongName()) > maxLongNameWidth {
				maxLongNameWidth = len(cmd.LongName())
			}
			// 渲染模板
			// tmpl, err := template.New("usage").Parse(cmd.UsageLine)
			// if err != nil {
			// 	panic(err)
			// }
			// fmt.Printf("Template Data: %+v\n", makeTmplData(cmd))
			// var renderedUsage strings.Builder
			// err = tmpl.Execute(&renderedUsage, makeTmplData(cmd))
			// if err != nil {
			// 	panic(err)
			// }

			// // 打印格式化输出
			// fmt.Printf("bigCmd.Commands: [%02d] Long name: %-*s\n%*sshort name: %s\n%*sshort: %s\n\n",
			// 	i, maxLongNameWidth, cmd.LongName(),
			// 	maxLongNameWidth+13, "", cmd.Name(),
			// 	maxLongNameWidth+18, "", renderedUsage.String())
		}
		fmt.Println("info maxLongNameWidth:", maxLongNameWidth)

		for i, cmd := range bigCmd.Commands {
			// Render the template
			tmpl, err := template.New("usage").Parse(cmd.UsageLine)
			if err != nil {
				panic(err)
			}

			var renderedUsage strings.Builder
			err = tmpl.Execute(&renderedUsage, makeTmplData(cmd))
			if err != nil {
				panic(err)
			}

			// Print formatted output
			fmt.Printf("bigCmd.Commands: [%02d] Long name: %-*s\n%*sshort name: %s\n%*sshort: %s, \n%*ssub commands num: %d\n\n",
				i, maxLongNameWidth, cmd.LongName(),
				maxLongNameWidth+13, "", cmd.Name(),
				maxLongNameWidth+18, "", renderedUsage.String(),
				maxLongNameWidth+7, "", len(cmd.Commands))

			if cmd.Name() != args[0] {
				continue
			}
			if len(cmd.Commands) > 0 {
				// test sub commands
				bigCmd = cmd
				args = args[1:]
				if len(args) == 0 {
					PrintUsage(os.Stderr, bigCmd)
					SetExitStatus(2)
					Exit()
				}
				if args[0] == "help" {
					// Accept 'go mod help' and 'go mod help foo' for 'go help mod' and 'go help mod foo'.
					Help(os.Stdout, append(strings.Split(cmdName, " "), args[1:]...))
					return
				}
				cmdName += " " + args[0]
				continue BigCmdLoop
			}
			if !cmd.Runnable() {
				continue
			}
			cmd.Flag.Usage = func() { cmd.Usage() }
			if cmd.CustomFlags {
				args = args[1:]
			} else {
				cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
			}

			buildCommandText(cmd)
			cmd.Run(cmd, args)
			Exit()
			return
		}
		helpArg := ""
		if i := strings.LastIndex(cmdName, " "); i >= 0 {
			helpArg = " " + cmdName[:i]
		}
		fmt.Fprintf(os.Stderr, "%s %s: unknown command\nRun '%s help%s' for usage.\n", CommandEnv.Exec, cmdName, CommandEnv.Exec, helpArg)
		SetExitStatus(2)
		Exit()
	}
}

// Sort sorts the commands
func Sort() {
	sort.Slice(RootCommand.Commands, func(i, j int) bool {
		return SortLessFunc(RootCommand.Commands[i], RootCommand.Commands[j])
	})
}

// SortLessFunc used for sort commands list, can be override from outside
var SortLessFunc = func(i, j *Command) bool {
	return i.Name() < j.Name()
}
