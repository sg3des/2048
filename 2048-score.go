package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/sg3des/fizzgui"
)

var leaderboardFilename = "leaderboard"

//Header struct contains information of current score and best results
type Header struct {
	con2048 *fizzgui.Container
	wgt2048 *fizzgui.Widget

	conCurr      *fizzgui.Container
	wgtCurrName  *fizzgui.Widget
	wgtCurrScore *fizzgui.Widget

	conBest      *fizzgui.Container
	wgtBestName  *fizzgui.Widget
	wgtBestScore *fizzgui.Widget

	conLB  *fizzgui.Container
	Names  [10]*fizzgui.Widget
	Scores [10]*fizzgui.Widget

	curr User
	best User
	LeaderBoard
}

type User struct {
	Score int
	Name  string
}

func NewHeader() *Header {
	s := new(Header)

	s.loadLeaderBoard()

	colGrey := fizzgui.Color(220, 220, 220, 255)

	s.con2048 = fizzgui.NewContainer("score", "1", "0", "33.3%", "100")
	s.con2048.Style.BackgroundColor = fizzgui.Color(236, 196, 0, 255)
	s.wgt2048 = s.newWdiget(s.con2048, "2048", "100%", NumsFont)

	s.conCurr = fizzgui.NewContainer("currScore", "33.3%", "0", "33.3%", "100")
	s.conCurr.Style.BackgroundColor = fizzgui.Color(187, 173, 160, 255)
	s.wgtCurrName = s.newWdiget(s.conCurr, "SCORE", "50%", TextFontSmall)
	s.wgtCurrName.Layout.Padding.B = 0
	s.wgtCurrName.Layout.Margin.B = 0
	s.wgtCurrName.Style.TextColor = colGrey
	s.wgtCurrScore = s.newWdiget(s.conCurr, "0", "0", TextFont)
	s.wgtCurrScore.Layout.Padding.T = 0
	s.wgtCurrScore.Layout.Margin.T = 0

	s.conBest = fizzgui.NewContainer("bestScore", "66.6%", "0", "33.3%", "100")
	s.conBest.Style.BackgroundColor = fizzgui.Color(187, 173, 160, 255)

	s.wgtBestName = s.conBest.NewButton("BEST", s.ShowLeaderBoard) //(s.conBest, "BEST", "50%", TextFontSmall)
	s.wgtBestName.Layout.SetWidth("100%")
	s.wgtBestName.Layout.SetHeight("50%")
	s.wgtBestName.Layout.Padding.B = 0
	s.wgtBestName.Layout.Margin.B = 0
	s.wgtBestName.Font = TextFontSmall
	s.wgtBestName.Style = fizzgui.NewStyle(colGrey, mgl32.Vec4{0, 0, 0, 0}, mgl32.Vec4{0, 0, 0, 0}, 0)
	s.wgtBestName.StyleHover = fizzgui.NewStyle(colGrey.Add(mgl32.Vec4{0.1, 0.1, 0.1, 0.1}), mgl32.Vec4{0, 0, 0, 0}, mgl32.Vec4{0, 0, 0, 0}, 0)
	s.wgtBestName.StyleActive = fizzgui.NewStyle(colGrey, mgl32.Vec4{0, 0, 0, 0}, mgl32.Vec4{0, 0, 0, 0}, 0)

	s.wgtBestScore = s.newWdiget(s.conBest, strconv.Itoa(s.best.Score), "0", TextFont)
	s.wgtBestScore.Layout.Padding.T = 0
	s.wgtBestScore.Layout.Margin.T = 0

	return s
}

func (*Header) newWdiget(c *fizzgui.Container, text, h string, f *fizzgui.Font) (wgt *fizzgui.Widget) {
	wgt = c.NewText(text)
	wgt.Font = f
	wgt.Layout.SetWidth("100%")
	wgt.Layout.SetHeight(h)
	wgt.TextAlign = fizzgui.TALIGN_CENTER
	wgt.Style.BackgroundColor = fizzgui.Color(0, 0, 0, 0)
	wgt.Style.TextColor = fizzgui.Color(255, 255, 255, 255)

	return
}

type LeaderBoard struct {
	Users []User
}

