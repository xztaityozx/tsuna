package process

import (
	"fmt"
	"github.com/mitchellh/go-ps"
	"github.com/samber/lo"
	"os"
)

// FindByName はプロセス名でプロセスを探して返す
func FindByName(name string) (*os.Process, error) {
	processes, err := ps.Processes()
	if err != nil {
		return nil, err
	}

	if proc, ok := lo.Find(processes, func(item ps.Process) bool {
		return item.Executable() == name
	}); ok {
		return os.FindProcess(proc.Pid())
	} else {
		return nil, fmt.Errorf("名前が %s なプロセスは見つかりませんでした", name)
	}
}
