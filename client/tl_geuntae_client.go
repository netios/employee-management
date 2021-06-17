package main

// SEND JSON
import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/lxn/walk"

	// GET IP

	// GUI

	"math/rand"
	"os"
	"sort"
	"syscall"
	"time"
	"unsafe"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	. "github.com/lxn/walk/declarative"
	"github.com/tealeg/xlsx"
)

// SOUND

// 엑셀

// RANDOM

//"strconv"
//"fmt"
//"math/rand"
//"strings"
//"time"

//Call DLL (경고창 구현

const (
	server_url = "127.0.0.1"
	server_port = "4000"
	
	//d_userid = "1"
	//d_username = "user1"

	d_userid   = "2"
	d_username = "user2"

	//d_userid = "3"
	//d_username = "user3"

	//d_userid = "7"
	//d_username = "user4"

	sound_path = "vorbis/"
)

const (
	SexMale Sex = 1 + iota
	SexFemale
	SexHermaphrodite
)

type DurationField struct {
	p *time.Duration
}

func (*DurationField) CanSet() bool       { return true }
func (f *DurationField) Get() interface{} { return f.p.String() }
func (f *DurationField) Set(v interface{}) error {
	x, err := time.ParseDuration(v.(string))
	if err == nil {
		*f.p = x
	}
	return err
}
func (f *DurationField) Zero() interface{} { return "" }

type Sex byte

type Jobday struct {
	Name          string
	ArrivalDate   time.Time
	SpeciesId     int
	Speed         int
	Sex           Sex
	Weight        float64
	PreferredFood string
	Domesticated  bool
	Remarks       string
	Patience      time.Duration
}

type ReportDaily struct {
	Remarks string // 업무내용
	//Patience      string   //
	//Domesticated  bool
}

func (a *Jobday) PatienceField() *DurationField {
	return &DurationField{&a.Patience}
}

func KnownSpecies() []*Species {
	return []*Species{
		{1, "List1"},
		{2, "List2"},
		{3, "List3"},
		{4, "List4"},
		{5, "List5"},
	}
}

type Species struct {
	Id   int
	Name string
}

//-------------------- 테이블뷰

type Foo struct {
	Index int
	//VoicePcm    time.Time
	VoicePcm         string
	VoiceScore       string
	VoiceTextUser    string
	VoiceTextDefault string
	checked          bool
}

type FooModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*Foo
}

func NewFooModel() *FooModel {
	m := new(FooModel)
	m.ResetRows()
	return m
}

func main() {

	// var server_time = get_server_time

	// if server_time < 15:00 출근 : actioncode = '출근'
	// if server_time > 15:00 퇴근 : actioncode = '퇴근'
	GUI_main()

}

