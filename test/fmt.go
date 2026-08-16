package main

import (
	"fmt"
	"strings"
	"text/template"
)

type Command struct {
	LongName  string
	Name      string
	UsageLine string
	Exec      string // 用于模板替换
}

func main() {
	commands := []Command{
		{"vlessenc", "vlessenc", "{{.Exec}} vlessenc", "xray"},
	}

	maxLongNameWidth := 8 // 假设最大宽度为 8

	for i, cmd := range commands {
		// 渲染模板
		tmpl, err := template.New("usage").Parse(cmd.UsageLine)
		if err != nil {
			panic(err)
		}
		var renderedUsage strings.Builder
		err = tmpl.Execute(&renderedUsage, cmd)
		if err != nil {
			panic(err)
		}

		// 打印格式化输出
		fmt.Printf("bigCmd.Commands: [%02d] Long name: %-*s\n%*sshort name: %s\n%*sshort: %s\n\n",
			i, maxLongNameWidth, cmd.LongName,
			maxLongNameWidth+13, "", cmd.Name,
			maxLongNameWidth+18, "", renderedUsage.String())
	}
}
