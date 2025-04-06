package main

import (
	"fmt"
	"math"
	"os"

	"github.com/google/uuid"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		panic("missing query")
	}

	switch args[1] {
	case "init":
		{
			initDatabaseInternalStructure()
		}
	case "tables":
		{
			dataSet, err := getAuralisTables()
			if err != nil {
				panic(err)
			}

			fmt.Println()
			displayDataSet(dataSet)
		}
	case "columns":
		{
			dataSet, err := getAuralisColumns()
			if err != nil {
				panic(err)
			}

			fmt.Println()
			displayDataSet(dataSet)
		}
	case "create":
		{
			err := cretateTable(TableDescriptor{
				source: SchemeTable[string, string]{"dbo", "users"},
				columnDescriptors: []ColumnDescriptor{
					{
						name:     "id1",
						dataType: smallint,
						position: 1,
					},
					{
						name:     "id2",
						dataType: smallint,
						position: 2,
					},
					{
						name:     "name",
						dataType: varchar,
						position: 3,
					},
					{
						name:     "id3",
						dataType: uniqueidentifier,
						position: 4,
					},
					{
						name:     "id4",
						dataType: smallint,
						position: 5,
					},
				},
			})

			if err != nil {
				panic(err)
			}
		}
	case "insert":
		{
			u1, _ := uuid.Parse("92bd41cc-62b5-41c9-b542-f9737941407a")
			u2, _ := uuid.Parse("37eae8fd-c2e6-48dc-b33f-d6bc5c9e1425")

			dataSet := DataSet{
				columnDescriptors: []ColumnDescriptor{
					{name: "id1", dataType: smallint},
					{name: "id2", dataType: smallint},
					{name: "name", dataType: varchar},
					{name: "id3", dataType: uniqueidentifier},
					{name: "id4", dataType: smallint},
				},
				rows: []Row{
					{
						[]Cell{{math.MinInt16}, {math.MaxInt16}, {"test-111"}, {u1}, {100}},
					},
					{
						[]Cell{{math.MinInt16}, {math.MaxInt16}, {"test-2"}, {u2}, {200}},
					},
				},
			}

			err := writeIntoTable(SchemeTable[string, string]{"dbo", "users"}, dataSet)
			if err != nil {
				panic(err)
			}
		}
	default:
		{
			dataSet, err := ExecuteQuery(args[1])
			if err != nil {
				panic(err)
			}

			fmt.Println()
			displayDataSet(dataSet)
		}
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
			tableRow = append(tableRow, fmt.Sprintf("%v", dataCell.value))
		}
		t.AppendRow(tableRow)
	}

	t.Render()
}
