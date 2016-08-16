package wot_paser

import (
	"appengine"
	"appengine/urlfetch"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	application_id = "76d0d8bd88a4c872b17404c408332284"
	api_url        = "https://api.worldoftanks.ru/wot"
	account        = User{Id: 2457506, Tanks: make(map[int]Tank)}
	exp_tanks      = make(map[int]ExpTank)
)

func toInit(c appengine.Context) {
	var data struct {
		Data []struct {
			Id         int     `json:"IDNum"`
			ExpFrag    float64 `json:"expFrag"`
			ExpDamage  float64 `json:"expDamage"`
			ExpSpot    float64 `json:"expSpot"`
			ExpDef     float64 `json:"expDef"`
			ExpWinRate float64 `json:"expWinRate"`
		} `json:"data"`
	}
	if len(account.Tanks) == 0 {
		getTanks(c)
		file, _ := ioutil.ReadFile("./expected_tank_values_27.json")
		json.Unmarshal(file, &data)

		for _, tank := range data.Data {
			exp_tanks[tank.Id] = ExpTank{
				Id:         tank.Id,
				ExpFrag:    tank.ExpFrag,
				ExpDamage:  tank.ExpDamage,
				ExpSpot:    tank.ExpSpot,
				ExpDef:     tank.ExpDef,
				ExpWinRate: tank.ExpWinRate,
			}
		}
	}
}

func getAccountStats(c appengine.Context) {
	send_request(c, "/account/info/")
}

func getTanks(c appengine.Context) {
	var data struct {
		Data map[string][]struct {
			Id int `json:"tank_id"`
		} `json:"data"`
	}

	result := send_request(c, "/account/tanks/")
	json.Unmarshal(result, &data)

	for _, tank := range data.Data[strconv.Itoa(account.Id)] {
		account.Tanks[tank.Id] = Tank{}
	}
}

func getTankStats(c appengine.Context) {
	var data struct {
		Data map[string][]struct {
			All struct {
				Spotted                    int `json:"spotted"`
				AvgDamageBlocked           int `json:"avg_damage_blocked"`
				DirectHitsReceived         int `json:"direct_hits_received"`
				ExplosionHits              int `json:"explosion_hits"`
				PiercingsReceived          int `json:"piercings_received"`
				Piercings                  int `json:"piercings"`
				Xp                         int `json:"xp"`
				SurvivedBattles            int `json:"survived_battles"`
				DroppedCapturePoints       int `json:"dropped_capture_points"`
				HitsPercents               int `json:"hits_percents"`
				Draws                      int `json:"draws"`
				DamageReceived             int `json:"damage_received"`
				Frags                      int `json:"frags"`
				CapturePoints              int `json:"capture_points"`
				Hits                       int `json:"hits"`
				BattleAvgXp                int `json:"battle_avg_xp"`
				Losses                     int `json:"losses"`
				DamageDealt                int `json:"damage_dealt"`
				NoDamageDirectHitsReceived int `json:"no_damage_direct_hits_received"`
				Shots                      int `json:"shots"`
				ExplosionHitsReceived      int `json:"explosion_hits_received"`
				TankingFactor              int `json:"tanking_factor"`
				Wins                       int `json:"wins"`
				Battles                    int `json:"battles"`
			} `json:"all"`
			Id       int `json:"tank_id"`
			Mastery  int `json:"mark_of_mastery"`
			MaxExp   int `json:"max_xp"`
			MaxFrags int `json:"max_frags"`
		} `json:"data"`
	}

	params := fmt.Sprintf("&tank_id=%s", tanksIdList(account))

	result := send_request(c, "/tanks/stats/", params)

	json.Unmarshal(result, &data)

	for _, tank := range data.Data[strconv.Itoa(account.Id)] {
		account.Tanks[tank.Id] = Tank{
			Id:                         tank.Id,
			Wins:                       tank.All.Wins,
			Battles:                    tank.All.Battles,
			Mastery:                    tank.Mastery,
			MaxExp:                     tank.MaxExp,
			Spotted:                    tank.All.Spotted,
			AvgDamageBlocked:           tank.All.AvgDamageBlocked,
			DirectHitsReceived:         tank.All.DirectHitsReceived,
			ExplosionHits:              tank.All.ExplosionHits,
			PiercingsReceived:          tank.All.PiercingsReceived,
			Piercings:                  tank.All.Piercings,
			Xp:                         tank.All.Xp,
			SurvivedBattles:            tank.All.SurvivedBattles,
			DroppedCapturePoints:       tank.All.DroppedCapturePoints,
			HitsPercents:               tank.All.HitsPercents,
			Draws:                      tank.All.Draws,
			DamageReceived:             tank.All.DamageReceived,
			Frags:                      tank.All.Frags,
			CapturePoints:              tank.All.CapturePoints,
			Hits:                       tank.All.Hits,
			BattleAvgXp:                tank.All.BattleAvgXp,
			Losses:                     tank.All.Losses,
			DamageDealt:                tank.All.DamageDealt,
			NoDamageDirectHitsReceived: tank.All.NoDamageDirectHitsReceived,
			Shots: tank.All.Shots,
			ExplosionHitsReceived: tank.All.ExplosionHitsReceived,
			TankingFactor:         tank.All.TankingFactor,
			MaxFrags:              tank.MaxFrags,
		}
	}
}

func send_request(c appengine.Context, uri string, params ...string) []byte {
	var (
		query string = ""
	)

	if len(params) > 0 {
		query = params[0]
	}

	tr := &urlfetch.Transport{Context: c, Deadline: time.Duration(30) * time.Second}

	req, _ := http.NewRequest(
		"GET",
		fmt.Sprint(api_url, uri, fmt.Sprintf("?application_id=%s&account_id=%d%s", application_id, account.Id, query)),
		strings.NewReader(""))

	response, _ := tr.RoundTrip(req)

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	return body
}

func getTankInfo(c appengine.Context) {
	var (
		data struct {
			Data map[string]struct{
				Name string `json:"name"`
			} `json:"data"`
		}
		names = make(map[int]string)
	)
	params := fmt.Sprintf("&tank_id=%s", tanksIdList(account))

	result := send_request(c, "/encyclopedia/vehicles/", params)

	json.Unmarshal(result, &data)

	for id, info := range data.Data {

		tankId, _ := strconv.Atoi(id)
		names[tankId] = info.Name
	}

	for id, tank := range account.Tanks {
		tank.Name = names[id]
		account.Tanks[id] = tank
	}
}
