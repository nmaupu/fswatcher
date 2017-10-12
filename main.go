package main

import (
	"bytes"
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/rjeczalik/notify"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Constants to get app name and description
const (
	AppName      = "fswatcher"
	AppDesc      = "Watch filesystem events using inotify"
	AppDescOther = "Use %f for file/directory name, %e for event name in command arguments"
)

var (
	// AppVersion can be set at release
	AppVersion string
	// List of all events with their bare name and their associated description
	eventsList = map[string]string{
		"IN_ACCESS":        "File was accessed (read).",
		"IN_MODIFY":        "File was modified.",
		"IN_ATTRIB":        "Metadata changed, e.g., permissions, timestamps, extended attributes, link count (since Linux 2.6.25), UID, GID, etc. ",
		"IN_CLOSE_WRITE":   "File opened for writing was closed.",
		"IN_CLOSE_NOWRITE": "File not opened for writing was closed.",
		"IN_OPEN":          "File was opened.",
		"IN_MOVED_FROM":    "File moved out of watched directory.",
		"IN_MOVED_TO":      "File moved into watched directory.",
		"IN_CREATE":        "File/directory created in watched directory.",
		"IN_DELETE":        "File/directory deleted from watched directory.",
		"IN_DELETE_SELF":   "Watched file/directory was itself deleted.",
		"IN_MOVE_SELF":     "Watched file/directory was itself moved.",
	}
)

func main() {
	if AppVersion == "" {
		AppVersion = "master"
	}

	err := processOpts(executeCommand)
	if err != nil {
		log.Fatal(err)
	}
}

func getAppDesc() string {
	var buffer bytes.Buffer

	for k, v := range eventsList {
		buffer.WriteString(fmt.Sprintf("%s : %s\n", k, v))
	}

	return fmt.Sprintf("%s\n\nAvailable events:\n%s\n\n%s", AppDesc, buffer.String(), AppDescOther)
}

// Process all opts given to the program
func processOpts(f func(c string, e []string, t string) error) error {
	var err error
	app := cli.App(AppName, getAppDesc())
	app.Version("v version", fmt.Sprintf("%s version %s", AppName, AppVersion))

	var (
		cmd    = app.StringOpt("c command cmd", "", "Command to execute when receiving a matching event")
		events = app.StringsOpt("e event", nil, "Events to listen to")
		target = app.StringArg("TARGET", "", "file or directory to listen for provided events")
	)

	app.Action = func() {
		var msgs []string

		if *cmd == "" {
			msgs = append(msgs, "Please provide command to execute")
		}
		if len(*events) == 0 {
			msgs = append(msgs, "Please provide inotify events")
		}
		if len(msgs) > 0 {
			for i := range msgs {
				log.Println(msgs[i])
			}
			log.Fatal(fmt.Sprintf("See %s --help for usage.", AppName))
		} else {
			err = f(*cmd, *events, *target)
		}
	}

	app.Run(os.Args)
	return err
}

// Execute program and loop for more events
func executeCommand(cmd string, events []string, target string) error {
	done := make(chan bool)

	// Creating real notify.Event (which are integers) from given strings
	var evts []notify.Event
	for i := range events {
		evts = append(evts, getEventFromString(events[i]))
	}

	c := make(chan notify.EventInfo, 1)
	log.Printf("Listening events %s from %s", events, target)
	err := notify.Watch(target, c, evts...)
	if err != nil {
		return err
	}
	defer notify.Stop(c)

	go func() {
		for {
			select {
			case ei := <-c:
				go processEvent(cmd, ei)
			}
		}
	}()

	<-done
	return nil
}

func processEvent(cmd string, e notify.EventInfo) error {
	// Replacing %f and %e from initial command
	s := strings.Split(cmd, " ")
	var args []string

	ref, _ := regexp.Compile("%f")
	ree, _ := regexp.Compile("%e")
	for _, v := range s[1:] {
		val := string(ref.ReplaceAll([]byte(v), []byte(e.Path())))
		val = string(ree.ReplaceAll([]byte(val), []byte(getStringFromEvent(e.Event()))))
		args = append(args, val)
	}

	log.Printf("Event received %s - executing %s %s", e, s[0], args)
	_, err := exec.Command(s[0], args...).Output()
	if err != nil {
		return err
	}

	return nil
}

func getEventFromString(event string) notify.Event {
	switch event {
	case "IN_ACCESS":
		return notify.InAccess
	case "IN_MODIFY":
		return notify.InModify
	case "IN_ATTRIB":
		return notify.InAttrib
	case "IN_CLOSE_WRITE":
		return notify.InCloseWrite
	case "IN_CLOSE_NOWRITE":
		return notify.InCloseNowrite
	case "IN_OPEN":
		return notify.InOpen
	case "IN_MOVED_FROM":
		return notify.InMovedFrom
	case "IN_MOVED_TO":
		return notify.InMovedTo
	case "IN_CREATE":
		return notify.InCreate
	case "IN_DELETE":
		return notify.InDelete
	case "IN_DELETE_SELF":
		return notify.InDeleteSelf
	case "IN_MOVE_SELF":
		return notify.InMoveSelf
	}

	return 0
}

func getStringFromEvent(event notify.Event) string {
	switch event {
	case notify.InAccess:
		return "IN_ACCESS"
	case notify.InModify:
		return "IN_MODIFY"
	case notify.InAttrib:
		return "IN_ATTRIB"
	case notify.InCloseWrite:
		return "IN_CLOSE_WRITE"
	case notify.InCloseNowrite:
		return "IN_CLOSE_NOWRITE"
	case notify.InOpen:
		return "IN_OPEN"
	case notify.InMovedFrom:
		return "IN_MOVED_FROM"
	case notify.InMovedTo:
		return "IN_MOVED_TO"
	case notify.InCreate:
		return "IN_CREATE"
	case notify.InDelete:
		return "IN_DELETE"
	case notify.InDeleteSelf:
		return "IN_DELETE_SELF"
	case notify.InMoveSelf:
		return "IN_MOVE_SELF"
	}

	return ""
}
