package utils

import (
	"fmt"
	"io"
	"time"

	"Alice-Seahat-Healthcare/seahat-be/entity"

	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GenerateMedicalCertificate(file io.Writer, telemedicine entity.Telemedicine) error {
	caser := cases.Title(language.Und)
	name := caser.String(telemedicine.User.Name)
	dateOfBirth := telemedicine.User.DateOfBirth.Format("02 January 2006")
	gender := caser.String(telemedicine.User.Gender)
	diagnose := *telemedicine.Diagnose
	restAt := telemedicine.StartRestAt.Local()
	restAtFormated := restAt.Format("02 january 2006")
	restDuration := *telemedicine.RestDuration
	doctorname := caser.String(telemedicine.Doctor.Name)
	doctorSpecialization := caser.String(telemedicine.Doctor.Specialization.Name)
	restEndAt := restAt.Add(time.Duration(restDuration * 24 * int(time.Hour))).Format("02 january 2006")
	yearAge := GetAge(telemedicine.User.DateOfBirth)

	marginX := 13.6
	marginY := 13.6
	pdf := gofpdf.New("L", "mm", "A5", "")
	pdf.SetMargins(marginX, marginY, marginX)
	tpl1 := gofpdi.ImportPage(pdf, "assets/pdf/template.pdf", 1, "/MediaBox")
	pdf.AddPage()

	pdf.SetFillColor(124, 252, 0)
	gofpdi.UseImportedTemplate(pdf, tpl1, 0, 0, 210, 148)
	var newTab float64 = 53
	pdf.SetFont("arial", "B", 12)
	pdf.SetXY(marginX, 30)
	pdf.CellFormat(0, 0, "Informasi Pasien", "", 0, "C", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetXY(newTab, pdf.GetY()+15)
	pdf.Cell(0, 0, "Nama")
	pdf.SetXY(newTab+30, pdf.GetY())
	pdf.MultiCell(106, 0, fmt.Sprintf(":   %s", name), "", "L", false)
	pdf.SetXY(newTab, pdf.GetY()+8)
	pdf.Cell(0, 0, "Tanggal Lahir")
	pdf.SetXY(newTab+30, pdf.GetY())
	pdf.MultiCell(106, 0, fmt.Sprintf(":   %s", dateOfBirth), "", "L", false)
	pdf.SetXY(newTab, pdf.GetY()+8)
	pdf.Cell(0, 0, "Jenis Kelamin")
	pdf.SetXY(newTab+30, pdf.GetY())
	pdf.MultiCell(106, 0, fmt.Sprintf(":   %s", gender), "", "L", false)
	pdf.SetXY(newTab, pdf.GetY()+8)
	pdf.Cell(0, 0, "Umur")
	pdf.SetXY(newTab+30, pdf.GetY())
	pdf.MultiCell(106, 0, fmt.Sprintf(":   %d tahun", yearAge), "", "L", false)

	if restDuration == 0 {
		pdf.SetXY(marginX, pdf.GetY()+15)
		pdf.MultiCell(182, 5, fmt.Sprintf("Berdasarkan anamnesa yang dilakukan selama konsultasi, pasien didiagnosis mengalami %s. Bila obat telah habis dan keadaan pasien belum membaik, dimohon untuk segera mengunjungi fasilitas kesehatan terdekat", diagnose), "", "J", false)
		pdf.SetXY(pdf.GetX(), pdf.GetY()+15)
		pdf.SetFont("arial", "B", 12.5)
		pdf.MultiCell(182, 4, doctorname, "", "R", false)
		pdf.SetFont("arial", "", 10)
		pdf.SetXY(pdf.GetX(), pdf.GetY()+1)
		pdf.MultiCell(182, 4, doctorSpecialization, "", "R", false)
	} else {
		pdf.SetXY(marginX, pdf.GetY()+15)
		pdf.MultiCell(182, 5, fmt.Sprintf("Berdasarkan anamnesa yang dilakukan selama konsultasi, pasien didiagnosis mengalami %s. Saya sebagai dokter merekomendasi untuk istirahat selama durasi %d hari dari tanggal %s hingga tanggal %s. Apabila dalam durasi tersebut keadaan pasien belum membaik, dimohon untuk segera mengunjungi fasilitas kesehatan terdekat", diagnose, restDuration, restAtFormated, restEndAt), "", "J", false)
		pdf.SetXY(pdf.GetX(), pdf.GetY()+5)
		pdf.SetFont("arial", "B", 12.5)
		pdf.MultiCell(182, 4, doctorname, "", "R", false)
		pdf.SetFont("arial", "", 10)
		pdf.SetXY(pdf.GetX(), pdf.GetY()+1)
		pdf.MultiCell(182, 4, doctorSpecialization, "", "R", false)
	}

	err := pdf.Output(file)
	if err != nil {
		return err
	}

	return nil

}

func GetAge(dateOfBirth time.Time) int {
	curentYear, currentMonth, currentDay := time.Now().Local().Date()
	yearOfBirth, monthOfBirth, dayOfBirth := dateOfBirth.Date()
	yearAge := curentYear - yearOfBirth
	monthAge := currentMonth - monthOfBirth
	dayAge := currentDay - dayOfBirth
	if monthAge < 0 || dayAge < 0 && yearAge != 0 {
		yearAge = yearAge - 1
	}
	return yearAge
}

func GeneratePrescription(file io.Writer, telemedicine entity.Telemedicine) error {
	marginX := 13.6
	marginY := 26.1
	pageNumber := 1
	timeNow := time.Now().Local()

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginX, marginY, marginX)
	tpl1 := gofpdi.ImportPage(pdf, "assets/pdf/template_prescription.pdf", 1, "/MediaBox")
	pdf.AddPage()
	gofpdi.UseImportedTemplate(pdf, tpl1, 0, 0, 210, 297)
	prescriptionHeader(pdf, pageNumber, telemedicine, timeNow)

	for _, prescription := range telemedicine.Prescriptions {
		if int(pdf.GetY()) >= 250 {
			pdf.SetXY(0, 0)
			pdf.AddPage()
			pageNumber++
			gofpdi.UseImportedTemplate(pdf, tpl1, 0, 0, 210, 297)
			prescriptionHeader(pdf, pageNumber, telemedicine, timeNow)
		}

		pdf.SetFont("arial", "B", 11)
		pdf.SetXY(22, pdf.GetY()+9.3)
		pdf.MultiCell(210, 0, fmt.Sprintf("R/ %s - %d %s", prescription.Drug.Name, prescription.Quantity, prescription.Drug.SellingUnit), "", "L", false)
		pdf.SetXY(22, pdf.GetY()+4)
		pdf.SetFont("arial", "", 10)
		pdf.MultiCell(210, 0, prescription.Notes, "", "L", false)

	}

	err := pdf.Output(file)
	if err != nil {
		return err
	}

	return nil

}

func prescriptionHeader(pdf *gofpdf.Fpdf, pageNumber int, telemedicine entity.Telemedicine, timeNow time.Time) {
	marginX := 13.6
	caser := cases.Title(language.Und)
	pdf.SetFont("arial", "B", 16)
	pdf.SetXY(marginX, 19)
	pdf.SetXY(20.14, pdf.GetY()+1.8)
	pdf.MultiCell(106, 0, fmt.Sprintf("Dokter %s", caser.String(telemedicine.Doctor.Name)), "", "L", false)
	pdf.SetFont("arial", "", 12)
	pdf.SetXY(20, pdf.GetY()+6)
	pdf.MultiCell(106, 0, caser.String(telemedicine.Doctor.Specialization.Name), "", "L", false)
	pdf.SetFont("arial", "", 11)
	pdf.SetXY(20, pdf.GetY()+18)
	pdf.MultiCell(106, 0, fmt.Sprintf("Halaman %d", pageNumber), "", "L", false)
	pdf.SetXY(82.5, pdf.GetY())
	pdf.MultiCell(106, 0, timeNow.Format("02 January 2006"), "", "R", false)
	pdf.SetFont("arial", "B", 16)
	pdf.SetXY(0, pdf.GetY()+10.1)
	pdf.MultiCell(210, 0, "Resep Digital", "", "C", false)
	pdf.SetFont("arial", "", 11)
	pdf.SetXY(0, pdf.GetY()+9.6)
	pdf.MultiCell(210, 0, "Jika kondisi anda tidak membaik, silakan untuk mengunjungi fasilitas terdekat secepat mungkin", "", "C", false)
	pdf.SetXY(0, pdf.GetY()+7)
	pdf.SetFont("arial", "B", 11)
	pdf.MultiCell(210, 0, fmt.Sprintf("Pasien: %s", caser.String(telemedicine.User.Name)), "", "C", false)
	pdf.SetXY(0, pdf.GetY()+6)

}
