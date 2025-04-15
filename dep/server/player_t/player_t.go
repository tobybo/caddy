
package player_t

import (
	"log"

    "go.mongodb.org/mongo-driver/bson"

	"server/tool"
)

var BASIC_PLAYER_ID int = 1000000;
var PERMISSION_READ int32 = 1;
var PERMISSION_WRITE int32 = 2;

type Player struct {
	Pid int32 `bson:"_id"`;
	Name string `bson:"name"`;
	Pwd string `bson:"pwd"`;
	Permission int32 `bson:"permission"`;
}

func (p *Player) LoadFromDB(bsonData bson.M) {
	p.Pid = bsonData["_id"].(int32);
	p.Name = bsonData["name"].(string);
	p.Pwd = bsonData["pwd"].(string);
	p.Permission = bsonData["permission"].(int32);
}

func New(name string, pwd string) (*Player) {
	// 向数据库请求一个自增后的id
	var pid = tool.GetSn("player");
	log.Println("new player, pid:", pid);
	return &Player{Pid: pid, Name: name, Pwd: pwd, Permission: PERMISSION_READ};
}
