# JHTTP 待办事项

## 功能增强

### ✅ 1. 自动识别环境变量文件（已完成）

- ✅ 自动识别env.json文件名和位置，实现以下查找策略：
  - ✅ 查找名称为`http-client.private.env.json`或`http-client.env.json`的文件
  - ✅ 默认在.http文件所在目录查找
  - ✅ 如果没有找到，逐级向上层目录查找，直到找到为止
  - ✅ 当自动找到环境文件时，在命令行输出警告信息，提示用户已自动使用了该文件
  
- ✅ 实现情况：
  - 修改了`environment/env.go`，添加了`FindEnvFile()`函数
  - 更新了`cmd/jhttp/main.go`，加入了自动查找环境文件的逻辑
  - 改进了`cli/cli.go`中的帮助文档，说明支持自动查找环境文件
  - 只有在上级目录找到环境文件时才显示警告信息
  - 支持多种环境文件名格式
  - 实现了完整的向上查找逻辑
  - 添加了测试用例验证功能正常

### 2. 增加并行测试和测试数量支持

- 增加选项，允许多个请求并行执行，加快测试速度
- 添加压测功能，可以指定请求的重复次数和并发数
- 收集和统计执行结果，包括成功率、平均响应时间、最大响应时间等

- 命令行参数扩展：
  ```
  --parallel         启用并行模式执行请求
  --concurrent <n>   设置并发数量
  --repeat <n>       设置每个请求的重复次数
  --stats            显示统计信息
  ```

- 具体实现方向：
  ```go
  // 在executor包中添加
  func (e *Executor) ExecuteParallel(httpFile *models.HTTPFile, options ParallelOptions) ([]*models.HTTPResponse, *models.Statistics, error) {
    // 使用goroutines并发执行请求
    // 收集结果和统计信息
  }
  ```

- 优先级：中
- 预计工作量：高

### 3. 美化JSON输出

- 对JSON响应体进行格式化，提高可读性
- 支持彩色输出，区分不同的JSON元素类型
- 提供选项控制是否格式化和彩色显示

- 命令行参数扩展：
  ```
  --pretty           美化输出（默认开启）
  --no-pretty        禁用美化输出
  --color            彩色输出
  --no-color         禁用彩色输出
  ```

- 具体实现方向：
  ```go
  // 在执行器或格式化工具中添加
  func FormatJSON(input string, pretty bool, color bool) string {
    // 解析JSON
    // 根据选项格式化
    // 添加缩进和颜色
    // 返回格式化后的字符串
  }
  ```

- 优先级：中
- 预计工作量：中

## 其他改进

### 代码重构和优化

- 改进错误处理，提供更详细的错误信息和故障排除建议
- 增加日志系统，可配置不同级别的日志输出
- 优化内存使用，特别是处理大型请求和响应时
- 增加单元测试和集成测试，提高代码质量和可靠性

### 文档完善

- 添加API文档，方便其他开发者使用项目作为库
- 完善示例和用例，展示常见使用场景
- 添加常见问题解答(FAQ)部分
- 提供性能优化指南

### 用户体验提升

- 添加进度显示和动画，特别是在处理多个请求时
- 增加交互模式，允许用户在运行时选择要执行的请求
- 添加保存历史记录功能，记录之前执行的请求和结果
- 支持请求模板和代码片段 