func GUI_main() {

	walk.FocusEffect, _ = walk.NewBorderGlowEffect(walk.RGB(0, 63, 255))
	walk.InteractionEffect, _ = walk.NewDropShadowEffect(walk.RGB(63, 63, 63))
	walk.ValidationErrorEffect, _ = walk.NewBorderGlowEffect(walk.RGB(255, 0, 0))

	//boldFont, _ := walk.NewFont("Consolas", 12, walk.FontBold)
	model := NewFooModel()

	var mw *walk.MainWindow
	var te *walk.TextEdit
	var tv *walk.TableView

	//Jobday := new(Jobday)
	reportdaily := new(ReportDaily)

	rt_str := fmt.Sprintf("(%s)님 출/퇴근 체크앱", d_username)

	if fileNotExists("tl_geuntae_client.ico") {
		msgbox("경고창","현재 경로에 tl_geuntae_client.ico 파일이 없습니다.")
		os.Exit(0)
	}
	
	if _, err := (MainWindow{
		AssignTo: &mw,
		Icon: "tl_geuntae_client.ico",
		Title:    rt_str,
		MinSize:  Size{1100, 900},
		Size:     Size{1100, 900},
		Layout:   VBox{},
		Children: []Widget{
			PushButton{
				Text: "▲ 출근 인증 버튼",
				OnClicked: func() {
					te.SetText(check_in())
				},
			},
			PushButton{
				Text: "▼ 퇴근 인증 버튼",
				OnClicked: func() {
					te.SetText(check_out())
				},
			},
			PushButton{
				Text: "▣ 휴가 연차 신청",
				OnClicked: func() {
					te.SetText(check_out())
				},
			},
			/*
				PushButton{
					Text:      "▣ 시트갱신",
					OnClicked: model.ResetRows,
				},
			*/
			PushButton{
				Text: "◆ 일일 업무 보고",
				OnClicked: func() {
					if cmd, err := Dialog_daily(mw, reportdaily); err != nil { // 팝업창을 띄워서
						fmt.Println(err)
					} else if cmd == walk.DlgCmdOK {
						te.SetText(fmt.Sprintf("%+v", reportdaily.Remarks))
						//te.SetText("--일일보고를 기록했습니다--") // 메시지창에 작업 결과를 찍는다.
						daily_put_info(reportdaily.Remarks)
					}

					model.ResetRows()
				},
			},
			/*
				PushButton{
					Text: "Edit Jobday",
					OnClicked: func() {
						if cmd, err := RunJobdayDialog(mw, Jobday); err != nil {
							fmt.Println(err)
						} else if cmd == walk.DlgCmdOK {
							te.SetText(fmt.Sprintf("%+v", Jobday))
						}
					},
				},
			*/
			TextEdit{
				//Size:     Size{1100, 100},
				TextColor: walk.RGB(0, 80, 0),
				AssignTo:  &te,
			},
			TableView{
				AssignTo:              &tv,
				//AlternatingRowBGColor: walk.RGB(239, 239, 239),
				//CheckBoxes:            true,
				ColumnsOrderable: true,
				MultiSelection:   true,
				Columns: []TableViewColumn{
					{Title: "NO", Width: 70},
					{Title: "날짜", Width: 150},
					{Title: "보고내용", Alignment: AlignNear, Width: 594}, // AlignNear, AlignCenter, AlignFar
					{Title: "비고", Alignment: AlignNear, Width: 250},   // AlignNear, AlignCenter, AlignFar
				},
				StyleCell: func(style *walk.CellStyle) {
					item := model.items[style.Row()]
					//style.Font = boldFont
					if item.checked {
						if style.Row()%2 == 0 {
							style.BackgroundColor = walk.RGB(159, 215, 255)
						} else {
							style.BackgroundColor = walk.RGB(143, 199, 239)
						}
					}
					/*
						// 화면출력 스타일을 설정한다. (1: 막대그래프, 2:점수, 3:파일명)
						switch style.Col() {

							case 0,1,2,3:

								if item.VoiceScore >= 80.0 {
									style.TextColor = walk.RGB(0, 191, 0)
									//style.Font = boldFont
								} else if item.VoiceScore < 40.0 {
									style.TextColor = walk.RGB(255, 0, 0)
									//style.Font = boldFont
								}

						}
					*/
				},
				Model: model,
				OnItemActivated: func() {
					//fmt.Printf("Activate: %v\n", tv.CurrentIndex())
					fmt.Print(model.items[tv.CurrentIndex()].VoicePcm)
					//fmt.Printf("\n")
					//runCmd("s16le","1", "8000", model.items[tv.CurrentIndex()].VoicePcm)
				},
			},
			/*
							TableView{
								//AssignTo:              &tv,
								AssignTo:              &tv,
								AlternatingRowBGColor: walk.RGB(239, 239, 239),
								//CheckBoxes:            true,
								ColumnsOrderable:      true,
								MultiSelection:        true,

								Columns: []TableViewColumn{
									{Title: "#", Width: 70},
									{Title: "파일명", Width: 360},
									{Title: "점수", Width: 50 },
									{Title: "TEXT",Alignment: AlignNear, Width: 600},  // AlignNear, AlignCenter, AlignFar
									{Title: "기준TEXT",Alignment: AlignNear, Width: 600},  // AlignNear, AlignCenter, AlignFar
								},
								StyleCell: func(style *walk.CellStyle) {
									item := model.items[style.Row()]
											//style.Font = boldFont
									if item.checked {
										if style.Row()%2 == 0 {
											style.BackgroundColor = walk.RGB(159, 215, 255)
										} else {
											style.BackgroundColor = walk.RGB(143, 199, 239)
										}
									}

									// 화면출력 스타일을 설정한다. (1: 막대그래프, 2:점수, 3:파일명)
									switch style.Col() {

										case 0,1,2,3,4:

											if item.VoiceScore >= 80.0 {
												style.TextColor = walk.RGB(0, 191, 0)
												style.Font = boldFont
											} else if item.VoiceScore < 40.0 {
												style.TextColor = walk.RGB(255, 0, 0)
												style.Font = boldFont
											}

									}
								},
								Model: model,
				//				OnSelectedIndexesChanged: func() {
				//					if tv.SelectedIndexes() != nil {
				//						runCmd("s16le","1", "8000", value)
				//					}
				//					fmt.Printf("SelectedIndexes: %v\n", tv.SelectedIndexes())
				//				},
								OnItemActivated: func() {
									//fmt.Printf("Activate: %v\n", tv.CurrentIndex())
									fmt.Print(model.items[tv.CurrentIndex()].VoicePcm)
									//fmt.Printf("\n")
									//runCmd("s16le","1", "8000", model.items[tv.CurrentIndex()].VoicePcm)
								},

							},
			*/
		},
	}).Run(); err != nil {
		fmt.Println(err)
	}
}

