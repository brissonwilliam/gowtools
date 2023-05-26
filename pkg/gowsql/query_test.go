package helpers

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	validInsertMultipleValuesKeys               = "(:int_b,:int_c)"
	validInsertMultipleValuesKeysPartial        = "(:int_c)"
	invalidInsertMultipleValuesKeysUnknownField = "(:int_d)"
	invalidInsertMultipleValuesKeysEmpty        = "()"
	invalidInsertMultipleValuesKeysEmptyString  = ""
	invalidInsertMultipleValuesKeysNoColon      = "(something, something_else)"
)

type buildInsertMultipleValuesTestStruct struct {
	A int32
	B int32 `db:"int_b"`
	C int32 `db:"int_c"`
}

func TestBuildMultipleNamedValuesSucceeds(t *testing.T) {
	s := givenRandomTestStructs(4)
	q, a, e := BuildMultipleNamedValues(
		validInsertMultipleValuesKeys, s,
	)

	assert.NoError(t, e)
	assert.Equal(t, "(?,?),(?,?),(?,?),(?,?)", q)
	expectedArgs := []interface{}{
		s[0].(buildInsertMultipleValuesTestStruct).B,
		s[0].(buildInsertMultipleValuesTestStruct).C,
		s[1].(buildInsertMultipleValuesTestStruct).B,
		s[1].(buildInsertMultipleValuesTestStruct).C,
		s[2].(buildInsertMultipleValuesTestStruct).B,
		s[2].(buildInsertMultipleValuesTestStruct).C,
		s[3].(buildInsertMultipleValuesTestStruct).B,
		s[3].(buildInsertMultipleValuesTestStruct).C,
	}
	assert.Equal(t, expectedArgs, a)
}

func TestBuildMultipleNamedValuesPartial(t *testing.T) {
	s := givenRandomTestStructs(4)
	q, a, e := BuildMultipleNamedValues(
		validInsertMultipleValuesKeysPartial, s,
	)

	assert.NoError(t, e)
	assert.Equal(t, "(?),(?),(?),(?)", q)
	expectedArgs := []interface{}{
		s[0].(buildInsertMultipleValuesTestStruct).C,
		s[1].(buildInsertMultipleValuesTestStruct).C,
		s[2].(buildInsertMultipleValuesTestStruct).C,
		s[3].(buildInsertMultipleValuesTestStruct).C,
	}
	assert.Equal(t, expectedArgs, a)
}

func TestBuildMultipleNamedValuesErrorUnknownField(t *testing.T) {
	s := givenRandomTestStructs(4)
	q, a, e := BuildMultipleNamedValues(
		invalidInsertMultipleValuesKeysUnknownField, s,
	)

	assert.Error(t, e)
	assert.Equal(t, "", q)
	assert.Nil(t, a)
}

func TestBuildMultipleNamedValuesErrorEmpty(t *testing.T) {
	s := givenRandomTestStructs(4)
	q, a, e := BuildMultipleNamedValues(
		invalidInsertMultipleValuesKeysEmpty, s,
	)

	assert.Error(t, e)
	assert.Equal(t, "Did not find named parameter in keys", e.Error())
	assert.Equal(t, "", q)
	assert.Nil(t, a)
}

func TestBuildMultipleNamedValuesErrorEmptyString(t *testing.T) {
	s := givenRandomTestStructs(4)
	q, a, e := BuildMultipleNamedValues(
		invalidInsertMultipleValuesKeysEmptyString, s,
	)

	assert.Error(t, e)
	assert.Equal(t, "Did not find named parameter in keys", e.Error())
	assert.Equal(t, "", q)
	assert.Nil(t, a)
}

func TestBuildMultipleNamedValuesErrorKeysNoColon(t *testing.T) {
	s := givenRandomTestStructs(4)
	q, a, e := BuildMultipleNamedValues(
		invalidInsertMultipleValuesKeysNoColon, s,
	)

	assert.Error(t, e)
	assert.Equal(t, "Did not find named parameter in keys", e.Error())
	assert.Equal(t, "", q)
	assert.Nil(t, a)
}

func TestPrefixColumns(t *testing.T) {
	assert.Equal(t, "", PrefixColumns("", "x"))
	assert.Equal(t, "x.a1,x.a2,x.a3", PrefixColumns("a1,a2,a3", "x"))
}

func givenARandomTestStruct() buildInsertMultipleValuesTestStruct {
	return buildInsertMultipleValuesTestStruct{
		A: rand.Int31(),
		B: rand.Int31(),
		C: rand.Int31(),
	}
}

func givenRandomTestStructs(count uint) (outStructs []interface{}) {
	for i := uint(0); i < count; i++ {
		s := givenARandomTestStruct()
		outStructs = append(outStructs, s)
	}
	return outStructs
}

func TestBuildNamedColumns(t *testing.T) {
	namedCols := BuildNamedColumns("a,b,c")
	assert.Equal(t, ":a,:b,:c", namedCols)

	namedCols = BuildNamedColumns("a")
	assert.Equal(t, ":a", namedCols)

	namedCols = BuildNamedColumns("")
	assert.Equal(t, "", namedCols)
}

func TestBuildOnDuplicateClause(t *testing.T) {
	t.Run("With override mode", func(t *testing.T) {
		m := OnDuplicateOverride

		s := BuildOnDuplicateClause("a,b,c", m)
		assert.Equal(t, "\t\ta = VALUES(a),\n\t\tb = VALUES(b),\n\t\tc = VALUES(c)", s)

		s = BuildOnDuplicateClause("a", m)
		assert.Equal(t, "\t\ta = VALUES(a)", s)

		s = BuildOnDuplicateClause("", OnDuplicateOverride)
		assert.Equal(t, "", s)
	})
	t.Run("With increment mode", func(t *testing.T) {
		m := OnDuplicateIncrement

		s := BuildOnDuplicateClause("a,b,c", m)
		assert.Equal(t, "\t\ta = a + VALUES(a),\n\t\tb = b + VALUES(b),\n\t\tc = c + VALUES(c)", s)

		s = BuildOnDuplicateClause("a", m)
		assert.Equal(t, "\t\ta = a + VALUES(a)", s)

		s = BuildOnDuplicateClause("", OnDuplicateOverride)
		assert.Equal(t, "", s)
	})
	t.Run("With invalid mode defaults to override", func(t *testing.T) {
		s := BuildOnDuplicateClause("a,b,c", 500)
		assert.Equal(t, "\t\ta = VALUES(a),\n\t\tb = VALUES(b),\n\t\tc = VALUES(c)", s)
	})
}
