package query

import (
	"bytes"
	"strings"
	"sync"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Copied from golint
var commonInitialisms = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
var commonInitialismsReplacer *strings.Replacer

func init() {
	var commonInitialismsForReplacer []string
	for _, initialism := range commonInitialisms {
		caser := cases.Title(language.AmericanEnglish)
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, caser.String(strings.ToLower(initialism)))
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
}

type safeMap struct {
	m map[string]string
	l *sync.RWMutex
}

func newSafeMap() *safeMap {
	return &safeMap{l: new(sync.RWMutex), m: make(map[string]string)}
}

func (s *safeMap) Get(key string) string {
	s.l.RLock()
	defer s.l.RUnlock()
	return s.m[key]
}

func (s *safeMap) Set(key string, value string) {
	s.l.Lock()
	defer s.l.Unlock()
	s.m[key] = value
}

// Namer is a function type which is given a string and return a string
type Namer func(string) string

// NamingStrategy represents naming strategies
type NamingStrategy struct {
	DB     Namer
	Table  Namer
	Column Namer
}

// TheNamingStrategy is being initialized with defaultNamingStrategy
var TheNamingStrategy = &NamingStrategy{
	DB:     defaultNamer,
	Table:  defaultNamer,
	Column: defaultNamer,
}

// AddNamingStrategy sets the naming strategy
func AddNamingStrategy(ns *NamingStrategy) {
	if ns.DB == nil {
		ns.DB = defaultNamer
	}
	if ns.Table == nil {
		ns.Table = defaultNamer
	}
	if ns.Column == nil {
		ns.Column = defaultNamer
	}
	TheNamingStrategy = ns
}

// DBName alters the given name by DB
func (ns *NamingStrategy) DBName(name string) string {
	return ns.DB(name)
}

// TableName alters the given name by Table
func (ns *NamingStrategy) TableName(name string) string {
	return ns.Table(name)
}

// ColumnName alters the given name by Column
func (ns *NamingStrategy) ColumnName(name string) string {
	return ns.Column(name)
}

// ToDBName convert string to db name
func ToDBName(name string) string {
	return TheNamingStrategy.DBName(name)
}

// ToTableName convert string to table name
func ToTableName(name string) string {
	return TheNamingStrategy.TableName(name)
}

// ToColumnName convert string to db name
func ToColumnName(name string) string {
	return TheNamingStrategy.ColumnName(name)
}

var smap = newSafeMap()

func defaultNamer(name string) string {
	const (
		lower = false
		upper = true
	)
	if v := smap.Get(name); v != "" {
		return v
	}
	if name == "" {
		return ""
	}
	var (
		value                                    = commonInitialismsReplacer.Replace(name)
		buf                                      = bytes.NewBufferString("")
		lastCase, currCase, nextCase, nextNumber bool
	)
	for i, v := range value[:len(value)-1] {
		nextCase = bool(value[i+1] >= 'A' && value[i+1] <= 'Z')
		nextNumber = bool(value[i+1] >= '0' && value[i+1] <= '9')
		if i > 0 {
			if currCase == bool(upper) {
				if lastCase == bool(upper) && (nextCase == bool(upper) || nextNumber == bool(upper)) {
					buf.WriteRune(v)
				} else {
					if value[i-1] != '_' && value[i+1] != '_' {
						buf.WriteRune('_')
					}
					buf.WriteRune(v)
				}
			} else {
				buf.WriteRune(v)
				if i == len(value)-2 && (nextCase == bool(upper) && nextNumber == bool(lower)) {
					buf.WriteRune('_')
				}
			}
		} else {
			currCase = upper
			buf.WriteRune(v)
		}
		lastCase = currCase
		currCase = nextCase
	}
	buf.WriteByte(value[len(value)-1])
	s := strings.ToLower(buf.String())
	smap.Set(name, s)
	return s
}
