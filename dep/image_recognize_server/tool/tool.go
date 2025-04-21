package tool

import (
    "log"

	"image_recognize_server/mongo"
	"image_recognize_server/define"

	"go.mongodb.org/mongo-driver/bson"
)

func IncSn(what string) int32 {
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

func GetSn(what string) int32 {
	filter := bson.M{"_id": what};
	res, err := mongo.QueryCollection(define.DB_NAME, "uniq", filter);
	if err != nil {
		log.Fatal(err)
	}
	if len(res) == 0 {
		return 0;
	}
	sn := res[0]["sn"].(int32);
	return sn;
}
