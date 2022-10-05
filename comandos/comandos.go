package comandos

import (  
    "fmt"
	"os"
	"bytes"
	"encoding/gob"
	"io"
	"math/rand"
	"strconv"
	"time"
	"strings"
)

type MBR struct {  
    Mbr_tamano 			[16]byte
	Mbr_fecha_creacion 	[16]byte
	Mbr_dsk_signature  	[16]byte
	Msk_fit				[4]byte
	Mbr_partition_1 	Partition
	Mbr_partition_2		Partition
	Mbr_partition_3		Partition
	Mbr_partition_4		Partition
}

type Partition struct {  
    Part_status [4]byte
	Part_type 	[4]byte
	Part_fit  	[4]byte
	Part_start 	[16]byte
	Part_size	[16]byte
	Part_name	[16]byte
}

type EBR struct {  
    Part_status [4]byte
	Part_start 	[16]byte
	Part_size  	[16]byte
	Part_next	[16]byte
	Part_name 	[16]byte

}

type Comandos struct {  
    Numero int
}


func (cmd *Comandos) Imprimir() {
	fmt.Println(cmd.Numero)
	cmd.Numero = cmd.Numero + 1
}

func (cmd *Comandos) Mkdisk(size int, fit byte, unit byte, path string){
	limite := 0
	bloque := make([]byte, 1024)

	if (unit == 'M'){
		size = size * 1024
	}

	for j := 0; j < 1024; j++ {
		bloque[j] = 0
	}

	disco, err := os.Create(path)
	if err != nil {
		msg_error(err)
	}

	for limite < size {
		_, err := disco.Write(bloque)
		if err != nil {
			msg_error(err)
		}
		limite++
	}

	mbr := MBR{}
	size = size * 1024
	currentTime := time.Now()
	date := currentTime.Format("02-01-2006")
	signature := rand.Intn(999999)
	
	empty_partition := getEmptyPartition()
	
	copy(mbr.Mbr_tamano[:], strconv.Itoa(size))
	copy(mbr.Mbr_fecha_creacion[:], date)
	copy(mbr.Mbr_dsk_signature[:], strconv.Itoa(signature))
	copy(mbr.Msk_fit[:], string(fit))

	mbr.Mbr_partition_1 = empty_partition
	mbr.Mbr_partition_2 = empty_partition
	mbr.Mbr_partition_3 = empty_partition
	mbr.Mbr_partition_4 = empty_partition

	mbr_bytes := struct_to_bytes(mbr)
	pos, err := disco.Seek(0, os.SEEK_SET)
	if err != nil {
		msg_error(err)
	}

	_, err = disco.WriteAt(mbr_bytes, pos)
		if err != nil {
		msg_error(err)
	}

	disco.Close()

	mostrar(mbr)
}


func getEmptyPartition() Partition{
	part := Partition{}
	copy(part.Part_status[:], "0")
	copy(part.Part_type[:], "P")
	copy(part.Part_fit[:], "F")
	copy(part.Part_start[:], strconv.Itoa(-1))
	copy(part.Part_size[:], strconv.Itoa(0))
	copy(part.Part_name[:], "")
	return part
}


