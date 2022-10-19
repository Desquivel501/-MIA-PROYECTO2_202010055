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
	Msk_fit				[1]byte
	Mbr_partition_1 	Partition
	Mbr_partition_2		Partition
	Mbr_partition_3		Partition
	Mbr_partition_4		Partition
}

type Partition struct {  
    Part_status [1]byte
	Part_type 	[1]byte
	Part_fit  	[1]byte
	Part_start 	[16]byte
	Part_size	[16]byte
	Part_name	[16]byte
}

type EBR struct {  
    Part_status [1]byte
	Part_fit  	[1]byte
	Part_start 	[16]byte
	Part_size  	[16]byte
	Part_next	[16]byte
	Part_name 	[16]byte
}

type Mounted struct {
	Id string
	Part Partition
	Master MBR
	Path string
	Id_disco int
}

type SuperBlock struct {
	S_filesystem_type    [1]byte
	S_inodes_count       [16]byte
	S_blocks_count       [16]byte
	S_free_blocks_count  [16]byte
	S_free_inodes_count  [16]byte
	S_mtime              [16]byte
	S_mnt_count          [16]byte
	S_magic              [16]byte
	S_inode_size         [16]byte
	S_block_size         [16]byte
	S_first_ino          [16]byte
	S_first_blo          [16]byte
	S_bm_inode_start     [16]byte
	S_bm_block_start     [16]byte
	S_inode_start        [16]byte
	S_block_start        [16]byte
}

type Inodo struct {
	I_uid     [3]byte
	I_gid     [3]byte
	I_size    [16]byte
	I_atime   [16]byte
	I_ctime   [16]byte
	I_mtime   [16]byte
	I_block   [16][16]byte
	I_type    [1]byte
	I_perm    [3]byte
}

type Carpeta struct{
	B_content  [4]Content
}

type Content struct{
	B_name   [16]byte
	B_inodo  [16]byte
}

type Archivo struct{
	B_content   [64]byte
}


