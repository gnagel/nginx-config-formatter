package formatter

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Fmt struct {
	ConfigFile   string
	CreateBackup bool
	InPlace      bool
	Verbose bool
}

func (self *Fmt) Run() error {
	// Read in the file
	if self.Verbose {
		log.Println("Reading source file: ", self.ConfigFile)
	}
	data, err := ioutil.ReadFile(self.ConfigFile)
	if nil != err {
		err = errors.Wrap(err, fmt.Sprintf("Failed to open file: %v", self.ConfigFile))
		log.Println("Failed to read source file: ", self.ConfigFile, err)
		return err
	}

	// Perform the auto-formatting
	input := string(data)
	output := FormatBody(input)

	// If we are not writing the file then just dump it to stdout
	if !self.InPlace {
		log.Println(output)
	}

	// If nothing changed, then just return
	if input == output || !self.InPlace {
		if self.Verbose {
			log.Println("No changes to source file")
		}
		return nil
	}

	// Otherwise the content changed,
	// so we need to write the new content back to the original location
	// and created a backup (if needed).
	if self.CreateBackup {
		fileName := fmt.Sprintf("%v.%v.bak", strings.ReplaceAll(self.ConfigFile, ".conf", ""), time.Now().Unix())
		log.Println("Creating backup file: ", fileName)

		err = ioutil.WriteFile(fileName, []byte(input), 0644)
		if nil != err {
			err = errors.Wrap(err, fmt.Sprintf("Failed to create backup file: %v", self.ConfigFile))
			return err
		}
	}

	// Overwrite the original
	err = ioutil.WriteFile(self.ConfigFile, []byte(output), 0644)
	if nil != err {
		log.Println("Error writing config file: ", self.ConfigFile)
		err = errors.Wrap(err, fmt.Sprintf("Failed to write to original file: %v", self.ConfigFile))
		return err
	}

	return nil
}

func FormatBody(body string) string {
	body = EscapeBlocks(body)
	lines := strings.Split(body, "\n")
	lines = CleanLines(lines)
	lines = MoveOpeningBracket(lines)
	lines = ScrubBlankLines(lines)
	lines = IndentLines(lines)
	body = strings.Join(lines, "\n")
	body = UnescapeBlocks(body)
	return body
}
