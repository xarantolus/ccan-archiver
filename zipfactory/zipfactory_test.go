package zipfactory

import "testing"

func TestGenerateFilename(t *testing.T) {

	table := map[string]string{
		"https://www.clonkx.de/endeavour/Freeware.c4k": "c4k",
	}

	for key, value := range table {

		if res := getURLExtension(key); res != value {
			t.Errorf("`%s`=`%s`, expected `%s`", key, res, value)
		}
	}
}