type Comandos struct {  
    Mounted_list []Mounted
	Id_disco int
	Consola []string
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
		cmd.msg_error(err)
	}

	for limite < size {
		_, err := disco.Write(bloque)
		if err != nil {
			cmd.msg_error(err)
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

	mbr_bytes := cmd.struct_to_bytes(mbr)
	pos, err := disco.Seek(0, os.SEEK_SET)
	if err != nil {
		cmd.msg_error(err)
	}

	_, err = disco.WriteAt(mbr_bytes, pos)
		if err != nil {
		cmd.msg_error(err)
	}

	// cmd.mostrar(disco, mbr)
	cmd.AddConsola("[MIA]@Proyecto2:~$ Creado disco: " + path)
	disco.Close()
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
		cmd.msg_error(err)
		return
	}

	empty_mbr := MBR{}
	size_mbr := len(cmd.struct_to_bytes(empty_mbr))

	buff_mbr := make([]byte, size_mbr)
	_, err = disco.ReadAt(buff_mbr, 0)
	if err != nil && err != io.EOF {
		disco.Close()
		cmd.msg_error(err)
		return
	}

	mbr := cmd.bytes_to_mbr(buff_mbr)

	if(existeNombre(mbr, name)){
		cmd.AddConsola("[MIA]@Proyecto2:~$ Ya existe una particion con ese nombre")
		// cmd.AddConsola("[MIA]@Proyecto2:~$ Ya existe una particion con ese nombre")
		disco.Close()
		return
	}


	if(type_ == 'L'){
		if !existeExtendida(mbr){
			cmd.AddConsola("[MIA]@Proyecto2:~$ No existe particion extendida")
			disco.Close()
			return
		}
		cmd.crearLogica(disco, mbr, size, fit, unit, name)
		return
	}

	inicio_libre := size_mbr + 1
	if(string(mbr.Mbr_partition_1.Part_status[:]) != "0"){
		ini := bytes_to_int(mbr.Mbr_partition_1.Part_start[:])
		size_ := bytes_to_int(mbr.Mbr_partition_1.Part_size[:])
        inicio_libre = ini + size_ + 1
    }

	if(string(mbr.Mbr_partition_2.Part_status[:]) != "0"){
		ini := bytes_to_int(mbr.Mbr_partition_2.Part_start[:])
		size_ := bytes_to_int(mbr.Mbr_partition_2.Part_size[:])
        inicio_libre = ini + size_ + 1
    }

	if(string(mbr.Mbr_partition_3.Part_status[:]) != "0"){
        ini := bytes_to_int(mbr.Mbr_partition_3.Part_start[:])
		size_ := bytes_to_int(mbr.Mbr_partition_3.Part_size[:])
        inicio_libre = ini + size_ + 1
    }

	if(string(mbr.Mbr_partition_4.Part_status[:]) != "0"){
        ini := bytes_to_int(mbr.Mbr_partition_4.Part_start[:])
		size_ := bytes_to_int(mbr.Mbr_partition_4.Part_size[:])
        inicio_libre = ini + size_ + 1
    }

	nueva_part := Partition{}

	copy(nueva_part.Part_status[:], "1")
	copy(nueva_part.Part_type[:], string(type_))
	copy(nueva_part.Part_fit[:], string(fit))
	copy(nueva_part.Part_start[:], strconv.Itoa(inicio_libre))
	copy(nueva_part.Part_size[:], strconv.Itoa(size))
	copy(nueva_part.Part_name[:], name)
    
	if(strings.Contains(string(mbr.Mbr_partition_1.Part_status[:]),"0")){
        mbr.Mbr_partition_1 = nueva_part
    } else if(strings.Contains(string(mbr.Mbr_partition_2.Part_status[:]),"0")){
        mbr.Mbr_partition_2 = nueva_part
    } else if(strings.Contains(string(mbr.Mbr_partition_3.Part_status[:]),"0")){
        mbr.Mbr_partition_3 = nueva_part
    } else if(strings.Contains(string(mbr.Mbr_partition_4.Part_status[:]),"0")){
        mbr.Mbr_partition_4 = nueva_part
    } else{
		cmd.AddConsola("[MIA]@Proyecto2:~$ No se puede crear otra particion")
		disco.Close()
		return
	}

	mbr_bytes := cmd.struct_to_bytes(mbr)
	disco.WriteAt(mbr_bytes, 0)

	if(type_ == 'E'){
		ebr := getEmptyEBR()
		ebr_bytes := cmd.struct_to_bytes(ebr)
		disco.WriteAt(ebr_bytes, int64(inicio_libre))
		cmd.AddConsola("[MIA]@Proyecto2:~$ Creada particion extendida " + name)	
	}else{
		cmd.AddConsola("[MIA]@Proyecto2:~$ Creada particion primaria " + name)	
	}
	
	// cmd.mostrar(disco, mbr)

	disco.Close()
}	

