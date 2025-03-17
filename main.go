package main

import (
	"encoding/binary"
	"fmt"
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
						dataType: uniqueidentifier,
					},
					{
						name:     "id3",
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
			c1 := Cell{ColumnDescriptor{name: "id1", dataType: smallint}, 65535}
			u, _ := uuid.Parse("92bd41cc-62b5-41c9-b542-f9737941407a")
			c2 := Cell{ColumnDescriptor{name: "id2", dataType: uniqueidentifier}, u}
			c3 := Cell{ColumnDescriptor{name: "id3", dataType: smallint}, 100}

			err := writeIntoTable("dbo", "users", c1, c2, c3)
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
	fmt.Printf("INFO: query row result:\n\n %v \n", cells)
	fmt.Printf("| ")

	for i := range cells {
		switch cells[i].columnDescriptor.dataType {
		case smallint:
			{
				fmt.Printf("%d", binary.BigEndian.Uint16(cells[i].value.([]byte)))
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
