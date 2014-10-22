package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"text/template"

	"github.com/codegangsta/cli"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

const helpTemplate = `
Usage: tomoe [options]

Options:
{{range .Flags}}  {{.}}
{{end}}`

const statusTmpl = "@tomo_e {{if .Place}}{{.Place}}で{{else}}どこでもいいので{{end}}{{if .Thing}}{{.Thing}}が{{end}}おいしいお店教えてください！{{if .People}}人数は{{.People}}人です！{{end}}"

func main() {
	cli.AppHelpTemplate = helpTemplate[1:]

	app := cli.NewApp()
	app.Name = "Tomoe API Client"
	app.HideHelp = true
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "place, p", Usage: "Set place"},
		cli.StringFlag{Name: "thing, t", Usage: "Set thing"},
		cli.IntFlag{Name: "people, n", Usage: "Set number of people"},
		cli.BoolFlag{Name: "dry-run", Usage: "Show preview"},
		cli.HelpFlag,
	}

	app.Action = doMain

	err := app.Run(os.Args)

	if err != nil {
		os.Exit(1)
	}
}

func doMain(c *cli.Context) {
	tmpl, err := template.New("status").Parse(statusTmpl)

	if err != nil {
		log.Fatal(err)
	}

	var doc bytes.Buffer
	err = tmpl.Execute(&doc, map[string]interface{}{
		"Place":  c.String("place"),
		"Thing":  c.String("thing"),
		"People": c.Int("people"),
	})

	if err != nil {
		log.Fatal(err)
	}

	status := doc.String()
	if c.Bool("dry-run") {
		fmt.Println(status)
		return
	}

	q := url.Values{}
	q.Set("status", status)

	cmd := exec.Command("sh", "-c", "open https://twitter.com/intent/tweet?"+q.Encode())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
