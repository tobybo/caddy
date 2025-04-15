package tool

import (
    "log"

	"server/mongo"
	"server/define"

	"go.mongodb.org/mongo-driver/bson"
)

func GetSn(what string) int32 {
	// var res struct {
	//     sn int32;
	// }
	res, err := mongo.FindAndModify(define.DB_NAME, "uniq", what, bson.M{"$inc": bson.M{"sn": 1}});
	if err != nil {
		log.Fatal(err)
	}
	sn := res["sn"].(int32);
	return sn;
}
