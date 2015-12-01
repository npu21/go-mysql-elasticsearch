package river

import (
	"strings"

	"github.com/siddontang/go-mysql/schema"
)

// If you want to sync MySQL data into elasticsearch, you must set a rule to let use know how to do it.
// The mapping rule may thi: schema + table <-> index + document type.
// schema and table is for MySQL, index and document type is for Elasticsearch.
type Rule struct {
	Schema string `toml:"schema"`
	Table  string `toml:"table"`
	Index  string `toml:"index"`
	Type   string `toml:"type"`
	Parent string `toml:"parent"`

	// Default, a MySQL table field name is mapped to Elasticsearch field name.
	// Sometimes, you want to use different name, e.g, the MySQL file name is title,
	// but in Elasticsearch, you want to name it my_title.
	SingleFieldMapping map[string]string `toml:"field"`
	FieldMapping       []*FieldMapping

	// MySQL table information
	TableInfo *schema.Table
}

type FieldMapping struct {
	Mysql   string `toml:"mysql"`
	Elastic string `toml:"elastic"`
	Type    string `toml:"type"`
}

func newDefaultRule(schema string, table string) *Rule {
	r := new(Rule)

	r.Schema = schema
	r.Table = table
	r.Index = table
	r.Type = table
	r.FieldMapping = []*FieldMapping{}
	r.SingleFieldMapping = make(map[string]string)

	return r
}

func (r *Rule) prepare() error {
	if r.FieldMapping == nil {
		r.FieldMapping = []*FieldMapping{}
	}

	if r.SingleFieldMapping != nil {
		for k, v := range r.SingleFieldMapping {
			composedField := strings.Split(v, ",")
			field := FieldMapping{
				Mysql:   k,
				Elastic: composedField[0],
			}
			if 0 == len(field.Elastic) {
				field.Elastic = field.Mysql
			}
			if 2 == len(composedField) {
				field.Type = composedField[1]
			}
			r.FieldMapping = append(r.FieldMapping, &field)
		}
	}

	if len(r.Index) == 0 {
		r.Index = r.Table
	}

	if len(r.Type) == 0 {
		r.Type = r.Index
	}

	return nil
}
