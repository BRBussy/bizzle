package XLSXStandardBank

import (
	"fmt"
	"github.com/BRBussy/bizzle/internal/pkg/budget/statement"
	"github.com/rs/zerolog/log"
	"github.com/tealeg/xlsx"
)

type Parser struct {
}

func (p Parser) ParseStatement(request *statement.ParseStatementRequest) (*statement.ParseStatementResponse, error) {
	// parse file
	excelFile, err := xlsx.OpenBinary(request.Statement)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse file")
		return nil, ErrUnableToParseFile{}
	}

	transactionsSheet, found := excelFile.Sheet["transactions"]
	if !found {
		err = ErrTransactionsSheetNotFound{}
		log.Error().Err(err)
		return nil, err
	}
	for _, row := range transactionsSheet.Rows {
		for _, cell := range row.Cells {
			fmt.Printf("%s\n", cell.String())
		}
	}

	return &statement.ParseStatementResponse{}, nil
}
