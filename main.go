package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	if err := run(os.Stdin); err != nil {
		log.Fatal(err)
	}
}

func run(r io.Reader) error {
	log.Println("starting recorder")
	dec := json.NewDecoder(r)
	var msg Message
	for {
		err := dec.Decode(&msg)
		if err == io.EOF {
			log.Println("reached end-of-file, exiting")
			return nil
		} else if err != nil {
			return err
		}
		fmt.Printf("%#v\n", msg)
	}
}

type Message struct {
	Timestamp float64 `json:"timestamp"`
	StationID string  `json:"station_id"`
	Channel   int64   `json:"channel"`
	Frequency float64 `json:"freq"`
	Level     float64 `json:"level"`
	Error     int64   `json:"error"`
	Mode      string  `json:"mode"`
	Label     string  `json:"label"`
	BlockID   string  `json:"block_id"`
	//Ack string|false
	Tail     string          `json:"tail"`
	Flight   string          `json:"flight"`
	MsgNo    string          `json:"msgno"`
	Text     string          `json:"text"`
	End      bool            `json:"end"`
	DepA     string          `json:"depa"`
	DstA     string          `json:"dsta"`
	ETA      string          `json:"eta"`
	GtOut    string          `json:"gtout"`
	GtIn     string          `json:"gtin"`
	WlOff    string          `json:"wloff"`
	Won      string          `json:"wlin"` // n.b. this appears to be mislabeled in the source: https://github.com/TLeconte/acarsdec/blob/7920079b8e005c6c798bd478a513211daf9bbd25/output.c#L293-L294
	Sublabel string          `json:"sublabel"`
	MFI      string          `json:"mfi"`
	AssStat  string          `json:"assstat"`
	Libacars json.RawMessage `json:"libacars"`
	App      struct {
		Name string `json:"name"`
		Ver  string `json:"ver"`
	} `json:"app"`
}
