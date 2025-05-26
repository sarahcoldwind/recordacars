package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"recordacars/db"
)

func main() {
	if err := run(context.Background(), os.Stdin); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, r io.Reader) error {
	log.Println("setting up database")
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("could not set up database pool: %w", err)
	}
	defer pool.Close()

	queries := db.New(pool)

	dec := json.NewDecoder(r)
	var (
		msg    Message
		params db.InsertACARSMessageParams
	)
	for {
		err := dec.Decode(&msg)
		if err == io.EOF {
			log.Println("reached end-of-file, exiting")
			return nil
		} else if err != nil {
			return err
		}
		log.Println("decoded ACARS message")
		populateParams(&params, &msg)
		if err := queries.InsertACARSMessage(ctx, params); err != nil {
			return err
		}
	}
}

type Message struct {
	Timestamp float64 `json:"timestamp"`
	StationID *string `json:"station_id"`
	Channel   int32   `json:"channel"`
	Freq      float64 `json:"freq"`
	Level     float64 `json:"level"`
	Error     int32   `json:"error"`
	Mode      string  `json:"mode"`
	Label     string  `json:"label"`
	BlockID   *string `json:"block_id"`
	//Ack string|false
	Tail     *string         `json:"tail"`
	Flight   *string         `json:"flight"`
	Msgno    *string         `json:"msgno"`
	Text     *string         `json:"text"`
	End      bool            `json:"end"`
	Depa     *string         `json:"depa"`
	Dsta     *string         `json:"dsta"`
	ETA      *string         `json:"eta"`
	Gtout    *string         `json:"gtout"`
	Gtin     *string         `json:"gtin"`
	Wloff    *string         `json:"wloff"`
	Won      *string         `json:"wlin"` // n.b. this appears to be mislabeled in the source: https://github.com/TLeconte/acarsdec/blob/7920079b8e005c6c798bd478a513211daf9bbd25/output.c#L293-L294
	Sublabel *string         `json:"sublabel"`
	MFI      *string         `json:"mfi"`
	Assstat  *string         `json:"assstat"`
	Libacars json.RawMessage `json:"libacars"`
	App      struct {
		Name string `json:"name"`
		Ver  string `json:"ver"`
	} `json:"app"`
}

// populateParams populates the given db.InsertACARSMessageParams struct using
// the values from msg.
func populateParams(params *db.InsertACARSMessageParams, msg *Message) {
	params.Timestamp = pgtype.Timestamptz{
		Time:             time.UnixMicro(int64(msg.Timestamp * 10e6)).UTC(),
		InfinityModifier: pgtype.Finite,
		Valid:            true,
	}
	params.StationID = pgText(msg.StationID)
	params.Channel = msg.Channel
	params.Freq = msg.Freq
	params.Level = msg.Level
	params.Error = msg.Error
	params.Mode = msg.Mode
	params.Label = msg.Label
	params.BlockID = pgText(msg.BlockID)
	params.Tail = pgText(msg.Tail)
	params.Flight = pgText(msg.Flight)
	params.Msgno = pgText(msg.Msgno)
	params.Text = pgText(msg.Text)
	params.End = msg.End
	params.Depa = pgText(msg.Depa)
	params.Dsta = pgText(msg.Dsta)
	params.Eta = pgText(msg.ETA)
	params.Gtout = pgText(msg.Gtout)
	params.Gtin = pgText(msg.Gtin)
	params.Wloff = pgText(msg.Wloff)
	params.Won = pgText(msg.Won)
	params.Sublabel = pgText(msg.Sublabel)
	params.Mfi = pgText(msg.MFI)
	params.Assstat = pgText(msg.Assstat)
	params.Libacars = msg.Libacars
	params.AppName = msg.App.Name
	params.AppVer = msg.App.Ver
}

func pgText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func pgFloat8(f *float64) pgtype.Float8 {
	if f == nil {
		return pgtype.Float8{}
	}
	return pgtype.Float8{Float64: *f, Valid: true}
}
