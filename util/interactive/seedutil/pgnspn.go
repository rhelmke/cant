package seedutil

import (
	"cant/models/pgn"
	"cant/models/spn"
	"cant/util/database"
	"cant/util/globals"
	"cant/util/parsers/pgnparser"
	"fmt"

	"github.com/BTBurke/clt"
)

// PGNSPNSession is a wrapper for clt.InteractiveSession
type PGNSPNSession struct {
	is       *clt.InteractiveSession
	spns     *spn.SPNs
	pgns     *pgn.PGNs
	filePath string
	parser   *pgnparser.PGNFile
}

// NewPGNSPNSession creates a new clt.InteractiveSession-Wrapper
func NewPGNSPNSession(filePath string) (*PGNSPNSession, error) {
	parser, err := pgnparser.Open(filePath)
	if err != nil {
		return nil, err
	}
	return &PGNSPNSession{is: clt.NewInteractiveSession(), pgns: &pgn.PGNs{}, spns: &spn.SPNs{}, filePath: filePath, parser: parser}, nil
}

// Welcome greets the user
func (sess *PGNSPNSession) Welcome() error {
	resp := sess.is.Say("We are about to seed the database with PGN and SPN data").
		Warn(clt.SStyled("This will wipe all PGN/SPN-related data!", clt.Bold, clt.Red)).
		AskYesNo("Are you sure you want to continue?", "n")
	if clt.IsNo(resp) {
		return fmt.Errorf("abort. shutting down")
	}
	return nil
}

// WriteToDB invokes the parser provided by package util/parsers/pgnparser and saves the data to the MySQL DB
func (sess *PGNSPNSession) WriteToDB() error {
	defer sess.parser.Close()
	spinner := clt.NewProgressSpinner("Reading " + clt.SStyled("'"+sess.parser.FilePath+"'", clt.Bold))
	spinner.Start()
	pgns, err := sess.parser.ScanAllPGN()
	if err != nil {
		spinner.Fail()
		return err
	}
	if err := sess.parser.Rewind(); err != nil {
		return err
	}
	spns, err := sess.parser.ScanAllSPN()
	if err != nil {
		spinner.Fail()
		return err
	}
	spinner.Success()
	spinner = clt.NewProgressSpinner("Wiping tables")
	spinner.Start()
	tables := []string{"pgn", "spn"}
	for _, table := range tables {
		if err := database.WipeTable(globals.DB, table); err != nil {
			spinner.Fail()
			return err
		}
	}
	spinner.Success()
	spinner = clt.NewProgressSpinner("Writing PGNs")
	spinner.Start()
	if err := pgns.SaveAll(); err != nil {
		spinner.Fail()
		return err
	}
	spinner.Success()
	spinner = clt.NewProgressSpinner("Writing SPNs")
	spinner.Start()
	if err := spns.SaveAll(); err != nil {
		spinner.Fail()
		return err
	}
	spinner.Success()
	sess.is.Say(clt.SStyled("Seeded PGN dataset with %d entries.", clt.Bold), len(pgns))
	sess.is.Say(clt.SStyled("Seeded SPN dataset with %d entries.", clt.Bold), len(spns))
	return nil
}
