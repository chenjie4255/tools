package routine

import (
	"context"
	"fmt"
	"github.com/chenjie4255/errors"
	"github.com/chenjie4255/tools/errcode"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hash/fnv"
)

type weeklyTable struct {
	dbName     string
	tableName  string
	maxHash    uint32
	mgoSession *mongo.Client
}

func NewWeeklyTable(dbName string, tableName string, maxHash uint32, mgoSession *mongo.Client) WeeklyTable {
	if maxHash == 0 {
		panic("maxHash cannot be zero")
	}

	//*mgoSession.Database(dbName).Collection(tableName)

	//mgoSession. (dbName).C(tableName).EnsureIndexKey("offset", "partition")
	//mgoSession.DB(dbName).C(tableName).EnsureIndexKey("index")

	return &weeklyTable{dbName, tableName, maxHash, mgoSession}
}

func hash2Int(str string, max uint32) uint32 {
	a := fnv.New32a()
	a.Write([]byte(str))
	val := a.Sum32()
	return val % max
}

func (t *weeklyTable) PartitionCount() uint32 {
	return t.maxHash
}

func (t *weeklyTable) UniqueName() string {
	return fmt.Sprintf("%s_%s", t.dbName, t.tableName)
}

func (t *weeklyTable) AddJob(job Job, offsets []WeekOffset) ([]string, error) {
	partition := hash2Int(job.UID, t.maxHash)
	var ret []string
	for _, offset := range offsets {
		idx := WeeklyTableIndex{}
		idx.Offset = uint32(offset)
		idx.Partition = partition

		coll := t.mgoSession.Database(t.dbName).Collection(t.tableName)

		index := fmt.Sprintf("%d_%d", offset, partition)

		ops := options.Update().SetUpsert(true)

		if _, err := coll.UpdateOne(context.Background(), bson.M{"offset": offset, "partition": partition},
			bson.M{"$addToSet": bson.M{"jobs": job}, "$setOnInsert": bson.M{"index": index}}, ops); err != nil {
			return nil, err
		}

		ret = append(ret, index)
	}

	return ret, nil
}

func (t *weeklyTable) ScanCellsPartitions(offset, from, to uint64) (Cells, error) {

	c := t.mgoSession.Database(t.dbName).Collection(t.tableName)
	ret := []Cell{}
	cur, err := c.Find(context.Background(), bson.M{"offset": offset, "partition": bson.M{"$gte": from, "$lt": to}}, nil)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	if err := cur.All(context.Background(), &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (t *weeklyTable) RemoveJob(uid string, indexes []string) error {
	c := t.mgoSession.Database(t.dbName).Collection(t.tableName)
	if len(indexes) == 0 {
		return errors.NewWithTag("indexes cannot be empty", errcode.ParamError)
	}

	query := bson.M{"index": bson.M{"$in": indexes}}
	updator := bson.M{"$pull": bson.M{"jobs": bson.M{"uid": uid}}}
	if _, err := c.UpdateMany(context.Background(), query, updator); err != nil {
		return err
	}

	return nil
}
