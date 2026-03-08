package party

import "log/slog"

// XP values taken for D&D 5e
var XpThresholds = []int{
	0,   // placeholder (no level 0)
	300, // 1
	900,
	2700,
	6500,
	14000,
	23000,
	34000,
	48000,
	64000,
	85000,
	100000,
	120000,
	140000,
	165000,
	195000,
	225000,
	265000,
	305000,
	355000,
}

// DistributeExperience distributes XP and checks for level-ups
func (p *Party) DistributeExperience(xp int) {
	if len(p.ActiveMembers) == 0 {
		slog.Warn("did not distribute experience", "xp", xp)
		return
	}

	splitXP := xp / len(p.ActiveMembers)
	extraXP := xp % len(p.ActiveMembers)

	for i := range p.ActiveMembers {
		p.ActiveMembers[i].XP += splitXP
		slog.Info("gained xp", "member", p.ActiveMembers[i].Name, "amount", splitXP)

		p.ActiveMembers[i].UpdateLevel()
	}

	slog.Info("distributed xp", "total", xp, "per_member", splitXP, "unclaimed", extraXP)
}