func (lb LeaderBoard) Len() int {
	return len(lb.Users)
}

func (lb LeaderBoard) Less(i, j int) bool {
	return lb.Users[i].Score < lb.Users[j].Score
}

func (lb LeaderBoard) Swap(i, j int) {
	lb.Users[i], lb.Users[j] = lb.Users[j], lb.Users[i]
}

func (lb LeaderBoard) UsersList() string {
	var list []string
	for _, u := range lb.Users {
		list = append(list, fmt.Sprintf("%-20s %.10d", u.Name, u.Score))
		log.Println(u)
	}

	return strings.Join(list, "\r\n")
}

func (s *Header) loadLeaderBoard() {
	s.conLB = fizzgui.NewContainer("leaderboard", "10%", "10%", "80%", "80%")
	s.conLB.Zorder = 3
	s.conLB.Hidden = true

	title := s.conLB.NewText("Leader Board")
	title.TextAlign = fizzgui.TALIGN_CENTER
	title.Layout.SetWidth("100%")

	for i := 0; i < 10; i++ {
		wgtName := s.conLB.NewText("")
		wgtName.Font = TextFontSmall
		wgtName.Layout.SetWidth("50%")

		wgtScore := s.conLB.NewText("")
		wgtScore.Font = TextFontSmall
		wgtScore.Layout.SetWidth("45%")
		wgtScore.TextAlign = fizzgui.TALIGN_RIGHT

		s.Names[i] = wgtName
		s.Scores[i] = wgtScore
	}

	closeBtn := s.conLB.NewButton("Close", s.CloseLeaderBoard)
	closeBtn.Layout.SetWidth("50%")
	closeBtn.Layout.PositionFixed = true
	closeBtn.Layout.HAlign = fizzgui.HAlignCenter
	closeBtn.Layout.VAlign = fizzgui.VAlignBottom
	closeBtn.Font = TextFontSmall

	f, err := os.Open(leaderboardFilename)
	if err != nil {
		return
	}

	err = gob.NewDecoder(f).Decode(&s.LeaderBoard)
	if err != nil {
		log.Println(err)
		return
	}

	sort.Sort(sort.Reverse(s.LeaderBoard))

	for _, u := range s.LeaderBoard.Users {
		s.best = u
		break
	}
}

func (s *Header) ShowLeaderBoard(_ *fizzgui.Widget) {
	if !s.conLB.Hidden {
		s.conLB.Hidden = true
		return
	}
	s.conLB.Hidden = false

	for i, u := range s.LeaderBoard.Users {
		if i > 9 {
			return
		}
		s.Names[i].Text = fmt.Sprintf("%2d   %-20s", i+1, u.Name)
		s.Scores[i].Text = fmt.Sprintf("%d", u.Score)
	}
}

func (s *Header) CloseLeaderBoard(_ *fizzgui.Widget) {
	s.conLB.Hidden = true
}

func (s *Header) AddScore(score int) {
	s.curr.Score += score
	s.UpdateCurr()
}

func (s *Header) NewGame() {
	if s.curr.Score > 0 {
		s.writeLeaderBoard()
	}

	if s.curr.Score > s.best.Score {
		s.best = s.curr
		s.UpdateBest()
	}

	s.curr.Score = 0
	s.UpdateCurr()
}

func (s *Header) writeLeaderBoard() {
	s.LeaderBoard.Users = append(s.LeaderBoard.Users, s.curr)
	sort.Sort(sort.Reverse(s.LeaderBoard))
	if len(s.LeaderBoard.Users) > 9 {
		s.LeaderBoard.Users = s.LeaderBoard.Users[:10]
	}

	f, err := os.OpenFile(leaderboardFilename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("failed open %s for store result", leaderboardFilename)
		return
	}
	defer f.Close()

	err = gob.NewEncoder(f).Encode(s.LeaderBoard)
	if err != nil {
		log.Println("failed encode data,", err)
	}
}

func (s *Header) UpdateBest() {
	s.wgtBestScore.Text = strconv.Itoa(s.best.Score)
}

func (s *Header) UpdateCurr() {
	s.wgtCurrScore.Text = strconv.Itoa(s.curr.Score)
}
