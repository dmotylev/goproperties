// properties: This package implements read/write operations on .properties file
//
// .properties is a file extension for files mainly used in Java related
// technologies to store the configurable parameters of an application.
// They can also be used for storing strings for Internationalization and
// localization; these are known as Property Resource Bundles.
//
// Each parameter is stored as a pair of strings, one storing the name of
// the parameter (called the key), and the other storing the value.
//
// Each line in a .properties file normally stores a single property.
// Several formats are possible for each line, including key=value,
// key = value, key:value, and key value.
//
// .properties files can use the number sign (#) or the exclamation mark (!)
// as the first non blank character in a line to denote that all text following
// it is a comment. The backwards slash is used to escape a character.
// An example of a properties file is provided below.
// <code>
// # You are reading the ".properties" entry.
// ! The exclamation mark can also mark text as comments.
// website = http://en.wikipedia.org/
// language = English
// # The backslash below tells the application to continue reading
// # the value onto the next line.
// message = Welcome to \
//           Wikipedia!
// # Add spaces to the key
// key\ with\ spaces = This is the value that could be looked up with the \
// key "key with spaces".
// # Empty lines are skipped
//
// # Unicode
// unicode=\u041f\u0440\u0438\u0432\u0435\u0442, \u0421\u043e\u0432\u0430!
// # Comment
// </code>
//
// In the example above, website would be a key, and its corresponding
// value would be http://en.wikipedia.org/. While the number sign and the
// exclamation mark marks text as comments, it has no effect when it is part
// of a property. Thus, the key message has the value Welcome to Wikipedia!
// and not Welcome to Wikipedia. Note also that all of the whitespace in
// front of Wikipedia! is excluded completely.
//
// The encoding of a .properties file is ISO-8859-1, also known as Latin-1.
// All non-Latin-1 characters must be entered by using Unicode escape characters,
// e. g. \uHHHH where HHHH is a hexadecimal index of the character in the Unicode
// character set. This allows for using .properties files as resource bundles for
// localization. A non-Latin-1 text file can be converted to a correct .properties
// file by using the native2ascii tool that is shipped with the JDK or by using
// a tool, such as po2prop, that manages the transformation from a bilingual
// localization format into .properties escaping.
//
// From Wikipedia, the free encyclopedia
// http://en.wikipedia.org/wiki/.properties
package properties

import (
	"errors"
	"io"
	"unicode/utf8"
)

// ErrMalformedUtf8Encoding means that it was not possible to convert \uXXXX
// string to utf8 rune.
var ErrMalformedUtf8Encoding error = errors.New("malformed \\uxxxx encoding")

// Reads key value pairs from reader and returns map[string]string
func Load(src io.Reader) (props map[string]string, err error) {
	err = nil
	lr := newLineReader(src)
	props = make(map[string]string)
	for {
		s, e := lr.readLine()
		if e == io.EOF {
			break
		}
		if e != nil {
			return nil, e
		}

		keyLen := 0
		precedingBackslash := false
		hasSep := false
		valueStart := len(s)

		for keyLen < len(s) {
			c := s[keyLen]

			if (c == '=' || c == ':') && !precedingBackslash {
				valueStart = keyLen + 1
				hasSep = true
				break
			}
			if (c == ' ' || c == '\t' || c == '\f') && !precedingBackslash {
				valueStart = keyLen + 1
				break
			}
			if c == '\\' {
				precedingBackslash = !precedingBackslash
			} else {
				precedingBackslash = false
			}

			keyLen++
		}

		for valueStart < len(s) {
			c := s[valueStart]
			if c != ' ' && c != '\t' && c != '\f' {
				if !hasSep && (c == '=' || c == ':') {
					hasSep = true
				} else {
					break
				}
			}
			valueStart++
		}
		key, err := decodeString(s[0:keyLen])
		if err != nil {
			return nil, err
		}
		value, err := decodeString(s[valueStart:len(s)])
		if err != nil {
			return nil, err
		}
		props[key] = value
	}
	return props, err
}

