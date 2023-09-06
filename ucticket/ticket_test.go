package ucticket

import (
	"sync"
	"testing"

	"github.com/astaxie/beego/orm"
	"github.com/stretchr/testify/assert"
)

const (
	Biztag_Test BiztagType = "user"
)

func TestMain(m *testing.M) {
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", "kaifa:kaifa@123@tcp(192.168.35.172:3306)/test?charset=utf8mb4", 2, 2)
	m.Run()
}

func initMongoEnv(usePreload bool) error {
	config := TicketConfig{
		DataSrc:       "mongodb://10.255.255.170:27017/?replicaSet=rsumsdata&readPreference=secondarypreferred&w=majority&journal=false&wtimeoutMS=5000&ssl=false",
		ConnCount:     5,
		Step:          100,
		TableName:     "ticket",
		ScopeList:     []string{Biztag_Test.String()},
		DatabaseName:  TICKET_DB_ALIASE_NAME,
		UsePreload:    usePreload,
		PreloadFactor: 0.5,
	}
	return InitTicketMongo(config)
}

func initMySQLEnv(usePreload bool) error {
	config := TicketConfig{
		DataSrc:       "kaifa:kaifa@123@tcp(192.168.35.172:3306)/test?charset=utf8mb4",
		ConnCount:     5,
		Step:          100,
		TableName:     "ticket",
		ScopeList:     []string{Biztag_Test.String()},
		UsePreload:    usePreload,
		PreloadFactor: 0.5,
	}
	return InitTicketDB(config)
}

func Test_TicketMongo(t *testing.T) {
	bizTag := Biztag_Test
	err := initMongoEnv(false)
	assert.Nil(t, err)

	newId, err := bizTag.GetGlobalId()
	assert.Nil(t, err)
	assert.True(t, newId > 0)
	t.Logf("init ticket mongo success, newId: %d", newId)

	// concurrent bizTag
	wg := new(sync.WaitGroup)
	time := 98
	for i := 0; i < time; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := bizTag.GetGlobalId()
			assert.Nil(t, err)
		}()
	}
	wg.Wait()

	newId2, err2 := bizTag.GetGlobalId()
	assert.Nil(t, err2)
	assert.Equal(t, newId+99, newId2)

	t.Logf("newId: %d, newId2： %d", newId, newId2)
}

func Test_TicketMySQL(t *testing.T) {
	bizTag := Biztag_Test
	err := initMySQLEnv(false)
	assert.Nil(t, err)

	newId, err := bizTag.GetGlobalId()
	assert.Nil(t, err)
	assert.True(t, newId > 0)
}

func BenchmarkGlobalId(b *testing.B) {
	b.StopTimer()
	initMySQLEnv(false)
	// assert.Nil(b, err)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Biztag_Test.GetGlobalId()
	}
}

func BenchmarkGlobalIdUseSecondaryCache(b *testing.B) {
	b.StopTimer()
	initMySQLEnv(true)
	// assert.Nil(b, err)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Biztag_Test.GetGlobalId()
	}
}

func BenchmarkMongoGlobalId(b *testing.B) {
	b.StopTimer()
	err := initMongoEnv(false)
	assert.Nil(b, err)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Biztag_Test.GetGlobalId()
	}
}

func BenchmarkMongoGlobalIdUseSecondaryCache(b *testing.B) {
	b.StopTimer()
	err := initMongoEnv(true)
	assert.Nil(b, err)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Biztag_Test.GetGlobalId()
	}
}
