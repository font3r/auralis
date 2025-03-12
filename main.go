package main

import (
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
						dataType: uniqueidentifier,
					},
					{
						name:     "id2",
						dataType: uniqueidentifier,
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
			c1 := Cell{ColumnDescriptor{name: "id1", dataType: uniqueidentifier}, u1}

			u2, _ := uuid.Parse("e1856542-6776-4935-b258-f350f293cf14")
			c2 := Cell{ColumnDescriptor{name: "id2", dataType: uniqueidentifier}, u2}

			err := writeIntoTable("dbo", "users", c1, c2)
			if err != nil {
				panic(err)
			}
		}
	case "select":
		{
			cells, err := readTableRow("dbo", "users")
			if err != nil {
				panic(err)
			}

			fmt.Printf("INFO: query result %s", cells)
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
