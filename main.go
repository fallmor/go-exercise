package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	CsvFilename := flag.String("csv", "problems.csv", "entrez un fichier csv contenant le quiz")
	duree := flag.Int("t", 30, "entrez timer")
	flag.Parse()
	_ = CsvFilename
	_ = duree

	f, err := os.Open(*CsvFilename)
	if err != nil {
		exit(fmt.Sprintf("Impossible d'ouvrir le fichier %s", *CsvFilename))
	}
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Impossible de lire le fichier %s", *CsvFilename))
	}
	aff := getProblem(lines)
	correct := 0
	montime := time.NewTimer(time.Duration((*duree)) * time.Second)
	mychan := make(chan string)
	for p, pa := range aff {
		go func() {
			var reponse string
			fmt.Printf("Voici la question %d qui est %s\n", p+1, pa.question)
			fmt.Scanf("%s\n", &reponse)
			mychan <- reponse
		}()
		select {
		case <-montime.C:
			fmt.Printf("Voici la question %d qui est %s\n", p+1, pa.question)
		case reponse := <-mychan:
			if reponse == pa.reponse {
				correct++
			}
		}
	}
	fmt.Printf("Vous avez trouve %d sur %d questions \n", correct, len(aff))
}

type problem struct {
	question string
	reponse  string
}

// parser le fichier excel pour avoir un struct pour manipuler les données
//getProblem prend en entrée un slice de deux dimensions et renvoie un struct
func getProblem(lines [][]string) []problem {
	rep := make([]problem, len(lines))
	for i, line := range lines {
		rep[i] = problem{
			question: line[0],
			reponse:  line[1],
		}
	}
	return rep
}
func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