/*
	bg, err := walk.NewSolidColorBrush(walk.RGB(0, 255, 0))
	if err != nil {
		return nil, err
	}
	w.SetBackground(bg)
*/

// 출근 처리
func check_in() string {

	var rt_str string

	if get_info("employee_cnt") == "false" {

		put_info("출근")
		gettime := get_info("employee_time_in")
		rt_str = fmt.Sprintf("출근시간 | %s | (%s)님의 출근 인증 완료함.", gettime, d_username)
		//playsound(sound_random_index("a"))

	} else {

		gettime := get_info("employee_time_in")
		rt_str = fmt.Sprintf("출근시간 | %s | (%s)님은 이미 출근한 상태입니다.\r\n________________________________________________________\r\n\r\n오전 09:30분 이후부터 지각처리임에 유념해주세요.", gettime, d_username)

	}
	return rt_str

}

// 퇴근 처리
func check_out() string {

	var rt_str string

	if get_info("employee_cnt") == "false" {

		rt_str = fmt.Sprintf("확인요망 | (%s)님 출근 인증을 먼저하세요.", d_username)

	} else {

		put_info("퇴근")
		gettime := get_info("employee_time_out")
		rt_str = fmt.Sprintf("퇴근시간 | %s | (%s)님의 퇴근을 등록함.\r\n________________________________________________________\r\n\r\n일반근무퇴근: 오후 18:00 ~ 18:30 까지 입니다.\r\n오전반차퇴근: 오후 14:00 ~ 14:30 입니다.", gettime, d_username)
		//playsound(sound_random_index("c"))

	}
	return rt_str

}

// 일일 보고 처리
func check_daybyday() string {

	var rt_str string

	if get_info("employee_cnt") == "false" {

		rt_str = fmt.Sprintf("확인요망 | (%s)님 출근 인증을 먼저하세요.", d_username)

	} else {

		// ---- 작업할 내용을 이곳에서 처리한다.

		// 메시지창 진행상태 출력
		gettime := get_info("employee_time_out")
		rt_str = fmt.Sprintf("일일보고 | %s | (%s)님의 일일보고를 등록함.\r\n________________________________________________________\r\n\r\n", gettime, d_username)
		// 음성플레이

	}
	return rt_str

}

