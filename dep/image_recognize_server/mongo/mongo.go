
package mongo

import (
    "context"
    "fmt"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

	//"server/player_t"
)

var gDb *mongo.Client = nil;

// host = "mongodb://localhost:27017"
func do_connect(host string) (*mongo.Client, error) {
    // 创建一个上下文
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // 连接到 MongoDB
    clientOptions := options.Client().ApplyURI(host)
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    // 检查连接
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB!")

	return client, nil
}

func Connect(host string) {
	client, err := do_connect(host)
	if err != nil {
        log.Fatal(err)
	}
	gDb = client
}

func SaveCollection(db_name string, collection string, data interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// bsonData, err := bson.Marshal(data)
	// if err != nil {
	//     log.Fatalf("Error marshaling bsonData data: %v", err)
	// }
	coll := gDb.Database(db_name).Collection(collection)
	_, err := coll.InsertOne(ctx, data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("save success.")
}

func UpdateCollection(db_name string, collection string, filter interface{}, data interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := gDb.Database(db_name).Collection(collection)
	chgs := make(map[string]interface{})
	chgs["$set"] = data
	_, err := coll.UpdateOne(ctx, filter, chgs)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("save success.")
}

func QueryCollection(db_name string, collection string, filter bson.M) (results []bson.M, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := gDb.Database(db_name).Collection(collection)
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		log.Fatal("Failed to find documents: ", err)
	}
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result);
	}
	return results, nil
}

func FindAndModify(db_name string, collectionName string, id string, chgs interface{}) (bson.M, error) {
    // 创建一个上下文
    ctx := context.TODO()

    // 获取指定的集合
    collection := gDb.Database(db_name).Collection(collectionName)

    // 定义查询条件，这里假设我们只更新第一条记录
    filter := bson.M{
		"_id": id,
	} // 可以根据需求修改查询条件

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true);
	var result bson.M;
	err := collection.FindOneAndUpdate(ctx, filter, chgs, opts).Decode(&result);
	if err != nil {
		log.Fatal(err)
	}

	//log.Println(result);

    return result, nil
}
