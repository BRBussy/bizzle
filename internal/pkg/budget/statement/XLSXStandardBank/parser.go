package XLSXStandardBank

import (
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

	// find transactions sheet
	transactionsSheet, found := excelFile.Sheet["transactions"]
	if !found {
		err = ErrTransactionsSheetNotFound{}
		log.Error().Err(err)
		return nil, err
	}

	// variables for validation and parsing of sheet
	reasonsInvalid := make([]string, 0)
	colHeaderIndex := make(map[string]int)

	// check for minimum no rows
	if len(transactionsSheet.Rows) < 3 {
		reasonsInvalid = append(reasonsInvalid, "less than 3 rows in sheet")
	} else {
		// index columns
		for headerCellIdx := range transactionsSheet.Rows[0].Cells {
			colHeaderIndex[transactionsSheet.Rows[0].Cells[headerCellIdx].Value] = headerCellIdx
		}
	}

	// check all required column headers are present
	for _, requiredColumnHeader := range RequiredColumnHeaders {
		if _, found := colHeaderIndex[requiredColumnHeader.String()]; !found {
			reasonsInvalid = append(reasonsInvalid, "missing column with header "+requiredColumnHeader.String())
		}
	}

	if len(reasonsInvalid) > 0 {
		return nil, ErrSheetInvalid{Reasons: reasonsInvalid}
	}

	return &statement.ParseStatementResponse{}, nil
}
