package main

import (
    "testing"
)

func TestCleanInput(t *testing.T) {
    cases := []struct {
        input           string
        expected        []string
    }{
        {
            input:      " hello world ",
            expected:   []string{"hello", "world"},
        },
        {
            input:      " cool test    for cool people  ",
            expected:   []string{"cool", "test", "for", "cool", "people"},
        },
        {
            input:      "single",
            expected:   []string{"single"},
        },
    }

    for _, c := range cases {
        actual := cleanInput(c.input)
        if len(actual) != len(c.expected) {
            t.Errorf("actual length: %v != expected length: %v", len(actual), len(c.expected))
        }

        for i := range actual {
            word := actual[i]
            expectedWord := c.expected[i]
            if word != expectedWord {
                t.Errorf("word: %v != expectedWord: %v", word, expectedWord)
            }
        }
    }
}
