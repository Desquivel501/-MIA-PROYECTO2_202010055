package analizador

import "MIA_PROYECTO2_202010055/comandos"
import (  
    "fmt"
	"strings"
	"strconv"
	
)

type analizador struct {  
    codigo   string
	cmd comandos.Comandos
}

func New(codigo string) analizador{
	var cmd comandos.Comandos 
	a := analizador {codigo, cmd}
	return a
}

func (a analizador) Imprimir() {
	fmt.Println(a.codigo)
}

func Split_txt(texto string) []string{
	name := ""
	path := ""
	var splited []string
	first := 0
	last := 0

	first = strings.Index(texto, "-path=\"")
	if (first > -1){
		last = strings.Index(texto[first+7:], "\"") + first + 7
		path = texto[first:last+1]
		texto = strings.Replace(texto, path, "", 1)
	}

	first = strings.Index(texto, "-name=\"")
	if (first > -1){
		last = strings.Index(texto[first+7:], "\"") + first + 7
		name = texto[first:last+1]
		texto = strings.Replace(texto, name, "", 1)
	}

	splited = strings.Split(texto, " ")
	if (path != ""){splited = append(splited, path)}
	if (name != ""){splited = append(splited, name)}

	return splited
}

func (a *analizador) Analizar(texto string){
	
	texto = strings.Replace(texto, "\n", "", 1)

	cmd_list := Split_txt(texto)
	var parametros []string
	comando := ""

	for i := 0; i < len(cmd_list); i++ {
        if(i == 0){
            comando = cmd_list[0]
        }else{
            parametros = append(parametros, cmd_list[i])
        }
    }
	a.Identificar(comando, parametros)
}


func (a *analizador) Identificar(comando string, parametros []string){
	if comando == "mkdisk"{
		fmt.Println("Comando mkdisk")
		size := -1
		fit := "FF"
		unit := "M"
		path := ""
		
		for i := 0; i < len(parametros); i++ {
			param := parametros[i]

			if (strings.Index(param, "-size=") == 0) {
				param = strings.Replace(param, "-size=", "", 1)
				if s, err := strconv.Atoi(param); err == nil {
					size = s
				}
				fmt.Println("Size: ",size)
			}

			if (strings.Index(param, "-fit=") == 0) {
				param = strings.Replace(param, "-fit=", "", 1)
				fit = param
				fmt.Println("Fit: ",fit)
			}

			if (strings.Index(param, "-unit=") == 0) {
				param = strings.Replace(param, "-unit=", "", 1)
				unit = param
				fmt.Println("Unit: ",unit)
			}

			if (strings.Index(param, "-path=") == 0) {
				param = strings.Replace(param, "-path=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				path = param
				fmt.Println("Path: ",path)
			}
		}

		if path == ""{
			fmt.Println("[MIA]@Proyecto2:~$ No se ha ingresado la ruta")
			return
		}

		if (fit != "FF") && (fit != "BF") && (fit != "WF"){
			fmt.Println(fit)
			fmt.Println("[MIA]@Proyecto2:~$ Fit incorrecto")
			return
		}

		if (unit != "M") && (unit != "K"){
			fmt.Println("[MIA]@Proyecto2:~$ Dimensional incorrecto")
			return
		}

		if (size <= 0){
			fmt.Println("[MIA]@Proyecto2:~$ Tamaño del disco incorrecto")
			return
		}

		a.cmd.Mkdisk(size, fit[1], unit[0], path)
	}

	if comando == "rmdisk"{
		fmt.Println("Comando rmdisk")
		path := ""
		
		for i := 0; i < len(parametros); i++ {
			param := parametros[i]
			if (strings.Index(param, "-path=") == 0) {
				param = strings.Replace(param, "-path=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				path = param
				fmt.Println("Path: ",path)
			}
		}
	}

	if comando == "fdisk"{
		fmt.Println("Comando fdisk")
		size := -1
		fit := "FF"
		unit := "M"
		path := ""
		name := ""
		type_ := ""
		
		for i := 0; i < len(parametros); i++ {
			param := parametros[i]

			if (strings.Index(param, "-size=") == 0) {
				param = strings.Replace(param, "-size=", "", 1)
				if s, err := strconv.Atoi(param); err == nil {
					size = s
				}else{
					fmt.Println(err)
				}
				fmt.Println("Size: ",size)
			}

			if (strings.Index(param, "-fit=") == 0) {
				param = strings.Replace(param, "-fit=", "", 1)
				fit = param
				fmt.Println("Fit: ",fit)
			}

			if (strings.Index(param, "-unit=") == 0) {
				param = strings.Replace(param, "-unit=", "", 1)
				unit = param
				fmt.Println("Unit: ",unit)
			}

			if (strings.Index(param, "-path=") == 0) {
				param = strings.Replace(param, "-path=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				path = param
				fmt.Println("Path: ",path)
			}

			if (strings.Index(param, "-name=") == 0) {
				param = strings.Replace(param, "-name=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				name = param
				fmt.Println("Name: ",name)
			}

			if (strings.Index(param, "-type=") == 0) {
				param = strings.Replace(param, "-type=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				type_ = param
				fmt.Println("type: ",type_)
			}
		}
		
		if path == ""{
			fmt.Println("[MIA]@Proyecto2:~$ No se ha ingresado la ruta")
			return
		}

		if (fit != "FF") && (fit != "BF") && (fit != "WF"){
			fmt.Println(fit)
			fmt.Println("[MIA]@Proyecto2:~$ Fit incorrecto")
			return
		}

		if (unit != "M") && (unit != "K"){
			fmt.Println("[MIA]@Proyecto2:~$ Dimensional incorrecto")
			return
		}

		if (size <= 0){
			fmt.Println("[MIA]@Proyecto2:~$ Tamaño de la particion incorrecta incorrecto")
			return
		}
		
		a.cmd.Fdisk(size, fit[1], unit[0], path, type_[0],name)
	}
}
