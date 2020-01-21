package XLSXStandardBank

import (
	"github.com/BRBussy/bizzle/internal/pkg/budget/statement"
	testifyAssert "github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestParser_ParseStatement(t *testing.T) {
	assert := testifyAssert.New(t)
	xlsxStandardBankParser := Parser{}

	dat, err := ioutil.ReadFile("/Users/bernardbussy/Google Drive (brbitzbussy@gmail.com)/Personal/2020/Budget/statement.xlsx")
	if err != nil {
		assert.FailNow(
			"failed to open file",
			err.Error(),
		)
		return
	}

	if _, err := xlsxStandardBankParser.ParseStatement(
		&statement.ParseStatementRequest{
			Statement: dat,
		},
	); err != nil {
		assert.FailNow(
			"failed to parse statement",
			err.Error(),
		)
		return
	}
}
