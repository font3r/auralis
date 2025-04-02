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
			CreateSchemaStore()
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
			c1 := Cell{ColumnDescriptor{name: "id1", dataType: smallint}, math.MinInt16}
			c2 := Cell{ColumnDescriptor{name: "id2", dataType: smallint}, math.MaxInt16}
			u, _ := uuid.Parse("92bd41cc-62b5-41c9-b542-f9737941407a")
			c3 := Cell{ColumnDescriptor{name: "id3", dataType: uniqueidentifier}, u}
			c4 := Cell{ColumnDescriptor{name: "id4", dataType: smallint}, 100}

			err := writeIntoTable("dbo", "users", c1, c2, c3, c4)
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

func displayQueryResultSet(cells []Cell) {
	fmt.Printf("INFO: query row result:\n\n%v\n", cells)
	fmt.Printf("| ")

	for i := range cells {
		switch cells[i].columnDescriptor.dataType {
		case smallint:
			{
				fmt.Printf("%d", int16(binary.BigEndian.Uint16(cells[i].value.([]byte))))
			}
		case uniqueidentifier:
			{
				fmt.Printf("%x", cells[i].value)
			}
		default:
			{
				fmt.Printf("%s", cells[i].value)
			}
		}

		fmt.Printf(" | ")
	}

	fmt.Println()
	fmt.Println()
}
