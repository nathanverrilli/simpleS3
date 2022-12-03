package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/pflag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var xLogFile *os.File
var xLogBuffer *bufio.Writer
var xLog log.Logger

var nFlags *pflag.FlagSet

func initLog() {
	var err error
	var logWriters []io.Writer

	xLogFile, err = os.OpenFile("simpleS3.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if nil != err {
		_, _ = fmt.Fprintf(os.Stderr, "error opening file: %v", err)
	}
	xLogBuffer = bufio.NewWriter(xLogFile)
	logWriters = append(logWriters, os.Stderr)
	logWriters = append(logWriters, xLogBuffer)

	xLog.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	xLog.SetOutput(io.MultiWriter(logWriters...))
}

// wordSepNormalizeFunc all options are lowercase, so
// ... lowercase they shall be
func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	from := []string{"-", "_"}
	to := "."
	if nil != f {
		for _, sep := range from {
			name = strings.Replace(name, sep, to, -1)
		}
	}
	return pflag.NormalizedName(strings.ToLower(name))
}

/* standard flags */

var FlagOutfile string
var FlagOrganization string
var FlagHelp bool
var FlagQuiet bool
var FlagVerbose bool
var FlagDebug bool

/* program specific flags */

func initFlags() {
	var err error

	nFlags = pflag.NewFlagSet("default", pflag.ContinueOnError)
	nFlags.SetNormalizeFunc(wordSepNormalizeFunc)

	nFlags.StringVarP(&FlagOutfile, "outfile", "",
		"temp.html", "html report file to append")
	err = nFlags.MarkHidden("outfile")
	if nil != err {
		xLog.Printf("huh? Failure to mark %s hidden?", "FlagOutfile")
	}

	nFlags.StringVarP(&FlagOrganization, "organization", "",
		"P3ID Technologies", "organization for email sends")
	// nFlags.MarkHidden("organization") // don't show this in usage

	nFlags.BoolVarP(&FlagHelp, "help", "h",
		false, "Display help message and usage information")
	nFlags.BoolVarP(&FlagDebug, "debug", "d",
		true, "Enable additional informational and operational logging output for debug purposes")
	nFlags.BoolVarP(&FlagQuiet, "quiet", "q",
		false, "Suppress output to stdout and stderr (output still goes to logfile)")
	nFlags.BoolVarP(&FlagVerbose, "verbose", "v",
		false, "Supply additional run messages; use --debug for more information")

	// Fetch and load the program flags
	err = nFlags.Parse(os.Args[1:])
	if nil != err {
		_, _ = fmt.Fprintf(os.Stderr, "\n%s\n", nFlags.FlagUsagesWrapped(75))
		xLog.Fatalf("\nerror parsing flags because: %s\n%s %s\n%s\n\t%v\n",
			err.Error(),
			"  common issue: 2 hyphens for long-form arguments,",
			"  1 hyphen for short-form argument",
			"  Program arguments are: ",
			os.Args)
	}

	// do quietness setup first
	// only write to logfile not stderr
	// for debug and verbose messages
	if FlagQuiet {
		xLog.SetOutput(xLogBuffer)
		// messages only to logfile, not stderr
	}

	if FlagDebug || FlagVerbose {
		xLog.Println("/**** program flags ********/")
		nFlags.VisitAll(logFlag)
		xLog.Println("/**** end program flags ****/")
	}

	// next simplest
	if FlagHelp {
		_, thisCmd := filepath.Split(os.Args[0])
		_, _ = fmt.Fprint(os.Stdout, "\n", "usage for ", thisCmd, ":\n")
		_, _ = fmt.Fprintf(os.Stdout, "%s\n", nFlags.FlagUsagesWrapped(75))
		UsageMessage()
		os.Exit(0)
	}

	if FlagVerbose {
		xLog.Print("Verbose mode active (all debug and informative messages)")
	}

}

func logFlag(flag *pflag.Flag) {
	xLog.Printf(" flag '%s' has value '%v' with default is '%s'",
		flag.Name, flag.Value, flag.DefValue)
}

// UsageMessage - describe capabilities and extended usage notes
func UsageMessage() {
	var sb strings.Builder
	sb.WriteString(" SIMPLES3")
	sb.WriteString(" Simple utility to demonstrate:")
	sb.WriteString("\t 1.\tAuthorization to S3\n")
	sb.WriteString("\t 2.\tReading the available files in the pseudo-directory\n")
	sb.WriteString("\t 3.\tDownloading one the file present in the bucket")

	_, err := fmt.Fprintf(os.Stdout, "\n%s\n", sb.String())
	if nil != err {
		xLog.Printf("error writing UsageMessage to stdout because %s\n",
			err.Error())
		myFatal()
	}
}
