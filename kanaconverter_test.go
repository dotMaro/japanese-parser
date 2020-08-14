package main

import "testing"

func TestKataToHira(t *testing.T) {
	tc := []struct {
		in     string
		expect string
	}{
		{"あ", "あ"},
		{"ア", "あ"},
		{"ハイ", "はい"},
		{"はイ", "はい"},
		{"ハい", "はい"},
		{"貫ク", "貫く"},
		{"アカタハナマラサ", "あかたはなまらさ"},
		{"ジポヅブプ", "じぽづぶぷ"},
	}

	for _, c := range tc {
		res := KataToHira(c.in)
		if res != c.expect {
			t.Errorf("%v should convert to %v, not %v", c.in, c.expect, res)
		}
	}
}
