package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shellus/jhttp/internal/cli"
	"github.com/shellus/jhttp/internal/environment"
	"github.com/shellus/jhttp/internal/executor"
	"github.com/shellus/jhttp/internal/parser"
)

const (
	exitSuccess = 0
	exitFailure = 1
)

func main() {
	// 设置程序名称
	progName := filepath.Base(os.Args[0])
	if progName == "" {
		progName = "jhttp"
	}

	// 解析命令行参数
	opts, err := cli.ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		cli.PrintUsage(os.Stderr, progName)
		os.Exit(exitFailure)
	}

	// 显示版本信息并退出
	if opts.ShowVersion {
		fmt.Printf("%s 版本 0.1.0\n", progName)
		os.Exit(exitSuccess)
	}

	// 显示帮助信息并退出
	if opts.ShowHelp {
		cli.PrintUsage(os.Stdout, progName)
		os.Exit(exitSuccess)
	}

	// 检查是否提供了HTTP文件
	if opts.HTTPFile == "" {
		fmt.Fprintln(os.Stderr, "错误: 必须提供一个.http文件")
		cli.PrintUsage(os.Stderr, progName)
		os.Exit(exitFailure)
	}

	// 解析HTTP文件
	httpFile, err := parser.ParseFile(opts.HTTPFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "解析文件错误: %v\n", err)
		os.Exit(exitFailure)
	}

	// 列出请求并退出
	if opts.ListRequests {
		fmt.Printf("文件 '%s' 中的请求:\n", opts.HTTPFile)
		if len(httpFile.Requests) == 0 {
			fmt.Println("  没有找到请求")
		} else {
			for i, req := range httpFile.Requests {
				displayName := req.Name
				if displayName == "" {
					displayName = fmt.Sprintf("%s %s", req.Method, req.URL)
				}
				fmt.Printf("%3d. %s\n", i+1, displayName)
				if req.Description != "" {
					// 对描述进行处理，确保多行描述缩进对齐
					descLines := strings.Split(req.Description, "\n")
					for _, line := range descLines {
						fmt.Printf("     描述: %s\n", line)
					}
				}
			}
		}
		os.Exit(exitSuccess)
	}

	// 处理环境变量文件
	var envVars map[string]string

	if opts.EnvFile != "" && opts.Env != "" {
		// 用户指定了环境文件路径和环境名称
		envVars, err = environment.LoadEnvFile(opts.EnvFile, opts.Env)
		if err != nil {
			fmt.Fprintf(os.Stderr, "加载环境变量错误: %v\n", err)
			os.Exit(exitFailure)
		}

		if opts.Verbose {
			fmt.Printf("已从 '%s' 加载环境 '%s' 中的 %d 个变量\n", opts.EnvFile, opts.Env, len(envVars))
		}
	} else if opts.Env != "" {
		// 用户只指定了环境名称，尝试自动查找环境文件
		envFilePath, found, needWarning := environment.FindEnvFile(opts.HTTPFile)
		if found {
			envVars, err = environment.LoadEnvFile(envFilePath, opts.Env)
			if err != nil {
				fmt.Fprintf(os.Stderr, "加载环境变量错误: %v\n", err)
				os.Exit(exitFailure)
			}

			if needWarning {
				fmt.Fprintf(os.Stderr, "警告: 自动使用了上级目录中的环境文件 '%s'\n", envFilePath)
			}

			if opts.Verbose {
				fmt.Printf("已从 '%s' 加载环境 '%s' 中的 %d 个变量\n", envFilePath, opts.Env, len(envVars))
			}
		} else {
			fmt.Fprintf(os.Stderr, "错误: 未找到环境文件，但指定了环境名称 '%s'\n", opts.Env)
			fmt.Fprintln(os.Stderr, "请使用 --env-file 参数指定环境文件路径，或确保在.http文件所在目录或上级目录有环境文件")
			os.Exit(exitFailure)
		}
	}

	// 将环境变量添加到HTTP文件
	if envVars != nil {
		if httpFile.EnvironmentVars == nil {
			httpFile.EnvironmentVars = make(map[string]map[string]string)
		}
		httpFile.EnvironmentVars[opts.Env] = envVars
	}

	// 创建执行器
	exec := executor.NewExecutor(opts.Verbose)

	// 执行HTTP请求
	responses, err := exec.ExecuteFile(httpFile, opts.RequestName, opts.Env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "执行请求错误: %v\n", err)
		os.Exit(exitFailure)
	}

	// 如果没有启用详细模式，但指定了单个请求，打印响应结果
	if !opts.Verbose && opts.RequestName != "" && len(responses) > 0 {
		executor.PrintResponse(responses[0])
	} else if !opts.Verbose {
		fmt.Printf("成功执行 %d 个HTTP请求\n", len(responses))
		for i, resp := range responses {
			fmt.Printf("\n请求 #%d: %s\n", i+1, resp.Request.Name)
			fmt.Printf("状态: %s\n", resp.Status)
			fmt.Printf("耗时: %d ms\n", resp.Time)
		}
	}

	// 如果指定了输出文件，将响应保存到文件
	if opts.OutputFile != "" && len(responses) > 0 {
		resp := responses[0]
		if err := os.WriteFile(opts.OutputFile, resp.Body, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "写入输出文件错误: %v\n", err)
			os.Exit(exitFailure)
		}
		fmt.Printf("响应已保存到文件: %s\n", opts.OutputFile)
	}

	os.Exit(exitSuccess)
}
