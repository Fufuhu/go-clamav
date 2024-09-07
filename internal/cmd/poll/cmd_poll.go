package poll

import (
	"fmt"
	"github.com/spf13/cobra"
)

type CommandPoll struct {
}

func (p *CommandPoll) Run(cmd *cobra.Command, args []string) {
	fmt.Println("poll called")
}
