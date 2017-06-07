// Copyright (c) 2017 Ernest Micklei
// 
// MIT License
// 
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
// 
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
// 
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package proto

// Group represents a (proto2 only) group.
// https://developers.google.com/protocol-buffers/docs/reference/proto2-spec#group_field
type Group struct {
	Name     string
	Optional bool
	Sequence int
	Elements []Visitee
}

// Accept dispatches the call to the visitor.
func (g *Group) Accept(v Visitor) {
	v.VisitGroup(g)
}

// addElement is part of elementContainer
func (g *Group) addElement(v Visitee) {
	g.Elements = append(g.Elements, v)
}

// elements is part of elementContainer
func (g *Group) elements() []Visitee {
	return g.Elements
}

// parse expects:
// groupName "=" fieldNumber { messageBody }
func (g *Group) parse(p *Parser) error {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != tIDENT {
		if !isKeyword(tok) {
			return p.unexpected(lit, "group name", g)
		}
	}
	g.Name = lit
	tok, lit = p.scanIgnoreWhitespace()
	if tok != tEQUALS {
		return p.unexpected(lit, "group =", g)
	}
	i, err := p.s.scanInteger()
	if err != nil {
		return p.unexpected(lit, "group sequence number", g)
	}
	g.Sequence = i
	tok, lit = p.scanIgnoreWhitespace()
	if tok != tLEFTCURLY {
		return p.unexpected(lit, "group opening {", g)
	}
	parseMessageBody(p, g)
	return nil
}
