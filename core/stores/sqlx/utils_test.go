package sqlx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscape(t *testing.T) {
	s := "a\x00\n\r\\'\"\x1ab"

	out := escape(s)

	assert.Equal(t, `a\x00\n\r\\\'\"\x1ab`, out)
}

func TestDesensitize(t *testing.T) {
	datasource := "user:pass@tcp(111.222.333.44:3306)/any_table?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	datasource = desensitize(datasource)
	assert.False(t, strings.Contains(datasource, "user"))
	assert.False(t, strings.Contains(datasource, "pass"))
	assert.True(t, strings.Contains(datasource, "tcp(111.222.333.44:3306)"))
}

func TestDesensitize_WithoutAccount(t *testing.T) {
	datasource := "tcp(111.222.333.44:3306)/any_table?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	datasource = desensitize(datasource)
	assert.True(t, strings.Contains(datasource, "tcp(111.222.333.44:3306)"))
}

func TestFormatForPrint(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		args   []interface{}
		expect string
	}{
		{
			name:   "no args",
			query:  "select user, name from table where id=?",
			expect: `select user, name from table where id=?`,
		},
		{
			name:   "one arg",
			query:  "select user, name from table where id=?",
			args:   []interface{}{"kevin"},
			expect: `select user, name from table where id=? ["kevin"]`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := formatForPrint(test.query, test.args...)
			assert.Equal(t, test.expect, actual)
		})
	}
}
