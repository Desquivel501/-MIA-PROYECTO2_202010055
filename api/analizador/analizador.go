package analizador

import "MIA_PROYECTO2_202010055/comandos"
import (  
    "fmt"
	"strings"
	"strconv"
	"bufio"
	"os"
	"net/http"
    "github.com/gin-gonic/gin"
)

type Response1 struct{
    Resultado string `json:"resultado"`
	Reporte string `json:"reporte"`
}

type Response2 struct{
    Instrucciones string `json:"instrucciones"`
}


type LoginRequest struct{
    Usuario string `json:"name"`
	Contrasenia string `json:"password"`
	Particion string `json:"part"`
}

type LoginResponse struct{
    Usuario string `json:"name"`
	Error string `json:"error"`
}


type LogoutRequest struct{
	Mensaje string `json:"mensaje"`
}



type analizador struct {  
    codigo   string
	cmd comandos.Comandos
}

func New(codigo string) analizador{
	var cmd comandos.Comandos 
	cmd.Id_disco = 1
	cmd.Graph = "graph G{}"
	a := analizador {codigo, cmd}
	return a
}

func (a analizador) Imprimir() {
	fmt.Println(a.codigo)
}

func Split_txt(texto string) []string{
	name := ""
	path := ""
	grp := ""
	pwd := ""
	usuario := ""
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

	first = strings.Index(texto, "-grp=\"")
	if (first > -1){
		last = strings.Index(texto[first+6:], "\"") + first + 6
		grp = texto[first:last+1]
		texto = strings.Replace(texto, grp, "", 1)
	}

	first = strings.Index(texto, "-pwd=\"")
	if (first > -1){
		last = strings.Index(texto[first+6:], "\"") + first + 6
		pwd = texto[first:last+1]
		texto = strings.Replace(texto, pwd, "", 1)
	}

	first = strings.Index(texto, "-usuario=\"")
	if (first > -1){
		last = strings.Index(texto[first+10:], "\"") + first + 10
		usuario = texto[first:last+1]
		texto = strings.Replace(texto, usuario, "", 1)
	}

	splited = strings.Split(texto, " ")
	if (path != ""){splited = append(splited, path)}
	if (name != ""){splited = append(splited, name)}
	if (grp != ""){splited = append(splited, grp)}
	if (pwd != ""){splited = append(splited, pwd)}
	if (usuario != ""){splited = append(splited, usuario)}

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
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ No se ha ingresado la ruta")
			return
		}

		if (fit != "FF") && (fit != "BF") && (fit != "WF"){
			fmt.Println(fit)
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ Fit incorrecto")
			return
		}

		if (unit != "M") && (unit != "K"){
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ Dimensional incorrecto")
			return
		}

		if (size <= 0){
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ Tamaño del disco incorrecto")
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
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ No se ha ingresado la ruta")
			return
		}

		if name == ""{
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ No se ha ingresado el nombre de la particion")
			return
		}

		if (fit != "FF") && (fit != "BF") && (fit != "WF"){
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ Fit incorrecto")
			return
		}

		if (unit != "M") && (unit != "K"){
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ Dimensional incorrecto")
			return
		}

		if (size <= 0){
			fmt.Println("[MIA]@Proyecto2:~$ Tamaño de la particion incorrecta incorrecto")
			return
		}
		
		a.cmd.Fdisk(size, fit[1], unit[0], path, type_[0],name)
	}

	if comando == "mount"{
		fmt.Println("Comando rmdisk")
		path := ""
		name := ""
		
		for i := 0; i < len(parametros); i++ {
			param := parametros[i]
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
		}

		if path == ""{
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ No se ha ingresado la ruta")
			return
		}

		if name == ""{
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ No se ha ingresado el nombre de la particion")
			return
		}

		a.cmd.Mount(path,name)
	}

	if comando == "show"{
		a.cmd.ShowMount()
	}

	if comando == "users"{
		id := ""
		
		for i := 0; i < len(parametros); i++ {
			param := parametros[i]
			if (strings.Index(param, "-id=") == 0) {
				param = strings.Replace(param, "-id=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				id = param
				fmt.Println("ID: ",id)
			}
		}

		a.cmd.GetUsers(id)
	}


	if comando == "rmusr"{
		fmt.Println("")
		fmt.Println("rmusr")
		usuario := ""
		
		for i := 0; i < len(parametros); i++ {
			param := parametros[i]
			if (strings.Index(param, "-usuario=") == 0) {
				param = strings.Replace(param, "-usuario=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				usuario = param
				fmt.Println("Usuario: ",usuario)
			}
		}

		a.cmd.Rmusr(usuario)
	}

	if comando == "rmgrp"{
		name := ""
		
		for i := 0; i < len(parametros); i++ {
			param := parametros[i]
			if (strings.Index(param, "-name=") == 0) {
				param = strings.Replace(param, "-name=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				name = param
				fmt.Println("Name: ",name)
			}
		}

		a.cmd.Rmgrp(name)
	}


	if comando == "mkfs"{
		fmt.Println("Comando rmdisk")
		id := ""
		
		for i := 0; i < len(parametros); i++ {
			param := parametros[i]
			if (strings.Index(param, "-id=") == 0) {
				param = strings.Replace(param, "-id=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				id = param
				fmt.Println("ID: ",id)
			}
		}

		if id == ""{
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ No se ha ingresado el ID de la particion")
			return
		}
		a.cmd.Mkfs(id)
	}

	if comando == "login"{
		id := ""
		usuario := ""
		password := ""
		
		for i := 0; i < len(parametros); i++ {
			param := parametros[i]
			if (strings.Index(param, "-id=") == 0) {
				param = strings.Replace(param, "-id=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				id = param
				fmt.Println("ID: ",id)
			}
			if (strings.Index(param, "-usuario=") == 0) {
				param = strings.Replace(param, "-usuario=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				usuario = param
				fmt.Println("Usuario: ",id)
			}
			if (strings.Index(param, "-password=") == 0) {
				param = strings.Replace(param, "-password=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				password = param
				fmt.Println("Password: ",id)
			}
		}

		if usuario == ""{
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ No se ha ingresado el nombre de usuario")
			return
		}

		if password == ""{
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ No se ha ingresado la contraseña")
			return
		}

		if id == ""{
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ No se ha ingresado el ID de la particion")
			return
		}
		
		a.LoginFunction(usuario, id, password)


		// a.cmd.GetUsers(id)
	}

	// if comando == "pause"{
    // 	fmt.Scanln() 
	// }

	if comando == "rep"{
		name := ""
		id := ""
		
		for i := 0; i < len(parametros); i++ {
			param := parametros[i]
			if (strings.Index(param, "-name=") == 0) {
				param = strings.Replace(param, "-name=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				name = param
				fmt.Println("Name: ",name)
			}

			if (strings.Index(param, "-id=") == 0) {
				param = strings.Replace(param, "-id=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				id = param
				fmt.Println("Id: ",id)
			}

		}

		if(name == "mbr"){
            a.cmd.RepDisco(id);
        }
	
	}

	if comando == "exec"{
		fmt.Println("Comando exec")
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

		if path == ""{
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ No se ha ingresado la ruta")
			return
		}

		a.AnalizarScript(path)
	}


	if comando == "mkusr"{
		grupo := ""
		usuario := ""
		password := ""
		
		for i := 0; i < len(parametros); i++ {
			param := parametros[i]
			if (strings.Index(param, "-grp=") == 0) {
				param = strings.Replace(param, "-grp=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				grupo = param
				fmt.Println("Grupo: ",grupo)
			}
			if (strings.Index(param, "-usuario=") == 0) {
				param = strings.Replace(param, "-usuario=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				usuario = param
				fmt.Println("Usuario: ",usuario)
			}
			if (strings.Index(param, "-pwd=") == 0) {
				param = strings.Replace(param, "-pwd=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				password = param
				fmt.Println("Password: ",password)
			}
		}

		if usuario == ""{
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ No se ha ingresado el nombre de usuario")
			return
		}

		if password == ""{
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ No se ha ingresado la contraseña")
			return
		}

		if grupo == ""{
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ No se ha ingresado el grupo")
			return
		}
		
		a.cmd.Mkusr(usuario, password, grupo)


		// a.cmd.GetUsers(id)
	}


	if comando == "mkgrp"{
		name := ""
		
		for i := 0; i < len(parametros); i++ {
			param := parametros[i]
			if (strings.Index(param, "-name=") == 0) {
				param = strings.Replace(param, "-name=", "", 1)
				param = strings.Replace(param, "\"", "", 2)
				name = param
				fmt.Println("Name: ",name)
			}
		}

		if name == ""{
			a.cmd.AddConsola("[MIA]@Proyecto2:~$ No se ha ingresado el nombre del grupo")
			return
		}
		
		a.cmd.Mkgrp(name)


		// a.cmd.GetUsers(id)
	}

}

func (a *analizador) LoginFunction(name string, id string, password string) string{
	a.cmd.Login(name, password, id)
	return ""
}


func (a *analizador) AnalizarScript(path string){
	
	file, _ := os.Open(path)
    // if err != nil {
    //     a.cmd.AddConsola("[MIA]@Proyecto2:~$ Error al leer el archivo:"+ err)
    // }
    // defer file.Close()

	scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        fmt.Println(scanner.Text())
		a.Analizar(scanner.Text())
    }

}


func (a *analizador) PostConsola(c *gin.Context){

	var Consola []string

	a.cmd.Consola = Consola
	a.cmd.Graph = "graph G{}"

	var texto Response2
	if err := c.BindJSON(&texto); err != nil {
        return
    }

	split := strings.Split(texto.Instrucciones, "\n")
	for _, com := range split {
		a.Analizar(com)
	}

	fmt.Println(a.cmd.Graph)

	res := Response1 {a.cmd.GetConsola(), a.cmd.Graph}
	c.IndentedJSON(http.StatusCreated, res)
}

func (a *analizador) PostLogin(c *gin.Context){

	var texto LoginRequest
	if err := c.BindJSON(&texto); err != nil {
        return
    }
	
	res_error := a.cmd.Login(texto.Usuario, texto.Contrasenia, texto.Particion)

	fmt.Println(res_error )

	// fmt.Println(texto.Usuario)
	// fmt.Println(texto.Contrasenia)
	// fmt.Println(texto.Particion)

	res := LoginResponse {texto.Usuario, res_error}
	c.IndentedJSON(http.StatusCreated, res)
}



func (a *analizador) PostLogout(c *gin.Context){

	var texto LogoutRequest
	if err := c.BindJSON(&texto); err != nil {
        return
    }
	
	res_error := a.cmd.Logout()

	fmt.Println(res_error )

	// fmt.Println(texto.Usuario)
	// fmt.Println(texto.Contrasenia)
	// fmt.Println(texto.Particion)

	res := LogoutRequest {res_error}
	c.IndentedJSON(http.StatusCreated, res)
}