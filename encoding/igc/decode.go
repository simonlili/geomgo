// Package igc 实现 an IGC 分析器.
package igc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/chengxiaoer/geomGo"
)

var (
	// ErrInvalidCharacter 将返回当遇到无效字符时.
	ErrInvalidCharacter = errors.New("invalid character")
	// ErrInvalidCharactersBeforeARecord 将返回在 A record之前遇到无效字符时。
	ErrInvalidCharactersBeforeARecord = errors.New("invalid characters before A record")
	// ErrInvalidBRecord 遇到无效的B记录时将返回.
	ErrInvalidBRecord = errors.New("invalid B record")
	// ErrInvalidHRecord 遇到无效的H记录时将返回.
	ErrInvalidHRecord = errors.New("invalid H record")
	// ErrInvalidIRecord 遇到无效的I记录时将返回.
	ErrInvalidIRecord = errors.New("invalid I record")
	// ErrEmptyLine 当遇到无效的空 line 时将返回.
	ErrEmptyLine = errors.New("empty line")
	// ErrMissingARecord 当找不到 A record 时将返回.
	ErrMissingARecord = errors.New("missing A record")
	// ErrOutOfRange 当值超过范围时将返回
	ErrOutOfRange = errors.New("out of range")

	hRegexp = regexp.MustCompile(`H([FP])([A-Z]{3})(.*?):(.*?)\s*\z`)
)

// An Errors is 当遍历每行出错时.
type Errors map[int]error

// A Header is an IGC header.
type Header struct {
	Source   string
	Key      string
	KeyExtra string
	Value    string
}

// A T 代表IGC文件解析.
type T struct {
	Headers    []Header
	LineString *geom.LineString
}

func (es Errors) Error() string {
	var ss []string
	for lineno, e := range es {
		ss = append(ss, fmt.Sprintf("%d: %s", lineno, e.Error()))
	}
	return strings.Join(ss, "\n")
}

// parseDec 解析十进制值在s[start:stop]中.
func parseDec(s string, start, stop int) (int, error) {
	result := 0
	neg := false
	if s[start] == '-' {
		neg = true
		start++
	}
	for i := start; i < stop; i++ {
		if c := s[i]; '0' <= c && c <= '9' {
			result = 10*result + int(c) - '0'
		} else {
			return 0, ErrInvalidCharacter
		}
	}
	if neg {
		result = -result
	}
	return result, nil
}

// parseDecInRange函数 解析十进制值 s[start:stop], 如果超出[min,max)时返回错误.
func parseDecInRange(s string, start, stop, min, max int) (int, error) {
	if result, err := parseDec(s, start, stop); err != nil {
		return result, err
	} else if result < min || max <= result {
		return result, ErrOutOfRange
	} else {
		return result, nil
	}
}

// parser  包含这个解析器的状态
type parser struct {
	headers           []Header
	coords            []float64
	year, month, day  int
	startAt           time.Time
	lastDate          time.Time
	ladStart, ladStop int
	lodStart, lodStop int
	tdsStart, tdsStop int
	bRecordLen        int
}

// newParser函数 创建一个新的解析器.
func newParser() *parser {
	return &parser{bRecordLen: 35}
}

// parseB方法 从 line 中解析一个 B record,同时更新 解析器的状态.
func (p *parser) parseB(line string) error {

	if len(line) != p.bRecordLen {
		return ErrInvalidBRecord
	}

	var err error

	var hour, minute, second, nsec int
	if hour, err = parseDecInRange(line, 1, 3, 0, 24); err != nil {
		return err
	}
	if minute, err = parseDecInRange(line, 3, 5, 0, 60); err != nil {
		return err
	}
	if second, err = parseDecInRange(line, 5, 7, 0, 60); err != nil {
		return err
	}
	if p.tdsStart != 0 {
		var decisecond int
		decisecond, err = parseDecInRange(line, p.tdsStart, p.tdsStop, 0, 10)
		if err != nil {
			return err
		}
		nsec = decisecond * 1e8
	}
	date := time.Date(p.year, time.Month(p.month), p.day, hour, minute, second, nsec, time.UTC)
	if date.Before(p.lastDate) {
		p.day++
		date = time.Date(p.year, time.Month(p.month), p.day, hour, minute, second, nsec, time.UTC)
	}

	if p.startAt.IsZero() {
		p.startAt = date
	}

	var latDeg, latMilliMin int
	if latDeg, err = parseDecInRange(line, 7, 9, 0, 90); err != nil {
		return err
	}
	// 特殊情况: latMilliMin 应当在[0, 60000)的范围内， but a number of flight recorders generate latMilliMins of 60000
	// FIXME check what happens in negative (S, W) hemispheres(半球)
	if latMilliMin, err = parseDecInRange(line, 9, 14, 0, 60000+1); err != nil {
		return err
	}
	lat := float64(60000*latDeg+latMilliMin) / 60000.
	if p.ladStart != 0 {
		var lad int
		if lad, err = parseDec(line, p.ladStart, p.ladStop); err == nil {
			lat += float64(lad) / 6000000.
		} else {
			return err
		}
	}
	switch c := line[14]; c {
	case 'N':
	case 'S':
		lat = -lat
	default:
		return ErrInvalidCharacter
	}

	var lngDeg, lngMilliMin int
	if lngDeg, err = parseDecInRange(line, 15, 18, 0, 180); err != nil {
		return err
	}
	if lngMilliMin, err = parseDecInRange(line, 18, 23, 0, 60000+1); err != nil {
		return err
	}
	lng := float64(60000*lngDeg+lngMilliMin) / 60000.
	if p.lodStart != 0 {
		var lod int
		if lod, err = parseDec(line, p.lodStart, p.lodStop); err == nil {
			lng += float64(lod) / 6000000.
		} else {
			return err
		}
	}
	switch c := line[23]; c {
	case 'E':
	case 'W':
		lng = -lng
	default:
		return ErrInvalidCharacter
	}

	var pressureAlt, ellipsoidAlt int
	if pressureAlt, err = parseDec(line, 25, 30); err != nil {
		return err
	}
	if ellipsoidAlt, err = parseDec(line, 30, 35); err != nil {
		return err
	}

	p.coords = append(p.coords, lng, lat, float64(ellipsoidAlt), float64(date.UnixNano())/1e9, float64(pressureAlt))
	p.lastDate = date

	return nil

}

