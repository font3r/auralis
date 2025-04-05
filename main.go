package main

import (
	"encoding/binary"
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
	case "create":
		{
			err := cretateTable(TableDescriptor{
				source: SchemeTable[string, string]{"dbo", "users"},
				columnDescriptors: []ColumnDescriptor{
					{
						name:     "id1",
						dataType: smallint,
					},
					{
						name:     "id2",
						dataType: smallint,
					},
					{
						name:     "name",
						dataType: varchar,
					},
					{
						name:     "id3",
						dataType: uniqueidentifier,
					},
					{
						name:     "id4",
						dataType: smallint,
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
			datSet, err := ExecuteQuery(args[1])
			if err != nil {
				panic(err)
			}

			fmt.Println()
			displayDataSet(datSet)
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
		var formattedCell string
		for i, dataCell := range dataRow.cells {
			switch dataSet.columnDescriptors[i].dataType {
			case smallint:
				formattedCell = fmt.Sprintf("%d", int16(binary.BigEndian.Uint16(dataCell.value.([]byte))))
			case uniqueidentifier:
				formattedCell = fmt.Sprintf("%x", dataCell.value)
			default:
				formattedCell = fmt.Sprintf("%s", dataCell.value)
			}

			tableRow = append(tableRow, formattedCell)
		}
		t.AppendRow(tableRow)
	}

	t.Render()
}
