package routine

import (
	"github.com/chenjie4255/tools/mongohelper"
	"github.com/chenjie4255/tools/testenv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_stringHash(t *testing.T) {

	val := hash2Int("1232111123", 22)
	t.Logf("1232: %d", val)
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("short test")
	}

	env := testenv.GetIntegratedTestEnv()
	if env.MongoHost == "" {
		t.Skip("the environment configuration is not ready yet")
	}

	mgoSession, _ := mongohelper.NewClient(env.MongoHost, "", "", "admin")
	Convey("add job", t, func() {
		tb := NewWeeklyTable("test", "weekly_jobs", 1024, mgoSession)
		mgoSession.Database("test").Drop(nil)

		addJob := func(offset int32, uid string) (string, error) {
			j := Job{}
			j.UID = uid
			j.Data = []byte(`111222333`)
			ret, err := tb.AddJob(j, []WeekOffset{WeekOffset(offset)})
			return ret[0], err
		}

		j1, err := addJob(1, "1")
		j2, err := addJob(2, "2")
		j3, err := addJob(3, "3")
		j4, err := addJob(4, "4")
		j5, err := addJob(5, "5")

		idx := []string{j1, j2, j3, j4, j5}

		t.Logf("new added jobs' index: %+v", idx)

		Convey("check jobs", func() {
			checks, err := tb.ScanCellsPartitions(1, 0, 1024)
			So(err, ShouldBeNil)

			jobs := checks.Jobs()
			So(jobs, ShouldHaveLength, 1)

			checks, err = tb.ScanCellsPartitions(3, 0, 1024)
			So(err, ShouldBeNil)

			jobs = checks.Jobs()
			So(jobs, ShouldHaveLength, 1)
		})

		Convey("remove jobs", func() {
			err = tb.RemoveJob("1", []string{j1})
			if err != nil {
				t.Errorf("failed to remove job, %s", err)
			}

			checks, err := tb.ScanCellsPartitions(1, 0, 1024)
			So(err, ShouldBeNil)
			So(checks.Jobs(), ShouldHaveLength, 0)
		})
	})
}
