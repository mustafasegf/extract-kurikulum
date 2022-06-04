package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Matkul struct {
	Title       string   `json:"ltitle"`
	Description string   `json:"description"`
	Code        string   `json:"code"`
	SKS         int      `json:"sks"`
	Depedencies []string `json:"depedencies"`
}

func main() {
	b, err := os.ReadFile("./silabus.txt")
	s := strings.Replace(string(b), "&", "dan", -1)
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`\n\nMata\sKuliah\n\n`)
	rawMatkuls := re.Split(s, -1)[1:]

	re = regexp.MustCompile(`(?sU)(?P<Title>.*)SKS\n\n(?P<SKS>\d).*Prasyarat\n+(?P<Depedencies>.*)\n+Deskripsi\n+(?P<Description>.*)\n+Capaian`)
	matkuls := []Matkul{}

	reCode := regexp.MustCompile(`\s+\((CS.*)\).*`)
	for _, rawMatkul := range rawMatkuls {
		match := re.FindStringSubmatch(rawMatkul)

		mapMatkul := map[string]string{}
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				mapMatkul[name] = strings.TrimSpace(match[i])
			}
		}

		title := strings.ReplaceAll(mapMatkul["Title"], "\n", " ")
		code := reCode.FindStringSubmatch(title)[1]
		title = reCode.ReplaceAllString(title, "")

		sks, _ := strconv.Atoi(mapMatkul["SKS"])
		deps := strings.Split(strings.ReplaceAll(mapMatkul["Depedencies"], "\n", " "), ",")
		for i, dep := range deps {
			deps[i] = strings.TrimSpace(dep)
		}

		matkul := Matkul{
			Title:       title,
			Description: strings.ReplaceAll(mapMatkul["Description"], "\n", " "),
			Code:        code,
			SKS:         sks,
			Depedencies: deps,
		}

		matkuls = append(matkuls, matkul)
	}

	b, err = json.MarshalIndent(matkuls, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
