# 日志记录器

## 概念介绍

### spec

`spec` 定义了日志的输出格式，默认的日志输出格式如下所示：

```
spec = "%{color}%{level}[%{time}] [%{module}]%{color:reset} => %{message}"
```
【颜色】【日志级别】【时间】【日志模块】【颜色重置】【日志内容】

以默认的日志格式打印日志信息会得到如下图所示的日志内容结构：
<img src="pics/1.PNG">

可以发现，这样的日志结构是按照【日志级别】`=>`【时间】`=>`【日志模块】`=>`【日志内容】顺序输出的，我们在默认的 `spec` 中还发现了两个特殊的标签：`%{color}` 和 `%{color:reset}`，这两个标签一般情况下是成对出现的，`%{color}` 标签后面的内容都会以彩色的样式被输出，而 `%{color:reset}` 后面的内容则恢复成终端默认的颜色形式被输出。这与上面图片展示的结果是一致的，【日志级别】【时间】【日志模块】都在 `%{color}` 标签后，且都在 `%{color:reset}` 标签前，所以这几个内容都以彩色的形式被输出，而【日志内容】在 `%{color:reset}` 标签之后，所以这部分内容以默认颜色被输出。

所以，我们可以通过设计 `spec` 定义我们想要的日志格式，例如下面给出的一些样例：

- 禁止彩色输出：
    ```
    spec = "%{level}[%{time}] [%{module}] => %{message}"
    ```
    <img src="pics/2.PNG">
- 不显示日志级别，但是以彩色形式输出
    ```
    spec = "%{color}[%{time}] [%{module}]%{color:reset} => %{message}"
    ```
    <img src="pics/3.PNG">
- 不显示日志级别，但是以彩色形式打印日志，且显示打印日志代码的位置
    ```
    spec = "%{color}[%{time}] [%{module}] %{location}%{color:reset} => %{message}"
    ```
    <img src="pics/4.PNG">

    **所以说，设置 `%{location}` 标签可以在打印日志的时候输出产生日志信息的位置，这在调试代码时很有用。**
- 对默认 `spec` 进行修改，【日志内容】前的 `=>` 符号换成其他符号，比如冒号
    ```
    spec = "%{color}%{level}[%{time}] [%{module}]%{color:reset}: %{message}"
    ```
    <img src="pics/5.PNG">

### terminal 与 json

`terminal` 与 `json` 规定了打印日志内容的样式，`terminal` 规定日志是以一般终端的样式被输出的，而 `json` 规定日志是以 `JSON` 的样式被输出的，下面给出了两种样式的案例：

`terminal` 样式
<img src="pics/6.PNG">

`json` 样式
<img src="pics/7.PNG">

### module

`module` 表示日志模块，一个区块链项目由很多个模块组成，比如共识模块、加密模块、存储模块和通信模块等，为了让人们更加清楚地通过日志来分析每个模块的运行情况，所以在日志中引入了**日志模块**这一概念。

我们在前面介绍 `spec` 概念的时候，给出的例子里，日志模块的值等于 `blockchain`。

### filter-level

`filter-level` 表示过滤的日志级别，这个概念被提出的原因是：区块链项目代码在运转过程中，会产生大量的日志信息，可是有些日志信息是我们在调试代码时需要看到的，但是在真正部署运行的时候，这些日志信息是不需要的，怎么将这部分日志信息给过滤掉呢？

为了解决上述的问题，**日志级别**这一概念就被提了出来。日志级别按照严重程度逐渐升高的顺序可以被分为 `DEBUG`、`INFO`、`WARN`、`ERROR`、`PANIC` 五个级别。在调试代码时，我们可能会将所有的调试信息在 `DEBUG` 这一级别输出。那么在项目运行时，这些日志信息我们不想将它们输出来，这个时候就可以定义 `filter-level`，例如，我们将 `filter-level` 定义为 `INFO`，这样，所有 `DEBUG` 级别的日志都将不会输出，而只输出 `INFO`、`WARN`、`ERROR`、`PANIC` 这四个级别的日志信息。

