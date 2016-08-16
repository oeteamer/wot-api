package wot_paser

import (
	"strconv"
	"strings"
	//"math"
	"log"
	"math"
)

type User struct {
	Id    int
	Tanks map[int]Tank
}

type Tank struct {
	Id                         int
	Wins                       int
	Battles                    int
	Mastery                    int
	MaxExp                     int
	Spotted                    int
	AvgDamageBlocked           int
	DirectHitsReceived         int
	ExplosionHits              int
	PiercingsReceived          int
	Piercings                  int
	Xp                         int
	SurvivedBattles            int
	DroppedCapturePoints       int
	HitsPercents               int
	Draws                      int
	DamageReceived             int
	Frags                      int
	CapturePoints              int
	Hits                       int
	BattleAvgXp                int
	Losses                     int
	DamageDealt                int
	NoDamageDirectHitsReceived int
	Shots                      int
	ExplosionHitsReceived      int
	TankingFactor              int
	MaxFrags                   int
	WN8                        int
	Eff                        int
	Name                       string
}

type ExpTank struct {
	Id         int
	ExpFrag    float64
	ExpDamage  float64
	ExpSpot    float64
	ExpDef     float64
	ExpWinRate float64
}

func tanksIdList(u User) string {
	var list []string
	for id := range u.Tanks {
		list = append(list, strconv.Itoa(id))
	}

	return strings.Join(list, ",")
}

func calculateStats(user User, expTanks map[int]ExpTank) {

	for id, tank := range user.Tanks {
		avgDAMAGE := float64(tank.DamageDealt) / float64(tank.Battles)
		avgSPOT := float64(tank.Spotted) / float64(tank.Battles)
		avgFRAG := float64(tank.Frags) / float64(tank.Battles)
		avgDEF := float64(tank.DroppedCapturePoints) / float64(tank.Battles)
		avgWIN := float64(tank.Wins) / float64(tank.Battles)

		rDAMAGE := avgDAMAGE / expTanks[id].ExpDamage
		rSPOT := avgSPOT / expTanks[id].ExpSpot
		rFRAG := avgFRAG / expTanks[id].ExpFrag
		rDEF := avgDEF / expTanks[id].ExpDef
		rWIN := avgWIN / expTanks[id].ExpWinRate

		rWINc := math.Max(0, (rWIN-0.71)/(1-0.71))
		rDAMAGEc := math.Max(0, (rDAMAGE-0.22)/(1-0.22))
		rFRAGc := math.Min(rDAMAGEc+0.2, math.Max(0, (rFRAG-0.12)/(1-0.12)))
		rSPOTc := math.Min(rDAMAGEc+0.1, math.Max(0, (rSPOT-0.38)/(1-0.38)))
		rDEFc := math.Min(rDAMAGEc+0.1, math.Max(0, (rDEF-0.10)/(1-0.10)))

		WN8 := 980*rDAMAGEc + 210*rDAMAGEc*rFRAGc + 155*rFRAGc*rSPOTc + 75*rDEFc*rFRAGc + 145*math.Min(1.8, rWINc)

		log.Print(tank.Name)
		log.Print(WN8)
	}
}
