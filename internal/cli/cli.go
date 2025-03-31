package cli

import (
	"flag"
	"fmt"
	"io"
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
	opts := &Options{}

	// 创建一个新的FlagSet
	fs := flag.NewFlagSet("jhttp", flag.ContinueOnError)
	fs.SetOutput(io.Discard) // 禁止FlagSet自己输出错误信息

	// 定义标志
	fs.StringVar(&opts.EnvFile, "env-file", "", "指定环境变量文件路径")
	fs.StringVar(&opts.Env, "env", "", "指定使用的环境名称")
	fs.StringVar(&opts.RequestName, "request", "", "指定要执行的请求名称")
	fs.StringVar(&opts.OutputFile, "output", "", "指定响应输出文件")
	fs.BoolVar(&opts.Verbose, "verbose", false, "输出详细信息")
	fs.BoolVar(&opts.ShowVersion, "version", false, "显示版本信息")
	fs.BoolVar(&opts.ShowHelp, "help", false, "显示帮助信息")
	fs.BoolVar(&opts.ListRequests, "list", false, "列出所有请求名称")

	// 解析参数
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	// 获取剩余的位置参数
	remaining := fs.Args()
	if len(remaining) > 0 {
		opts.HTTPFile = remaining[0]
	}

	return opts, nil
}

// PrintUsage 打印使用说明
func PrintUsage(w io.Writer, progName string) {
	fmt.Fprintf(w, "用法: %s [选项] <http-file>\n\n", progName)
	fmt.Fprintf(w, "描述:\n")
	fmt.Fprintf(w, "  %s 是一个命令行工具，用于执行IntelliJ IDEA格式的.http文件。\n\n", progName)
	fmt.Fprintf(w, "选项:\n")
	fmt.Fprintf(w, "  --env-file <file>     指定环境变量文件路径\n")
	fmt.Fprintf(w, "  --env <n>          指定使用的环境名称（支持自动查找环境文件）\n")
	fmt.Fprintf(w, "  --request <n>      指定要执行的请求名称\n")
	fmt.Fprintf(w, "  --output <file>       指定响应输出文件\n")
	fmt.Fprintf(w, "  --verbose             输出详细信息\n")
	fmt.Fprintf(w, "  --version             显示版本信息\n")
	fmt.Fprintf(w, "  --help                显示帮助信息\n")
	fmt.Fprintf(w, "  --list                列出所有请求名称\n\n")
	fmt.Fprintf(w, "请求名称格式说明:\n")
	fmt.Fprintf(w, "  请求名称以'###'开头定义，例如：### 获取用户信息\n")
	fmt.Fprintf(w, "  紧随其后的注释行（以'#'开头）会被保存为请求的描述，而不会成为请求名称的一部分\n")
	fmt.Fprintf(w, "  使用--request参数时，需要使用完整的请求名（不包含注释内容）\n")
	fmt.Fprintf(w, "  如遇到请求无法匹配的情况，请使用--list选项查看实际的请求名称\n\n")
	fmt.Fprintf(w, "示例:\n")
	fmt.Fprintf(w, "  %s example.http\n", progName)
	fmt.Fprintf(w, "  %s --env 开发环境 example.http           # 自动查找环境文件\n", progName)
	fmt.Fprintf(w, "  %s --env-file env.json --env 开发环境 example.http\n", progName)
	fmt.Fprintf(w, "  %s --request \"获取用户信息\" example.http\n", progName)
}
