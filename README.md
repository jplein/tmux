# Tmux

A go library for working with tmux

## Running a tmux shell command

To run a tmux shell command and collect its output:

```
output, err = tmux.Command("list-sessions", "-F", "#{session_name}");
```

## Using a tmux Runner

Commands are easy to work with, but you can get better performance by using a Runner.

A Runner starts a `tmux -C` process, and writes to it to send commands, and reads from it to capture the output of the commands.

To initialize, use something like this:

```
var r *tmux.Runner = &tmux.Runner{}

if err = r.Init(); err != nil {
  return err
}
```

When you're done with the Runner, close it, to stop the `tmux -C` process:

```
if err = r.Close(); err != nil {
  os.Stderr.Write([]byte(fmt.Sprintf("Error closing tmux runner: %s", err.Error())))
}
```

To run a command and collect its output, use the Run method:

```
output, err = r.Run("list-sessions -F '#{session_name}'")
```

The Runner type also has many other functions for tasks like starting a new
tmux session, getting the active window, etc. For details, run `go doc tmux.Runner`.