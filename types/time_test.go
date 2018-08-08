package types

import (
	"testing"
	"time"
	"encoding/json"
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
