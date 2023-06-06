package logging

import (
	"fmt"
	"io"
	"regexp"
	"sync"
)

// formatRegexp 匹配日志格式的正则表达式。这里用到了非捕获匹配规则 `?:`，非捕获匹配规则的简单解释如下：
// 给定一个待匹配字符串："%{color:reset}%{level:info}"，我们利用 formatRegexp 来匹配该字符串，得到如
// 下结果：
// [0 13] [2 6] [8 12] [14 26] [16 20] [22 25]，每个中括号里的数字分别对应匹配到的不同字符串的左右边界；
// 如果去掉正则表达式里的 `?:`，改用捕获匹配规则，得到的匹配结果会是下面这样的：
// [0 13] [2 6] [7 12] [8 12] [14 26] [16 20] [21 25] [22 25]，会发现多了 [7 12] 和 [21 25]，也就是
// 说，捕获匹配也将字符串的子串 ":reset" 和 ":info" 也单独匹配出来了。
var formatRegexp = regexp.MustCompile(`%{(color|level|time|module|location|message)(?::(.*?))?}`)

var defaultFormat = "%{color}%{level}[%{time}] [%{module}]%{color:reset} => %{message}"

type Formatter interface {
	Format(w io.Writer, e Entry)
}

func ParseFormat(spec string) ([]Formatter, error) {
	cursor := 0
	formatters := make([]Formatter, 0)

	matches := formatRegexp.FindAllStringSubmatchIndex(spec, -1)
	for _, match := range matches {
		start, end := match[0], match[1]
		verbStart, verbEnd := match[2], match[3]
		formatStart, formatEnd := match[4], match[5]

		if start > cursor {
			formatters = append(formatters, StringFormatter{value: spec[cursor:start]})
		}

		verb := spec[verbStart:verbEnd]
		var format string // 与 color 和 time 相关。
		if formatStart > 0 {
			format = spec[formatStart:formatEnd]
		}

		formatter, err := NewFormatter(verb, format)
		if err != nil {
			return nil, err
		}

		formatters = append(formatters, formatter)
		cursor = end
	}

	if cursor != len(spec) {
		formatters = append(formatters, StringFormatter{value: spec[cursor:]})
	}

	return formatters, nil
}

type MultiFormatter struct {
	mutex      sync.RWMutex
	formatters []Formatter
}

func NewMultiFormatter(formatters ...Formatter) *MultiFormatter {
	return &MultiFormatter{formatters: formatters}
}

func (m *MultiFormatter) Format(w io.Writer, e Entry) {
	m.mutex.RLock()

	switch e.FormatSelector {
	case "json":
		fmt.Fprint(w, "{")
		for i, formatter := range m.formatters {
			if _, ok := formatter.(StringFormatter); ok {
				continue
			}
			formatter.Format(w, e)
			if i != len(m.formatters)-1 {
				if _, ok := formatter.(ColorFormatter); !ok {
					fmt.Fprint(w, ",")
				}
			}
		}
		fmt.Fprint(w, "}")
	case "terminal", "":
		for _, formatter := range m.formatters {
			formatter.Format(w, e)
		}
	}

	fmt.Fprint(w, "\n")

	m.mutex.RUnlock()
}

func (m *MultiFormatter) SetFormatters(formatters []Formatter) {
	m.mutex.Lock()
	m.formatters = formatters
	m.mutex.Unlock()
}

func (m *MultiFormatter) clone() *MultiFormatter {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	cpy := &MultiFormatter{
		formatters: make([]Formatter, len(m.formatters)),
	}
	copy(cpy.formatters, m.formatters)
	return cpy
}

type StringFormatter struct {
	value string
}

func (sf StringFormatter) Format(w io.Writer, e Entry) {
	fmt.Fprint(w, sf.value)
}

func NewFormatter(verb, format string) (Formatter, error) {
	switch verb {
	case "color":
		return newColorFormatter(format)
	case "level":
		return LevelFormatter{}, nil
	case "time":
		return TimeFormatter{layout: format}, nil
	case "module":
		return ModuleFormatter{}, nil
	case "location":
		return LocationFormatter{}, nil
	case "message":
		return MessageFormatter{}, nil
	default:
		return nil, fmt.Errorf("unknown verb: (%s)", verb)
	}
}

type ColorFormatter struct {
	reset bool
}

func newColorFormatter(f string) (ColorFormatter, error) {
	switch f {
	case "reset":
		return ColorFormatter{reset: true}, nil
	case "":
		return ColorFormatter{}, nil
	default:
		return ColorFormatter{}, fmt.Errorf("invalid color option: %s", f)
	}
}

func (cf ColorFormatter) Format(w io.Writer, e Entry) {

	switch {
	case cf.reset:
		fmt.Fprint(w, ResetColor())
	default:
		fmt.Fprint(w, e.Level.SpecifiedColor().Color())
	}
}

type LevelFormatter struct{}

func (lf LevelFormatter) Format(w io.Writer, e Entry) {
	switch e.FormatSelector {
	case "terminal", "":
		fmt.Fprint(w, e.Level.CapitalString())
	case "json":
		fmt.Fprintf(w, "\"level\":\"%s\"", e.Level.CapitalString())
	}
}

type TimeFormatter struct {
	layout string
}

func (tf TimeFormatter) Format(w io.Writer, e Entry) {
	var t string
	if tf.layout != "" {
		t = e.Time.Format(tf.layout)
	} else {
		t = e.Time.Format("2006-01-02 15:04:05.000 CST")
	}
	switch e.FormatSelector {
	case "terminal", "":
		fmt.Fprint(w, t)
	case "json":
		fmt.Fprintf(w, "\"time\":\"%s\"", t)
	}
}

type ModuleFormatter struct{}

func (ModuleFormatter) Format(w io.Writer, e Entry) {
	switch e.FormatSelector {
	case "terminal", "":
		fmt.Fprint(w, e.Module)
	case "json":
		fmt.Fprintf(w, "\"module\":\"%s\"", e.Module)
	}
}

type LocationFormatter struct{}

func (LocationFormatter) Format(w io.Writer, e Entry) {
	switch e.FormatSelector {
	case "terminal", "":
		fmt.Fprint(w, e.Call.String())
	case "json":
		fmt.Fprintf(w, "\"location\":\"%s\"", e.Call.String())
	}
}

type MessageFormatter struct{}

func (MessageFormatter) Format(w io.Writer, e Entry) {
	switch e.FormatSelector {
	case "terminal", "":
		msg := fmt.Sprintf("%s %s", e.Message, normalize(e.KeyValues...))
		fmt.Fprint(w, msg)
	case "json":
		fmt.Fprintf(w, "\"message\":\"%s\"", e.Message)
		if len(e.KeyValues) > 0 {
			fmt.Fprint(w, ", ")
			if len(e.KeyValues)%2 != 0 {
				e.KeyValues = append(e.KeyValues, "ERROR_LOG")
			}
			for i := 0; i < len(e.KeyValues)-1; i += 2 {
				if i < len(e.KeyValues)-2 {
					fmt.Fprintf(w, "\"%s\":\"%s\", ", toString(e.KeyValues[i]), toString(e.KeyValues[i+1]))
				} else {
					fmt.Fprintf(w, "\"%s\":\"%s\"", toString(e.KeyValues[i]), toString(e.KeyValues[i+1]))
				}
			}
		}
	}
}
