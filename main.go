package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func DosyaOku(tür string) (string, error) {
	dosyaAdi := "templates/stil/" + tür + ".txt"
	dosya, err := os.ReadFile(dosyaAdi)
	if err != nil {
		return "File Read Error ", err
	}
	return string(dosya), err
}

func YazdirVeBul(satirlar []string, kelime string, output http.ResponseWriter) {
	for h := 1; h < 9; h++ {
		for i := 0; i < len(kelime); i++ {
			for satirIndex, satir := range satirlar {
				if satirIndex == (int(kelime[i])-32)*9+h {

					fmt.Fprint(output, satir)
				}
			}
		}
		fmt.Fprintln(output)
	}
}

func asciiHandler(output http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodGet {
		http.Error(output, "Allowed method", http.StatusMethodNotAllowed)
		return
	}

	if request.URL.Path != "/sayfa_cekme" {
		http.Error(output, "Page not found", http.StatusNotFound)
		return
	}

	metin := request.FormValue("metin")
	tür := request.FormValue("secim")

	if tür == "" && metin == "" {
		http.Error(output, " Text and Type Cannot Be Blank", http.StatusBadRequest)
		return
	}

	if metin == "" {
		http.Error(output, " Text  Cannot be left blank", http.StatusBadRequest)
		return
	}
	if tür == "" {
		http.Error(output, "Type Cannot be left blank", http.StatusBadRequest)
		return
	}

	for _, harf := range metin {
		if harf >= 1 && harf <= 127 {
			continue
		} else {
			http.Error(output, "Invalid character in text", http.StatusBadRequest)
			return
		}
	}

	ifade := strings.Split(metin, "\n")

	satirlar, err := DosyaOku(tür)
	if err != nil {
		http.Error(output, "File Read Error", http.StatusInternalServerError)
		return
	}

	for i, kelime := range ifade {
		if kelime == "" {
			if i != 0 {
				fmt.Fprintln(output)
			}
			continue
		}

		YazdirVeBul(strings.Split(satirlar, "\n"), kelime, output)

		if err != nil {
			http.Error(output, fmt.Sprintf("ERROR", err), http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	http.HandleFunc("/sayfa_cekme", asciiHandler)
	http.Handle("/", http.FileServer(http.Dir("templates")))
	port := ":8080"
	fmt.Printf("The server is listening on port %s...\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
