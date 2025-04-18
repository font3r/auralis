package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		panic("missing query")
	}

	initDatabaseInternalStructure()

	dataSet, err := ExecuteQuery(args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	if dataSet != nil {
		displayDataSet(dataSet)
	} else {
		log.Printf("NO RESULT\n")
	}
}

func displayDataSet(dataSet *DataSet) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetAutoIndex(true)

	style := table.StyleDefault
	style.Format.Header = text.FormatDefault
	t.SetStyle(style)

	tableHeader := table.Row{}
	for _, cd := range dataSet.columnDescriptors {
		tableHeader = append(tableHeader, cd.name)
	}
	t.AppendHeader(tableHeader)

	for _, dataRow := range dataSet.rows {
		tableRow := table.Row{}
		for _, dataCell := range dataRow.cells {
			tableRow = append(tableRow, fmt.Sprintf("%v", dataCell))
		}
		t.AppendRow(tableRow)
	}

	t.Render()
}
