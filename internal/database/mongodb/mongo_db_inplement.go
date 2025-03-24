package mongodb

import (
	"context"
	"fmt"
	// "go.mongodb.org/mongo-driver/bson"

	"github.com/adjivas/eir/pkg/factory"
	mongo_driver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDbConnector struct {
	Mongodb *mongo_driver.Database
}

func NewMongoDbConnector(mongo *factory.Mongodb) MongoDbConnector {
	uri := mongo.Url
	name := mongo.Name
 
	// Create a new client and connect to the server
	client, err := mongo_driver.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// Send a ping to confirm a successful connection
	// var result bson.M
	db := client.Database(name)
	// if err := db.RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
	// 	panic(err)
	// }
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return MongoDbConnector {
		Mongodb: db,
	}
}

// func (m MongoDbConnector) createEquipementStatus(db *mongo.Database) (err error) {
// 	jsonSchema := bson.M{
// 		"bsonType": "object",
// 		"required": []string{"pei", "equipement_status"},
// 		"properties": bson.M{
// 			"pei": bson.M{
// 				"bsonType": "string",
// 				"pattern": "^(imei-[0-9]{15}|imeisv-[0-9]{16}|mac((-[0-9a-fA-F]{2}){6})(-untrusted)?|eui((-[0-9a-fA-F]{2}){8}))$",
// 				"description": "Data type representing the PEI of the UE",
// 			},
// 			"supi": bson.M{
// 				"bsonType": "string",
// 				"pattern": "^(imsi-[0-9]{5,15}|nai-.+|gci-.+|gli-.+)$",
// 				"description": "Data type representing the SUPI of the subscriber",
// 			},
// 			"gpsi": bson.M{
// 				"bsonType": "string",
// 				"pattern": "^(msisdn-[0-9]{5,15}|extid-[^@]+@[^@]+)$",
// 				"description": "Data type representing the GPSI of the subscriber",
// 			},
// 			"equipement_status": bson.M{
// 				"bsonType": "string",
// 				"enum": []string{"WHITELISTED", "BLACKLISTED", "GREYLISTED"},
// 				"description": "Indicates the PEI is white, black or grey listed",
// 			},
// 		},
// 	}
// 	validator := bson.M{
// 		"$jsonSchema": jsonSchema,
// 	}
// 	opts := options.CreateCollection().SetValidator(validator)
//
// 	if err := db.CreateCollection(context.TODO(), "policyData.ues.eirData", opts); err != nil {
// 		return fmt.Errorf("Can't creates EquipementStatus[%s]", err.Error())
// 	}
// 	return nil
// }
//
// func (m MongoDbConnector) dropEquipementStatus() {
// 	if err := m.Mongodb.Collection("policyData.ues.eirData").Drop(context.TODO()); err != nil {
// 		fmt.Errorf("Can't drops EquipementStatus[%s]", err.Error())
// 	}
// }
func (m MongoDbConnector) dropEquipementStatus(db *mongo_driver.Database) {
	if err := db.Collection("policyData.ues.eirData").Drop(context.TODO()); err != nil {
		fmt.Errorf("Can't drops EquipementStatus[%s]", err.Error())
	}
}
