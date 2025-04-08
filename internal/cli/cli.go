package cli

import (
	"fmt"
	"io"
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

	// 预处理参数以检查是否显示版本或帮助
	showVersion := false
	showHelp := false
	for _, arg := range args {
		if arg == "-v" || arg == "-version" || arg == "--version" {
			showVersion = true
		}
		if arg == "-h" || arg == "-help" || arg == "--help" {
			showHelp = true
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
			
			// 捕获内置help和version标志的状态
			opts.ShowHelp = c.Bool("help") || c.Bool("h")
			opts.ShowVersion = c.Bool("version") || c.Bool("v")
			
			return nil
		},
	}

	// 如果不是显示版本或帮助，则禁用输出
	if !showVersion && !showHelp {
		app.Writer = io.Discard
		app.ErrWriter = io.Discard
	}

	// 解析参数
	err := app.Run(append([]string{"jhttp"}, processedArgs...))
	
	// 如果显示版本或帮助，我们已经完成任务了，所以返回错误来中止主程序
	if showVersion || showHelp {
		return opts, fmt.Errorf("显示信息完成")
	}
	
	return opts, err
}

// PrintUsage 打印使用说明
func PrintUsage(w io.Writer, progName string) {
	fmt.Fprintf(w, "用法: %s [选项] <http-file> [request-name]\n\n", progName)
	fmt.Fprintf(w, "选项:\n")
	fmt.Fprintf(w, "  -env-file, -ef <file>  指定环境变量文件路径\n")
	fmt.Fprintf(w, "  -env, -e <n>           指定使用的环境名称 (默认: local)\n")
	fmt.Fprintf(w, "  -request, -r <n>       指定要执行的请求名称\n")
	fmt.Fprintf(w, "  -output, -o <file>     指定响应输出文件\n")
	fmt.Fprintf(w, "  -verbose, -vvv               输出详细信息\n")
	fmt.Fprintf(w, "  -version, -v           显示版本信息\n")
	fmt.Fprintf(w, "  -help, -h              显示帮助信息\n")
	fmt.Fprintf(w, "  -list, -l              列出所有请求名称\n\n")
	
	fmt.Fprintf(w, "请求名称格式说明:\n")
	fmt.Fprintf(w, "  请求名称以'###'开头定义，例如：### 获取用户信息\n")
	fmt.Fprintf(w, "  紧随其后的注释行（以'#'开头）会被保存为请求的描述，而不会成为请求名称的一部分\n")
	fmt.Fprintf(w, "  使用-request参数或第二个位置参数指定请求名称时，需要使用完整的请求名（不包含注释内容）\n")
	fmt.Fprintf(w, "  如遇到请求无法匹配的情况，请使用-list选项查看实际的请求名称\n\n")
	
	fmt.Fprintf(w, "示例:\n")
	fmt.Fprintf(w, "  %s example.http\n", progName)
	fmt.Fprintf(w, "  %s example.http \"获取用户信息\"    # 使用位置参数指定请求名称\n", progName)
	fmt.Fprintf(w, "  %s -e local example.http      # 使用短参数名\n", progName)
	fmt.Fprintf(w, "  %s -env-file env.json -env prod example.http\n", progName)
	fmt.Fprintf(w, "  %s -r \"获取用户信息\" example.http   # 指定请求名称\n", progName)
}
