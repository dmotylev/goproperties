// properties: This package implements access to .properties file
// See further details on http://en.wikipedia.org/wiki/.properties
package properties

import "io"
import "os"
import "log"

// Reads key value pairs from reader and returns map[string]string
func Load(src io.Reader) (props map[string]string, err os.Error) {
	err = nil;
	lr := newLineReader(src);
	for {
		s, e := lr.readLine();
		if e == os.EOF {
			break
		}
		if e != nil {
			err = e;
			break;
		}


		log.Stdout(s);
	}
	return nil, err;
}

/* Read in a "logical line" from an InputStream/Reader, skip all comment
 * and blank lines and filter out those leading whitespace characters
 * (\u0020, \u0009 and \u000c) from the beginning of a "natural line".
 * Method returns the char length of the "logical line" and stores
 * the line in "buffer".
 */
type lineReader struct {
	reader		io.Reader;
	buffer		[]byte;
	lineBuffer	[]byte;
	limit		int;
	offset		int;
	exhausted	bool;
}

func newLineReader(r io.Reader) *lineReader {
	n := new(lineReader);
	n.reader = r;
	n.buffer = make([]byte, 1024);
	n.lineBuffer = make([]byte, 1024);
	n.limit = 0;
	n.offset = 0;
	n.exhausted = false;
	return n;
}

func (lr *lineReader) readLine() (line string, e os.Error) {
	if lr.exhausted {
		return "", os.EOF
	}
	nextCharIndex := 0;
	char := byte(0);

	skipLF := false;
	skipWhiteSpace := true;
	appendedLineBegin := false;
	isNewLine := true;
	isCommentLine := false;
	precedingBackslash := false;

	for {
		if lr.offset >= lr.limit {
			lr.limit, e = io.ReadFull(lr.reader, lr.buffer);
			lr.offset = 0;
			if e == os.EOF {
				lr.exhausted = true;
				if isCommentLine {
					return "", os.EOF
				}
				return string(lr.lineBuffer[0:nextCharIndex]), nil;
			}
			if e == io.ErrUnexpectedEOF {
				if isCommentLine {
					return "", os.EOF
				}
				continue;
			}
			if e != nil {
				lr.exhausted = true;
				return "", e;
			}
		}

		char = lr.buffer[lr.offset];
		lr.offset++;

		if skipLF {
			skipLF = false;
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
			skipWhiteSpace = false;
			appendedLineBegin = false;
		}

		if isNewLine {
			isNewLine = false;
			if char == '#' || char == '!' {
				isCommentLine = true;
				continue;
			}
		}

		if char != '\n' && char != '\r' {
			lr.lineBuffer[nextCharIndex] = char;
			nextCharIndex++;
			if nextCharIndex == len(lr.lineBuffer) {
				newBuffer := make([]byte, len(lr.lineBuffer)*2);
				for i, x := range lr.lineBuffer {
					newBuffer[i] = x
				}
				lr.lineBuffer = newBuffer;
			}
			//flip the preceding backslash flag
			precedingBackslash = char == '\\' && !precedingBackslash;
		} else {
			// reached EOL
			if isCommentLine || nextCharIndex == 0 {
				isCommentLine = false;
				isNewLine = true;
				skipWhiteSpace = true;
				nextCharIndex = 0;
				continue;
			}
			if lr.offset >= lr.limit {
				lr.limit, e = io.ReadFull(lr.reader, lr.buffer);
				lr.offset = 0;
				if e == os.EOF || e == io.ErrUnexpectedEOF {
					lr.exhausted = true;
					return string(lr.lineBuffer[0:nextCharIndex]), nil;
				}
				if e != nil {
					lr.exhausted = true;
					return "", e;
				}
			}
			if precedingBackslash {
				nextCharIndex--;
				//skip the leading whitespace characters in following line
				skipWhiteSpace = true;
				appendedLineBegin = true;
				precedingBackslash = false;
				if char == '\r' {
					skipLF = true
				}
			} else {
				return string(lr.lineBuffer[0:nextCharIndex]), nil
			}
		}
	}

	return "", nil;
}
