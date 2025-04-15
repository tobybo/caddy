package player_mng

import (
	"log"

	"server/player_t"
	"server/mongo"
	"server/define"

    "go.mongodb.org/mongo-driver/bson"
)

var gPlys map[string]*player_t.Player = make(map[string]*player_t.Player);

func CheckPlayerExists(name string) bool {
	_, ok := gPlys[name];
	return ok;
}

func CreatePlayer(name string, pwd string) {
	gPlys[name] = player_t.New(name, pwd);
	mongo.SaveCollection(define.DB_NAME, "player", gPlys[name]);
}

func LoadPlys() {
	filter := bson.M{};
	results, err := mongo.QueryCollection(define.DB_NAME, "player", filter);
	if err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		ply := player_t.Player{};
		ply.LoadFromDB(result);
		gPlys[ply.Name] = &ply;
	}
	log.Println("load player success.");
}

func ViewPlys() {
	log.Println("view player:");
	log.Println(gPlys);
}

func GetPlayer(name string) (*player_t.Player, bool) {
	ply, exists := gPlys[name];
	return ply, exists;
}
