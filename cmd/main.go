package main

import (
	"encoding/csv"
	"fmt"
	"github.com/qPyth/parses-task/internal/parsers"
	"github.com/qPyth/parses-task/internal/types"
	"io"
	"log"
	"os"
)

type Convertable interface {
	ToStringSlice() []string
}

func main() {
	parser := parsers.NewParser()
	category := "all"
	country := "russia"
	persons, err := parser.ParseTopInstagram(category, country)
	if err != nil {
		log.Fatalf(err.Error())
	}
	file, err := os.Create(fmt.Sprintf("top50-%s-%s.csv", category, country))
	if err != nil {
		log.Fatalf(err.Error())
	}
	header := []string{"Rank", "Influencer", "Category", "Followers", "Country", "Eng. (Auth.)", "Eng. (Avg.)"}

	err = ExportToCsv(persons, header, file)
	if err != nil {
		log.Fatalf(err.Error())
	}

}

func ExportToCsv(data []types.Influencer, header []string, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	err := writer.Write(header)
	if err != nil {
		return fmt.Errorf("write header error: %w", err)
	}
	for _, p := range data {
		if err = writer.Write(p.ToStringSlice()); err != nil {
			return fmt.Errorf("data write error: %w\n", err)
		}
	}

	return nil
}
