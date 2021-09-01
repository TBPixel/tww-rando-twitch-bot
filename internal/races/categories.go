package races

import "strings"

const (
	Standard   = "Standard Race"
	SpoilerLog = "Spoiler Log Race"
	Fun        = "Fun Race"
	Custom     = "Custom"
)

type ExamplePerma struct {
	Preset      string
	Perma       string
	Description string
}

var examplePermas = []ExamplePerma{
	{
		Preset:      "beginner",
		Description: "Start with Sword, 2 DRM, Puzzle Secret Caves, Free Gifts, Mail & Misc.",
		Perma:       "MS45LjAAQQAFAwIADzDAAYAcQIFBATAA",
	},
	{
		Preset:      "co-op",
		Description: "5 DRM, Puzzle Secret Caves, Free Gifts, Tingle Chests, Short Sidequests, Mail, Island Puzzles, Submarines & Misc.",
		Perma:       "MS45LjAAQQAVCyYAD3DABAAAAAAAAQAA",
	},
	{
		Preset:      "s1",
		Description: "Start with Sword, 3 DRM, Puzzle Secret Caves, Great Fairies, Free Gifts, Tingle Chests, Short Sidequests, Mail & Misc.",
		Perma:       "MS45LjAAQQAXAwYCDxDAAgAAAAAAAQAA",
	},
	{
		Preset:      "s3",
		Description: "Start with Sword, 4 DRM, Puzzle Secret Caves, Great Fairies, Free Gifts, Tingle Chests, Short Sidequests, Mail.",
		Perma:       "MS45LjAAQQAXAwQATjDAAwgAAAAAAQAA",
	},
	{
		Preset:      "s4",
		Description: "3 DRM, Puzzle Secret Caves, Island Puzzles, Free Gifts, Mail, Submarines & Misc.",
		Perma:       "MS45LjAAQQAFCyIAD3DAAgAAAAAAAQAA",
	},
	{
		Preset:      "allsanity",
		Description: "Everything enabled.",
		Perma:       "MS45LjAAQQD//3+CD3BABQAAAAAAAAAA",
	},
	{
		Preset:      "preset-a",
		Description: "Long Sidequests.",
		Perma:       "MS45LjAAQQA3AyYCD1DAAgAAAAAAAAAA",
	},
	{
		Preset:      "preset-b",
		Description: "Triforce Charts, Big Octos and Gunboats.",
		Perma:       "MS45LjAAQQAXYyaCD1DAAgAAAAAAAAAA",
	},
	{
		Preset:      "preset-c",
		Description: "Swordless.",
		Perma:       "MS45LjAAQQAXAyYCD5DAAgAAAAAAAAAA",
	},
	{
		Preset:      "preset-d",
		Description: "Lookout Platforms and Rafts.",
		Perma:       "MS45LjAAQQAXByYCD1DAAgAAAAAAAAAA",
	},
	{
		Preset:      "preset-e",
		Description: "4 Dungeon Race Mode and Key-Lunacy.",
		Perma:       "MS45LjAAQQAXA2YCD1DAAwAAAAAAAAAA",
	},
	{
		Preset:      "preset-f",
		Description: "Combat Secret Caves, Submarines.",
		Perma:       "MS45LjAAQQAfCyYCD1DAAgAAAAAAAAAA",
	},
}

func Presets() []string {
	var presets []string
	for _, e := range examplePermas {
		presets = append(presets, e.Preset)
	}

	return presets
}

func ExamplePermaByPreset(preset string) *ExamplePerma {
	for _, e := range examplePermas {
		if e.Preset != strings.ToLower(preset) {
			continue
		}

		return &e
	}

	return nil
}
