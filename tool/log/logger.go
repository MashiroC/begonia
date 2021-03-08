package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/MashiroC/begonia/tool/qconv"
	"github.com/gookit/color"
)

// 定义一个接口
//type Connector interface {
//	Fatal(msg ...interface{})
//	Error(msg ...interface{})
//	Warn(msg ...interface{})
//	Debug(msg ...interface{})
//	Info(msg ...interface{})    // 打印
//	Output(msg string)              // 将信息输出到文件中
//	WithMsg(msg ...interface{}) // 将输入的信息添加抬头（例如添加打印时间等）
//	GetFile() *os.File              // 文件输出的配置
//}

type Level uint8

type Fields map[string]interface{}

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "Debug"
	case LevelInfo:
		return "Info"
	case LevelWarn:
		return "Warn"
	case LevelError:
		return "Error"
	case LevelFatal:
		return "Fatal"
	case LevelPanic:
		return "Panic"
	default:
		return ""
	}
}

// 外部可调用Logger
var Logger Log

type Log struct {
	Writer   io.Writer // 输出
	TimeNow  time.Time // 时间
	Message  string    // 信息
	Data     Fields
	FilePath string // 文件输出路径
	Skip     int    // caller的深度

	initialized   bool // 该日志对象是否初始化
	isShowCaller  bool // 是否查找路径
	isOutToCenter bool //
	onlyStdOut    bool // 只在终端输出
	color         color.Color
	level         Level
	cxt           context.Context
	mu            sync.Mutex
	content       string
}

// 构造函数
func DefaultNewLogger() *Log {
	return &Log{
		Writer:       os.Stdout, // 默认终端输出
		Skip:         4,
		mu:           sync.Mutex{},
		isShowCaller: true,
	}
}

// 初始化logger
func InitLogger() {
	if Logger.initialized {
		return
	}
	Logger = Log{
		Writer:       os.Stdout, // 默认终端输出
		Skip:         4,
		mu:           sync.Mutex{},
		isShowCaller: true,
	}
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// 默认为当前下的路径
	Logger.FilePath = dir + "/begonia.log"
	Logger.onlyStdOut = false
}

// 格式化logger
func (l *Log) Formatter() {
	var title string
	switch l.level {
	case LevelDebug:
		l.color = color.Green
		title = l.color.Sprintf("%s", "Debug")
	case LevelInfo:
		l.color = color.Green
		title = l.color.Sprintf("%s", "Info")
	case LevelWarn:
		l.color = color.Magenta
		title = l.color.Sprintf("%s", "Warn")
	case LevelError:
		l.color = color.Red
		title = l.color.Sprintf("%s", "Error")
	case LevelFatal:
		l.color = color.Red
		title = l.color.Sprintf("%s", "Fatal")
	case LevelPanic:
		l.color = color.Red
		title = l.color.Sprintf("%s", "Panic")
	}
	// [Error] [2006-01-02 - 15:04:05] test
	buf := bytes.NewBufferString(fmt.Sprintf("[%s] [%s] %s",
		title, time.Now().Format("2006-01-02 - 15:04:05"), l.Message))

	// 如果有字段
	if l.Data != nil {
		buf.WriteString(" | ")
		format, _ := json.Marshal(l.Data)
		buf.Write(format)
	}
	// 如果要求显示路径
	if l.isShowCaller || l.level > LevelInfo {
		for _, v := range l.GetCaller() {
			buf.WriteByte('\n')
			buf.WriteString(v)
		}
	}
	l.mu.Lock()
	l.content = fmt.Sprintln(buf.String())
	l.Message = ""
	l.mu.Unlock()
	return
}

// 添加level
func (l *Log) SetLevel(lev Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = lev
	return
}

// 添加上下文（虽然目前没有什么用）
func (l *Log) SetContext(ctx context.Context) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cxt = ctx
	return
}

// 添加字段
func (l *Log) SetFields(fields Fields) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.Data == nil {
		l.Data = make(Fields, len(fields))
	}
	for k, v := range fields {
		l.Data[k] = v
	}
	return
}

// 输出
func (l *Log) SetWriter(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Writer = w
	return
}

// 添加信息
func (l *Log) SetMsg(msg ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Message = fmt.Sprint(l.Message, msg)
	return
}

