package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Alphabet map[string]int

type Cipher struct {
	Name          string
	Desc          string
	CaseSensitive bool
	Letters       Alphabet
}

var (
	ciphers = make(map[string]Cipher)

	ciphername = flag.String("c", "aq36", "Cipher to use")
	listFiles  = flag.Bool("l", false, "List available ciphers")
	viewCipher = flag.Bool("v", false, "View info about a cipher")
	textQuery  = flag.Bool("q", false, "Query the database with words")
	numQuery   = flag.Bool("n", false, "Query the databse with a number")
	noSave     = flag.Bool("x", false, "Don't save words to the database")

	cfgdir string

	//go:embed ciphers
	content embed.FS
)

func main() {
	cfgdir, _ = os.UserConfigDir()
	cfgdir += "/gomatria/"

	addAlphabets()
	flag.Parse()

	args := flag.Args()
	values := make([]int, 0)

	if *listFiles {
		listAlphas()
	}

	ciph, exist := ciphers[*ciphername]

	if !exist {
		fmt.Printf("Cipher \"%s\" does not exist.\n", *ciphername)
		os.Exit(1)
	}

	if *viewCipher {
		seeCipher(ciph)
	}

	if *numQuery {
		num, errConv := strconv.ParseInt(args[0], 10, 64)
		if errConv != nil {
			fmt.Println("Error: The -n flag only allows numbers.")
			os.Exit(1)
		}
		queryDB(int(num), ciph)
		os.Exit(0)
	}

	for _, arg := range args {

		if !ciph.CaseSensitive {
			arg = strings.ToUpper(arg)
		}

		val := AqCalc(arg, ciph)
		values = append(values, val)

		fmt.Printf("%s = %v\n", arg, val)

		if !*noSave {
			saveDB(ciph, arg)
		}
	}

	if *textQuery {
		for _, num := range removeDuplicate(values) {
			fmt.Println("")
			queryDB(num, ciph)
		}
	}
}

func removeDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func AqCalc(text string, ciph Cipher) (result int) {
	for _, letter := range strings.Split(text, "") {
		value, ok := ciph.Letters[letter]
		if ok {
			result += value
		} else {
			switch letter {
			case "0":
				result += 0
			case "1":
				result += 1
			case "2":
				result += 2
			case "3":
				result += 3
			case "4":
				result += 4
			case "5":
				result += 5
			case "6":
				result += 6
			case "7":
				result += 7
			case "8":
				result += 8
			case "9":
				result += 9
			default:
				result += value
			}
		}
	}

	return
}

func readCipher(filename string) Cipher {
	contents, err := os.ReadFile(filename)
	handleError(err)

	var ciph Cipher
	err = json.Unmarshal(contents, &ciph)
	handleError(err)

	return ciph
}

func addAlphabets() {
	errCreate := os.Mkdir(cfgdir, os.ModePerm)

	if errCreate != nil && !os.IsExist(errCreate) {
		log.Fatal(errCreate)
	}

	basefiles, errReadBase := content.ReadDir("ciphers")
	handleError(errReadBase)

	for _, bfile := range basefiles {
		name := bfile.Name()
		jsony, _ := content.ReadFile("ciphers/" + name)

		var ciph Cipher
		_ = json.Unmarshal(jsony, &ciph)
		ciphers[ciph.Name] = ciph
	}

	customfiles, errReadCustom := os.ReadDir(cfgdir)
	handleError(errReadCustom)

	for _, file := range customfiles {
		name := file.Name()
		if !file.IsDir() && strings.Contains(name, ".json") {
			ciph := readCipher(cfgdir + name)
			ciphers[ciph.Name] = ciph
		}
	}
}

func listAlphas() {
	fmt.Println("Available ciphers:")

	var names []string

	for ciph := range ciphers {
		names = append(names, ciph)
	}

	sort.Strings(names)

	for _, name := range names {
		fmt.Println(name)
	}

	os.Exit(0)
}

func seeCipher(ciph Cipher) {
	fmt.Printf("Cipher %s:\n", ciph.Name)
	fmt.Println(ciph.Desc)
	fmt.Println("\nAlphabet used:")

	strs := make([]string, 0)

	for letter := range ciph.Letters {
		strs = append(strs, letter)
	}

	sort.Strings(strs)

	for _, letter := range strs {
		fmt.Printf("%s = %v\n", letter, ciph.Letters[letter])
	}

	os.Exit(0)
}

func saveDB(ciph Cipher, text string) {
	db, err := sql.Open("sqlite3", cfgdir+"gomatria.db")
	handleError(err)
	defer db.Close()

	entry := text

	entryval := AqCalc(text, ciph)

	dbname := fmt.Sprintf("%s_%v", ciph.Name, entryval)

	table := "CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE, entry VARCHAR(250) UNIQUE);"
	table = fmt.Sprintf(table, dbname)
	insert := "INSERT OR IGNORE INTO %s (entry) VALUES (@val);"
	insert = fmt.Sprintf(insert, dbname)

	_, errTable := db.Exec(table)
	handleError(errTable)

	_, errInsert := db.Exec(insert, sql.Named("val", entry))
	handleError(errInsert)
}

func queryDB(num int, ciph Cipher) {
	db, err := sql.Open("sqlite3", cfgdir+"gomatria.db")
	handleError(err)
	defer db.Close()

	dbname := fmt.Sprintf("%s_%v", ciph.Name, num)

	rows, errQuery := db.Query(fmt.Sprintf("SELECT entry FROM %s;", dbname))

	if errQuery != nil {
		fmt.Printf("There are no entires with value %v for cipher %s.\n", num, ciph.Name)
		os.Exit(1)
	}

	entries := make([]string, 0)

	for rows.Next() {
		var entry string
		errScan := rows.Scan(&entry)
		handleError(errScan)
		entries = append(entries, entry)
	}

	sort.Strings(entries)

	if len(entries) == 0 {
		fmt.Printf("There are no entires with value %v for cipher %s.\n", num, ciph.Name)
		os.Exit(1)
	}

	fmt.Printf("Results for %v in cipher %s:\n", num, ciph.Name)

	for _, entry := range entries {
		fmt.Println("-", entry)
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
