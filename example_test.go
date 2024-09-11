package freelist_test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Snawoot/freelist"
)

func Example() {
	const jsonStream = `
	{"Name": "Ed", "Text": "Knock knock."}
	{"Name": "Sam", "Text": "Who's there?"}
	{"Name": "Ed", "Text": "Go fmt."}
	{"Name": "Sam", "Text": "Go fmt who?"}
	{"Name": "Ed", "Text": "Go fmt yourself!"}
`
	type Message struct {
		Name, Text string
	}
	dec := json.NewDecoder(strings.NewReader(jsonStream))
	var messages []*Message
	var allocator freelist.Freelist[Message]
	for {
		m := allocator.Alloc()
		if err := dec.Decode(m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		messages = append(messages, m)
	}
	for _, m := range messages {
		fmt.Printf("%s: %s\n", m.Name, m.Text)
		allocator.Free(m)
	}
	messages = nil // make sure pointers released
}
