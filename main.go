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

var bot *tgbotapi.BotAPI

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	data, err := ioutil.ReadFile("D:/botToken.txt")
	bot, err = tgbotapi.NewBotAPI(string(data))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		userName := update.Message.From.UserName
		if !(userName == "EnsPzr") {
			MesajGonder(update.Message.Chat.ID,"Yetkiniz Bulunmamaktadır.")
			continue
		}
		gelenMesaj := update.Message.Text
		if (gelenMesaj == "/ekranresmi") {
			DosyalariSil()
			n := screenshot.NumActiveDisplays()
			for i := 0; i < n; i++ {
				EkranResmiAlVeGonder(update.Message.Chat.ID, i);
			}
		} else if (gelenMesaj == "/kilitle") {
			sonuc := KomutCalistir("/C", "kilitle.bat");
			if (sonuc == nil) {
				MesajGonder(update.Message.Chat.ID, "Kilitlendi...");
			} else {
				MesajGonder(update.Message.Chat.ID, "Kilitleme İşlemi Sırasında Hata Oluştu."+sonuc.Error());
			}
		} else if (gelenMesaj == "/kapat") {
			sonuc := KomutCalistir("/C", "shutdown", "/s");
			if (sonuc == nil) {
				MesajGonder(update.Message.Chat.ID, "30 Saniye İçerisinde Bilgisayar Kapatılıyor...");
			} else {
				MesajGonder(update.Message.Chat.ID, "Kapatma İşlemi Sırasında Hata Oluştu."+err.Error());
			}
		} else if (strings.HasPrefix(gelenMesaj, "/kapat ")) {
			sure := strings.Replace(gelenMesaj, "/kapat ", "", -1)
			sonuc := KomutCalistir("/C", "shutdown", "/s", "/t", sure);
			if (sonuc == nil) {
				MesajGonder(update.Message.Chat.ID, sure+" Saniye İçerisinde Bilgisayar Kapatılıyor...");
			} else {
				MesajGonder(update.Message.Chat.ID, "Kapatma İşlemi Sırasında Hata Oluştu."+err.Error());
			}
		} else if (gelenMesaj == "/yenidenbaslat") {
			sonuc := KomutCalistir("/C", "shutdown", "/r");
			if (sonuc == nil) {
				MesajGonder(update.Message.Chat.ID, "Bilgisayar Yeniden Başlatılıyor...");
			} else {
				MesajGonder(update.Message.Chat.ID, "Yeniden Başlatma İşlemi Sırasında Hata Oluştu."+err.Error());
			}
		} else if (gelenMesaj == "/iptal") {
			sonuc := KomutCalistir("/C", "shutdown", "/a");
			if (sonuc == nil) {
				MesajGonder(update.Message.Chat.ID, "İşlemler İptal Edildi...");
			} else {
				MesajGonder(update.Message.Chat.ID, "İşlem İptali Sırasında Hata Oluştu."+err.Error());
			}
		} else {
			MesajGonder(update.Message.Chat.ID, "Komut Bulunamadı...");
		}

	}
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {
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

func DosyalariSil() {
	files, err := ioutil.ReadDir("./resimler")
	if err != nil {
		fmt.Println(err)
	}

	for _, f := range files {
		fmt.Println("%s", f.Name())
		os.Remove("./resimler/" + f.Name())
	}
	fmt.Println("Resimler Silindi")
}

func MesajGonder(chatId int64, metin string) {
	if metin == "" {
		return
	}
	msg := tgbotapi.NewMessage(chatId, metin)
	bot.Send(msg)
}

func ResimUploadEt(chatId int64, resimYolu string) string {
	resim := tgbotapi.NewPhotoUpload(chatId, resimYolu)
	bot.Send(resim)
	return resim.FileID;
}

func ResimGonder(chatId int64, resimId string) {
	msg := tgbotapi.NewPhotoShare(chatId, resimId)
	bot.Send(msg)
}

func EkranResmiAlVeGonder(chatId int64, ekranNo int) {
	bounds := screenshot.GetDisplayBounds(ekranNo)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		MesajGonder(chatId, "Ekran resim alımı sırasında hata=>"+err.Error());
	}
	fileName := "./resimler/" + fmt.Sprintf("%d_%dx%d.png", ekranNo, bounds.Dx(), bounds.Dy())
	file, err := os.Create(fileName)
	if (err != nil) {
		MesajGonder(chatId, "Ekran resim gönderimi sırasında hata=>"+err.Error());
	} else {
		defer file.Close()
		png.Encode(file, img)
		resimId := ResimUploadEt(chatId, fileName);
		ResimGonder(chatId, resimId)
	}
}

func KomutCalistir(komutlar ...string) (cevap error) {
	var komut string;
	for i := 0; i < len(komutlar); i++ {
		komut = komut + " " + komutlar[i]
	}
	cmd := exec.Command("cmd.exe", komut)
	if err := cmd.Run();
		err != nil {
		return err;
	}
	return nil;
}
