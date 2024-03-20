package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
)

type User struct {
	Name string
	Age  int
}

func main() {
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
		"user": &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: "User",
				Fields: graphql.Fields{
					"name": &graphql.Field{
						Type: graphql.String,
					},
					"age": &graphql.Field{
						Type: graphql.Int,
					},
				},
			}),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return User{
					Name: "Raj",
					Age:  68,
				}, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: fields,
	}

	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
	}

	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		log.Fatalf("Failed to create new schema: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		query, err := ioutil.ReadAll(r.Body)

		if err != nil {
			fmt.Println(err)
			w.Write([]byte(err.Error()))
		}

		fmt.Println(string(query))

		params := graphql.Params{Schema: schema, RequestString: string(query)}
		res := graphql.Do(params)

		if len(res.Errors) > 0 {
			log.Printf("failed to execute graphql operation, errors: %+v \n", res.Errors)
		}

		rJson, _ := json.Marshal(res)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(rJson))
	})

	fmt.Println("Started Server")
	http.ListenAndServe("localhost:8000", nil)
}