func (cmd *Comandos) crearLogica(disco *os.File, master MBR, size int, fit byte, unit byte, name string){
	extendida :=  getExtendida(master)

	inicio_part := bytes_to_int(extendida.Part_start[:])
	size_part := bytes_to_int(extendida.Part_size[:])
	ebr := cmd.leerEBR(disco, int64(inicio_part))

	inicio_libre := inicio_part
	fin_libre := inicio_part + size_part - 1

	if(size < len(cmd.struct_to_bytes(ebr))){
		cmd.AddConsola("[MIA]@Proyecto2:~$ No hay espacio para crear la particion")
		return
	}

	if(strings.Contains(string(ebr.Part_status[:]), "0")){
		copy(ebr.Part_status[:], "1")
		copy(ebr.Part_fit[:], string(fit))
		copy(ebr.Part_start[:], strconv.Itoa(inicio_libre + len(cmd.struct_to_bytes(ebr)) ))
		copy(ebr.Part_size[:], strconv.Itoa(size))
		copy(ebr.Part_next[:], strconv.Itoa(-1))
		copy(ebr.Part_name[:], name)

		ebr_bytes := cmd.struct_to_bytes(ebr)
		disco.WriteAt(ebr_bytes, int64(inicio_libre))
		cmd.AddConsola("[MIA]@Proyecto2:~$ Creada particion logica " + name)	

	}else{
		pos := 0
		for{
			ebr_ini := bytes_to_int(ebr.Part_start[:])
			ebr_size := bytes_to_int(ebr.Part_size[:])
			inicio_libre = ebr_ini + ebr_size - 1
			
			if (fin_libre - inicio_libre + 1 < size){
				cmd.AddConsola("[MIA]@Proyecto2:~$ No hay espacio para crear la particion")
				return
			}

			if(strings.Contains(string(ebr.Part_next[:]), "-1")){
				break
			}else{
				pos = bytes_to_int(ebr.Part_next[:])
				ebr = cmd.leerEBR(disco, int64(pos))
			}
		}

		nuevo_ebr := getEmptyEBR()
		copy(nuevo_ebr.Part_status[:], "1")
		copy(nuevo_ebr.Part_fit[:], string(fit))
		copy(nuevo_ebr.Part_start[:], strconv.Itoa(inicio_libre + len(cmd.struct_to_bytes(ebr))))
		copy(nuevo_ebr.Part_size[:], strconv.Itoa(size))
		copy(nuevo_ebr.Part_next[:], strconv.Itoa(-1))
		copy(nuevo_ebr.Part_name[:], name)

		copy(ebr.Part_next[:], strconv.Itoa(inicio_libre))

		pos = bytes_to_int(ebr.Part_start[:]) - len(cmd.struct_to_bytes(ebr))
		ebr_bytes := cmd.struct_to_bytes(ebr)
		disco.WriteAt(ebr_bytes, int64(pos))

		ebr_bytes = cmd.struct_to_bytes(nuevo_ebr)
		disco.WriteAt(ebr_bytes, int64(inicio_libre))

		cmd.AddConsola("[MIA]@Proyecto2:~$ Creada particion logica " + name)
	}

	// cmd.AddConsola("[MIA]@Proyecto2:~$ Creado disco", path)
	// cmd.mostrar(disco, master)
	disco.Close()
}


func (cmd *Comandos) Mount(path string, name string){
	disco, err := os.OpenFile(path, os.O_RDWR, 0660)
	if err != nil {
		cmd.msg_error(err)
		return
	}

	empty_mbr := MBR{}
	size_mbr := len(cmd.struct_to_bytes(empty_mbr))

	buff_mbr := make([]byte, size_mbr)
	_, err = disco.ReadAt(buff_mbr, 0)
	if err != nil && err != io.EOF {
		disco.Close()
		cmd.msg_error(err)
		return
	}
	

	master := cmd.bytes_to_mbr(buff_mbr)

	mount_ := cmd.getId(path)
	mount_.Path = path
	mount_.Master = master

	part_name := master.Mbr_partition_1.Part_name;
	if(strings.Contains(string(part_name[:]), name)){
		mount_.Part = master.Mbr_partition_1
		cmd.Mounted_list = append(cmd.Mounted_list, mount_)
		disco.Close()
		cmd.AddConsola("[MIA]@Proyecto2:~$ Particion " + string(part_name[:]) + " montada con id: " + mount_.Id)
		return
    }

    part_name = master.Mbr_partition_2.Part_name;
    if(strings.Contains(string(part_name[:]), name)){
        mount_.Part = master.Mbr_partition_2
		cmd.Mounted_list = append(cmd.Mounted_list, mount_)
		disco.Close()
		cmd.AddConsola("[MIA]@Proyecto2:~$ Particion " + string(part_name[:]) + " montada con id: " + mount_.Id)
		return
    }
	
    part_name = master.Mbr_partition_3.Part_name;
    if(strings.Contains(string(part_name[:]), name)){
        mount_.Part = master.Mbr_partition_3
		cmd.Mounted_list = append(cmd.Mounted_list, mount_)
		disco.Close()
		cmd.AddConsola("[MIA]@Proyecto2:~$ Particion " + string(part_name[:]) + " montada con id: " + mount_.Id)
		return
    }

    part_name = master.Mbr_partition_4.Part_name;
    if(strings.Contains(string(part_name[:]), name)){
		mount_.Part = master.Mbr_partition_4
		cmd.Mounted_list = append(cmd.Mounted_list, mount_)
		disco.Close()
		cmd.AddConsola("[MIA]@Proyecto2:~$ Particion " + string(part_name[:]) + " montada con id: " + mount_.Id)
		return
    }

	if (existeExtendida(master)){
		extendida := getExtendida(master)
		inicio_part := bytes_to_int(extendida.Part_start[:])
		ebr := cmd.leerEBR(disco, int64(inicio_part))
		
		part_name = ebr.Part_name;
		if(strings.Contains(string(part_name[:]), name)){
			mount_.Part = master.Mbr_partition_4
			cmd.Mounted_list = append(cmd.Mounted_list, mount_)
			disco.Close()
			cmd.AddConsola("[MIA]@Proyecto2:~$ Particion " + string(part_name[:]) + " montada con id: " + mount_.Id)
			return
		}
	
		for{
			cmd.AddConsola("1")
			if(strings.Contains(string(ebr.Part_next[:]), "-1")){
				break
			}
			inicio_ebr := bytes_to_int(ebr.Part_next[:])
			ebr = cmd.leerEBR(disco, int64(inicio_ebr))
			
			part_name = ebr.Part_name;
			if(strings.Contains(string(part_name[:]), name)){
				mount_.Part = master.Mbr_partition_4
				cmd.Mounted_list = append(cmd.Mounted_list, mount_)
				disco.Close()
				cmd.AddConsola("[MIA]@Proyecto2:~$ Particion " + string(part_name[:]) + " montada con id: " + mount_.Id)
				return
			}
		}
	}
	cmd.AddConsola("[MIA]@Proyecto2:~$ Ha ocurrido un error al montar la particion")
	disco.Close()
}

