package encoder

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Danyiyk/FunCodex/internal/utils"
)

func getEncodedLetter(charMap map[string]string, letter string, bit int) string {
	for k, v := range charMap {
		if k == letter || v == letter {
			if bit == 0 {
				return k
			} else {
				return v
			}
		}
	}

	return ""
}

func Encode(str string, hiddenText string) string {
	splittedMainString := strings.Split(str, "")

	splittedString := strings.Split(hiddenText, "")
	returnValue := ""
	hiddenTextBinary := [][utils.BitsPerCharacter]int{}

	// Check if length match
	_, maxCharacters := utils.GetAvailableSpace(str)
	hiddenTextLength := len(splittedString)

	for i := range maxCharacters {
		if i >= hiddenTextLength {
			hiddenTextBinary = append(hiddenTextBinary, utils.IntegerToBinary(utils.EmptyCharBit, utils.BitsPerCharacter))
			break
		}

		letter := splittedString[i]
		found := utils.ArrayFind(utils.EncodableCharacters, strings.ToLower(letter))

		if found != -1 {
			hiddenTextBinary = append(hiddenTextBinary, utils.IntegerToBinary(found, utils.BitsPerCharacter))
		} else {
			fmt.Printf("[Warning] Invalid character in hidden string: %s\n", letter)
			returnValue = ""
			return returnValue
		}
	}

	slices.Reverse(hiddenTextBinary)

	if hiddenTextLength > maxCharacters {
		returnValue = ""
		return returnValue
	}

	var parsing [utils.BitsPerCharacter]int
	var parsingPosition int
	cursor := len(hiddenTextBinary)

	for i := range splittedMainString {
		letter := splittedMainString[i]
		bitValue, letterCharmap := utils.GetBit(letter)

		if cursor > 0 && parsingPosition < 1 {
			cursor--
			parsing = hiddenTextBinary[cursor]
			parsingPosition = utils.BitsPerCharacter
		}

		if bitValue == -1 || cursor == 0 && parsingPosition < 1 {
			returnValue = returnValue + letter
			continue
		}

		parsingPosition--
		encodedLetter := getEncodedLetter(letterCharmap, letter, parsing[utils.BitsPerCharacter-1-parsingPosition])

		returnValue = returnValue + encodedLetter
	}

	return returnValue
}
