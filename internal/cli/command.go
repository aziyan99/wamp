package cli

import (
	"fmt"
	"os"
	"strings"
)

type Command struct {
	Name        string
	Short       string
	Long        string
	Run         func(cmd *Command, args []string)
	SubCommands map[string]*Command
	Flags       map[string]*string
}

func NewCommand(name, short, long string, run func(cmd *Command, args []string)) *Command {
	cmd := &Command{
		Name:        name,
		Short:       short,
		Long:        long,
		Run:         run,
		SubCommands: make(map[string]*Command),
		Flags:       make(map[string]*string),
	}

	return cmd
}

func (c *Command) AddCommand(cmd *Command) {
	c.SubCommands[cmd.Name] = cmd
}

func (c *Command) AddCommands(cmds ...*Command) {
	for i := 0; i < len(cmds); i++ {
		c.SubCommands[cmds[i].Name] = cmds[i]
	}
}

func (c *Command) AddFlag(name, short, defaultValue, description string) {
	c.Flags[name] = &defaultValue
}

func (c *Command) Execute() {
	args := os.Args[1:]

	// Check for subcommands
	if len(args) > 0 {
		if cmd, ok := c.SubCommands[args[0]]; ok {
			// Re-slice os.Args to pass to subcommand
			os.Args = append([]string{os.Args[0]}, args[1:]...)
			cmd.Execute()
			return
		}
	}

	// Parse flags
	parsedArgs := []string{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if len(arg) > 2 && arg[:2] == "--" {
			flagName := arg[2:]
			if _, ok := c.Flags[flagName]; ok {
				if flagName == "ssl" {
					*c.Flags[flagName] = "true"
				} else if i+1 < len(args) {
					*c.Flags[flagName] = args[i+1]
					i++ // Skip the flag value
				}
			}
		} else if len(arg) > 1 && arg[0] == '-' {
			flagName := arg[1:]
			if _, ok := c.Flags[flagName]; ok {
				if flagName == "s" {
					*c.Flags["ssl"] = "true"
				} else if i+1 < len(args) {
					*c.Flags[flagName] = args[i+1]
					i++ // Skip the flag value
				}
			}
		} else {
			parsedArgs = append(parsedArgs, arg)
		}
	}

	c.Run(c, parsedArgs)

	os.Exit(0)
}

func AddHelpCommand(cmd *Command) {
	helpCmd := NewCommand(
		"help",
		"Prints help information",
		"This command prints help information for a specific command.",
		func(c *Command, args []string) {
			if len(args) == 0 {
				fmt.Println(cmd.Long)
				fmt.Println("\nAvailable Commands:")

				// Find the longest command name
				longestName := 0
				for _, subCmd := range cmd.SubCommands {
					if len(subCmd.Name) > longestName {
						longestName = len(subCmd.Name)
					}
					for _, subSubCmd := range subCmd.SubCommands {
						if len(subCmd.Name+":"+subSubCmd.Name) > longestName {
							longestName = len(subCmd.Name + ":" + subSubCmd.Name)
						}
					}
				}

				for _, subCmd := range cmd.SubCommands {
					padding := longestName - len(subCmd.Name)
					if padding < 0 {
						padding = 0
					}
					fmt.Printf("  %s%s  %s\n", subCmd.Name, strings.Repeat(" ", padding), subCmd.Short)
					if len(subCmd.SubCommands) > 0 {
						for _, subSubCmd := range subCmd.SubCommands {
							name := subCmd.Name + " " + subSubCmd.Name
							padding := longestName - len(name)
							if padding < 0 {
								padding = 0
							}
							fmt.Printf("    %s%s  %s\n", name, strings.Repeat(" ", padding), subSubCmd.Short)
						}
						fmt.Println()
					}
				}
				return
			}

			// try to find the command
			var cmdToHelp *Command
			for _, subCmd := range cmd.SubCommands {
				if subCmd.Name == args[0] {
					cmdToHelp = subCmd
					break
				}
				for _, subSubCmd := range subCmd.SubCommands {
					if subSubCmd.Name == args[0] {
						cmdToHelp = subSubCmd
						break
					}
				}
			}

			if cmdToHelp == nil {
				fmt.Printf("Unknown command: %s\n", args[0])
				return
			}

			fmt.Println(cmdToHelp.Long)
			fmt.Println("\nUsage:")
			fmt.Printf("  %s %s [flags]\n", cmd.Name, cmdToHelp.Name)

			if len(cmdToHelp.Flags) > 0 {
				fmt.Println("\nFlags:")
				// Find the longest flag name
				longestFlagName := 0
				for name := range cmdToHelp.Flags {
					if len(name) > longestFlagName {
						longestFlagName = len(name)
					}
				}
				for name, flag := range cmdToHelp.Flags {
					padding := longestFlagName - len(name)
					if padding < 0 {
						padding = 0
					}
					fmt.Printf("  --%s%s  %s (default: \"%s\")\n", name, strings.Repeat(" ", padding), "", *flag)
				}
			}
		},
	)
	cmd.AddCommand(helpCmd)
}

func AddHelpCommands(cmds ...*Command) {
	for i := 0; i < len(cmds); i++ {
		AddHelpCommand(cmds[i])
	}
}
