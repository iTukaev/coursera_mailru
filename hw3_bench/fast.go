package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
	"io"
	"os"
	"strconv"
	"strings"
)

type User struct {
	Browsers []string `json:"browsers"`
	Company string `json:"-"`
	Country string `json:"-"`
	Email string `json:"email"`
	Job string `json:"-"`
	Name string `json:"name"`
	Phone string `json:"-"`
}

var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	buf := bufio.NewScanner(file)
	//r := regexp.MustCompile("@")
	seenBrowsers := make(map[string]struct{}, 100)
	uniqueBrowsers := 0
	usersBuf := bytes.NewBufferString("")
	user := User{}
	i := -1
	email := ""

	for buf.Scan() {
		i++

		line := buf.Bytes()
		err = easyjson.Unmarshal(line, &user)
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		for _, browser := range user.Browsers {

			if ok := strings.Contains(browser, "Android"); ok {
				isAndroid = true
				_, ok = seenBrowsers[browser]
				if !ok {
					uniqueBrowsers++
					seenBrowsers[browser] = struct{}{}
				}
			}
			if ok := strings.Contains(browser, "MSIE"); ok {
				isMSIE = true
				_, ok = seenBrowsers[browser]
				if !ok {
					uniqueBrowsers++
					seenBrowsers[browser] = struct{}{}
				}

			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}
		email = strings.ReplaceAll(user.Email, "@", " [at] ")
		usersBuf.WriteRune('[')
		usersBuf.WriteString(strconv.Itoa(i))
		usersBuf.WriteRune(']')
		usersBuf.WriteRune(' ')
		usersBuf.WriteString(user.Name)
		usersBuf.WriteRune(' ')
		usersBuf.WriteRune('<')
		usersBuf.WriteString(email)
		usersBuf.WriteRune('>')
		usersBuf.WriteRune('\n')
	}

	fmt.Fprintln(out, "found users:\n" + usersBuf.String())
	fmt.Fprintln(out, "Total unique browsers", uniqueBrowsers) //len(seenBrowsers))
	usersBuf.Reset()
}

func easyjson9f2eff5fDecodeCourseraHw3BenchJsoner(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "email":
			out.Email = string(in.String())
		case "name":
			out.Name = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9f2eff5fEncodeCourseraHw3BenchJsoner(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"browsers\":"
		out.RawString(prefix[1:])
		if in.Browsers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Browsers {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9f2eff5fEncodeCourseraHw3BenchJsoner(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9f2eff5fEncodeCourseraHw3BenchJsoner(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9f2eff5fDecodeCourseraHw3BenchJsoner(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9f2eff5fDecodeCourseraHw3BenchJsoner(l, v)
}