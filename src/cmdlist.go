package main


import (
    "fmt"
    "os"
    "flag"
    "time"
)


// ---------------------------------------------------------------------------------------------------
// List commands

type ListCommand struct {
    Ctx             *Context
    setName         *string
    beforeDate      *string
    afterDate       *string
    flagDetailed    *bool
    firstResult     *int
    maxResults      *int
}


// Attempt to parse a date string and return it as a heap allocated time.Time.
// If the string is empty, returns nil.  If there was an error parsing the string,
// the program dies
func parseDateString(dateString string) *time.Time {
    if (dateString != "") {
        parsedTime, err := time.ParseInLocation(DateFormat, dateString, time.Local)
        if (err != nil) {
            die("Invalid date: " + err.Error())
        }

        heapAllocTime := new(time.Time)
        *heapAllocTime = parsedTime

        return heapAllocTime
    } else {
        return nil
    }

}



// Get list identifier arguments
func (lc *ListCommand) genListIdentifierArgsFromCommandLine() ListIdentifierArgs {
    set := *(lc.setName)
    if (set == "") {
        set = lc.Ctx.Provider.Set
    }

    args := ListIdentifierArgs{
        Set: set,
        From: parseDateString(*(lc.afterDate)),
        Until: parseDateString(*(lc.beforeDate)),
    }

    return args
}


// List the identifiers from a provider
func (lc *ListCommand) listIdentifiers() {
    var deletedCount int = 0

    args := lc.genListIdentifierArgsFromCommandLine()

    lc.Ctx.Session.ListIdentifiers(args, *(lc.firstResult), *(lc.maxResults), func(res ListIdentifierResult) bool {
        if (res.Deleted) {
            deletedCount += 1
        } else {
            fmt.Printf("%s\n", res.Identifier)
        }
        return true
    })

    if (deletedCount > 0) {
        fmt.Fprintf(os.Stderr, "oaipmh: %d deleted record(s) not displayed.\n", deletedCount)
    }

}

// List the identifiers in detail from a provider
func (lc *ListCommand) listIdentifiersInDetail() {
    args := lc.genListIdentifierArgsFromCommandLine()

    lc.Ctx.Session.ListIdentifiers(args, *(lc.firstResult), *(lc.maxResults), func(res ListIdentifierResult) bool {
        if (res.Deleted) {
            fmt.Printf("D ")
        } else {
            fmt.Printf(". ")
        }
        fmt.Printf("%-20s ", res.Sets[0])
        fmt.Printf("%-20s  ", res.Datestamp)
        fmt.Printf("%s\n", res.Identifier)
        return true
    })
}

func (lc *ListCommand) Flags(fs *flag.FlagSet) *flag.FlagSet {
    lc.setName = fs.String("s", "", "The set to retrieve")
    lc.beforeDate = fs.String("B", "", "List metadata records that have been updated before this date (YYYY-MM-DD).")
    lc.afterDate = fs.String("A", "", "List metadata records that have been updated after this date (YYYY-MM-DD).")
    lc.flagDetailed = fs.Bool("l", false, "List metadata in detail.")
    lc.firstResult = fs.Int("f", 0, "The first result to return.")
    lc.maxResults = fs.Int("c", 100000, "Maximum number of results to return.")

    return fs
}

func (lc *ListCommand) Run(args []string) {
    if *(lc.flagDetailed) {
        lc.listIdentifiersInDetail()
    } else {
        lc.listIdentifiers()
    }
}
