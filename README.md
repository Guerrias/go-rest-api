# Minimal API - Golang

This project is an API build with the programing language Golang and using MongoDB for data storage.

## Tools Used
* Programing Language : Go v1.16
* Database : MongoDB v4.2.3
* Package Installed: 
  * Gorilla "github.com/gorilla/mux" v1.8.0
  * Mongo Driver for Go:  "go.mongodb.org/mongo-driver" v1.5.4
* Test's Interface : Postman
* Editor : Visual Studio
* Version Control : Github

## Data Structure
The data stored is a representation of employees' information. The structure is as follow:

type Employee struct{<br>
     "_id": automatically generated,<br>
     "firstName": "example first name",<br>
     "lastName": "example last name",<br>
     "position": "example position"<br>
}

> The attribute "_id" is unmutable.

## Operability
The API is running on port 8000.
MongoDB is running on port 27017.

### Routes
* Get : "http://localhost:8000/api/employees?lastname=lastnameValue&firstname=firstnameValue&position=positionValue" => List all the employees in the database.
* Get : "http://localhost:8000/api/employees/{id}" => List the information of the employee which ID matches the parameter
* Post : "http://localhost:8000/api/employees" => Add an entry in the database
* Update : "http://localhost:8000/api/employees/{id}" => Update the entry in the database which ID matches the parameter
* Delete : "http://localhost:8000/api/employees/{id}" => Delete an entry in the database which ID matches the parameter

#### Note
The parameters for the end point to read all the employees are optional. When present, the result will be filtered based on the attributes passed in parameter. The filter algorithm is based on regex patterns rather than exact values.

For instance, "position=soft" might have the same result as "position=sof" 
