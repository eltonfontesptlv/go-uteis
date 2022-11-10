package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)

type Component struct {
	sql             string
	quantidade      int
	quantidadeTotal int
}

func (c *Component) initReadDir(dir string) ([]fs.FileInfo, error) {
	fmt.Println("Iniciando a leitura dos arquivos...")

	files, err := ioutil.ReadDir(dirFilesInput)
	if err != nil {
		return nil, err
	}

	fmt.Println(len(files), "arquivo(s) encontrado(s)")
	fmt.Println("")

	return files, nil
}

func (c *Component) trimSpace() {
	c.sql = strings.TrimSpace(c.sql)
}

func (c *Component) writeFile(nameFile string) error {
	if c.sql == "" {
		return errors.New("Nenhuma query criada")
	}
	c.trimSpace()
	data := []byte(c.sql)
	err := os.WriteFile(dirFilesOutput+"/"+nameFile, data, 0777)
	if err != nil {
		return err
	}
	return nil
}

func (c *Component) readFiles(files []os.FileInfo) error {

	for l, file := range files {
		arq, err := os.Open(dirFilesInput + "/" + file.Name())
		if err != nil {
			panic(err)
		}
		defer arq.Close()

		records, err := csv.NewReader(arq).ReadAll()
		if err != nil {
			return err
		}

		c.quantidade = 0
		for k, record := range records {
			if k == 0 || record[42] == "" {
				continue
			}
			c.sql += fmt.Sprintf("UPDATE payments SET acquirer_reference='%s' WHERE psp_reference_id='%s';\n", record[42], record[2])
			c.quantidadeTotal++
			c.quantidade++
		}

		fmt.Println("Arquivo", (l + 1), "com", c.quantidade, "registro(s)")
	}
	fmt.Println("Foram registrada(s)", c.quantidadeTotal, "querys")
	return nil

}

const dirFilesInput = "./files/input"
const dirFilesOutput = "./files/output"

func main() {

	component := Component{}
	files, err := component.initReadDir(dirFilesInput)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(-1)
	}

	err = component.readFiles(files)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}

	err = component.writeFile("query.sql")
	if err != nil {
		fmt.Println("Error writing file:", err)
	}
}
