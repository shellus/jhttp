package cli

import (
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

// Options 包含命令行解析后的选项
type Options struct {
	HTTPFile     string // HTTP文件路径
	EnvFile      string // 环境变量文件
	Env          string // 环境名称
	RequestName  string // 请求名称（可选）
	OutputFile   string // 输出文件（可选）
	Verbose      bool   // 详细输出
	ShowVersion  bool   // 显示版本信息
	ShowHelp     bool   // 显示帮助信息
	ListRequests bool   // 列出所有请求
}

// ParseArgs 解析命令行参数
func ParseArgs(args []string) (*Options, error) {
	opts := &Options{
		Env: "local", // 默认使用local环境
	}

	// 检查是否有帮助或版本标志
	skipExecution := false
	for _, arg := range args {
		if arg == "-h" || arg == "--help" || arg == "-help" || 
		   arg == "-v" || arg == "--version" || arg == "-version" {
			skipExecution = true
			break
		}
	}

	// 预处理参数，以支持参数位置的灵活性
	var processedArgs []string
	var positionalArgs []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			processedArgs = append(processedArgs, arg)
		} else {
			positionalArgs = append(positionalArgs, arg)
		}
	}

	// 将位置参数附加到处理后的参数列表中
	processedArgs = append(processedArgs, positionalArgs...)

	// 创建应用
	app := &cli.App{
		Name:      "jhttp",
		Usage:     "执行IntelliJ IDEA格式的.http文件",
		UsageText: "jhttp [选项] <http-file> [request-name]",
		Version:   "0.1.0",
		Description: `一个命令行工具，用于执行IntelliJ IDEA格式的.http文件
请求名称格式说明:
  请求名称以'###'开头定义，例如：### 获取用户信息
  紧随其后的注释行（以'#'开头）会被保存为请求的描述，而不会成为请求名称的一部分
  使用-request参数或第二个位置参数指定请求名称时，需要使用完整的请求名（不包含注释内容）
  如遇到请求无法匹配的情况，请使用-list选项查看实际的请求名称

示例:
  jhttp example.http
  jhttp example.http "获取用户信息"    # 使用位置参数指定请求名称
  jhttp -e local example.http      # 使用短参数名
  jhttp -env-file env.json -env prod example.http
  jhttp -r "获取用户信息" example.http   # 指定请求名称`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "env-file",
				Aliases:     []string{"ef"},
				Usage:       "指定环境变量文件路径",
				Destination: &opts.EnvFile,
			},
			&cli.StringFlag{
				Name:        "env",
				Aliases:     []string{"e"},
				Usage:       "指定使用的环境名称（默认：local）",
				Value:       "local",
				Destination: &opts.Env,
			},
			&cli.StringFlag{
				Name:        "request",
				Aliases:     []string{"r"},
				Usage:       "指定要执行的请求名称",
				Destination: &opts.RequestName,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "指定响应输出文件",
				Destination: &opts.OutputFile,
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"vvv"},
				Usage:       "输出详细信息",
				Destination: &opts.Verbose,
			},
			&cli.BoolFlag{
				Name:        "list",
				Aliases:     []string{"l"},
				Usage:       "列出所有请求名称",
				Destination: &opts.ListRequests,
			},
		},
		Action: func(c *cli.Context) error {
			// 处理位置参数
			if c.NArg() > 0 {
				opts.HTTPFile = c.Args().Get(0)
				if c.NArg() > 1 && opts.RequestName == "" {
					opts.RequestName = c.Args().Get(1)
				}
			}
			
			return nil
		},
		// 设置HideHelp和HideVersion为false，确保帮助和版本标志被正常处理
		HideHelp:    false,
		HideVersion: false,
	}

	// 如果是显示帮助或版本信息，确保输出到控制台
	if !skipExecution {
		app.ExitErrHandler = func(context *cli.Context, err error) {
			// 不退出，让我们自己处理错误
		}
	}

	// 解析参数
	err := app.Run(append([]string{"jhttp"}, processedArgs...))
	
	// 如果是显示帮助或版本，在显示后退出应用
	if skipExecution {
		os.Exit(0)
	}
	
	return opts, err
}