func (cmd *Comandos) ShowMount(){
	cmd.AddConsola("Particiones montadas:")
	for _, m := range cmd.Mounted_list{
		path := slicePath(m.Path)
		cmd.AddConsola("- ID: " + m.Id +", Nombre Disco: " + path[len(path) - 1])
	}
}

func (cmd *Comandos) GetMount(id string) (Mounted, int){
	empty := Mounted {}
	for _, m := range cmd.Mounted_list{
		if(m.Id == id){
			return m, 1
		}
	}
	return empty, 0
}

func (cmd *Comandos) Mkfs(id string){
	part, err_ := cmd.GetMount(id)
	if(err_ == 0){
		cmd.AddConsola("[MIA]@Proyecto2:~$ La particion no ha sido montada")
		return
	}

	currentTime := time.Now()
	date := currentTime.Format("02-01-2006")

	empty_i := Inodo{}
	SI := len(cmd.struct_to_bytes(empty_i))

	empty_b := Carpeta{}
	SB := len(cmd.struct_to_bytes(empty_b))

	empty_sb:= SuperBlock{}
	SS := len(cmd.struct_to_bytes(empty_sb))

	size_part := bytes_to_int(part.Part.Part_size[:])
	n := cmd.getNumeroEstructuras(size_part)

	part_start := bytes_to_int(part.Part.Part_start[:])

	super := SuperBlock{}
	copy(super.S_filesystem_type[:], "2")
	copy(super.S_inodes_count[:], strconv.Itoa(n))
	copy(super.S_blocks_count[:], strconv.Itoa(n*3))
	copy(super.S_free_blocks_count[:], strconv.Itoa(n - 2))
	copy(super.S_free_inodes_count[:], strconv.Itoa(n*3 - 2))
	copy(super.S_mtime[:], date)
	copy(super.S_mnt_count[:], strconv.Itoa(1))
	copy(super.S_inode_size[:], strconv.Itoa(SI))
	copy(super.S_block_size[:], strconv.Itoa(SB))
	copy(super.S_first_ino[:], strconv.Itoa(2))
	copy(super.S_first_blo[:], strconv.Itoa(2))
	copy(super.S_bm_inode_start[:], strconv.Itoa(part_start + SS))
	copy(super.S_bm_block_start[:], strconv.Itoa(part_start + SS+n))
	copy(super.S_inode_start[:], strconv.Itoa(part_start + SS + n + 3*n))
	copy(super.S_block_start[:], strconv.Itoa(part_start + SS + n + 3*n + SI*n ))

	disco, err := os.OpenFile(part.Path, os.O_RDWR, 0660)
	if err != nil {
		cmd.msg_error(err)
		return
	}
	
	super_bytes := cmd.struct_to_bytes(super)
	disco.WriteAt(super_bytes, int64(part_start))

	// bm_inodos := make([]byte, n)
	// disco.WriteAt(bm_inodos, int64(part_start + SS))

	// bm_bloques := make([]byte, 3*n)
	// disco.WriteAt(bm_bloques, int64(part_start + SS + n))


	//-----------------------------------------------------	USERS.TXT

	users := Inodo{}
	copy(users.I_uid[:], "1")
	copy(users.I_gid[:], "1")
	copy(users.I_size[:], strconv.Itoa(0))
	copy(users.I_atime[:], date)
	copy(users.I_ctime[:], date)
	copy(users.I_type[:], "1")
	copy(users.I_perm[:], "664")

	for i := 0; i < 16; i += 1{
		copy(users.I_block[i][:], strconv.Itoa(-1))
	}
	copy(users.I_block[0][:], strconv.Itoa(1))

	users_txt := "1,G,root\n1,U,root,123\n"
	start := part_start + SS + n + 3*n + SI*n 

	users_arr := []byte(users_txt)
	
	block := Archivo{}
	copy(block.B_content[:], users_arr)
	block_bytes := cmd.struct_to_bytes(block)

	disco.WriteAt(block_bytes, int64(start + SB))
	
	// cmd.AddConsola(start + SB)


	//-----------------------------------------------------	BLOQUE CARPETA

	empty_cont := Content{}
	copy(empty_cont.B_name[:], "")
	copy(empty_cont.B_inodo[:], strconv.Itoa(-1))

	users_cont := Content{}
	copy(users_cont.B_name[:], "users.txt")
	copy(users_cont.B_inodo[:], strconv.Itoa(1))

	root_dir := Carpeta{}

	var cont_arr [4]Content
	cont_arr[0] = users_cont
	cont_arr[1] = empty_cont
	cont_arr[2] = empty_cont
	cont_arr[3] = empty_cont
	root_dir.B_content = cont_arr

	carpeta_bytes := cmd.struct_to_bytes(root_dir)
	disco.WriteAt(carpeta_bytes, int64(start))

	
	//-----------------------------------------------------	INODO ROOT
	root := Inodo{}
	copy(root.I_uid[:], "1")
	copy(root.I_gid[:], "1")
	copy(root.I_size[:], strconv.Itoa(0))
	copy(root.I_atime[:], date)
	copy(root.I_ctime[:], date)
	copy(root.I_type[:], "0")
	copy(root.I_perm[:], "664")

	for i := 0; i < 16; i += 1{
		copy(root.I_block[i][:], strconv.Itoa(-1))
	}
	copy(root.I_block[0][:], strconv.Itoa(0))


	//-----------------------------------------------------	ESCRIBIR EN DISCO
	
	bm_inodos := make([]byte, n)
	disco.WriteAt(bm_inodos, int64(part_start + SS))

	bm_bloques := make([]byte, 3*n)
	disco.WriteAt(bm_bloques, int64(part_start + SS + n))

	root_bytes := cmd.struct_to_bytes(root)
	disco.WriteAt(root_bytes, int64(part_start + SS))

	users_bytes := cmd.struct_to_bytes(users)
	disco.WriteAt(users_bytes, int64(part_start + SS + SI))


}

