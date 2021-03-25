package log

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
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

// 因为传递时，是只能string
type Fields map[string]string

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

// 外部可调用Logger
var Logger Log
var colors = []color.Color{color.Green, color.Red}
var levelString = []string{"Debug", "Info", "Warn", "Error", "Fatal", "Panic"}

func (l Level) String() string {
	return levelString[l]
}

type Log struct {
	Writer  io.Writer // 输出
	TimeNow time.Time // 时间
	Skip    int       // caller的深度
	Data    Fields    //  字段

	initialized   bool // 该日志对象是否初始化
	isShowCaller  bool // 是否查找路径
	isOutToCenter bool //
	//onlyStdOut    bool // 只在终端输出
	color   color.Color
	level   Level
	cxt     context.Context
	mu      sync.Mutex
	content string
}

// 构造函数
func DefaultNewLogger() *Log {
	return &Log{
		Writer:       os.Stdout, // 默认终端输出
		Skip:         4,
		mu:           sync.Mutex{},
		isShowCaller: false,
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
		initialized:  true,
		isShowCaller: false,
	}
}

// 添加上下文（虽然目前没有什么用）
func (l *Log) SetContext(ctx context.Context) *Log {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cxt = ctx
	return l
}

// 添加字段
func (l *Log) SetFields(fields Fields) *Log {
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
func (l *Log) SetWriter(w io.Writer) *Log {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Writer = w
	return l
}

// 获取当前时间
func (l *Log) SetNowTime() *Log {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.TimeNow = time.Now()
	return l
}

func (l *Log) SetCaller() *Log {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.isShowCaller {
		return l
	}
	l.isShowCaller = true
	return l

}
func (l *Log) OutCaller() *Log {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.isShowCaller {
		return l
	}
	l.isShowCaller = false
	return l

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

func (l *Log) GetLevel() Level {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

// Json序列化
func (l *Log) JsonFormat(msg ...interface{}) map[string]string {
	l.mu.Lock()
	defer l.mu.Unlock()
	data := make(Fields, len(l.Data)+3)
	data["msg"] = fmt.Sprintf("%s", msg...)
	data["time"] = l.TimeNow.String()
	if len(l.Data) > 0 {
		for k, v := range l.Data {
			if _, ok := data[k]; !ok {
				data[k] = v
			}
		}
	}
	return data
}

func (l *Log) Output(level Level, msg string) error {
	l.setLevel(level)
	if len(msg) != 0 {
		l.SetFields(Fields{"msg": msg})
	}
	l.formatter()
	// 打印到控制台
	_, err := l.Writer.Write(qconv.Qs2b(l.content))
	if err != nil {
		return err
	}
	// 不止是终端输出
	return nil
	//if len(l.FilePath) == 0 {
	//	return "", nil
	//}
	// 获取到要输出的文件路径
	//file := l.GetFile()
	//defer file.Close()
	//n, _ := file.Seek(0, os.SEEK_END)
	//// 写入文件
	//return file.WriteAt(qconv.Qs2b(l.content), n)

}
func (l *Log) OutPutAndContent() string {
	l.Output(l.level, "")
	return l.content
}

//
//// 默认为当前位置的begonia.log里
//func (l *Log) GetFile() *os.File {
//	if len(l.FilePath) == 0 {
//		return nil
//	}
//	//fmt.Println(path)
//	if _, err := os.Stat(l.FilePath); err != nil {
//		// 文件不存在,创建
//		f, err := os.Create(l.FilePath)
//		//defer f.Close()  // 记得关闭
//		if err != nil {
//			panic(err)
//		}
//		return f
//	}
//	// 打开该文件，追加模式
//	f, err := os.OpenFile(l.FilePath, os.O_WRONLY, os.ModeAppend)
//
//	if err != nil {
//		panic(err)
//	}
//
//	return f
//}
func (l *Log) Print(level Level, v ...interface{}) {
	l.Output(level, fmt.Sprint(v...))
}

func (l *Log) Debug(v ...interface{}) {
	l.Output(LevelDebug, fmt.Sprint(v...))
}
func (l *Log) Info(v ...interface{}) {
	l.Output(LevelInfo, fmt.Sprint(v...))
}
func (l *Log) Warn(v ...interface{}) {
	l.Output(LevelWarn, fmt.Sprint(v...))
}
func (l *Log) Error(v ...interface{}) {
	l.Output(LevelError, fmt.Sprint(v...))
}
func (l *Log) Fatal(v ...interface{}) {
	l.Output(LevelFatal, fmt.Sprint(v...))
}
func (l *Log) Panic(v ...interface{}) {
	l.Output(LevelPanic, fmt.Sprint(v...))
}

// 添加level
func (l *Log) setLevel(lev Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = lev
	return
}

// 解析fields
func (l *Log) parsedFields() string {
	buf := strings.Builder{}
	for k, v := range l.Data {
		field := fmt.Sprintf(" [%s=%s]", k, v)
		buf.WriteString(field)
	}
	return buf.String()
}

// 格式化logger
func (l *Log) formatter() *Log {
	var title string
	if l.level > LevelInfo {
		title = colors[1].Sprintf("%s", l.level.String())
	} else {
		title = colors[0].Sprintf("%s", l.level.String())
	}

	// [Error] [2006-01-02 - 15:04:05] test
	buf := bytes.NewBufferString(fmt.Sprintf("[%s] [%s] %s",
		l.SetNowTime().TimeNow.Format("2006-01-02 - 15:04:05"), title, l.parsedFields()))

	// 如果要求显示路径
	if l.isShowCaller || l.level > LevelInfo {
		for _, v := range l.GetCaller() {
			buf.WriteByte('\n')
			buf.WriteString(v)
		}
	}
	l.mu.Lock()
	l.content = fmt.Sprintln(buf.String())
	l.mu.Unlock()
	return l
}
