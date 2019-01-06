package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nick-jones/brulee"
)

func main() {
	input := `
	score(politics) = 0
	score(sports) = 0
	score(film) = 0

	# politics

	when
		var(title) contains "conservatives"
		or var(title) contains "labour"
	then
		score(politics) += 10
		when
			var(title) contains "party"
		then
			score(politics) += 5
		done
	done

	when
		var(title) contains "brexit"
	then
		score(politics) += 5
	done

	when
		var(topic) in [
			"elections", "international relations", "national security",
			"economics", "liberalism", "socialism", "fascism", "eu"
		]
	then
		score(politics) += 5
	done

	# sports

	when
		var(title) contains "football"
		or var(title) matches /(snow|skate|wake|kite)board(ing|er)/
		or var(title) contains "cricket"
		or (var(title) contains "running" and var(title) does not contain "zombies")
	then
		score(sports) += 10
	done

	# film

	when
		var(title) matches /(film|movie)/
	then
		score(film) += 10
	done

	when
		var(title) contains "zombies"
	then
		score(film) += 2
	done
	`

	program := brulee.MustCompile(strings.NewReader(input))
	program.Dump(os.Stdout)

	articles := []struct {
		title string
		topic string
	}{
		{
			title: "Corbyn tells May to strike new Brexit deal the Labour party can back",
			topic: "EU",
		},
		{
			title: "The greatest games that decided college football championships",
			topic: "football",
		},
		{
			title: "That Snowboarder Guy: An Interview with Mathieu Crepel",
			topic: "snowboarding",
		},
		{
			title: "Running From Zombies: the best film of 2019?",
			topic: "movies",
		},
	}
	for _, a := range articles {
		vars := map[string]string{
			"title": strings.ToLower(a.title),
			"topic": strings.ToLower(a.topic),
		}
		scores, err := program.Run(vars)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Title: %s, Topic: %s\n", a.title, a.topic)
		fmt.Print("Scores:")
		for name, score := range scores {
			fmt.Printf(" %s=%d", name, score)
		}
		fmt.Println()
		fmt.Println("+-----+")
	}
}
