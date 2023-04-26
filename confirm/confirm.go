// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package confirm

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/manifoldco/promptui"
)

// Confirm is send the prompt and get result
func Confirm(prompt string) (bool, error) {
	var yesRx = regexp.MustCompile("^(?:y(?:es)?)$")
	var noRx = regexp.MustCompile("^(?:n(?:o)?)$")
	promptLabel := fmt.Sprintf("%s Yes [y/yes], No [n/no]", prompt)
	validate := func(input string) error {
		input = strings.ToLower(input)
		if !yesRx.MatchString(input) && !noRx.MatchString(input) {
			return errors.New("invalid input, please enter 'y', 'yes', 'n', or 'no'")
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ .  }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | green }} ",
	}

	promptObj := promptui.Prompt{
		Label:     promptLabel,
		Templates: templates,
		Validate:  validate,
	}
	result, err := promptObj.Run()
	if err != nil {
		return false, err
	}
	result = strings.ToLower(result)
	if yesRx.MatchString(result) {
		return true, nil
	}
	return false, nil
}
