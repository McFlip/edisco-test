package ingestemail

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
)

type Email struct {
	From, To, Date, Subject string
}

func ParseEml(emlStr string) (Email, error) {
	r := strings.NewReader(emlStr)
	m, err := mail.ReadMessage(r)
	if err != nil {
		return Email{}, err
	}

	header := m.Header
	from := header.Get("From")
	to := header.Get("To")
	date := header.Get("Date")
	subject := header.Get("Subject")

	// clean up the date string so that in can be type-cast into a timestamp with timezone object in Postgres
	date, _, _ = strings.Cut(date, " (")

	myEmail := Email{
		From:    from,
		To:      to,
		Date:    date,
		Subject: subject,
	}

	return myEmail, nil
}

func Ingest(inDir, outFile string) error {
	f, err := os.OpenFile(outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	ok := filepath.WalkDir(inDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		// NOTE: I'm only parsing email at first. ReadPST is able to output only email but files are just numbers with no extension.
		// readpst -a ".fubar" -D -b -S -t e
		// if d.IsDir() || !strings.HasSuffix(d.Name(), ".eml") {
		// 	return nil
		// }
		if d.IsDir() {
			return nil
		}
		emlBs, err := os.ReadFile(path)
		if err != nil {
			slog.Log(context.Background(), slog.LevelError, fmt.Sprintf("Error reading .eml file: %s", err))
			return nil
		}

		eml, err := ParseEml(string(emlBs))
		if err != nil {
			slog.Log(context.Background(), slog.LevelError, fmt.Sprintf("Error parsing email: %s", err))
			return nil
		}

		emlJson, err := json.Marshal(eml)
		if err != nil {
			slog.Log(context.Background(), slog.LevelError, fmt.Sprintf("Error marshalling email: %s", err))
			return nil
		}
		emlJson = append(emlJson, byte('\n'))
		_, err = f.Write(emlJson)
		if err != nil {
			slog.Log(context.Background(), slog.LevelError, fmt.Sprintf("Error writing json: %s", err))
			return nil
		}
		return nil
	})
	return ok
}
