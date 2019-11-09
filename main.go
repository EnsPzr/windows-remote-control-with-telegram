package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kardianos/service"
	"github.com/kbinani/screenshot"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) run() {
	data, err := ioutil.ReadFile("D:/botToken.txt")
	bot, err := tgbotapi.NewBotAPI(string(data))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		userName:=update.Message.From.UserName
		if(userName == "EnsPzr"){
			gelenMesaj:=update.Message.Text
			if(gelenMesaj =="/ekranresmi"){
				DosyalariSil()
				n := screenshot.NumActiveDisplays()
				for i := 0; i < n; i++ {
					bounds := screenshot.GetDisplayBounds(i)
					fmt.Println("lahmacun",bounds.Dx(), bounds.Dy())
					img, err := screenshot.CaptureRect(bounds)
					if err != nil {
						panic(err)
					}
					fileName :="resimler/" +fmt.Sprintf("%d_%dx%d.png", i, bounds.Dx(), bounds.Dy())
					file, _ := os.Create(fileName)
					defer file.Close()
					png.Encode(file, img)
					resim:=tgbotapi.NewPhotoUpload(update.Message.Chat.ID,fileName)
					bot.Send(resim)
					resimId:=resim.FileID
					fmt.Println("resimId=>",resimId)
					msg := tgbotapi.NewPhotoShare(update.Message.Chat.ID, resimId)
					bot.Send(msg)
				}
			} else if(gelenMesaj=="/kilitle"){
				app := "rundll32.exe user32.dll,LockWorkStation"
				if err := exec.Command("cmd", "/C", app).Run(); err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Kilitleme İşlemi Sırasında Hata Oluştu."+err.Error())
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Kilitlendi...")
					bot.Send(msg)
				}
			} else if(gelenMesaj =="/kapat"){
				if err := exec.Command("cmd", "/C", "shutdown", "/s").Run(); err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Kapatma İşlemi Sırasında Hata Oluştu."+err.Error())
					bot.Send(msg)
				} else{
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "30 Saniye İçerisinde Bilgisayar Kapatılıyor...")
					bot.Send(msg)
				}
			} else if(strings.HasPrefix(gelenMesaj,"/kapat ")){
				sure:= strings.Replace(gelenMesaj,"/kapat ","",-1)
				fmt.Println("Buraya Girdi"+sure)
				if err := exec.Command("cmd", "/C", "shutdown", "/s","/t",sure).Run(); err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Kapatma İşlemi Sırasında Hata Oluştu."+err.Error())
					bot.Send(msg)
				} else{
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, sure +" Saniye İçerisinde Bilgisayar Kapatılıyor...")
					bot.Send(msg)
				}
			} else if(gelenMesaj=="/yenidenbaslat"){
				if err := exec.Command("cmd", "/C", "shutdown", "/r").Run(); err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Yeniden Başlatma İşlemi Sırasında Hata Oluştu."+err.Error())
					bot.Send(msg)
				} else{
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Bilgisayar Yeniden Başlatılıyor...")
					bot.Send(msg)
				}
			}else if(gelenMesaj=="/iptal"){
				if err := exec.Command("cmd", "/C", "shutdown", "/a").Run(); err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "İşlem İptali Sırasında Hata Oluştu."+err.Error())
					bot.Send(msg)
				} else{
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "İşlemler İptal Edildi...")
					bot.Send(msg)
				}
			}	else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Komut Bulunamadı")
				bot.Send(msg)
			}
		}else{
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Yetkiniz Bulunmamaktadır")
			bot.Send(msg)
		}
	}
}
func (p *program) Stop(s service.Service) error {
	return nil
}
func main(){
	svcConfig := &service.Config{
		Name:        "uzaktankontrolbot",
		DisplayName: "Uzaktan Kontrol Telegram Bot",
		Description: "Bilgisayarı Uzaktan Kontrol Etmeye Yarayan Bot",
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	var argsWithoutProgName = os.Args[1:]
	for _, param := range argsWithoutProgName {
		if param == "install" {
			err = s.Install()
			if err != nil {
				fmt.Println("hata")
				fmt.Println("Elhamdülillah")
			}
			return
		} else if param == "uninstall" {
			err = s.Uninstall()
			if err != nil {
				fmt.Println("hata")
			} else {
				fmt.Println("Elhamdülillah")
			}
			return
		}
	}
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func DosyalariSil(){
	files, err := ioutil.ReadDir("./resimler")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println("%s",f.Name())
		os.Remove("resimler/"+f.Name())
	}
	fmt.Println("Resimler Silindi")
}