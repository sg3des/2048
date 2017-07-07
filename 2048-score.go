package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

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

	s.con2048 = fizzgui.NewContainer("score", "0", "0", "33.3%", "100")
	s.con2048.Style.BackgroundColor = fizzgui.Color(236, 196, 0, 255)
	s.wgt2048 = s.newWdiget(s.con2048, "2048", "100%", NumsFont)

	s.conCurr = fizzgui.NewContainer("currScore", "33.3%", "0", "33.3%", "100")
	s.conCurr.Style.BackgroundColor = fizzgui.Color(187, 173, 160, 255)
	s.wgtCurrName = s.newWdiget(s.conCurr, "SCORE", "50%", TextFontSmall)
	s.wgtCurrName.Layout.Padding.B = 0
	s.wgtCurrName.Layout.Margin.B = 0
	s.wgtCurrName.Style.TextColor = fizzgui.Color(220, 220, 220, 255)
	s.wgtCurrScore = s.newWdiget(s.conCurr, "0", "0", TextFont)
	s.wgtCurrScore.Layout.Padding.T = 0
	s.wgtCurrScore.Layout.Margin.T = 0

	s.conBest = fizzgui.NewContainer("bestScore", "66.6%", "0", "33.3%", "100")
	s.conBest.Style.BackgroundColor = fizzgui.Color(187, 173, 160, 255)
	s.wgtBestName = s.newWdiget(s.conBest, "BEST", "50%", TextFontSmall)
	s.wgtBestName.Layout.Padding.B = 0
	s.wgtBestName.Layout.Margin.B = 0
	s.wgtBestName.Style.TextColor = fizzgui.Color(220, 220, 220, 255)
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

func (lb LeaderBoard) ToWrite() (data []byte) {
	for _, u := range lb.Users {
		data = append(data, []byte(fmt.Sprintf("%d %s\r\n", u.Score, u.Name))...)
	}

	return data
}

func (s *Header) loadLeaderBoard() {
	f, err := os.Open(leaderboardFilename)
	if err != nil {
		return
	}

	err = gob.NewDecoder(f).Decode(&s.LeaderBoard)
	if err != nil {
		log.Println(err)
	}
	sort.Sort(s.LeaderBoard)

	for i := len(s.LeaderBoard.Users) - 1; i > 0; i-- {
		s.best = s.LeaderBoard.Users[i]
		break
	}
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
	sort.Sort(s.LeaderBoard)
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

	// ioutil.WriteFile(leaderboardFilename, s.LeaderBoard.ToWrite(), 0644)
}

func (s *Header) UpdateBest() {
	s.wgtBestScore.Text = strconv.Itoa(s.best.Score)
}

func (s *Header) UpdateCurr() {
	s.wgtCurrScore.Text = strconv.Itoa(s.curr.Score)
}
