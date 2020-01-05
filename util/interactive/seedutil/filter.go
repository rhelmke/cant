package seedutil

import (
	"fmt"

	"cant/can/flowroutines/filter/drop"
	"cant/can/flowroutines/filter/manipulate/awgn60928dummy"
	"cant/can/flowroutines/filter/manipulate/awgn65096"
	"cant/can/flowroutines/filter/manipulate/encrypt"
	"cant/can/flowroutines/filter/manipulate/fakeview58880"
	"cant/can/flowroutines/filter/manipulate/fakeview59136"
	"cant/can/flowroutines/filter/manipulate/limit65096"
	"cant/can/flowroutines/filter/manipulate/lowres129025"
	"cant/can/flowroutines/filter/manipulate/tuxawgn65096"
	"cant/can/flowroutines/filter/manipulate/tuxlimit65096"
	"cant/can/flowroutines/filter/manipulate/tuxlowres129025"
	filterModel "cant/models/filter"
	"cant/util/database"
	"cant/util/globals"

	"github.com/BTBurke/clt"
)

// FilterSession is a wrapper for clt.InteractiveSession
type FilterSession struct {
	is    *clt.InteractiveSession
	filts filterModel.Filters
}

// NewFilterSession creates a new clt.InteractiveSession-Wrapper
func NewFilterSession() (*FilterSession, error) {
	return &FilterSession{is: clt.NewInteractiveSession(), filts: filterModel.Filters{}}, nil
}

// Welcome greets the user
func (sess *FilterSession) Welcome() error {
	resp := sess.is.Say("We are about to seed the database with Filter data").
		Warn(clt.SStyled("This will wipe all Filter-related data!", clt.Bold, clt.Red)).
		AskYesNo("Are you sure you want to continue?", "n")
	if clt.IsNo(resp) {
		return fmt.Errorf("abort. shutting down")
	}
	return nil
}

// WriteToDB writes all filters to the db
func (sess *FilterSession) WriteToDB() error {
	spinner := clt.NewProgressSpinner("Wiping tables")
	spinner.Start()
	tables := []string{"filter", "filter_for"}
	for _, table := range tables {
		if err := database.WipeTable(globals.DB, table); err != nil {
			spinner.Fail()
			return err
		}
	}
	spinner.Success()
	spinner = clt.NewProgressSpinner("Writing Filters")
	spinner.Start()
	filt := filterModel.Filter{ID: drop.UniqIdentifier(), Name: drop.GetName()}
	for _, pgn := range drop.SupportedPGNs() {
		filt.For = append(filt.For, filterModel.ForPGN{PGN: pgn, Enabled: false})
	}
	sess.filts = append(sess.filts, filt)
	filt = filterModel.Filter{ID: awgn65096.UniqIdentifier(), Name: awgn65096.GetName()}
	for _, pgn := range awgn65096.SupportedPGNs() {
		filt.For = append(filt.For, filterModel.ForPGN{PGN: pgn, Enabled: false})
	}
	sess.filts = append(sess.filts, filt)
	filt = filterModel.Filter{ID: awgn60928dummy.UniqIdentifier(), Name: awgn60928dummy.GetName()}
	for _, pgn := range awgn60928dummy.SupportedPGNs() {
		filt.For = append(filt.For, filterModel.ForPGN{PGN: pgn, Enabled: false})
	}
	sess.filts = append(sess.filts, filt)
	filt = filterModel.Filter{ID: limit65096.UniqIdentifier(), Name: limit65096.GetName()}
	for _, pgn := range limit65096.SupportedPGNs() {
		filt.For = append(filt.For, filterModel.ForPGN{PGN: pgn, Enabled: false})
	}
	sess.filts = append(sess.filts, filt)
	filt = filterModel.Filter{ID: lowres129025.UniqIdentifier(), Name: lowres129025.GetName()}
	for _, pgn := range lowres129025.SupportedPGNs() {
		filt.For = append(filt.For, filterModel.ForPGN{PGN: pgn, Enabled: false})
	}
	sess.filts = append(sess.filts, filt)
	filt = filterModel.Filter{ID: fakeview59136.UniqIdentifier(), Name: fakeview59136.GetName()}
	for _, pgn := range fakeview59136.SupportedPGNs() {
		filt.For = append(filt.For, filterModel.ForPGN{PGN: pgn, Enabled: false})
	}
	sess.filts = append(sess.filts, filt)
	filt = filterModel.Filter{ID: fakeview58880.UniqIdentifier(), Name: fakeview58880.GetName()}
	for _, pgn := range fakeview58880.SupportedPGNs() {
		filt.For = append(filt.For, filterModel.ForPGN{PGN: pgn, Enabled: false})
	}
	sess.filts = append(sess.filts, filt)
	filt = filterModel.Filter{ID: tuxawgn65096.UniqIdentifier(), Name: tuxawgn65096.GetName()}
	for _, pgn := range tuxawgn65096.SupportedPGNs() {
		filt.For = append(filt.For, filterModel.ForPGN{PGN: pgn, Enabled: false})
	}
	sess.filts = append(sess.filts, filt)
	filt = filterModel.Filter{ID: tuxlimit65096.UniqIdentifier(), Name: tuxlimit65096.GetName()}
	for _, pgn := range tuxlimit65096.SupportedPGNs() {
		filt.For = append(filt.For, filterModel.ForPGN{PGN: pgn, Enabled: false})
	}
	sess.filts = append(sess.filts, filt)
	filt = filterModel.Filter{ID: tuxlowres129025.UniqIdentifier(), Name: tuxlowres129025.GetName()}
	for _, pgn := range tuxlowres129025.SupportedPGNs() {
		filt.For = append(filt.For, filterModel.ForPGN{PGN: pgn, Enabled: false})
	}
	sess.filts = append(sess.filts, filt)
	filt = filterModel.Filter{ID: encrypt.UniqIdentifier(), Name: encrypt.GetName()}
	for _, pgn := range encrypt.SupportedPGNs() {
		filt.For = append(filt.For, filterModel.ForPGN{PGN: pgn, Enabled: false})
	}
	sess.filts = append(sess.filts, filt)

	if err := sess.filts.SaveAll(); err != nil {
		spinner.Fail()
		return err
	}
	spinner.Success()
	return nil
}
