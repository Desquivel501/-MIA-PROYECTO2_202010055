package main  

import "MIA_PROYECTO2_202010055/analizador"


func main() {  
    a := analizador.New("")
    a.Analizar("fdisk -type=E -path=/home/Disco2.dk -name=Part3 -unit=K -size=200")
    // a.Analizar("fdisk -type=E -path=/home/Disco2.dk -name=Part3 -unit=K -size=200")
    // a.Analizar("fdisk -type=E -path=/home/Disco2.dk -name=Part3 -unit=K -size=200")
    
}