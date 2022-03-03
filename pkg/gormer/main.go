package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var Viper = viper.New()

func main() {
	var configFilePath = flag.String("c", "./", "config file path")
	flag.Parse()

	Viper.SetConfigName("gormer")        // config file name without file type
	Viper.SetConfigType("yaml")          // config file type
	Viper.AddConfigPath(*configFilePath) // config file path
	if err := Viper.ReadInConfig(); err != nil {
		panic(err)
	} else if err = Viper.UnmarshalKey("database.tables", &tableInfo); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", tableInfo)

	var (
		dsn         = Viper.GetString("database.dsn")
		projectPath = Viper.GetString("project.base")
		goMod       = Viper.GetString("project.go_mod")
		gormPath    = Viper.GetString("project.gorm")
		daoPath     = Viper.GetString("project.dao")
		modelPath   = Viper.GetString("project.model")
		// daoLogPackage = Viper.GetString("project.log.package")
	)

	if projectPath == "" || gormPath == "" || daoPath == "" || modelPath == "" || goMod == "" {
		fmt.Println("dsn,projectPath,gormPath,daoPath,modelPath,goMod 为必填参数，请检查")
		os.Exit(1)
	}

	// 创建文件夹（如果已存在会报错，不影响）
	for _, path := range []string{projectPath + gormPath, projectPath + daoPath, projectPath + modelPath} {
		os.MkdirAll(path, os.ModePerm)
	}

	var (
		tables   []string
		tMatcher map[string]TableInfo
		db       *sql.DB
	)
	if dsn != "" {
		// 连接mysql
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		fmt.Println("start to generate gorm structs")

		// 读取数据库中的表
		tables, err = getAllTables(db)
		if err != nil {
			fmt.Printf("getAllTables error %+v", err)
			os.Exit(1)
		}
	} else {
		for _, v := range tableInfo {
			tables = append(tables, v.Name)
		}
	}
	tMatcher = getTableMatcher()

	for _, table := range tables {
		// 不存在的表直接过滤
		if tMatcher != nil {
			if _, ok := tMatcher[table]; !ok {
				fmt.Printf("table %s ignored\n", table)
				continue
			}
		}

		// 1.生成结构
		structResult, err := Generate(db, table, tMatcher[table])
		if err != nil {
			fmt.Printf("Generate table %s error %+v\n", table, err)
			os.Exit(1)
		}

		// 检查字段的有效性
		fieldMap := make(map[string]FieldLevel)
		for _, column := range structResult.Columns {
			fieldMap[column.GormName] = column
		}
		fmt.Println(fieldMap)
		for i, gen := range tMatcher[table].GenQueries {
			if gen.Fields != "" {
				fields := strings.Split(gen.Fields, ",")
				for _, field := range fields {
					if fieldDetail, ok := fieldMap[strings.TrimSpace(field)]; !ok {
						fmt.Printf("Field %s not match table %s\n", strings.TrimSpace(field), table)
						os.Exit(1)
					} else {
						tMatcher[table].GenQueries[i].GenFields = append(tMatcher[table].GenQueries[i].GenFields, fieldDetail.FieldName)
					}
				}
			}
		}

		// 2.生成gormer file
		if gormPath[len(gormPath)-1] == '/' {
			gormPath = gormPath[:len(gormPath)-1]
		}
		dirs := strings.Split(gormPath, "/")
		header := fmt.Sprintf(gormerHeader, dirs[len(dirs)-1])
		err = parseToFile(projectPath+gormPath, tMatcher[table].Name, header, structResult, parseToTmpl, gormerTmpl)
		if err != nil {
			fmt.Printf("parseToFile error %+v\n", err)
			os.Exit(1)
		}

		if daoPath[len(daoPath)-1] == '/' {
			daoPath = daoPath[:len(daoPath)-1]
		}
		if modelPath[len(modelPath)-1] == '/' {
			modelPath = modelPath[:len(modelPath)-1]
		}

		// 3-1.生成dao file
		dirs = strings.Split(daoPath, "/")
		header = fmt.Sprintf(daoHeader, dirs[len(dirs)-1], goMod, gormPath, goMod, modelPath)
		err = parseToFile(projectPath+daoPath, tMatcher[table].Name, header, structResult, parseToTmpl, daoTmpl)
		if err != nil {
			fmt.Printf("parseToFile error %+v\n", err)
			os.Exit(1)
		}
		// 3-2.生成 dao ext file
		extFile := fmt.Sprintf("%s/%s_ext.go", projectPath+daoPath, tMatcher[table].Name)
		if _, err = os.Stat(extFile); err != nil {
			file, err := os.OpenFile(extFile, os.O_WRONLY+os.O_CREATE+os.O_TRUNC, os.ModePerm)
			if err != nil {
				fmt.Printf("OpenFile error %+v\n", err)
				os.Exit(1)
			}

			_, err = file.WriteString(fmt.Sprintf(daoExtHeader, dirs[len(dirs)-1]))
			if err != nil {
				file.Close()
				fmt.Printf("WriteString error %+v\n", err)
				os.Exit(1)
			}
			file.Close()
		}

		// 4-1.生成model file
		dirs = strings.Split(modelPath, "/")
		header = fmt.Sprintf(modelHeader, dirs[len(dirs)-1], goMod, gormPath)
		err = parseToFile(projectPath+modelPath, tMatcher[table].Name, header, structResult, parseToTmpl, modelTmpl)
		if err != nil {
			fmt.Printf("parseToFile error %+v\n", err)
			os.Exit(1)
		}

		// 4-2.生成 model ext file
		extFile = fmt.Sprintf("%s/%s_ext.go", projectPath+modelPath, tMatcher[table].Name)
		if _, err = os.Stat(extFile); err != nil {
			file, err := os.OpenFile(extFile, os.O_WRONLY+os.O_CREATE+os.O_TRUNC, os.ModePerm)
			if err != nil {
				fmt.Printf("OpenFile error %+v\n", err)
				os.Exit(1)
			}

			_, err = file.WriteString(fmt.Sprintf(modelExtHeader, dirs[len(dirs)-1]) + fmt.Sprintf(modelExtTmpl, structResult.StructName.UpperP))
			if err != nil {
				file.Close()
				fmt.Printf("WriteString error %+v\n", err)
				os.Exit(1)
			}
			file.Close()
		}
		fmt.Printf("Generate Table %s Finished\n", table)
	}

	// deprecated
	// 生成dao层统一的log
	// if daoPath[len(daoPath)-1] == '/' {
	// 	daoPath = daoPath[:len(daoPath)-1]
	// }
	// dirs := strings.Split(daoPath, "/")
	// header := fmt.Sprintf(daoBaseHeader, dirs[len(dirs)-1], daoLogPackage)
	// err = parseToFile(projectPath+daoPath, "gormer_base", header, StructLevel{LogOn: Viper.GetBool("project.log.mode")}, parseToTmpl, daoBaseTmpl)
	// if err != nil {
	// 	fmt.Printf("parseToFile error %+v\n", err)
	// 	os.Exit(1)
	// }

}

func parseToFile(filePath string, fileName string, fileHeader string, structResult StructLevel, parseFunc func(StructLevel, string) (string, error), text string) error {
	result, err := parseFunc(structResult, text)
	if err != nil {
		return errors.Wrapf(err, "parseToDaoTmpl structResult %v", structResult)
	}
	path := fmt.Sprintf("%s/%s.go", filePath, fileName)
	file, err := os.OpenFile(path, os.O_WRONLY+os.O_CREATE+os.O_TRUNC, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "OpenFile path %s", path)
	}
	defer file.Close()

	_, err = file.WriteString(fileHeader + result)
	if err != nil {
		return errors.Wrap(err, "WriteString to file")
	}

	// go fmt files
	exec.Command("go", "fmt", path).Run()
	return nil
}

func parseToTmpl(structData StructLevel, text string) (string, error) {
	tmpl, err := template.New("t").Funcs(template.FuncMap{"counter": counter}).Parse(text)
	// tmpl, err := template.New("t").Parse(text)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, structData)
	return buf.String(), nil
}

// counter 是为了去除一个数组最后一个分隔符的问题，如 1,2,3 不填最后的逗号
func counter() func() int {
	i := -1
	return func() int {
		i++
		return i
	}
}
