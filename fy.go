package fy

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tidwall/pretty"
)

const (
	Name    = "fy"
	Version = "v0.1.0"
	Target  = "target"
	From    = "from"
	Pretty  = "pretty"
	Color   = "color"
)

func rootCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   Name,
		Short: "命令行翻译工具",
		Long:  "调用第三方API实现翻译",
		Run: func(c *cobra.Command, args []string) {
			if len(args) == 0 {
				c.Help()
				return
			}
			target, err := c.Flags().GetString(Target)
			if err != nil {
				log.Fatalln(err)
			}

			from, err := c.Flags().GetString(From)
			if err != nil {
				log.Fatalln(err)
			}
			ptty, err := c.Flags().GetBool(Pretty)
			if err != nil {
				log.Fatalln(err)
			}
			color, err := c.Flags().GetBool(Color)
			if err != nil {
				log.Fatalln(err)
			}
			var buf strings.Builder
			for _, v := range args {
				buf.WriteString(v)
				buf.WriteByte(' ')
			}
			query := buf.String()
			baidu := NewBaidu(from, target)
			result, err := baidu.Translation(query)
			if err != nil {
				log.Fatalln(err)
			}
			resultJson, err := json.Marshal(result)
			if ptty {
				if color {
					resultJson = pretty.Color(pretty.PrettyOptions(resultJson, pretty.DefaultOptions), nil)
				} else {
					resultJson = pretty.Pretty(resultJson)
				}
			}
			fmt.Println(string(resultJson))
		},
		Version: Version,
	}
	cmd.Flags().StringP(Target, "t", "zh", "翻译的目标语言")
	cmd.Flags().StringP(From, "f", "auto", "翻译的原文语言")
	cmd.Flags().BoolP(Pretty, "p", true, "是否美化输出")
	cmd.Flags().BoolP(Color, "c", color, "输出是否添加颜色, windows下cmd不支持颜色输出")
	return cmd
}

func Execute() {
	r := rootCmd()
	pretty.DefaultOptions.SortKeys = true
	if err := r.Execute(); err != nil {
		log.Fatalln(err)
	}
}
