package pgnparser

import (
	"cant/models/pgn"
	"cant/models/spn"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// PGNFile represents a pgn/spn import file from cc isobus
type PGNFile struct {
	FilePath string
	fp       *os.File
	sc       *csv.Reader
}

type parseFunc func(string, ...interface{}) error

// Parser implements the actual parsing
type Parser struct {
	Attribute string
	run       func(string, ...interface{}) error
}

// Open creates a new *PGNFile
func Open(filePath string) (*PGNFile, error) {
	if _, err := os.Stat(filePath); err != nil {
		return nil, err
	}
	fp, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	sc := csv.NewReader(fp)
	sc.TrimLeadingSpace = true
	sc.FieldsPerRecord = 32
	// skip header
	if _, err := sc.Read(); err != nil {
		return nil, err
	}
	return &PGNFile{FilePath: filePath, fp: fp, sc: sc}, nil
}

// getAllPGNParsers return all parsing functions
func (pp *PGNFile) getAllParserFuncs() []parseFunc {
	parseUint32 := func(field string, v ...interface{}) error {
		if len(v) != 1 {
			return fmt.Errorf("expected 3 arguments")
		}
		if field == "" {
			*v[0].(*uint32) = uint32(0)
			return nil
		}
		tmp, err := strconv.ParseUint(field, 10, 32)
		if err != nil {
			return fmt.Errorf("not a number: %s", field)
		}
		*v[0].(*uint32) = uint32(tmp)
		return nil
	}
	parseString := func(field string, v ...interface{}) error {
		if len(v) != 1 {
			return fmt.Errorf("expected 1 argument")
		}
		*v[0].(*string) = field
		return nil
	}
	parseBool := func(field string, v ...interface{}) error {
		if len(v) != 1 {
			return fmt.Errorf("expected 1 argument")
		}
		switch strings.ToLower(field) {
		case "no", "", "0":
			field = "false"
		case "yes", "1":
			field = "true"
		default:
			return fmt.Errorf("can not interpret boolean expression: %s", field)
		}
		tmp, err := strconv.ParseBool(field)
		if err != nil {
			return err
		}
		*v[0].(*bool) = tmp
		return nil
	}
	parseInt := func(field string, v ...interface{}) error {
		if len(v) != 1 {
			return fmt.Errorf("expected 1 argument")
		}
		if field == "Variable" {
			*v[0].(*int) = -1
			return nil
		}
		tmp, err := strconv.Atoi(field)
		if err != nil {
			return err
		}
		*v[0].(*int) = tmp
		return nil
	}
	return []parseFunc{
		parseUint32,
		parseString,
		parseInt,
		parseBool,
	}
}

func (pp *PGNFile) getAllPGNParsers() []*Parser {
	fns := pp.getAllParserFuncs()
	return []*Parser{
		&Parser{Attribute: "id", run: fns[0]},
		&Parser{Attribute: "name", run: fns[1]},
		&Parser{Attribute: "edp", run: fns[2]},
		&Parser{Attribute: "dp", run: fns[2]},
		&Parser{Attribute: "pf", run: fns[2]},
		&Parser{Attribute: "ps", run: fns[1]},
		&Parser{Attribute: "multipacket", run: fns[3]},
		&Parser{Attribute: "pgn_dlc", run: fns[2]},
	}
}

func (pp *PGNFile) getAllSPNParsers() []*Parser {
	fns := pp.getAllParserFuncs()
	return []*Parser{
		&Parser{Attribute: "id", run: fns[0]},
		&Parser{Attribute: "name", run: fns[1]},
		&Parser{Attribute: "pgn", run: fns[0]},
	}
}

// ExecutePGN executes the PGN parser
func (p *Parser) ExecutePGN(field string, data *pgn.PGN) error {
	var res error
	switch p.Attribute {
	case "id":
		res = p.run(field, &data.ID)
	case "name":
		res = p.run(field, &data.Name)
	case "edp":
		res = p.run(field, &data.EDP)
	case "dp":
		res = p.run(field, &data.DP)
	case "pf":
		res = p.run(field, &data.PF)
	case "ps":
		res = p.run(field, &data.PS)
	case "multipacket":
		res = p.run(field, &data.Multipacket)
	case "pgn_dlc":
		res = p.run(field, &data.DLC)
	}
	return res
}

// ExecuteSPN executes the SPN parser
func (p *Parser) ExecuteSPN(field string, data *spn.SPN) error {
	var res error
	switch p.Attribute {
	case "id":
		res = p.run(field, &data.ID)
	case "name":
		res = p.run(field, &data.Name)
	case "pgn":
		res = p.run(field, &data.PGN)
	}
	return res
}

// ScanPGN scans a single PGN and returns it
func (pp *PGNFile) ScanPGN() (pgn.PGN, error) {
	record, err := pp.sc.Read()
	data := pgn.New()
	if err != nil {
		return data, err
	}
	parsers := pp.getAllPGNParsers()
	relevant := []int{0, 1, 4, 5, 6, 7, 8, 10}
	for i, parser := range parsers {
		field := record[relevant[i]]
		if field == "" && i == 0 {
			return data, nil
		} else if field == "" && i != 1 {
			field = "0"
		}
		if err := parser.ExecutePGN(field, &data); err != nil {
			return data, err
		}
	}
	return data, nil
}

// ScanSPN scans a single SPN and returns it
func (pp *PGNFile) ScanSPN() (spn.SPN, error) {
	record, err := pp.sc.Read()
	data := spn.New()
	if err != nil {
		return data, err
	}
	parsers := pp.getAllSPNParsers()
	relevant := []int{14, 15, 0}
	for i, parser := range parsers {
		field := record[relevant[i]]
		if field == "" && i == 6 {
			continue
		} else if field == "" && i != 15 {
			field = "0"
		}
		if err := parser.ExecuteSPN(field, &data); err != nil {
			return data, err
		}
	}
	return data, nil
}

// ScanAllSPN scans all SPNs within a *PGNFile. It returns an array of SPNs
func (pp *PGNFile) ScanAllSPN() (spn.SPNs, error) {
	spns := spn.SPNs{}
	for {
		data, err := pp.ScanSPN()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		spns = append(spns, data)
	}
	return spns, nil
}

// ScanAllPGN scans all PGNs within a *PGNFile. It returns an array of PGNs
func (pp *PGNFile) ScanAllPGN() (pgn.PGNs, error) {
	pgns := pgn.PGNs{}
	for {
		data, err := pp.ScanPGN()
		if err == io.EOF {
			break
		} else if err != nil {
			return pgns, err
		}
		if len(pgns) == 0 || data.ID != pgns[len(pgns)-1].ID {
			pgns = append(pgns, data)
		}
	}
	return pgns, nil
}

// Close the *PGNFile
func (pp *PGNFile) Close() {
	pp.fp.Close()
}

// Rewind the PGNFile
func (pp *PGNFile) Rewind() error {
	if _, err := pp.fp.Seek(0, 0); err != nil {
		return err
	}
	if _, err := pp.sc.Read(); err != nil {
		return err
	}
	return nil
}
