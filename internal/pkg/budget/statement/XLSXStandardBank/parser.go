package XLSXStandardBank

import (
	"fmt"
	"github.com/BRBussy/bizzle/internal/pkg/budget/statement"
	"github.com/rs/zerolog/log"
	"github.com/tealeg/xlsx"
	"strconv"
	"time"
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
	colHeaderIndex := make(map[ColumnHeader]int)

	// check for minimum no rows
	if len(transactionsSheet.Rows) < 3 {
		reasonsInvalid = append(reasonsInvalid, "less than 3 rows in sheet")
	} else {
		// index columns
		for headerCellIdx := range transactionsSheet.Rows[0].Cells {
			colHeaderIndex[ColumnHeader(transactionsSheet.Rows[0].Cells[headerCellIdx].Value)] = headerCellIdx
		}
	}

	// check all required column headers are present
	for _, requiredColumnHeader := range RequiredColumnHeaders {
		if _, found := colHeaderIndex[requiredColumnHeader]; !found {
			reasonsInvalid = append(reasonsInvalid, "missing column with header "+requiredColumnHeader.String())
		}
	}

	year, err := transactionsSheet.Rows[1].Cells[colHeaderIndex[DateColumnHeader]].Int()
	if err != nil {
		reasonsInvalid = append(reasonsInvalid, "could not find starting year")
	}

	if len(reasonsInvalid) > 0 {
		return nil, ErrSheetInvalid{Reasons: reasonsInvalid}
	}

	// sheet appears valid, parse the rest
	for rowIdx := range transactionsSheet.Rows[2:] {
		// check if this is a year row and update year if it is
		potentialYear, err := transactionsSheet.Rows[1].Cells[colHeaderIndex[DateColumnHeader]].Int()
		if err == nil {
			// assume this is a year row
			// update year
			year = potentialYear
			// go to next row
			continue
		}

		// otherwise assume this is an entry row, try and parse the date
		entryDate, err := time.Parse(
			"",
			fmt.Sprintf(
				"%s %s",
				"",
				strconv.Itoa(year),
			),
		)
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("could not parse date in row %i", rowIdx+1))
		}

		// parse
	}

	return &statement.ParseStatementResponse{}, nil
}
