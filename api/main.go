package main  

import "MIA_PROYECTO2_202010055/analizador"

import (
    // "bufio"
	// "fmt"
	// "os"
	// "stginrings"
    "github.com/gin-gonic/gin"
)



func main() {  
    a := analizador.New("")

    router := gin.Default()
    router.Use(CORSMiddleware())
    router.POST("/consola", a.PostConsola)
    router.Run("127.0.0.1:5000")
}


func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {

        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}




    // mkdisk -size=-2 -unit=M -path=/home/desquivel/Desktop/Disco3.dsk

    // mkdisk -size=2 -unit=M -path=/home/desquivel/Desktop/Disco1.dsk
    // fdisk -type=P -path=/home/desquivel/Desktop/Disco1.dsk -unit=K -name=Particion1 -size=500
    // mount -path=/home/desquivel/Desktop/Disco1.dsk -name=Particion1
    // mkfs -id=551A

    // mkdisk -size=5 -unit=M -path=/home/desquivel/Desktop/Disco2.dsk
    // fdisk -type=P -path=/home/desquivel/Desktop/Disco2.dsk -unit=K -name=Particion1 -size=300
    // fdisk -type=E -path=/home/desquivel/Desktop/Disco2.dsk -unit=K -name=Particion2 -size=1000
    // fdisk -type=L -path=/home/desquivel/Desktop/Disco2.dsk -unit=K -name=Particion3 -size=200
    // fdisk -type=L -path=/home/desquivel/Desktop/Disco2.dsk -unit=K -name=Particion4 -size=400
    // mount -path=/home/desquivel/Desktop/Disco2.dsk -name=Particion1
    // mount -path=/home/desquivel/Desktop/Disco2.dsk -name=Particion3 


    // exec -path=/home/desquivel/Desktop/script.txt


    // reader := bufio.NewReader(os.Stdin)
    // for true{
    //     fmt.Println("||========================================================||")
    //     fmt.Println("||                    [MIA] PROYECTO 2                    ||")
    //     fmt.Println("||========================================================||")
    //     fmt.Print("[MIA]@Proyecto2:~$ ")
        
    //     comando, _ := reader.ReadString('\n')

    //     if(strings.Contains(comando, "exit\n")){
    //         break
    //     }
    //     a.Analizar(comando)
        
    // }
    