package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/MashiroC/begonia/tool/qconv"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
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

// 颜色
const (
	green   = "\033[32m"
	white   = "\033[37m"
	yellow  = "\033[33m"
	red     = "\033[31m"
	blue    = "\033[34m"
	magenta = "\033[35m"
)

type Logger struct {
	Writer   io.Writer // 输出
	TimeNow  time.Time // 时间
	Message  string    // 信息
	Data     Fields
	FilePath string // 文件输出路径

	color   string
	level   Level
	cxt     context.Context
	mu      sync.Mutex
	caller  []string //调用栈的信息
	content string
}

// 构造函数
func NewLogger() *Logger {
	l := Logger{
		Writer: os.Stdout, // 默认终端输出
		mu:     sync.Mutex{},
	}
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	l.FilePath = dir + "/begonia.log"
	return &l
}

// 格式化logger
func (l *Logger) Formatter() *Logger {
	var title string
	switch l.level {
	case LevelDebug:
		l.color = blue
		title = "Debug"
	case LevelInfo:
		l.color = green
		title = "Info"
	case LevelWarn:
		l.color = magenta
		title = "Warn"
	case LevelError:
		l.color = red
		title = "Error"
	case LevelFatal:
		l.color = red
		title = "Fatal"
	case LevelPanic:
		l.color = red
		title = "Panic"
	}
	// [Error] [2006-01-02 - 15:04:05] test
	buf := bytes.NewBufferString(fmt.Sprintf("%s[%s]\033[0m [%s] %s",
		l.color, title, time.Now().Format("2006-01-02 - 15:04:05"), l.Message))

	// 如果有字段
	if l.Data != nil {
		buf.WriteString(" | ")
		format, _ := json.Marshal(l.Data)
		buf.Write(format)
	}
	// 如果要求显示路径
	if len(l.caller) != 0 {
		for _, v := range l.caller {
			buf.WriteByte('\n')
			buf.WriteString(v)
		}
	}
	l.mu.Lock()
	l.content = fmt.Sprintln(buf.String())
	l.mu.Unlock()
	return l
}

// 添加level
func (l *Logger) WithLevel(lev Level) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = lev
	return l
}

// 添加上下文（虽然目前没有什么用）
func (l *Logger) WithContext(ctx context.Context) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cxt = ctx
	return l
}

// 添加字段
func (l *Logger) WithFields(fields Fields) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.Data == nil {
		l.Data = make(Fields, len(fields))
	}
	for k, v := range fields {
		l.Data[k] = v
	}
	return l
}

// 输出
func (l *Logger) WithOut(w io.Writer) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Writer = w
	return l
}

// 添加信息
func (l *Logger) WithMsg(msg ...interface{}) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Message = fmt.Sprint(l.Message, msg)
	return l
}

// 获取当前时间
func (l *Logger) GetNowTime() *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.TimeNow = time.Now()
	return l
}

// 调用栈 skip > 1
func (l *Logger) WithCaller(skip int) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	//pc -> 函数指针  file -> 绝对路径  line -> 对应文件的行数  ok -> 是否获取成功
	pc, file, line, ok := runtime.Caller(skip - 1)
	if ok {
		fun := runtime.FuncForPC(pc)

		l.caller = []string{fmt.Sprintf("%s:  %d %s", file, line, fun.Name())}
	}
	return l
}

// 获取当前整个栈的信息
func (l *Logger) WithCallerFrames() *Logger {
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
	l.caller = callers
	return l
}

// Json序列化
func (l *Logger) JsonFormat() map[string]interface{} {
	l.mu.Lock()
	defer l.mu.Unlock()
	data := make(Fields, len(l.Data)+3)
	data["msg"] = l.Message
	data["callers"] = l.caller
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

func (l *Logger) Output(msg string) {
	l.WithMsg(msg).Formatter()
	// 获取到要输出的文件路径
	file := l.GetFile()
	defer file.Close()
	n, _ := file.Seek(0, os.SEEK_END)
	// 写入文件
	file.WriteAt(qconv.Qs2b(l.content), n)
	// 打印到控制台
	l.Writer.Write(qconv.Qs2b(l.content))
}

// 默认为当前位置的begonia.log里
func (l *Logger) GetFile() *os.File {

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
func (l *Logger) Print(v ...interface{}) {
	l.WithLevel(LevelInfo).Output(fmt.Sprint(v...))
}

func (l *Logger) Debug(v ...interface{}) {
	l.WithLevel(LevelDebug).Output(fmt.Sprint(v...))
}
func (l *Logger) Info(v ...interface{}) {
	l.WithLevel(LevelInfo).Output(fmt.Sprint(v...))
}
func (l *Logger) Warn(v ...interface{}) {
	l.WithLevel(LevelWarn).Output(fmt.Sprint(v...))
}
func (l *Logger) Error(v ...interface{}) {
	l.WithLevel(LevelError).Output(fmt.Sprint(v...))
}
func (l *Logger) Fatal(v ...interface{}) {
	l.WithLevel(LevelFatal).Output(fmt.Sprint(v...))
}
func (l *Logger) Panic(v ...interface{}) {
	l.WithLevel(LevelPanic).Output(fmt.Sprint(v...))
}
