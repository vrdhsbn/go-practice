package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"text/template"
)

// 構造体の定義
type Page struct {
	Title string
	Body  []byte
}

// ページを保存するメソッド
// 関数名の前の(p *Page)はレシーバである。
// Page構造体からこのメソッドを呼べるようにしているってこと。p1.save()みたいに呼べる。
func (p *Page) save() error {
	filename := "data/" + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

// ページを読み込むメソッド
func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	// こういうエラーハンドリングの書き方はよく出てくるのでイディオムとして覚えておく。
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// 要調査：構造体そのものではなくポインタで返す理由って？
	// ふつうにリターンすると構造体のコピー（実体とは別物）が作られてしまう？
	return &Page{Title: title, Body: body}, nil
}

// テンプレートを読み込んでグローバル変数に格納する
// これを使うことでrenderTemplate実行時に毎回ファイルを読み込むのを避けることができる
var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))

// テンプレートファイルを出力する関数
// レスポンス作成時にエラーがあったら500エラーを返す
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ハンドラ関数。
// ページが存在しなければ新規作成できるよう編集画面にリダイレクトする
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// 編集用画面を出力するハンドラ関数。
// titleに相当するファイルが存在しなければ新規ページを用意する
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// 編集した内容を保存するハンドラ関数。
// 保存時にエラーがあれば500エラーを返す。
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	// FormValueの値はstringなのでbyteスライスに変換する必要がある
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// URLのバリデーションのためのグローバル変数
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// ハンドラ関数を返すラッパ関数
// この形で返された関数はクロージャというらしい
// 引数fnに、これまで定義してきたsave/edit/viewハンドラが渡される
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URLのパスがvalidでなければ404エラーを返す
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// ToDo:練習として全体をGinで書き直してみたい
func main() {
	// それぞれのルートにアクセスした際のハンドラ関数を紐付ける
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	// サーバを起動（エラーが出たらログを出力するためにlogでラップしている）
	log.Fatal(http.ListenAndServe(":8080", nil))
}
