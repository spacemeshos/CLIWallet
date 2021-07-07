package repl

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
)

var emptyComplete = func(prompt.Document) []prompt.Suggest { return []prompt.Suggest{} }

func runPrompt(executor func(string), completer func(prompt.Document) []prompt.Suggest,
	firstTime func(), length uint16) {
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(prefix),
		prompt.OptionPrefixTextColor(prompt.LightGray),
		prompt.OptionMaxSuggestion(length),
		prompt.OptionShowCompletionAtStart(),
	)
	firstTime()
	p.Run()
}

// executes prompt waiting for an input with y or n
func yesOrNoQuestion(msg string) string {
	var input string
	for {
		input = prompt.Input(prefix+msg,
			emptyComplete,
			prompt.OptionPrefixTextColor(prompt.LightGray))

		if input == "y" || input == "n" {
			break
		}

		fmt.Println("invalid command.")
	}

	return input
}

func multipleChoice(names []string) int {
	var input string
	if len(names) == 0 {
		return 0
	}
	for {
		for n, ac := range names {
			fmt.Println(n+1, ac)
		}
		input = prompt.Input(prefix,
			emptyComplete,
			prompt.OptionPrefixTextColor(prompt.LightGray))

		if num, err := strconv.Atoi(input); err == nil {
			if num > 0 && num <= len(names) {
				return num
			}
		}

		s := strings.TrimSpace(input)
		if s == "quit" || s == "exit" {
			fmt.Println("Bye!")
			os.Exit(0)
			return 0
		}

		fmt.Println("invalid command.")

	}
}

// executes prompt waiting an input not blank
func inputNotBlank(msg string) string {
	var input string
	for {
		input = prompt.Input(prefix+msg,
			emptyComplete,
			prompt.OptionPrefixTextColor(prompt.LightGray))

		if strings.TrimSpace(input) != "" {
			break
		}

		fmt.Println("please enter a value.")
	}

	return input
}