// Decodes \t,\n,\r,\f and \uXXXX characters in string
func decodeString(in string) (string, error) {
	out := make([]byte, len(in))
	o := 0
	for i := 0; i < len(in); {
		if in[i] == '\\' {
			i++
			switch in[i] {
			case 'u':
				i++
				rune := 0
				for j := 0; j < 4; j++ {
					switch {
					case in[i] >= '0' && in[i] <= '9':
						rune = (rune << 4) + int(in[i]) - '0'
						break
					case in[i] >= 'a' && in[i] <= 'f':
						rune = (rune << 4) + 10 + int(in[i]) - 'a'
						break
					case in[i] >= 'A' && in[i] <= 'F':
						rune = (rune << 4) + 10 + int(in[i]) - 'A'
						break
					default:
						return "", ErrMalformedUtf8Encoding
					}
					i++
				}
				bytes := make([]byte, utf8.RuneLen(rune))
				bytesWritten := utf8.EncodeRune(bytes, rune)
				for j := 0; j < bytesWritten; j++ {
					out[o] = bytes[j]
					o++
				}
				continue
			case 't':
				out[o] = '\t'
				o++
				i++
				continue
			case 'r':
				out[o] = '\r'
				o++
				i++
				continue
			case 'n':
				out[o] = '\n'
				o++
				i++
				continue
			case 'f':
				out[o] = '\f'
				o++
				i++
				continue
			}
			out[o] = in[i]
			o++
			i++
			continue
		}
		out[o] = in[i]
		o++
		i++
	}

	return string(out[0:o]), nil
}

// Read in a "logical line" from an InputStream/Reader, skip all comment
// and blank lines and filter out those leading whitespace characters
// (\u0020, \u0009 and \u000c) from the beginning of a "natural line".
type lineReader struct {
	reader     io.Reader
	buffer     []byte
	lineBuffer []byte
	limit      int
	offset     int
	exhausted  bool
}

func newLineReader(r io.Reader) *lineReader {
	n := new(lineReader)
	n.reader = r
	n.buffer = make([]byte, 1024)
	n.lineBuffer = make([]byte, 1024)
	n.limit = 0
	n.offset = 0
	n.exhausted = false
	return n
}

// Returns the "logical line" from given reader
func (lr *lineReader) readLine() (line string, e error) {
	if lr.exhausted {
		return "", io.EOF
	}
	nextCharIndex := 0
	char := byte(0)

	skipLF := false
	skipWhiteSpace := true
	appendedLineBegin := false
	isNewLine := true
	isCommentLine := false
	precedingBackslash := false

	for {
		if lr.offset >= lr.limit {
			lr.limit, e = io.ReadFull(lr.reader, lr.buffer)
			lr.offset = 0
			if e == io.EOF {
				lr.exhausted = true
				if isCommentLine {
					return "", io.EOF
				}
				return string(lr.lineBuffer[0:nextCharIndex]), nil
			}
			if e == io.ErrUnexpectedEOF {
				if isCommentLine {
					return "", io.EOF
				}
				continue
			}
			if e != nil {
				lr.exhausted = true
				return "", e
			}
		}

		char = lr.buffer[lr.offset]
		lr.offset++

		if skipLF {
			skipLF = false
			if char == '\n' {
				continue
			}
		}

		if skipWhiteSpace {
			if char == ' ' || char == '\t' || char == '\f' {
				continue
			}
			if !appendedLineBegin && (char == '\r' || char == '\n') {
				continue
			}
			skipWhiteSpace = false
			appendedLineBegin = false
		}

		if isNewLine {
			isNewLine = false
			if char == '#' || char == '!' {
				isCommentLine = true
				continue
			}
		}

		if char != '\n' && char != '\r' {
			lr.lineBuffer[nextCharIndex] = char
			nextCharIndex++
			if nextCharIndex == len(lr.lineBuffer) {
				newBuffer := make([]byte, len(lr.lineBuffer)*2)
				for i, x := range lr.lineBuffer {
					newBuffer[i] = x
				}
				lr.lineBuffer = newBuffer
			}
			//flip the preceding backslash flag
			precedingBackslash = char == '\\' && !precedingBackslash
		} else {
			// reached EOL
			if isCommentLine || nextCharIndex == 0 {
				isCommentLine = false
				isNewLine = true
				skipWhiteSpace = true
				nextCharIndex = 0
				continue
			}
			if lr.offset >= lr.limit {
				lr.limit, e = io.ReadFull(lr.reader, lr.buffer)
				lr.offset = 0
				if e == io.EOF || e == io.ErrUnexpectedEOF {
					lr.exhausted = true
					return string(lr.lineBuffer[0:nextCharIndex]), nil
				}
				if e != nil {
					lr.exhausted = true
					return "", e
				}
			}
			if precedingBackslash {
				nextCharIndex--
				//skip the leading whitespace characters in following line
				skipWhiteSpace = true
				appendedLineBegin = true
				precedingBackslash = false
				if char == '\r' {
					skipLF = true
				}
			} else {
				return string(lr.lineBuffer[0:nextCharIndex]), nil
			}
		}
	}

	return "", nil
}
