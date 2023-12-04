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
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".eml") {
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
