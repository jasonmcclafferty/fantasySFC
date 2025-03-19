package models

var counties = []string{
	"Antrim",
	"Armagh",
	"Carlow",
	"Cavan",
	"Clare",
	"Cork",
	"Derry",
	"Donegal",
	"Down",
	"Dublin",
	"Fermanagh",
	"Galway",
	"Kerry",
	"Kildare",
	"Kilkenny",
	"Laois",
	"Leitrim",
	"Limerick",
	"Longford",
	"Louth",
	"Mayo",
	"Meath",
	"Monaghan",
	"Offaly",
	"Roscommon",
	"Sligo",
	"Tipperary",
	"Tyrone",
	"Waterford",
	"Westmeath",
	"Wexford",
	"Wicklow",
}

type County struct {
	Name  string
	Squad []Player
}

type Province struct {
	Name     string
	Counties []County
}

type Player struct {
	County County
	Price  int
	Pos    string
	Score  int
}

type Squad struct {
	GK []Player
	FB []Player
	HB []Player
	MF []Player
	HF []Player
	FF []Player
}

type ActiveTeam struct {
	Players  []Player
	Subs     []Player
	Captain  Player
	Setpiece Player
}