// 获取当前时间
func (l *Log) SetNowTime() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.TimeNow = time.Now()
	return
}
func (l *Log) SetOnlyStdOut() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.onlyStdOut{
		return
	}
	l.onlyStdOut = true
	return
}
func (l *Log) SetCaller() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.isShowCaller {
		return
	}
	l.isShowCaller = true
	return

}
func (l *Log) OutCaller() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.isShowCaller {
		return
	}
	l.isShowCaller = false
	return

}
func (l *Log) OutOnlyStdOut() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.onlyStdOut{
		return
	}
	l.onlyStdOut = false
	return
}

// 调用栈 skip > 1
//TODO:
func (l *Log) GetCaller() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	var callers []string
	// 判断是否已经设置
	//pc -> 函数指针  file -> 绝对路径  line -> 对应文件的行数  ok -> 是否获取成功
	pc, file, line, ok := runtime.Caller(l.Skip)
	if ok {
		fun := runtime.FuncForPC(pc)

		callers = []string{fmt.Sprintf("%s:  %d %s", file, line, fun.Name())}
	}
	return callers
}

// 获取当前整个栈的信息
func (l *Log) GetCallerFrames() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	// 最多15个（感觉已经很多）
	maxCallerDepth := 15
	minCallerDepth := 1
	callers := []string{}
	pcs := make([]uintptr, maxCallerDepth)
	depth := runtime.Callers(minCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	// 循环遍历获取栈
	frame, more := frames.Next()
	for more {
		callers = append(callers, fmt.Sprintf("%s:  %d %s",
			frame.File, frame.Line, frame.Function))
		frame, more = frames.Next()
	}
	return callers
}

// Json序列化
func (l *Log) JsonFormat() map[string]interface{} {
	l.mu.Lock()
	defer l.mu.Unlock()
	data := make(Fields, len(l.Data)+3)
	data["msg"] = l.Message
	data["time"] = l.TimeNow
	if len(l.Data) > 0 {
		for k, v := range l.Data {
			if _, ok := data[k]; !ok {
				data[k] = v
			}
		}
	}
	return data
}

func (l *Log) Output(msg string) (int, error) {
	l.SetMsg(msg)
	l.Formatter()
	// 打印到控制台
	l.Writer.Write(qconv.Qs2b(l.content))
	if len(l.FilePath) == 0 {
		return 0, nil
	}
	// 获取到要输出的文件路径
	file := l.GetFile()
	defer file.Close()
	n, _ := file.Seek(0, os.SEEK_END)
	// 写入文件
	return file.WriteAt(qconv.Qs2b(l.content), n)

}

// 默认为当前位置的begonia.log里
func (l *Log) GetFile() *os.File {
	if len(l.FilePath) == 0 {
		return nil
	}
	//fmt.Println(path)
	if _, err := os.Stat(l.FilePath); err != nil {
		// 文件不存在,创建
		f, err := os.Create(l.FilePath)
		//defer f.Close()  // 记得关闭
		if err != nil {
			panic(err)
		}
		return f
	}
	// 打开该文件，追加模式
	f, err := os.OpenFile(l.FilePath, os.O_WRONLY, os.ModeAppend)

	if err != nil {
		panic(err)
	}

	return f
}
func (l *Log) Print(v ...interface{}) {
	l.SetLevel(LevelInfo)
	l.Output(fmt.Sprint(v...))
}

func (l *Log) Debug(v ...interface{}) {
	l.SetLevel(LevelDebug)
	l.Output(fmt.Sprint(v...))
}
func (l *Log) Info(v ...interface{}) {
	l.SetLevel(LevelInfo)
	l.Output(fmt.Sprint(v...))
}
func (l *Log) Warn(v ...interface{}) {
	l.SetLevel(LevelWarn)
	l.Output(fmt.Sprint(v...))
}
func (l *Log) Error(v ...interface{}) {
	l.SetLevel(LevelError)
	l.Output(fmt.Sprint(v...))
}
func (l *Log) Fatal(v ...interface{}) {
	l.SetLevel(LevelFatal)
	l.Output(fmt.Sprint(v...))
}
func (l *Log) Panic(v ...interface{}) {
	l.SetLevel(LevelPanic)
	l.Output(fmt.Sprint(v...))
}
