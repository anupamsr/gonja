package statements

import (
	"fmt"

	"github.com/paradime-io/gonja/exec"
	"github.com/paradime-io/gonja/nodes"
	"github.com/paradime-io/gonja/parser"
	"github.com/paradime-io/gonja/tokens"
)

type AutoescapeStmt struct {
	Wrapper    *nodes.Wrapper
	Autoescape bool
}

func (stmt *AutoescapeStmt) Position() *tokens.Token { return stmt.Wrapper.Position() }
func (stmt *AutoescapeStmt) String() string {
	t := stmt.Position()
	return fmt.Sprintf("AutoescapeStmt(Line=%d Col=%d)", t.Line, t.Col)
}

func (stmt *AutoescapeStmt) Execute(r *exec.Renderer, tag *nodes.StatementBlock) error {
	sub := r.Inherit()
	sub.Autoescape = stmt.Autoescape

	err := sub.ExecuteWrapper(stmt.Wrapper)
	if err != nil {
		return err
	}

	return nil
}

func autoescapeParser(p *parser.Parser, args *parser.Parser) (nodes.Statement, error) {
	stmt := &AutoescapeStmt{}

	wrapper, _, err := p.WrapUntil("endautoescape")
	if err != nil {
		return nil, err
	}
	stmt.Wrapper = wrapper

	modeToken := args.Match(tokens.Name)
	if modeToken == nil {
		return nil, args.Error("A mode is required for autoescape statement.", nil)
	}
	if modeToken.Val == "true" {
		stmt.Autoescape = true
	} else if modeToken.Val == "false" {
		stmt.Autoescape = false
	} else {
		return nil, args.Error("Only 'true' or 'false' is valid as an autoescape statement.", nil)
	}

	if !args.Stream.End() {
		return nil, args.Error("Malformed autoescape statement args.", nil)
	}

	return stmt, nil
}

func init() {
	All.Register("autoescape", autoescapeParser)
}
