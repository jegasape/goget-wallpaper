package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"
)

func main() {
	base := "https://goget-wallpaper.jegasape.workers.dev"
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	nums := rand.Perm(90)[:12]

	rows := ""
	i := 0
	for row := 0; row < 3; row++ {
		cells := ""
		for col := 0; col < 4; col++ {
			n := nums[i] + 1
			cells += fmt.Sprintf("    <td><img src=\"%s?v=%s%d\" width=\"250\" /></td>\n", base, ts, n)
			i++
		}
		rows += fmt.Sprintf("  <tr>\n%s  </tr>\n", cells)
	}

	table := fmt.Sprintf("<div align=\"center\">\n<table>\n%s</table>\n</div>", rows)

	content, err := os.ReadFile("README.md")
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(`(?s)<!-- COLLAGE_START -->.*?<!-- COLLAGE_END -->`)
	newContent := re.ReplaceAllString(
		string(content),
		fmt.Sprintf("<!-- COLLAGE_START -->\n%s\n<!-- COLLAGE_END -->", table),
	)

	err = os.WriteFile("README.md", []byte(newContent), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Shuffled with timestamp %s\n", ts)
}
