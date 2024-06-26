package graphql

import (
	"Backend/ent"
	"Backend/ent/user"
	"context"
	"entgo.io/ent/dialect"
	"errors"
	"github.com/graphql-go/graphql"
)

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: UserType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: ResolveUser,
			},
			"users": &graphql.Field{
				Type:    graphql.NewList(UserType),
				Resolve: ResolveAllUsers,
			},
		},
	}),
	Mutation: graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"addUser": &graphql.Field{
				Type: UserType,
				Args: graphql.FieldConfigArgument{
					"user": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(UserInputType),
					},
				},
				Resolve: AddUser,
			},
			"deleteUser": &graphql.Field{
				Type: UserType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: deleteUser,
			}, "updateUser": &graphql.Field{
				Type: UserType,
				Args: graphql.FieldConfigArgument{
					"user": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(UserUpdateInputType),
					},
				},
				Resolve: UpdateUser,
			},
		},
	}),
})

var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"age": &graphql.Field{
			Type: graphql.Int,
		},
		"salary": &graphql.Field{
			Type: graphql.Float,
		},
	},
})
var UserInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UserInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"age": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"salary": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.Float),
		},
	},
})
var UserUpdateInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UserUpateInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"age": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
		"salary": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.Float),
		},
	},
})

func ResolveUser(params graphql.ResolveParams) (interface{}, error) {

	id, ok := params.Args["id"].(int)
	if !ok {
		return nil, errors.New("missing 'id' argument")
	}

	// Connect to the database
	client, err := ent.Open(dialect.MySQL, "root:root@tcp(127.0.0.1:3306)/my_db?parseTime=true")
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Create a context
	ctx := context.Background()

	foundUser, err := client.User.Query().
		Where(user.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	userData := map[string]interface{}{
		"id":     foundUser.ID,
		"name":   foundUser.Name,
		"age":    foundUser.Age,
		"salary": foundUser.Salary,
	}

	return userData, nil
}

func ResolveAllUsers(params graphql.ResolveParams) (interface{}, error) {

	client, err := ent.Open(dialect.MySQL, "root:root@tcp(127.0.0.1:3306)/my_db?parseTime=true")
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()

	allUsers, err := client.User.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	var userData []map[string]interface{}

	for _, u := range allUsers {
		userData = append(userData, map[string]interface{}{
			"id":     u.ID,
			"name":   u.Name,
			"age":    u.Age,
			"salary": u.Salary,
		})
	}

	return userData, nil
}
func AddUser(params graphql.ResolveParams) (interface{}, error) {

	userInput, ok := params.Args["user"].(map[string]interface{})
	if !ok {
		return nil, errors.New("missing 'user' argument")
	}

	name, ok := userInput["name"].(string)
	if !ok {
		return nil, errors.New("missing 'name' field in 'user'")
	}

	age, ok := userInput["age"].(int)
	if !ok {
		return nil, errors.New("missing 'age' field in 'user'")
	}

	salary, ok := userInput["salary"].(float64)
	if !ok {
		return nil, errors.New("missing 'salary' field in 'user'")
	}

	client, err := ent.Open(dialect.MySQL, "root:root@tcp(127.0.0.1:3306)/my_db?parseTime=true")
	if err != nil {
		return nil, err
	}

	defer client.Close()
	ctx := context.Background()

	newUser, err := client.User.Create().
		SetName(name).
		SetAge(age).
		SetSalary(salary).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func deleteUser(params graphql.ResolveParams) (interface{}, error) {
	id, ok := params.Args["id"].(int)
	if !ok {
		return nil, errors.New("missing 'id' argument")
	}

	client, err := ent.Open(dialect.MySQL, "root:root@tcp(127.0.0.1:3306)/my_db?parseTime=true")
	if err != nil {
		return nil, err
	}

	defer client.Close()
	ctx := context.Background()
	foundUser, err := client.User.Query().
		Where(user.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	err = client.User.DeleteOne(foundUser).Exec(ctx)
	if err != nil {
		return nil, err
	}

	userData := map[string]interface{}{
		"id":     foundUser.ID,
		"name":   foundUser.Name,
		"age":    foundUser.Age,
		"salary": foundUser.Salary,
	}
	return userData, nil
}
func UpdateUser(params graphql.ResolveParams) (interface{}, error) {
	userInput, ok := params.Args["user"].(map[string]interface{})
	if !ok {
		return nil, errors.New("missing 'user' argument")
	}
	id, ok := userInput["id"].(int)
	if !ok {
		return nil, errors.New("missing 'id' argument")
	}

	newName, newNameExists := userInput["name"].(string)
	newAge, newAgeExists := userInput["age"].(int)
	newSalary, newSalaryExists := userInput["salary"].(float64)

	client, err := ent.Open(dialect.MySQL, "root:root@tcp(127.0.0.1:3306)/my_db?parseTime=true")
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx := context.Background()

	foundUser, err := client.User.Query().
		Where(user.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	if newNameExists {
		foundUser = foundUser.Update().SetName(newName).SaveX(ctx)
	}
	if newAgeExists {
		foundUser = foundUser.Update().SetAge(newAge).SaveX(ctx)
	}
	if newSalaryExists {
		foundUser = foundUser.Update().SetSalary(newSalary).SaveX(ctx)
	}

	userData := map[string]interface{}{
		"id":     foundUser.ID,
		"name":   foundUser.Name,
		"age":    foundUser.Age,
		"salary": foundUser.Salary,
	}
	return userData, nil
}
