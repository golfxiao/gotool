package ucticket

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoTicketStore struct{}

func NewMongoTicketStore() *MongoTicketStore {
	return &MongoTicketStore{}
}

func (this *MongoTicketStore) LoadIDSegment(bizTag string) (segment *TicketSegment, err error) {
	return this.LoadIDSegmentWithNum(bizTag, int64(cfg.Step))
}

func (this *MongoTicketStore) LoadIDSegmentWithNum(bizTag string, num int64) (segment *TicketSegment, err error) {
	collection := mongoPool.Database(cfg.DatabaseName).Collection(cfg.TableName)
	after := options.After

	r := collection.FindOneAndUpdate(context.Background(),
		bizTagFilter(bizTag), incrNum(int(num)),
		&options.FindOneAndUpdateOptions{ReturnDocument: &after})
	if r.Err() != nil {
		return nil, r.Err()
	}

	segment = new(TicketSegment)
	err = r.Decode(segment)
	if err != nil {
		return nil, err
	}
	segment.Step = int(num)
	return segment, err
}

func (this *MongoTicketStore) InitScope(bizTag string, step int, maxId int64) (err error) {
	collection := mongoPool.Database(cfg.DatabaseName).Collection(cfg.TableName)
	_, err = collection.InsertOne(context.Background(), map[string]interface{}{
		"biz_tag": bizTag,
		"max_id":  maxId,
		"step":    step,
	})
	return
}

func bizTagFilter(bizTag string) *bson.M {
	return &bson.M{
		"biz_tag": bizTag,
	}
}

func incrNum(step int) bson.M {
	return bson.M{
		"$inc": &bson.M{
			"max_id": step,
		},
	}
}