func (cmd *Comandos) Fdisk(size int, fit byte, unit byte, path string, type_ byte, name string){

	if (unit == 'K'){
		size = size * 1024
	}

	if (unit == 'M'){
		size = size * 1024 * 1024
	}

	disco, err := os.OpenFile(path, os.O_RDWR, 0660)
	if err != nil {
		msg_error(err)
		return
	}

	empty_mbr := MBR{}
	size_mbr := len(struct_to_bytes(empty_mbr))

	buff_mbr := make([]byte, size_mbr)
	_, err = disco.ReadAt(buff_mbr, 0)
	if err != nil && err != io.EOF {
		msg_error(err)
	}

	mbr := bytes_to_mbr(buff_mbr)

	if(existeNombre(mbr, name)){
		fmt.Println("[MIA]@Proyecto2:~$ Ya existe una particion con ese nombre")
		return
	}

	inicio_libre := size_mbr + 1
	if(string(mbr.Mbr_partition_1.Part_status[:]) != "0"){
		ini, _ := strconv.Atoi(string(mbr.Mbr_partition_1.Part_start[:]))
		size_, _ := strconv.Atoi(string(mbr.Mbr_partition_1.Part_size[:]))
        inicio_libre = ini + size_ + 1
    }

	if(string(mbr.Mbr_partition_2.Part_status[:]) != "0"){
        ini, _ := strconv.Atoi(string(mbr.Mbr_partition_2.Part_start[:]))
		size_, _ := strconv.Atoi(string(mbr.Mbr_partition_2.Part_size[:]))
        inicio_libre = ini + size_ + 1
    }

	if(string(mbr.Mbr_partition_3.Part_status[:]) != "0"){
        ini, _ := strconv.Atoi(string(mbr.Mbr_partition_3.Part_start[:]))
		size_, _ := strconv.Atoi(string(mbr.Mbr_partition_3.Part_size[:]))
        inicio_libre = ini + size_ + 1
    }

	if(string(mbr.Mbr_partition_4.Part_status[:]) != "0"){
        ini, _ := strconv.Atoi(string(mbr.Mbr_partition_4.Part_start[:]))
		size_, _ := strconv.Atoi(string(mbr.Mbr_partition_4.Part_size[:]))
        inicio_libre = ini + size_ + 1
    }

	// fmt(Println(inicio_libre))
	nueva_part := Partition{}

	copy(nueva_part.Part_status[:], "1")
	copy(nueva_part.Part_type[:], string(type_))
	copy(nueva_part.Part_fit[:], string(fit))
	copy(nueva_part.Part_start[:], strconv.Itoa(inicio_libre))
	copy(nueva_part.Part_size[:], strconv.Itoa(size))
	copy(nueva_part.Part_name[:], name)
    
	fmt.Println(string(mbr.Mbr_partition_1.Part_status[:]))

	if(strings.Contains(string(mbr.Mbr_partition_1.Part_status[:]),"0")){
        mbr.Mbr_partition_1 = nueva_part
    } else if(strings.Contains(string(mbr.Mbr_partition_2.Part_status[:]),"0")){
        mbr.Mbr_partition_2 = nueva_part
    } else if(strings.Contains(string(mbr.Mbr_partition_3.Part_status[:]),"0")){
        mbr.Mbr_partition_3 = nueva_part
    } else if(strings.Contains(string(mbr.Mbr_partition_4.Part_status[:]),"0")){
        mbr.Mbr_partition_4 = nueva_part
    } else{
		fmt.Println("[MIA]@Proyecto2:~$ No se puede crear otra particion")
		return
	}

	mbr_bytes := struct_to_bytes(mbr)
	pos, err := disco.Seek(0, os.SEEK_SET)
	if err != nil {
		msg_error(err)
	}

	_, err = disco.WriteAt(mbr_bytes, pos)
		if err != nil {
		msg_error(err)
	}

	disco.Close()

	mostrar(mbr)

}	


func existeNombre(master MBR, nombre_s string) bool{

	part_name := master.Mbr_partition_1.Part_name;
    
	if(strings.Contains(string(part_name[:]), nombre_s)){
        return true;
    }
    part_name = master.Mbr_partition_2.Part_name;
    if(strings.Contains(string(part_name[:]), nombre_s)){
        return true;
    }
    part_name = master.Mbr_partition_3.Part_name;
    if(strings.Contains(string(part_name[:]), nombre_s)){
        return true;
    }
    part_name = master.Mbr_partition_4.Part_name;
    if(strings.Contains(string(part_name[:]), nombre_s)){
        return true;
    }
	return false
}


func mostrar(master MBR){
	fmt.Print("[MIA]@Proyecto2:~$ ")

    fmt.Print("MBR ->")
        
	fmt.Print("SIZE: ")
	fmt.Print(string(master.Mbr_tamano[:]))

    fmt.Print(", TIME: ")
	fmt.Print(string(master.Mbr_fecha_creacion[:]))

	fmt.Print(", SIGNATURE: ")
	fmt.Print(string(master.Mbr_dsk_signature[:]))

	fmt.Print(", FIT: ")
	fmt.Println(string(master.Msk_fit[:]))

	fmt.Println("PARTITIONS: ")
    fmt.Println("-- Name ",string(master.Mbr_partition_1.Part_name[:]), ", Size: ", string(master.Mbr_partition_1.Part_size[:]), ", Start: ", string(master.Mbr_partition_1.Part_start[:]));
	fmt.Println("-- Name ",string(master.Mbr_partition_2.Part_name[:]), ", Size: ", string(master.Mbr_partition_2.Part_size[:]), ", Start: ", string(master.Mbr_partition_2.Part_start[:]));
	fmt.Println("-- Name ",string(master.Mbr_partition_3.Part_name[:]), ", Size: ", string(master.Mbr_partition_3.Part_size[:]), ", Start: ", string(master.Mbr_partition_3.Part_start[:]));
	fmt.Println("-- Name ",string(master.Mbr_partition_4.Part_name[:]), ", Size: ", string(master.Mbr_partition_4.Part_size[:]), ", Start: ", string(master.Mbr_partition_4.Part_start[:]));

}


func msg_error(err error) {
	fmt.Println("Error: ", err)
}

func struct_to_bytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil && err != io.EOF {
		msg_error(err)
	}
	return buf.Bytes()
}

func bytes_to_mbr(s []byte) MBR {
	p := MBR{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil && err != io.EOF {
		msg_error(err)
	}
	return p
}

func bytes_to_ebr(s []byte) EBR {
	p := EBR{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil && err != io.EOF {
		msg_error(err)
	}
	return p
}