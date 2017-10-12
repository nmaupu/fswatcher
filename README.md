[![Build Status](https://travis-ci.org/nmaupu/fswatcher.svg?branch=master)](https://travis-ci.org/nmaupu/fswatcher)
[![Go Report Card](https://goreportcard.com/badge/github.com/nmaupu/fswatcher)](https://goreportcard.com/report/github.com/nmaupu/fswatcher)

# What is fswatcher

_fswatcher_ is a tool to listen for filesystem changes using inotify events.
See for more info : http://man7.org/linux/man-pages/man7/inotify.7.html

NOTE : It is currently only working for **unix based** OS !

# Building

```
glide install && make
```

# Usage

```
root@22f4580a3b27:/go/src/github.com/nmaupu/fswatcher# ./bin/fswatcher  --help

Usage: fswatcher [OPTIONS] TARGET

Watch filesystem events using inotify

Available events:
IN_ACCESS : File was accessed (read).
IN_MODIFY : File was modified.
IN_CLOSE_WRITE : File opened for writing was closed.
IN_OPEN : File was opened.
IN_MOVED_FROM : File moved out of watched directory.
IN_DELETE_SELF : Watched file/directory was itself deleted.
IN_ATTRIB : Metadata changed, e.g., permissions, timestamps, extended attributes, link count (since Linux 2.6.25), UID, GID, etc.
IN_CLOSE_NOWRITE : File not opened for writing was closed.
IN_MOVED_TO : File moved into watched directory.
IN_CREATE : File/directory created in watched directory.
IN_DELETE : File/directory deleted from watched directory.
IN_MOVE_SELF : Watched file/directory was itself moved.


Use %f for file/directory name, %e for event name in command arguments

Arguments:
  TARGET=""    file or directory to listen for provided events

Options:
  -v, --version             Show the version and exit
  -c, --command, --cmd=""   Command to execute when receiving a matching event
  -e, --event=[]            Events to listen to
```

## Usage example

```
# ./bin/fswatcher -c "/tmp/script.sh %f" -e IN_CLOSE_WRITE -e IN_MOVED_TO /tmp/test
2017/10/12 21:38:01 Listening events [IN_CLOSE_WRITE IN_MOVED_TO] from /tmp/test
```

In the command parameters, it's possible to get information on the event :
- `%e` can be used to get the event name
- `%f` can be used to get the file name/directory name source of the event

Now, pop an event !
```
echo 42 >> /tmp/test/toto
```

We can now see in the console :
```
2017/10/12 21:38:33 Event received notify.InCloseWrite: "/tmp/test/toto" - executing /tmp/script.sh [/tmp/test/toto]
```

# Dependencies

- https://github.com/rjeczalik/notify
- https://github.com/jawher/mow.cli
