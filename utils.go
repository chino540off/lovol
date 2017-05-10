package main

import (
  "fmt"
  "log"
  "os/exec"
  "strings"
  "syscall"
)

func printCommand(cmd *exec.Cmd) {
  log.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
  if err != nil {
    log.Printf("==> Error: %s\n", err.Error())
  }
}

func printOutput(outs []byte) {
  if len(outs) > 0 {
    log.Printf("==> Output: %s\n", string(outs))
  }
}

func execute(cmd string, args map[string]string) error {
  var _args []string

  for k, v := range args {
    _args = append(_args, fmt.Sprintf("%s=%s", k, v))
  }
  return _execute(cmd, _args...)
}

func _execute(cmd string, args...string) error {
  command := exec.Command(cmd, args...)
  printCommand(command)

  var waitStatus syscall.WaitStatus

  if err := command.Run(); err != nil {
    printError(err)
    // Did the command fail because of an unsuccessful exit code
    if exitError, ok := err.(*exec.ExitError); ok {
      waitStatus = exitError.Sys().(syscall.WaitStatus)
      printOutput([]byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
    }
    return err
  } else {
    // Command was successful
    waitStatus = command.ProcessState.Sys().(syscall.WaitStatus)
    printOutput([]byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
  }

  return nil
}