func (cmd *Comandos)createFile(disco *os.File, content string, b_libre int) []int {
	content_arr := []byte(content)
	chunks := chunkSlice(content_arr)
	fmt.Println(chunks)
	pos := 0
	var used_blocks []int

	for i := 0; i < len(chunks); i += 1 {
		block := Archivo{}
		copy(block.B_content[:], chunks[i])
		block_bytes := cmd.struct_to_bytes(block)
		pos = b_libre + len(block_bytes) * i
		disco.WriteAt(block_bytes, int64(pos))
		used_blocks = append(used_blocks, i+1)
	}
	return used_blocks
}

func chunkSlice(slice []byte) [][]byte {
	var chunks [][]byte
	for i := 0; i < len(slice); i += 64 {
		end := i + 64

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
 

func (cmd *Comandos) getId(disco string) Mounted{
	id_num := 65
	id_disco := 0

	for _, m := range cmd.Mounted_list{
		if (disco == m.Path){
			id_num = id_num + 1
			id_disco = m.Id_disco
		}
	}

	if (id_num == 65){
		id_disco = cmd.Id_disco
		cmd.Id_disco = cmd.Id_disco + 1
	}

	id_letra := string(rune(id_num))
	id_disco_str := strconv.Itoa(id_disco)

	id_str := "55" + id_disco_str + id_letra
	
	mount_ := Mounted{}
	mount_.Id = id_str
	mount_.Id_disco = id_disco
	return mount_
}

func slicePath(path string) []string {
	split := strings.Split(path, "/")
	return split
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

func existeExtendida(master MBR) bool{

	part_name := master.Mbr_partition_1.Part_type;
	if(strings.Contains(string(part_name[:]), "E")){
        return true;
    }
    part_name = master.Mbr_partition_2.Part_type;
    if(strings.Contains(string(part_name[:]), "E")){
        return true;
    }
    part_name = master.Mbr_partition_3.Part_type;
    if(strings.Contains(string(part_name[:]), "E")){
        return true;
    }
    part_name = master.Mbr_partition_4.Part_type;
    if(strings.Contains(string(part_name[:]), "E")){
        return true;
    }
	return false
}

func getExtendida(master MBR) Partition{
	if(strings.Contains(string(master.Mbr_partition_1.Part_type[:]), "E")) {return master.Mbr_partition_1}
    if(strings.Contains(string(master.Mbr_partition_2.Part_type[:]), "E")) {return master.Mbr_partition_2}
    if(strings.Contains(string(master.Mbr_partition_3.Part_type[:]), "E")) {return master.Mbr_partition_3}
	return master.Mbr_partition_4
}

func (cmd *Comandos)leerEBR(disco *os.File, pos int64) EBR{
	empty_ebr := EBR{}
	size_ebr := len(cmd.struct_to_bytes(empty_ebr))
	buff_ebr := make([]byte, size_ebr)
	disco.ReadAt(buff_ebr, pos)
	ebr := cmd.bytes_to_ebr(buff_ebr)
	return ebr
}

func (cmd *Comandos)mostrar(disco *os.File, master MBR){

	
    cmd.AddConsola("MIA]@Proyecto2:~$ MBR ->")
        
	// fmt.Print("SIZE: ")
	cmd.AddConsola("SIZE: " + string(master.Mbr_tamano[:]))

    // fmt.Print("TIME: ")
	cmd.AddConsola("TIME: " + string(master.Mbr_fecha_creacion[:]))

	// fmt.Print(", SIGNATURE: ")
	cmd.AddConsola("SIGNATURE: " + string(master.Mbr_dsk_signature[:]))

	// fmt.Print(", FIT: ")
	cmd.AddConsola("FIT: " + string(master.Msk_fit[:]))

	cmd.AddConsola("PARTITIONS: ")

    cmd.AddConsola("-- Name "+string(master.Mbr_partition_1.Part_name[:])+ ", Size: "+ string(master.Mbr_partition_1.Part_size[:])+ ", Start: "+ string(master.Mbr_partition_1.Part_start[:]));
	if(strings.Contains(string(master.Mbr_partition_1.Part_type[:]), "E")){
		cmd.mostrarExtendida(disco, master.Mbr_partition_1)
	}

	cmd.AddConsola("-- Name "+string(master.Mbr_partition_2.Part_name[:])+ ", Size: "+ string(master.Mbr_partition_2.Part_size[:])+ ", Start: "+ string(master.Mbr_partition_2.Part_start[:]));
	if(strings.Contains(string(master.Mbr_partition_2.Part_type[:]), "E")){
		cmd.mostrarExtendida(disco, master.Mbr_partition_2)
	}

	cmd.AddConsola("-- Name "+string(master.Mbr_partition_3.Part_name[:])+ ", Size: "+ string(master.Mbr_partition_3.Part_size[:])+ ", Start: "+ string(master.Mbr_partition_3.Part_start[:]));
	if(strings.Contains(string(master.Mbr_partition_3.Part_type[:]), "E")){
		cmd.mostrarExtendida(disco, master.Mbr_partition_3)
	}

	cmd.AddConsola("-- Name "+string(master.Mbr_partition_4.Part_name[:])+ ", Size: "+ string(master.Mbr_partition_4.Part_size[:])+ ", Start: "+ string(master.Mbr_partition_4.Part_start[:]));
	if(strings.Contains(string(master.Mbr_partition_4.Part_type[:]), "E")){
		cmd.mostrarExtendida(disco, master.Mbr_partition_4)
	}

	cmd.AddConsola(" ")
}

func (cmd *Comandos)mostrarExtendida(disco *os.File, extendida Partition){
	inicio_part := bytes_to_int(extendida.Part_start[:])
	ebr := cmd.leerEBR(disco, int64(inicio_part))

	cmd.AddConsola("---- Name "+ string(ebr.Part_name[:]) + ", Size: "+string(ebr.Part_size[:])+", Start: "+ string(ebr.Part_start[:]));

	for{
		if(strings.Contains(string(ebr.Part_next[:]), "-1")){
			break
		}
		inicio_ebr := bytes_to_int(ebr.Part_next[:])
		ebr = cmd.leerEBR(disco, int64(inicio_ebr))
		cmd.AddConsola("---- Name "+string(ebr.Part_name[:])+ ", Size: "+ string(ebr.Part_size[:])+ ", Start: "+ string(ebr.Part_start[:]));
	}
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


func getEmptyEBR() EBR{
	ebr := EBR{}
	copy(ebr.Part_status[:], "0")
	copy(ebr.Part_fit[:], "F")
	copy(ebr.Part_start[:], strconv.Itoa(-1))
	copy(ebr.Part_size[:], strconv.Itoa(0))
	copy(ebr.Part_next[:], strconv.Itoa(-1))
	copy(ebr.Part_name[:], "")
	return ebr
}


func (cmd *Comandos)getNumeroEstructuras(part_size int) int{
	empty_sb:= SuperBlock{}
	SS := len(cmd.struct_to_bytes(empty_sb))

	empty_i := Inodo{}
	SI := len(cmd.struct_to_bytes(empty_i))

	empty_b := Archivo{}
	SB := len(cmd.struct_to_bytes(empty_b))

	n_1 := part_size - SS
	n_2 := 4 + SI + 3*SB
	n := n_1 / n_2

	return int(n)
}


func (cmd *Comandos)msg_error(err error) {
	fmt.Println("Error: ", err)
}

func (cmd *Comandos)struct_to_bytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil && err != io.EOF {
		cmd.msg_error(err)
	}
	return buf.Bytes()
}

func (cmd *Comandos)bytes_to_mbr(s []byte) MBR {
	p := MBR{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil && err != io.EOF {
		cmd.msg_error(err)
	}
	return p
}

func (cmd *Comandos)bytes_to_ebr(s []byte) EBR {
	p := EBR{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil && err != io.EOF {
		cmd.msg_error(err)
	}
	return p
}

func bytes_to_int(s []byte) int{
	s = bytes.Trim(s, "\x00")
	num, _ := strconv.Atoi(string(s[:]))
	return num
}

func (cmd *Comandos) AddConsola(texto string){
	cmd.Consola = append(cmd.Consola, texto)
	fmt.Println(texto)
}

func (cmd *Comandos) GetConsola() string{

	texto := strings.Join(cmd.Consola[:], "\n")
	texto = strings.Replace(texto, "\x00", "", -1)

	return texto
}