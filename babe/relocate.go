package babe

import (
	"errors"
	"github.com/mrnavastar/assist/bytes"
	"slices"
	"strings"
)

func ParseRelocations(relocations []string) [][]string {
	var parsedRelocations [][]string
	for _, relocation := range relocations {
		parsedRelocations = append(parsedRelocations, strings.Split(strings.ReplaceAll(relocation, ".", "/"), ":"))
	}
	return parsedRelocations
}

func ParseRelocation(relocation string) [][]string {
	return ParseRelocations([]string{relocation})
}

// TODO: Dont modify regular strings - actually test that this works
// TODO: Make this function less hacky and gross
func RelocateClass(class *Class, relocations [][]string) bool {
	modified := false
	var skip []uint16

	for _, constant := range class.ConstantPool {
		switch info := constant.(type) {
		case *StringInfo:
			skip = append(skip, info.StringIndex-1)
		}
	}

	for i, constant := range class.ConstantPool {
		switch info := constant.(type) {
		case *Utf8Info:
			if slices.Contains(skip, uint16(i)) {
				continue
			}

			for _, relocation := range relocations {
				s := info.String()
				if strings.Contains(s, relocation[0]) {
					info.Set(strings.ReplaceAll(s, relocation[0], relocation[1]))
					modified = true
					break
				}
			}
		}
	}
	return modified
}

func RelocateJar(filename string, relocations [][]string) error {
	return ModifyJar(filename, func(member *JarMember) error {
		for _, relocation := range relocations {
			if strings.Contains(member.Name, relocation[0]) {
				member.Name = strings.ReplaceAll(member.Name, relocation[0], relocation[1])
				break
			}
		}

		class, err := member.GetAsClass()
		if err != nil {
			if errors.Is(err, ErrNotClass) {
				return nil
			}
			return err
		}

		if RelocateClass(&class, relocations) {
			member.Buffer = bytes.NewBuffer()
			class.Write(member.Buffer.Data)
		}
		return nil
	})
}
