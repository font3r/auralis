package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"

	"github.com/google/uuid"
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
				name:   "users",
				scheme: "dbo",
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
			u, _ := uuid.Parse("92bd41cc-62b5-41c9-b542-f9737941407a")

			dataSet := DataSet{
				columnDescriptors: []ColumnDescriptor{
					{name: "id1", dataType: smallint},
					{name: "id2", dataType: smallint},
					{name: "id3", dataType: uniqueidentifier},
					{name: "id4", dataType: smallint},
				},
				rows: []Row{
					{
						[]Cell{
							{math.MinInt16}, {math.MaxInt16}, {u}, {100},
						},
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
			cells, err := ExecuteQuery(args[1])
			if err != nil {
				panic(err)
			}

			displayQueryResultSet(cells)
		}
	}
}

func displayQueryResultSet(dataSet *DataSet) {
	fmt.Printf("INFO: query row result:\n\n%v\n", dataSet)
	fmt.Printf("| ")

	for _, row := range dataSet.rows {
		for cellIndex, cell := range row.cells {
			switch dataSet.columnDescriptors[cellIndex].dataType {
			case smallint:
				{
					fmt.Printf("%d", int16(binary.BigEndian.Uint16(cell.value.([]byte))))
				}
			case uniqueidentifier:
				{
					fmt.Printf("%x", cell.value)
				}
			default:
				{
					fmt.Printf("%s", cell.value)
				}
			}
			fmt.Printf(" | ")
		}
		fmt.Println()
	}

	fmt.Println()
	fmt.Println()
}