func getdate(in_str string) string {
	//var rt_str string
	runes := []rune(in_str)
	return string(runes[0:9])
}

//--------------------------- 일일 업무 저장 기록
func Dialog_daily(owner walk.Form, reportdaily *ReportDaily) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	//rt_str := fmt.Sprintf("[%s] %s 님 일일 업무 보고",getdate(get_info("employee_time_in")),d_username)

	return Dialog{
		AssignTo: &dlg,
		//Title:         Bind("'Jobday Details' + (Jobday.Name == '' ? '' : ' - ' + Jobday.Name)"),
		Title: Bind("tttttttttttttt"),
		//Title:         rt_str,
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "reportdaily",
			DataSource:     reportdaily,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{550, 150},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					VSpacer{
						ColumnSpan: 2,
						Size:       0,
					},
					Label{
						//ColumnSpan: 2,
						Text: "일일업무내용:",
					},
					TextEdit{
						//ColumnSpan: 2,
						MinSize: Size{100, 50},
						Text:    Bind("Remarks"),
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								fmt.Println(err)
								return
							}

							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}

//------------------------ 연차 사용 조절
func RunJobdayDialog(owner walk.Form, Jobday *Jobday) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	return Dialog{
		AssignTo:      &dlg,
		Title:         Bind("'Jobday Details' + (Jobday.Name == '' ? '' : ' - ' + Jobday.Name)"),
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "Jobday",
			DataSource:     Jobday,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{300, 300},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Name:",
					},
					LineEdit{
						Text: Bind("Name"),
					},

					Label{
						Text: "Arrival Date:",
					},
					DateEdit{
						Date: Bind("ArrivalDate"),
					},

					Label{
						Text: "Species:",
					},
					ComboBox{
						Value:         Bind("SpeciesId", SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         KnownSpecies(),
					},

					Label{
						Text: "Speed:",
					},
					Slider{
						Value: Bind("Speed"),
					},

					RadioButtonGroupBox{
						ColumnSpan: 2,
						Title:      "Sex",
						Layout:     HBox{},
						DataMember: "Sex",
						Buttons: []RadioButton{
							{Text: "Male", Value: SexMale},
							{Text: "Female", Value: SexFemale},
							{Text: "Hermaphrodite", Value: SexHermaphrodite},
						},
					},

					Label{
						Text: "Weight:",
					},
					NumberEdit{
						Value:    Bind("Weight", Range{0.01, 9999.99}),
						Suffix:   " kg",
						Decimals: 2,
					},

					Label{
						Text: "Preferred Food:",
					},
					ComboBox{
						Editable: true,
						Value:    Bind("PreferredFood"),
						Model:    []string{"Fruit", "Grass", "Fish", "Meat"},
					},

					Label{
						Text: "Domesticated:",
					},
					CheckBox{
						Checked: Bind("Domesticated"),
					},

					VSpacer{
						ColumnSpan: 2,
						Size:       8,
					},

					Label{
						ColumnSpan: 2,
						Text:       "Remarks:",
					},
					TextEdit{
						ColumnSpan: 2,
						MinSize:    Size{100, 50},
						Text:       Bind("Remarks"),
					},

					Label{
						ColumnSpan: 2,
						Text:       "Patience:",
					},
					LineEdit{
						ColumnSpan: 2,
						Text:       Bind("PatienceField"),
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								fmt.Println(err)
								return
							}

							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}

//func db_get_daily () {
// 디비에서 접속해서 자료를 읽어온다.
// 일간 정보를 읽어온다.
//select * from
//d_userid
//}

