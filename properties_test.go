// Copyright (c) 2012 The Goproperties Authors.
// 
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is furnished to do
// so, subject to the following conditions:
// 
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
// 
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package properties

import (
	"bytes"
	. "launchpad.net/gocheck"
	"math"
	"os"
	"testing"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type PropertiesSuite struct {
	p Properties
}

var _ = Suite(&PropertiesSuite{})

func (s *PropertiesSuite) SetUpTest(c *C) {
	r := bytes.NewReader([]byte(source))
	p := &Properties{}
	err := p.Load(r)
	if err != nil {
		c.Fatalf("failed to load properties: %s", err)
	}
	s.p = *p
}

func (s *PropertiesSuite) TestGeneric(c *C) {
	c.Assert(s.p["website"], Equals, "http://en.wikipedia.org/")
	c.Assert(s.p["language"], Equals, "English")
	c.Assert(s.p["message"], Equals, "Welcome to Wikipedia!")
	c.Assert(s.p["unicode"], Equals, "Привет, Сова!")
	c.Assert(s.p["key with spaces"], Equals, "This is the value that could be looked up with the key \"key with spaces\".")
}

func (s *PropertiesSuite) TestLoad(c *C) {
	p, err := Load(c.MkDir() + string(os.PathSeparator) + "nofile")
	c.Assert(p, NotNil)
	c.Assert(err, ErrorMatches, ".*no such file.*")
}

func (s *PropertiesSuite) TestGetString(c *C) {
	c.Assert(s.p.GetString("string", "not found"), Equals, "found")
	c.Assert(s.p.GetString("missed", "not found"), Equals, "not found")
}

func (s *PropertiesSuite) TestGetBool(c *C) {
	c.Assert(s.p.GetBool("bool", false), Equals, true)
	c.Assert(s.p.GetBool("missed", true), Equals, true)
}

func (s *PropertiesSuite) TestGetFloat(c *C) {
	c.Assert(s.p.GetFloat("float", math.MaxFloat64), Equals, math.SmallestNonzeroFloat64)
	c.Assert(s.p.GetFloat("missed", math.MaxFloat64), Equals, math.MaxFloat64)
}

func (s *PropertiesSuite) TestGetInt(c *C) {
	c.Assert(s.p.GetInt("int", math.MaxInt64), Equals, int64(math.MinInt64))
	c.Assert(s.p.GetInt("missed", math.MaxInt64), Equals, int64(math.MaxInt64))
	c.Assert(s.p.GetInt("hex", 0xCAFEBABE), Equals, int64(0xCAFEBABE))
}

func (s *PropertiesSuite) TestGetUint(c *C) {
	c.Assert(s.p.GetUint("uint", 42), Equals, uint64(math.MaxUint64))
	c.Assert(s.p.GetUint("missed", 42), Equals, uint64(42))
	c.Assert(s.p.GetUint("hex", 0xCAFEBABE), Equals, uint64(0xCAFEBABE))
}

const source = `
# You are reading the ".properties" entry.
! The exclamation mark can also mark text as comments.
website = http://en.wikipedia.org/
language = English
# The backslash below tells the application to continue reading
# the value onto the next line.
message = Welcome to \
          Wikipedia!
# Add spaces to the key
key\ with\ spaces = This is the value that could be looked up with the \
key "key with spaces".
# Empty lines are skipped


# Unicode
unicode=\u041f\u0440\u0438\u0432\u0435\u0442, \u0421\u043e\u0432\u0430!
# Comment
string=found
bool=true
float=4.940656458412465441765687928682213723651e-324
int=-9223372036854775808
uint=18446744073709551615
hex=0xCAFEBABE
`
