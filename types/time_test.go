package types

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestTime_MarshalJSON(t *testing.T) {
	now := time.Now()
	nowJson, err := json.Marshal(now)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("nowJson:", string(nowJson[:]))

	tim := Time(now)
	timJson, err := json.Marshal(tim)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("timJson:", string(timJson[:]))
}

func TestTime_GetDays(t *testing.T) {
	now := time.Now()

	today := Time(now)
	days := today.GetDays(now)
	if days != 0 {
		log.Fatal("expect 0; actual ", days)
	}

	yesterday := Time(now.Add(-time.Hour * 24))
	days = yesterday.GetDays(now)
	if days != -1 {
		log.Fatal("expect -1; actual ", days)
	}

	tomorrow := Time(now.Add(time.Hour * 24))
	days = tomorrow.GetDays(now)
	if days != 1 {
		log.Fatal("expect 1; actual ", days)
	}

	after := Time(now.Add(time.Hour * 24 * 7))
	days = after.GetDays(now)
	if days != 7 {
		log.Fatal("expect 1; actual ", days)
	}

	before := Time(now.Add(-time.Hour * 24 * 5))
	days = before.GetDays(now)
	if days != -5 {
		log.Fatal("expect 1; actual ", days)
	}
}
