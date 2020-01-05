package interactive

import (
	"cant/util/interactive/seedutil"
)

// SeedPGNSPN seeds given data into the mysql database
func SeedPGNSPN(filePath string) error {
	sess, err := seedutil.NewPGNSPNSession(filePath)
	if err != nil {
		return err
	}
	if err := sess.Welcome(); err != nil {
		return err
	}
	if err := sess.WriteToDB(); err != nil {
		return err
	}
	return nil
}

func SeedFilter() error {
	sess, err := seedutil.NewFilterSession()
	if err != nil {
		return err
	}
	if err := sess.Welcome(); err != nil {
		return err
	}
	if err := sess.WriteToDB(); err != nil {
		return err
	}
	return nil
}
