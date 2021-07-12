package employee

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Employee Structure (Model)
type Employee struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Position  string             `json:"position"`
}

var Ctx context.Context
var Collections *mongo.Collection

//var cancel context.CancelFunc

type EmployeeInterface interface {
	init()
	printError(error, http.ResponseWriter)
	GetEmployees(http.ResponseWriter, *http.Request)
	GetEmployee(http.ResponseWriter, *http.Request)
	CreateEmployee(http.ResponseWriter, *http.Request)
	UpdateEmployee(http.ResponseWriter, *http.Request)
	DeleteEmployee(http.ResponseWriter, *http.Request)
}

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ := mongo.Connect(Ctx, clientOptions)
	Ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
	Collections = client.Database("Keepcurrent").Collection("employees")
}

func printError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"message": "` + err.Error() + `"}`))
}

// Get all Employees
func GetEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var employees []Employee
	filter := bson.M{}

	if params := r.URL.Query(); len(params) != 0 {
		// create a query filter
		query := []bson.M{}
		for key, value := range params {
			query = append(query, bson.M{
				key: bson.M{
					"$regex": primitive.Regex{
						Pattern: value[0],
						Options: "i",
					},
				},
			})
		}
		filter = bson.M{"$or": query}
	}

	// get all the employees satisfying the filter requirements
	cursor, err := Collections.Find(Ctx, filter)
	if err != nil {
		printError(err, w)
		return
	}
	defer cursor.Close(Ctx)
	// placing all the employees returned by the Find method in a slice table
	for cursor.Next(Ctx) {
		var employee Employee
		cursor.Decode(&employee)
		employees = append(employees, employee)
	}
	if err := cursor.Err(); err != nil {
		printError(err, w)
		return
	}
	json.NewEncoder(w).Encode(employees)
}

// Get Single Employee
func GetEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // getting params
	var employee Employee
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	err := Collections.FindOne(Ctx, filter).Decode(&employee)

	if err != nil {
		printError(err, w)
		return
	}
	json.NewEncoder(w).Encode(employee)
}

//Create a new Employee
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var employee Employee
	json.NewDecoder(r.Body).Decode(&employee)
	result, err := Collections.InsertOne(Ctx, employee)
	if err != nil {
		printError(err, w)
		return
	}
	json.NewEncoder(w).Encode(result)
}

//update an Employee
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // getting params
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	var updateEmployee, oldEmployee Employee
	var result Employee
	json.NewDecoder(r.Body).Decode(&updateEmployee)

	Collections.FindOne(Ctx, filter).Decode(&oldEmployee)
	if updateEmployee.FirstName == "" {
		updateEmployee.FirstName = oldEmployee.FirstName
	}
	if updateEmployee.LastName == "" {
		updateEmployee.LastName = oldEmployee.LastName
	}
	if updateEmployee.Position == "" {
		updateEmployee.Position = oldEmployee.Position
	}
	update := bson.M{"$set": updateEmployee}
	err := Collections.FindOneAndUpdate(Ctx, filter, update).Decode(&result)
	if err != nil {
		printError(err, w)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "Successfully Updated"}`))
}

//delete an Employee
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // getting params
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var deleted Employee
	err := Collections.FindOneAndDelete(Ctx, bson.M{"_id": id}).Decode(&deleted)
	if err != nil {
		printError(err, w)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "Successfully Deleted"}`))
}
