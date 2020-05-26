package cavgocql

import (
	"strings"

	"goAgent/internal/sqlscanner"
)

func querySignature(query string) string {
	s := sqlscanner.NewScanner(query)
	for s.Scan() {
		if s.Token() != sqlscanner.COMMENT {
			break
		}
	}

	scanUntil := func(until sqlscanner.Token) bool {
		for s.Scan() {
			if s.Token() == until {
				return true
			}
		}
		return false
	}
	scanToken := func(tok sqlscanner.Token) bool {
		for s.Scan() {
			switch s.Token() {
			case tok:
				return true
			case sqlscanner.COMMENT:
			default:
				return false
			}
		}
		return false
	}

	switch s.Token() {
	case sqlscanner.DELETE:
		if !scanUntil(sqlscanner.FROM) {
			break
		}
		if !scanToken(sqlscanner.IDENT) {
			break
		}
		tableName := s.Text()
		for scanToken(sqlscanner.PERIOD) && scanToken(sqlscanner.IDENT) {
			tableName += "." + s.Text()
		}
		return "DELETE FROM " + tableName

	case sqlscanner.INSERT:
		if !scanUntil(sqlscanner.INTO) {
			break
		}
		if !scanToken(sqlscanner.IDENT) {
			break
		}
		tableName := s.Text()
		for scanToken(sqlscanner.PERIOD) && scanToken(sqlscanner.IDENT) {
			tableName += "." + s.Text()
		}
		return "INSERT INTO " + tableName

	case sqlscanner.SELECT:
		var level int
	scanLoop:
		for s.Scan() {
			switch tok := s.Token(); tok {
			case sqlscanner.LPAREN:
				level++
			case sqlscanner.RPAREN:
				level--
			case sqlscanner.FROM:
				if level != 0 {
					continue scanLoop
				}
				if !scanToken(sqlscanner.IDENT) {
					break scanLoop
				}
				tableName := s.Text()
				for scanToken(sqlscanner.PERIOD) && scanToken(sqlscanner.IDENT) {
					tableName += "." + s.Text()
				}
				return "SELECT FROM " + tableName
			}
		}

	case sqlscanner.TRUNCATE:
		if !scanUntil(sqlscanner.IDENT) {
			break
		}
		tableName := s.Text()
		for scanToken(sqlscanner.PERIOD) && scanToken(sqlscanner.IDENT) {
			tableName += "." + s.Text()
		}
		return "TRUNCATE " + tableName

	case sqlscanner.UPDATE:
		if !scanToken(sqlscanner.IDENT) {
			break
		}
		tableName := s.Text()
		for scanToken(sqlscanner.PERIOD) && scanToken(sqlscanner.IDENT) {
			tableName += "." + s.Text()
		}
		return "UPDATE " + tableName
	}

	fields := strings.Fields(query)
	if len(fields) == 0 {
		return ""
	}
	return strings.ToUpper(fields[0])
}
