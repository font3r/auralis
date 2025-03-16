package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/google/uuid"
)

const (
	schemeStore string = "scheme_store"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		panic("missing query")
	}

	switch args[1] {
	case "init":
		{
			createSchemaStore()
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
	case "select":
		{
			cells, err := readTableRow(SchemeTable[string, string]{"dbo", "users"})
			if err != nil {
				panic(err)
			}

			fmt.Printf("INFO: query row result %v\n", cells)
		}
	default:
		{
			tokens := Analyze("select id1, id2, id3 from users")
			fmt.Printf("INFO: tokens %v\n", tokens)

			query, err := ParseTokens(tokens)
			if err != nil {
				panic(err)
			}

			switch query := query.(type) {
			case SelectQuery:
				{
					fmt.Printf("INFO: parsed query %v\n", query)
					cells, err := readTableRow(query.source)
					if err != nil {
						panic(err)
					}

					fmt.Printf("INFO: query row result:\n %v \n", cells)
					for _, c := range cells {
						switch c.columnDescriptor.dataType {
						case smallint:
							{
								fmt.Printf("%d|", binary.BigEndian.Uint16(c.value.([]byte)))
							}
						case uniqueidentifier:
							{
								fmt.Printf("%x|", c.value)
							}
						default:
							{
								fmt.Printf("%s|", c.value)
							}
						}
					}
					fmt.Println()
				}
			default:
				panic("unsupported query")
			}
		}
	}

}

func createSchemaStore() error {
	schemeF, err := os.OpenFile(fmt.Sprintf("./data/%s", schemeStore), os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer schemeF.Close()

	return nil
}
