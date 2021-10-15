package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"unsafe"

	"github.com/FloatTech/ZeroBot-Plugin/utils/dl"
	ctrl "github.com/fumiama/ZeroBot-Hook/control"
	zero "github.com/fumiama/ZeroBot-Hook/hook"
	"github.com/fumiama/ZeroBot-Hook/hook/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

func Inita() {
	// -------------在此下书写插件内容-------------
	err := loadcfg("cfg.pb")
	if err != nil {
		panic(err)
	}
	for i, s := range table {
		index[s] = uint32(i)
	}
	err = os.MkdirAll(base, 0755)
	if err != nil {
		panic(err)
	}
	// 插件主体
	en := ctrl.Register("fortune", &ctrl.Options{
		DisableOnDefault: false,
		Help: "每日运势: \n" +
			"- 运势|抽签\n" +
			"- 设置底图[车万 DC4 爱因斯坦 星空列车 樱云之恋 富婆妹 李清歌 公主连结 原神 明日方舟 碧蓝航线 碧蓝幻想 战双 阴阳师]",
	})
	en.OnRegex(`^设置底图(.*)`).SetBlock(true).SecondPriority().
		Handle(func(ctx *zero.Ctx) {
			gid := ctx.Event.GroupID
			if gid <= 0 {
				// 个人用户设为负数
				gid = -ctx.Event.UserID
			}
			i, ok := index[ctx.State["regex_matched"].([]string)[1]]
			if ok {
				conf.Kind[gid] = i
				savecfg("cfg.pb")
				ctx.SendChain(message.Text("设置成功~"))
			} else {
				ctx.SendChain(message.Text("没有这个底图哦～"))
			}
		})
	en.OnFullMatchGroup([]string{"运势", "抽签"}).SetBlock(true).SecondPriority().
		Handle(func(ctx *zero.Ctx) {
			// 检查签文文件是否存在
			mikuji := base + "运势签文.json"
			if _, err := os.Stat(mikuji); err != nil && !os.IsExist(err) {
				ctx.SendChain(message.Text("正在下载签文文件，请稍后..."))
				err := dl.DownloadTo(site+"运势签文.json", mikuji)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				ctx.SendChain(message.Text("下载签文文件完毕"))
			}
			// 检查字体文件是否存在
			ttf := base + "sakura.ttf"
			if _, err := os.Stat(ttf); err != nil && !os.IsExist(err) {
				ctx.SendChain(message.Text("正在下载字体文件，请稍后..."))
				err := dl.DownloadTo(site+"sakura.ttf", ttf)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				ctx.SendChain(message.Text("下载字体文件完毕"))
			}
			// 获取该群背景类型，默认车万
			kind := "车万"
			gid := ctx.Event.GroupID
			if gid <= 0 {
				// 个人用户设为负数
				gid = -ctx.Event.UserID
			}
			fmt.Println("[fortune]gid:", ctx.Event.GroupID, "uid:", ctx.Event.UserID)
			if v, ok := conf.Kind[gid]; ok {
				kind = table[v]
			}
			// 检查背景图片是否存在
			folder := base + kind
			if _, err := os.Stat(folder); err != nil && !os.IsExist(err) {
				ctx.SendChain(message.Text("正在下载背景图片，请稍后..."))
				zipfile := kind + ".zip"
				zipcache := base + zipfile
				err := dl.DownloadTo(site+zipfile, zipcache)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				ctx.SendChain(message.Text("下载背景图片完毕"))
				err = unpack(zipcache, folder+"/")
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				ctx.SendChain(message.Text("解压背景图片完毕"))
				// 释放空间
				os.Remove(zipcache)
			}
			// 生成种子
			t, _ := strconv.ParseInt(time.Now().Format("20060102"), 10, 64)
			seed := ctx.Event.UserID + t
			// 随机获取背景
			background, err := randimage(folder+"/", seed)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			// 随机获取签文
			title, text, err := randtext(mikuji, seed)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			// 绘制背景
			d, err := draw(background, title, text)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			// 发送图片
			ctx.SendChain(message.Image("base64://" + helper.BytesToString(d)))
		})
	// -------------在此上书写插件内容-------------
}

// 以下勿动
// Hook 改变本插件的环境变量以加载插件
func Hook(botconf interface{}, apicallers interface{}, hooknew interface{},
	matlist interface{}, matlock interface{}, defen interface{},
	reg interface{}, del interface{},
	sndgrpmsg interface{}, sndprivmsg interface{}, getmsg interface{},
	parsectx interface{},
	custnode interface{}, pasemsg interface{}, parsemsgfromarr interface{},
) {
	zero.Hook(botconf, apicallers, hooknew, matlist, matlock, defen)
	rd := getdata(&reg)
	dd := getdata(&del)
	ctrl.Register = *(*(func(service string, o *ctrl.Options) *zero.Engine))(unsafe.Pointer(&rd))
	ctrl.Delete = *(*(func(service string)))(unsafe.Pointer(&dd))
	zero.HookCtx(sndgrpmsg, sndprivmsg, getmsg, parsectx)
	message.HookMsg(custnode, pasemsg, parsemsgfromarr)
	IsHooked = true
	// fmt.Printf("[plugin]set reg: %x, del: %x\n", ctrl.Register, ctrl.Delete)
}

// IsHooked 已经 hook 则不再重复 hook
var IsHooked bool

// 没有方法的interface
type eface struct {
	_type uintptr
	data  unsafe.Pointer
}

func getdata(ptr *interface{}) unsafe.Pointer {
	return (*eface)(unsafe.Pointer(ptr)).data
}