func get_info(cmd string) string {

	// cmd List
	// -------------------
	//  __ 공용 시간값을 받음
	//	employee_time     : 서버기준 현재시간 [시간반환]
	//
	//  __ 아래는 UserID 값을 요구함.(개인데티터를 받음)
	//	employee_cnt      : 출근했는지 여부 [0]=미출근, [1]=출근
	//	employee_time_in  : 출근한시간 [시간반환]
	//	employee_time_out : 퇴근한시간 [시간반환]

	var url string

	url = fmt.Sprintf("http://%s:%s/%s/%s", server_url, server_port, cmd, d_userid)

	// switch cmd {
	// case "employee_time":
	// url = fmt.Sprintf("http://%s:%s/%s",server_url,server_port,cmd)
	// default:
	// url = fmt.Sprintf("http://%s:%s/%s/%s",server_url,server_port,cmd,d_userid)
	// }

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	/*
	   if err == nil {
	       str := string(respBody)
	       println(str)
	   }
	*/
	return string(respBody)
}

func put_info(ctype string) {

	ip, err := findSystemIP()
	if err != nil {
		fmt.Println(err)
	}

	url_putdata := fmt.Sprintf("http://%s:%s/employee", server_url, server_port)

	resp, err := http.PostForm(url_putdata, url.Values{"UserId": {d_userid}, "Name": {d_username}, "ActCode": {ctype}, "Ip": {ip}})

	defer resp.Body.Close()

	// Response 체크.
	/* // 디버그용
	    respBody, err := ioutil.ReadAll(resp.Body)
	    if err != nil {
			fmt.Println(err)
		}

	    if err == nil {
	        str := string(respBody)
	        println(str)
	    }
	*/
}

// 일일보고
func daily_put_info(ctype string) {

	ip, err := findSystemIP()

	//ip := "test"

	if err != nil {
		fmt.Println(err)
	}

	url_putdata := fmt.Sprintf("http://%s:%s/report_daily", server_url, server_port)

	resp, err := http.PostForm(url_putdata, url.Values{"UserId": {d_userid}, "Name": {d_username}, "Report": {ctype}, "Comment": {ip}})

	defer resp.Body.Close()

	// Response 체크.
	/* // 디버그용
	    respBody, err := ioutil.ReadAll(resp.Body)
	    if err != nil {
			fmt.Println(err)
		}

	    if err == nil {
	        str := string(respBody)
	        println(str)
	    }
	*/
}

//---------------- SOUND
func sound_random_index(stype string) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return fmt.Sprintf("%s%s%d", sound_path, stype, r1.Intn(3))
}

//---------------- SOUND
func playsound(file_sound string) {

	// Open first sample File
	f, err := os.Open(file_sound)

	// Check for errors when opening the file
	if err != nil {
		fmt.Println("error:", err)
	}

	// Decode the .ogg File, if you have a .wav file, use wav.Decode(f)
	//s, format, err := vorbis.Decode(f)
	s, format, _ := wav.Decode(f)
	if err != nil {
		fmt.Println("error:", err)
	}
	//println(format.SampleRate.N(time.Second/10))
	// Init the Speaker with the SampleRate of the format and a buffer size of 1/10s
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	//speaker.Init(24000, 2400)

	// Channel, which will signal the end of the playback.
	playing := make(chan struct{})

	// Now we Play our Streamer on the Speaker
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		// Callback after the stream Ends
		close(playing)
	})))
	<-playing

}

func findSystemIP() (string, error) {
	// list of system network interfaces
	// https://golang.org/pkg/net/#Interfaces
	intfs, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	// mapping between network interface name and index
	// https://golang.org/pkg/net/#Interface
	for _, intf := range intfs {
		// skip down interface & check next intf
		if intf.Flags&net.FlagUp == 0 {
			continue
		}
		// skip loopback & check next intf
		if intf.Flags&net.FlagLoopback != 0 {
			continue
		}
		// list of unicast interface addresses for specific interface
		// https://golang.org/pkg/net/#Interface.Addrs
		addrs, err := intf.Addrs()
		if err != nil {
			return "", err
		}
		// network end point address
		// https://golang.org/pkg/net/#Addr
		for _, addr := range addrs {
			var ip net.IP
			// Addr type switch required as a result of IPNet & IPAddr return in
			// https://golang.org/src/net/interface_windows.go?h=interfaceAddrTable
			switch v := addr.(type) {
			// net.IPNet satisfies Addr interface
			// since it contains Network() & String()
			// https://golang.org/pkg/net/#IPNet
			case *net.IPNet:
				ip = v.IP
			// net.IPAddr satisfies Addr interface
			// since it contains Network() & String()
			// https://golang.org/pkg/net/#IPAddr
			case *net.IPAddr:
				ip = v.IP
			}
			// skip loopback & check next addr
			if ip == nil || ip.IsLoopback() {
				continue
			}
			// convert IP IPv4 address to 4-byte
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			// return IP address as string
			return ip.String(), nil
		}
	}
	return "", errors.New("no ip interface up")
}

