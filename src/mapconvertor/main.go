//go:generate goversioninfo -icon=resource/icon.ico -manifest=resource/mapconvertor.exe.manifest
package main

import (
	"datastruct"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reader"
	"reader/parsestif"
	"sendmail"
	"strconv"
	"strings"
	"time"
	"writer/writemap"
	"zipcompress"
)

func getConf(conf *datastruct.Conf) (*os.File, *log.Logger, *os.File, *log.Logger) {
	var infoLogger *log.Logger
	var errLogger *log.Logger
	rootPath := filepath.Dir(os.Args[0])
	year, month, day := time.Now().Date()
	name := fmt.Sprintf("%04d%02d%02d.log", year, month, day)
	nameError := fmt.Sprintf("%04d%02d%02d_error.log", year, month, day)
	errLogFile, err := os.OpenFile(filepath.Join(rootPath, nameError), os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		errLogger.Fatalln(err.Error())
	}
	errLogger = log.New(errLogFile, "ERROR ", log.LstdFlags|log.Lshortfile)
	logFile, err := os.OpenFile(filepath.Join(rootPath, name), os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		errLogger.Fatalln(err.Error())
	}
	infoLogger = log.New(logFile, "INFO ", log.LstdFlags|log.Lshortfile)

	configName := filepath.Join(rootPath, "config.json")
	if _, err = os.Stat(configName); os.IsNotExist(err) {
		defer errLogger.Fatalln("Please configure the path for input, output and backup.")
		conf.InputPath = "."
		conf.OutputPath = "."
		conf.BackupPath = "."
		configFile, err := os.Create(configName)
		if err != nil {
			errLogger.Fatalln(err.Error())
		}
		defer configFile.Close()
		encoder := json.NewEncoder(configFile)
		err = encoder.Encode(conf)
		if err != nil {
			errLogger.Fatalln(err.Error())
		}
	} else {
		configFile, err := os.Open(configName)
		if err != nil {
			errLogger.Fatalln(err.Error())
		}
		defer configFile.Close()
		decoder := json.NewDecoder(configFile)
		err = decoder.Decode(conf)
		if err != nil {
			errLogger.Fatalln(err.Error())
		}
	}

	return logFile, infoLogger, errLogFile, errLogger
}

func main() {
	var err error
	var buf strings.Builder

	conf := datastruct.Conf{}
	logFile, infoLogger, errLogFile, errLogger := getConf(&conf)
	defer func() {
		if buf.String() != "" {
			err := sendmail.SendMail(conf.Host, conf.MailFrom, conf.MailTo, conf.MailType, conf.Subject, buf.String())
			if err != nil {
				errLogger.Println(err.Error())
			}
		}
	}()
	defer logFile.Close()
	defer errLogFile.Close()

	rootDir, err := ioutil.ReadDir(conf.InputPath)
	if err != nil {
		errLogger.Fatalln(err.Error())
	}
	infoLogger.Println("Start converting......")
	listFile := datastruct.ListFile{}
	lotMap := []datastruct.MapContent{}
	waferMap := datastruct.WaferMap{}
	for _, fi := range rootDir {
		if fi.IsDir() {
			lotDir, err := ioutil.ReadDir(filepath.Join(conf.InputPath, fi.Name()))
			if err != nil {
				errLogger.Println(err.Error())
				buf.WriteString(err.Error())
				buf.WriteString("\n")
				continue
			}
			for _, lot := range lotDir {
				if lot.IsDir() {
					compressName := lot.Name()
					waferDir, err := ioutil.ReadDir(filepath.Join(conf.InputPath, fi.Name(), lot.Name()))
					if err != nil {
						errLogger.Println(err.Error())
						buf.WriteString(err.Error())
						buf.WriteString("\n")
						continue
					}
					listFile.WAFERCOUNT = 0
					listFile.TOTALPASS = 0
					listFile.TOTALPASS = 0
					listFile.LISTLINE = make([]datastruct.ListLine, len(waferDir))
					lotMap = make([]datastruct.MapContent, len(waferDir)+1)
					for i, wafer := range waferDir {
						if !wafer.IsDir() && !strings.HasSuffix(wafer.Name(), "lst") {
							infoLogger.Printf("Read map %s ...\n", wafer.Name())
							mapContent, err := reader.Read(filepath.Join(conf.InputPath, fi.Name(), lot.Name(), wafer.Name()))
							if err != nil {
								errLogger.Println(err.Error())
								buf.WriteString(err.Error())
								buf.WriteString("\n")
								continue
							}
							infoLogger.Printf("Parse map %s ...\n", wafer.Name())
							err = parsestif.Parser(mapContent, &waferMap)
							if err != nil {
								errLogger.Println(err.Error())
								buf.WriteString(err.Error())
								buf.WriteString("\n")
								continue
							}
							infoLogger.Printf("Parse map %s complete.\n", wafer.Name())
							infoLogger.Printf("Write map %s-%s.MAP ...\n", waferMap.LOT, waferMap.WAFER)
							// lotMap = append(lotMap, writemap.WriteMap(waferMap, &listFile))
							writemap.WriteMap(waferMap, &(lotMap[i]), &listFile, &(listFile.LISTLINE[i]))
							infoLogger.Printf("Write map %s-%s.MAP complete.\n", waferMap.LOT, waferMap.WAFER)
						}
					}
					// lotMap = append(lotMap, writemap.WriteList(listFile))
					writemap.WriteList(&(lotMap[len(waferDir)]), listFile)
					infoLogger.Printf("Write LIS file %s.LIS complete.\n", listFile.LOT)
					if listFile.WAFERCOUNT != 25 {
						infoLogger.Printf("Wafer Count %d, need to rename folder.", listFile.WAFERCOUNT)
						waferID, err := strconv.Atoi(listFile.LISTLINE[0].WAFERID)
						if err != nil {
							errLogger.Println(err.Error())
							buf.WriteString(err.Error())
							buf.WriteString("\n")
							continue
						}
						compressName = fmt.Sprintf("%sL%c", listFile.LOT, waferID+48)
					}
					infoLogger.Printf("Compress maps for Lot %s...", compressName)
					err = zipcompress.CompressMap(lotMap, fmt.Sprintf("%s/%s.ZIP", conf.OutputPath, compressName))
					if err != nil {
						errLogger.Println(err.Error())
						buf.WriteString(err.Error())
						buf.WriteString("\n")
						continue
					}
					infoLogger.Printf("Compress map %s complete.\n", compressName)
					backupPath := filepath.Join(conf.BackupPath, fi.Name())
					if _, err = os.Stat(backupPath); os.IsNotExist(err) {
						err = os.MkdirAll(backupPath, 0777)
						if err != nil {
							errLogger.Println(err.Error())
							buf.WriteString(err.Error())
							buf.WriteString("\n")
							continue
						}
					}
					err = os.Rename(filepath.Join(conf.InputPath, fi.Name(), lot.Name()), filepath.Join(backupPath, lot.Name()))
					if err != nil {
						errLogger.Println(err.Error())
						buf.WriteString(err.Error())
						buf.WriteString("\n")
						continue
					}
					infoLogger.Printf("Backup Lot %s complete.\n", lot.Name())
				}
			}
		}
	}
	infoLogger.Println("Convert Finished.")
}