结合前面介绍的 `module`，我们可以将共识模块的 `filter-level` 设置为 `INFO`，而将通信模块的 `filter-level` 设置为 `WARN`，这样给不同的日志模块设置不同的 `filter-level`，可以体现出我们所关心的侧重点是不同的。

## 使用教程

### 入门

通过下面的代码给出一个简单的使用例子：
```go
opt := logging.Option{
	Module:         "blockchain",
	FilterLevel:    logging.DebugLevel,
	Spec:           "%{color}%{level}[%{time}] [%{module}]%{color:reset}: %{message}",
	FormatSelector: "json",
	Writer:         os.Stdout,
}
logger, err := logging.NewLogger(opt)
if err != nil {
    panic(err)
}
logger.Info("info message", "key1", "value1", "key2", 99)
```
通过对上面代码的分析，可以看到，我们在实例化一个 `logger` 日志记录器时，需要先提供一个 `logging.Option`，这个“选项”实际上是为实例化日志记录器提供一些额外信息的。`Module` 表示日志模块的名称。`FilterLevel` 就是前面介绍的 `filter-level`，表示过滤的日志级别，这里 `FilterLevel` 被设置为 `DebugLevel`，就表示所有级别的日志都不会被过滤掉。`Spec` 就是前面介绍的 `spec`，用来定义日志格式。`FormatSelector` 表示日志被打印的样式，这里只能为 `FormatSelector` 赋值 `"json"` 或者 `"terminal"` 或者 `""`，赋其他的值都会报错。`Writer` 表示日志被打印的地方，这里设置成 `os.Stdout`，表示日志会被打印到终端，我们也可以自己创建一个日志文件，将日志打印到文件里。

上面一段代码执行后，会得到如下的效果：
<img src="pics/8.PNG">

### 添加日志模块

我们在区块链项目里，实例化了一个根日志记录器，我们想在这个日志记录器里添加新的日志模块，这样方便将来衍生出一个 `child` 日志记录器，方法如下代码所示：
```go
logger.SetModule("consensus", logging.InfoLevel)
logger.SetModule("p2p", logging.WarnLevel)
```
通过上面的代码，我们为根日志记录器添加了两个日志模块，分别是 `consensus` 和 `p2p`，而且，我们分别还为这两个日志模块设置了过滤级别：`INFO` 和 `WARN`。

### 衍生出一个 child 日志记录器

当我们实例化一个根日志记录器后，想要为其他模块生成一个日志记录器，可以通过以下代码来实现：
```go
logger.SetModule("consensus", logging.InfoLevel)
logger.SetModule("p2p", logging.WarnLevel)
consensusLogger := logger.DeriveChildLogger("consensus")

consensusLogger.Debug("debug message", "key", "value")
consensusLogger.Info("info message", "key", "value")
```
在上面的代码里，我们衍生出 `consensusLogger` 这个日志记录器，将来作为共识模块的日志记录器，由于我们之前在添加 `consensus` 模块时，指定了过滤的日志级别是 `INFO`，所以，可以预想的是，`Debug()` 方法将不起作用，也就是说，`consensusLogger` 无法输出 `DEBUG` 级别的日志。上述代码执行后的结果如下图所示：
<img src="pics/9.PNG">

### 更新日志记录器

我们从根日志记录器衍生出一个 `child` 日志记录器后，例如上一节的 `consensusLogger`，我们想禁止日志以彩色形式被输出，该怎么办呢？可以通过更新日志记录器的选项来实现，如下代码所示：
```go
logger.SetModule("consensus", logging.InfoLevel)
consensusLogger := logger.DeriveChildLogger("consensus")
consensusLogger.Update(logging.Option{
	Spec: "%{level}[%{time}] [%{module}]: %{message}",
})
consensusLogger.Info("info message", "key", "value")
```
执行上述代码，得到如下输出结果：
<img src="pics/10.PNG">

**从上面给的案例来看，我们可以通过重新定制 `logging.Option`，然后调用 `Update` 方法，就可以更新日志记录器的相关属性，实现重新定制日志记录器功能的目标。**