// parseH方法 从line 中解析一个 H record，并更新 解析器p 的状态.
func (p *parser) parseH(line string) error {
	if m := hRegexp.FindStringSubmatch(line); m != nil {
		p.headers = append(p.headers, Header{
			Source:   m[1],
			Key:      m[2],
			KeyExtra: m[3],
			Value:    m[4],
		})
	}
	switch {
	case strings.HasPrefix(line, "HFDTE"):
		return p.parseHFDTE(line)
	default:
		return nil
	}
}

// parseHFDTE 从line 中解析一个 HFDTE record，并更新 解析器p 的状态.
func (p *parser) parseHFDTE(line string) error {
	var err error
	var day, month, year int
	if len(line) != 11 {
		return ErrInvalidHRecord
	}
	if day, err = parseDecInRange(line, 5, 7, 1, 31+1); err != nil {
		return err
	}
	if month, err = parseDecInRange(line, 7, 9, 1, 12+1); err != nil {
		return err
	}
	if year, err = parseDec(line, 9, 11); err != nil {
		return err
	}
	// FIXME 检查无效数据
	p.day = day
	p.month = month
	if year < 70 {
		p.year = 2000 + year
	} else {
		p.year = 1970 + year
	}
	return nil
}

// parseI 从line 中解析一个 I record，并更新 解析器p 的状态..
func (p *parser) parseI(line string) error {
	var err error
	var n int
	if len(line) < 3 {
		return ErrInvalidIRecord
	}
	if n, err = parseDec(line, 1, 3); err != nil {
		return err
	}
	if len(line) < 7*n+3 {
		return ErrInvalidIRecord
	}
	for i := 0; i < n; i++ {
		var start, stop int
		if start, err = parseDec(line, 7*i+3, 7*i+5); err != nil {
			return err
		}
		if stop, err = parseDec(line, 7*i+5, 7*i+7); err != nil {
			return err
		}
		if start != p.bRecordLen+1 || stop < start {
			return ErrInvalidIRecord
		}
		p.bRecordLen = stop
		switch line[7*i+7 : 7*i+10] {
		case "LAD":
			p.ladStart, p.ladStop = start-1, stop
		case "LOD":
			p.lodStart, p.lodStop = start-1, stop
		case "TDS":
			p.tdsStart, p.tdsStop = start-1, stop
		}
	}
	return nil
}

// parseLine 从 line 中解析 single reord,同时更新解析器 p 的状态.
func (p *parser) parseLine(line string) error {
	switch line[0] {
	case 'B':
		return p.parseB(line)
	case 'H':
		return p.parseH(line)
	case 'I':
		return p.parseI(line)
	default:
		return nil
	}
}

// doParse函数 读取 r, 解析它所能找到的所有记录, 并更新解析器 p 的状态.
func doParse(r io.Reader) (*parser, Errors) {
	errors := make(Errors)
	p := newParser()
	s := bufio.NewScanner(r)
	foundA := false
	leadingNoise := false
	for lineno := 1; s.Scan(); lineno++ {
		line := s.Text()
		if len(line) == 0 {
			// errors[lineno] = ErrEmptyLine
		} else if foundA {
			if err := p.parseLine(line); err != nil {
				errors[lineno] = err
			}
		} else {
			if c := line[0]; c == 'A' {
				foundA = true
			} else if 'A' <= c && c <= 'Z' {
				// 所有记录必须以大写字符开头，才是有效的。.
				leadingNoise = true
				continue
			} else if i := strings.IndexRune(line, 'A'); i != -1 {
				// 剥去任何噪音。
				// 噪声必须包含至少一个不可打印字符 (like XOFF or a fragment of a Unicode BOM).
				for _, c := range line[:i] {
					if !(c == ' ' || ('A' <= c && c <= 'Z')) {
						foundA = true
						leadingNoise = true
						break
					}
				}
			}
		}
	}
	if !foundA {
		errors[1] = ErrMissingARecord
	} else if leadingNoise {
		errors[1] = ErrInvalidCharactersBeforeARecord
	}
	return p, errors
}

// Read 读取 a igc.T from r, 其中应包含IGC的记录.
func Read(r io.Reader) (*T, error) {
	p, errors := doParse(r)
	if len(errors) != 0 {
		return nil, errors
	}
	return &T{
		Headers:    p.headers,
		LineString: geom.NewLineStringFlat(geom.Layout(5), p.coords),
	}, nil
}
