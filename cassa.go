package main

import(
    "fmt"

    "github.com/gocassa/gocassa"
)

type Sale struct {
	id          string
	customerid  string
	sellerid    string
	price       int
	//created     int
}

func main() {
	keySpace, err := gocassa.ConnectToKeySpace("razvan", []string{"127.0.0.1"}, "", "")
	if err != nil {
		panic(err)
	}

	salesTable := keySpace.Table("sale", &Sale{}, gocassa.Keys{
		PartitionKeys: []string{"id"},
	})

	//fmt.Println(salesTable)

	err = salesTable.Set(Sale{
		id: "sale-1",
		customerid: "customer-1",
		sellerid: "seller-1",
		price: 42,
		//created: int(time.Now().Unix()),
	}).Run()

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("END")

	//result := Sale{}
	//if err := salesTable.Where(gocassa.Eq("id", "sale-1")).ReadOne(&result).Run(); err != nil {
	//	panic(err)
	//}
	//fmt.Println(result)
}