func (m *FooModel) ResetRows() {

	var score string
	var name string
	var named string
	var maxrow int

	returncode := 1

	excelFileName := "geuntae.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		//log.Fatal(err)
		returncode = 0
	}

	if returncode == 0 { // 파일이 없을때

		retstr01 := fmt.Sprintf("%s 파일을 찾을 수 없습니다.", excelFileName)
		retstr02 := fmt.Sprintf("프로그램 경로에 %s 파일을 복사해 주세요", excelFileName)
		m.items = make([]*Foo, 1)
		m.items[0] = &Foo{Index: 1, VoicePcm: retstr01, VoiceScore: "", VoiceTextUser: retstr02}

	} else { // 파일이 있을때

		for _, sheet := range xlFile.Sheets {
			maxrow = sheet.MaxRow
		}
		//m := &EnvModel{items: make([]EnvItem, maxrow-1)}
		m.items = make([]*Foo, maxrow-1)
		for _, sheet := range xlFile.Sheets {
			for i, row := range sheet.Rows {
				if i > 0 {
					for cellcount, cell := range row.Cells {
						text := cell.String()
						if cellcount == 0 {
							//tmp, err = strconv.Atoi(text)  // 점수
							score = text // 날짜
						}
						if cellcount == 1 {
							name = text // 보고내용
						}
						if cellcount == 2 {
							named = text // 바고
						}
					}
					m.items[i-1] = &Foo{Index: i, VoiceScore: score, VoiceTextUser: name, VoiceTextDefault: named}
					//					cellcount = 0
				}
			}
		}
	}

	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()
	m.Sort(m.sortColumn, m.sortOrder)
}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *FooModel) RowCount() int {
	return len(m.items)
}

// 아이템을 화면을
func (m *FooModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index

	case 1:
		return item.VoiceScore

	case 2:
		return item.VoiceTextUser

	case 3:
		return item.VoiceTextDefault

	}

	panic("unexpected col")
}

// Called by the TableView to retrieve if a given row is checked.
func (m *FooModel) Checked(row int) bool {
	return m.items[row].checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *FooModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked
	return nil
}

// Called by the TableView to sort the model.
func (m *FooModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case 0:
			return c(a.Index < b.Index)

		case 1:
			//return c(a.VoicePcm.Before(b.VoicePcm))
			return c(a.VoicePcm < b.VoicePcm)

		case 2:
			return c(a.VoiceScore < b.VoiceScore)

		case 3:
			return c(a.VoiceTextUser < b.VoiceTextUser)

		case 4:
			return c(a.VoiceTextDefault < b.VoiceTextDefault)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func msgbox(title, msgtext string) {
	var mod = syscall.NewLazyDLL("user32.dll")
	var proc = mod.NewProc("MessageBoxW")
	//var MB_YESNOCANCEL = 0x00000003
	var MB_OK = 0x00000000

	//ret, _, _ := proc.Call(0,
	ret, _, _ := proc.Call(0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(msgtext))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(MB_OK))
	// 디버그
	fmt.Printf("Return: %d\n", ret)
}

func fileNotExists(filename string) bool {
    _, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return true
    }
    return false
}
