package poll

import (
	"fmt"
	"github.com/Fufuhu/go-clamav/config"
	"github.com/Fufuhu/go-clamav/internal/cmd"
	"github.com/spf13/cobra"
)

type CommandPoll struct {
}

func (p *CommandPoll) Run(cmd *cobra.Command, args []string) {
	fmt.Println("poll called")
}

func NewCommand(cfg config.Configuration) cmd.CommandInterface {
	return &CommandPoll{}
}
