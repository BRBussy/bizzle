package XLSXStandardBank

import (
	"fmt"
	"github.com/BRBussy/bizzle/internal/pkg/budget"
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

	// return if the sheet is invalid
	if len(reasonsInvalid) > 0 {
		return nil, ErrSheetInvalid{Reasons: reasonsInvalid}
	}

	// sheet appears valid, parse the rest
	parsedBudgetEntries := make([]budget.Entry, 0)
	for rowIdx := range transactionsSheet.Rows[2:] {
		// check if this is a year row and update year if it is
		potentialYear, err := transactionsSheet.Rows[2:][rowIdx].Cells[colHeaderIndex[DateColumnHeader]].Int()
		if err == nil {
			// assume this is a year row
			// update year
			year = potentialYear
			// go to next row
			continue
		}

		// otherwise assume this is an entry row, try and parse the date
		entryDate, err := time.Parse(
			"02 Jan 2006",
			fmt.Sprintf(
				"%s %s",
				transactionsSheet.Rows[2:][rowIdx].Cells[colHeaderIndex[DateColumnHeader]].String(),
				strconv.Itoa(year),
			),
		)
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("could not parse date in row %d", rowIdx+3))
			continue
		}

		// try and parse the out and in amounts
		inAmountCell := transactionsSheet.Rows[2:][rowIdx].Cells[colHeaderIndex[InColumnHeader]]
		outAmountCell := transactionsSheet.Rows[2:][rowIdx].Cells[colHeaderIndex[OutColumnHeader]]
		var amount float64
		if inAmountCell.String() == "" && outAmountCell.String() == "" {
			// both are blank
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("both in and out amounts are blank in row %d", rowIdx+3))
			continue
		} else if inAmountCell.String() != "" && outAmountCell.String() == "" {
			// only in value set, try and parse
			amount, err = inAmountCell.Float()
			if err != nil {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("could not parse in amount in row %d", rowIdx+3))
				continue
			}
		} else if inAmountCell.String() == "" && outAmountCell.String() != "" {
			// only out value set, try and parse
			amount, err = outAmountCell.Float()
			if err != nil {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("could not parse out amount in row %d", rowIdx+3))
				continue
			}
		} else {
			// both are set
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("both in and out amount set in row %d", rowIdx+3))
			continue
		}

		parsedBudgetEntries = append(
			parsedBudgetEntries,
			budget.Entry{
				Date:        entryDate.Unix(),
				Description: transactionsSheet.Rows[2:][rowIdx].Cells[colHeaderIndex[DescriptionColumnHeader]].String(),
				Amount:      amount,
			},
		)
	}

	// return if parsing any rows failed
	if len(reasonsInvalid) > 0 {
		return nil, ErrSheetInvalid{Reasons: reasonsInvalid}
	}

	for _, item := range parsedBudgetEntries {
		fmt.Printf("%v\n", item)
	}

	return &statement.ParseStatementResponse{}, nil
}
