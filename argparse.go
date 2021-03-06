package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

/**
 * Package name: main
 * Project name: timesheet
 * Created by: Praveen Premaratne
 * Created on: 29/02/2020 17:52
 */

func (app *App) Parser() {
	app.Started = app.getDateTime()

	flag.BoolVar(&app.Help, "h", false,
		"HELP: This tool can be used to log time spent on a specific Jira ticket on a project.")
	flag.StringVar(&app.Ticket, "r", "",
		"REQUIRED: Jira ticket reference. E.g. DDSP-4")
	flag.StringVar(&app.TimeSpent, "t", "",
		"REQUIRED: The time spent as days (#d), hours (#h), or minutes (#m or #). E.g. 8h")
	flag.StringVar(&app.Started, "d", "",
		fmt.Sprintf("Default %s. The date on which the worklog effort was started in full date (YYYY-MM-DD) or"+
			" relative date (-N) format. eg: 2006-01-02 or -1.", app.getDate()))
	flag.StringVar(&app.Comment, "m", "",
		"OPTIONAL: A comment about the worklog")
	flag.StringVar(&app.Encode, "e", "",
		"HELP: Base64 encode the given credentials."+
			" Format: email:token;domain. e.g. example@example.com:abcThisIsFake;xyz.atlassian.net")
	flag.BoolVar(&app.TimeRemaining, "remaining", false,
		"HELP: Print how many hour can be book for the current day. -d is also available")
	flag.BoolVar(&app.History, "history", false,
		"HELP: Print the timesheet of the day -d is also available to change the week")
	flag.BoolVar(&app.PrintWeek, "week", false,
		"HELP: Print timesheet of the current week. -d is also available to change the week")
	flag.BoolVar(&app.PrintMonth, "month", false,
		"HELP: Print timesheet of the current month. -d is also available to change the week")
	flag.BoolVar(&app.Version, "v", false, "Print application version")
	flag.Parse()
	app.validate()
}

func (app *App) validate() {

	if len(os.Args[1:]) < 1 {
		fmt.Printf("no arguments are given\n\n")
		app.usage()
		os.Exit(1)
	}

	if app.Started != "" {
		if len(RelativeDateFormat.FindStringSubmatch(app.Started)) > 0 {
			relativeTime, err := app.GetDateFromRelative()
			if err != nil {
				panic(err)
			}
			app.Started = fmt.Sprintf("%sT%s", relativeTime.Format(YmdFormat), app.getTimeFixed())
		} else if DateFormat.MatchString(app.Started) {
			app.Started = fmt.Sprintf("%sT%s", app.Started, app.getTimeFixed())
		} else {
			panic("provided date didn't match expected format. try -h for help")
		}
	} else {
		app.Started = app.getDateTime()
	}

	if app.Help {
		app.usage()
		os.Exit(0)
	}

	if app.Version {
		fmt.Printf(SIGNATURE)
		fmt.Println("Version:", VERSION)
		os.Exit(0)
	}

	if app.TimeRemaining || app.PrintWeek || app.History || app.PrintMonth {
		return
	}

	if app.Encode != "" {
		app.CredentialEncode()
		os.Exit(0)
	}

	if app.Ticket == "" {
		panic(errors.New("please provide a ticket reference. -r"))
	}

	if app.TimeSpent == "" {
		panic(errors.New("no time given. -t"))
	}
}

func (app *App) usage() {
	fmt.Printf("timesheet (-r -t [-d] [-m]] [[-h] [-e] [-d]) ([-remaining] [-history])\n")
	flag.PrintDefaults()
	fmt.Printf("Example:\n" +
		"\ttimesheet -r DDSP-XXXX -t 8h -m \"Jenkins pipeline completed\"\n" +
		"\ttimesheet -r DDSP-XXXX -t 1h -m \"Investigated possible solutions\" -d 2020-03-05\n" +
		"\ttimesheet -remaining\n" +
		"\ttimesheet -remaining -d 2020-03-05\n" +
		"\ttimesheet -history\n" +
		"\ttimesheet -history -d -1\n")
}
