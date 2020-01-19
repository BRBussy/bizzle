package XLSXStandardBank

import (
	"bytes"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/BRBussy/bizzle/internal/pkg/budget/statement"
	"github.com/rs/zerolog/log"
)

type Parser struct {
}

func (p Parser) ParseStatement(request *statement.ParseStatementRequest) (*statement.ParseStatementResponse, error) {
	f, err := excelize.OpenReader(bytes.NewReader(request.Statement))
	if err != nil {
		log.Error().Err(err).Msg("unable to open file")
		return nil, err
	}

	fmt.Println(f)

	return &statement.ParseStatementResponse{}, nil
}
