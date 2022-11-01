package comandos

import (  
    "fmt"
	"os"
	"bytes"
	"encoding/gob"
	"io"
	"strconv"
	"time"
	"strings"
	"math/rand"
	// "math"
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
	Creado_home bool
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

type DiskData struct{
	Block_start int
	Block_size int
	Inode_start int
	Inode_size int
}


type Comandos struct {  
    Mounted_list []Mounted
	Id_disco int
	Consola []string
	Graph string
	Usuario string
	Root bool
	Part_id string
	Creado_home bool
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

	path_dir := slicePath(path)
	path_dir = path_dir[:len(path_dir)-1]
	path_dir_str := strings.Join(path_dir, "/")

	// fmt.Println(path_dir_str)
	err := os.MkdirAll(path_dir_str, os.ModePerm)
	if err != nil {
		// log.Println(err)
		fmt.Println("here")
		cmd.msg_error(err)
		return
	}

	disco, err := os.Create(path)
	if err != nil {
		cmd.msg_error(err)
		return
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
		return
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

func (cmd *Comandos) PushMount(id string, mount Mounted){
	for i, m := range cmd.Mounted_list{
		if(m.Id == id){
			cmd.Mounted_list[i] = mount
		}
	}
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
	copy(super.S_free_blocks_count[:], strconv.Itoa(n*3 - 2))
	copy(super.S_free_inodes_count[:], strconv.Itoa(n - 2))
	copy(super.S_mtime[:], date)
	copy(super.S_magic[:], "0xEF53")
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

	copy(users.I_block[0][:], strconv.Itoa(1))
	for i := 1; i < 16; i += 1{
		copy(users.I_block[i][:], strconv.Itoa(-1))
	}
	
	users_txt := "1,G,root\n1,U,root,123"
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
	
	copy(root.I_block[0][:], strconv.Itoa(0))
	
	for i := 1; i < 16; i += 1{
		copy(root.I_block[i][:], strconv.Itoa(-1))
	}
	
	//-----------------------------------------------------	ESCRIBIR EN DISCO
	
	bm_inodos := make([]byte, n)
	disco.WriteAt(bm_inodos, int64(part_start + SS))

	bm_bloques := make([]byte, 3*n)
	disco.WriteAt(bm_bloques, int64(part_start + SS + n))

	root_bytes := cmd.struct_to_bytes(root)
	disco.WriteAt(root_bytes, int64(part_start + SS + n + 3*n))

	users_bytes := cmd.struct_to_bytes(users)
	disco.WriteAt(users_bytes, int64(part_start + SS + n + 3*n + SI))

	cmd.AddConsola("[MIA]@Proyecto2:~$ La particion " + id + " ha sido formateada exitosamente")

	disco.Close()
}

func (cmd *Comandos) GetUsers(id string) (string, string) {
	part, err_ := cmd.GetMount(id)
	if(err_ == 0){
		cmd.AddConsola("[MIA]@Proyecto2:~$ La particion no ha sido montada")
		return "ERROR", "[MIA]@Proyecto2:~$ La particion no ha sido montada"
	}

	disco, err := os.OpenFile(part.Path, os.O_RDWR, 0660)
	if err != nil {
		cmd.msg_error(err)
		return "ERROR", "[MIA]@Proyecto2:~$ No se ha podido abrir el disco"
	}

	part_start := bytes_to_int(part.Part.Part_start[:])

	super_block := cmd.leerSuper(disco, int64(part_start))

	inodo_size := bytes_to_int(super_block.S_inode_size[:])
	block_size := bytes_to_int(super_block.S_block_size[:])
	block_start := bytes_to_int(super_block.S_block_start[:])

	inicio_root := bytes_to_int(super_block.S_inode_start[:])
	
	carpeta := cmd.leerCarpeta(disco, int64(block_start))

	contenido := carpeta.B_content[0]

	num_inodo_users := bytes_to_int(contenido.B_inodo[:])

	inodo_users := cmd.leerInodo(disco, int64(inicio_root + num_inodo_users*inodo_size))

	contenido_str := ""
	for i := 0; i < 16; i += 1{
		bloque_pos := bytes_to_int(inodo_users.I_block[i][:])

		if(bloque_pos == -1){
			continue
		}

		bloque := cmd.leerArchivo(disco, int64(block_start + bloque_pos*block_size))

		// fmt.Println(strconv.Itoa(bloque_pos) + "------" + strconv.Itoa(block_size) + " --- " +strconv.Itoa(block_start + bloque_pos*block_size))

		content := bytes.Trim(bloque.B_content[:], "\x00")
		
		contenido_str += string(content[:])
	}

	fmt.Println(contenido_str)
	disco.Close()
	return contenido_str, ""
}


func (cmd *Comandos)Login(usuario string, password string, id string) string{

	texto, error_msg := cmd.GetUsers(id)
	if texto == "ERROR"{
		return error_msg
	}
	 
	lineas := strings.Split(texto, "\n")

	for _, linea := range lineas {

		if linea == ""{
			continue
		}

		linea_array := strings.Split(linea, ",")

		if linea_array[0] == "0"{
			continue
		}

		if linea_array[1] == "G"{
			fmt.Println("Grupo: " + linea_array[2])
		}

		
		if linea_array[1] == "U"{
			fmt.Println("Usuario: " + linea_array[2] + ", Contraseña: " + linea_array[3])

			if(linea_array[2] == usuario){
				if(linea_array[3] == password){

					if(linea_array[2] == "root"){
						cmd.Root = true
					}
					cmd.Part_id = id
					cmd.Usuario = usuario

					return ""
				}
				return "[MIA]@Proyecto2:~$ Contraseña Incorrecta"
			}
		} 
	} 
	return "[MIA]@Proyecto2:~$ No se ha encontrado el usuario"
}


func (cmd *Comandos) Logout() string {

	username := cmd.Usuario
	if (username == ""){
		return "[MIA]@Proyecto2:~$ No se ha iniciado sesion"
	}

	cmd.Root = false
	cmd.Part_id = ""
	cmd.Usuario = ""
	return ""
}


func (cmd *Comandos) Mkusr(username string, password string, grupo string){

	if( cmd.Root == false){
		cmd.AddConsola("[MIA]@Proyecto2:~$ Usuario no es root")
		return
	}

	id := cmd.Part_id

	part, err_ := cmd.GetMount(id)
	if(err_ == 0){
		cmd.AddConsola("[MIA]@Proyecto2:~$ La particion no ha sido montada")
		return
	}

	disco, err := os.OpenFile(part.Path, os.O_RDWR, 0660)
	if err != nil {
		cmd.msg_error(err)
		// return "ERROR", "[MIA]@Proyecto2:~$ No se ha podido abrir el disco"
	}

	part_start := bytes_to_int(part.Part.Part_start[:])

	super_block := cmd.leerSuper(disco, int64(part_start))

	inodo_size := bytes_to_int(super_block.S_inode_size[:])
	block_size := bytes_to_int(super_block.S_block_size[:])
	block_start := bytes_to_int(super_block.S_block_start[:])

	inicio_root := bytes_to_int(super_block.S_inode_start[:])
	
	carpeta := cmd.leerCarpeta(disco, int64(block_start))

	contenido := carpeta.B_content[0]

	num_inodo_users := bytes_to_int(contenido.B_inodo[:])

	inodo_users := cmd.leerInodo(disco, int64(inicio_root + num_inodo_users*inodo_size))

	contenido_str := ""
	for i := 0; i < 16; i += 1{
		
		bloque_pos := bytes_to_int(inodo_users.I_block[i][:])
		bloque := cmd.leerArchivo(disco, int64(block_start + bloque_pos*block_size))

		content := bytes.Trim(bloque.B_content[:], "\x00")

		contenido_str += string(content[:])
	}

	lineas := strings.Split(contenido_str, "\n")

	for _, linea := range lineas {
		if linea == ""{
			continue
		}

		// fmt.Println("---------------" + linea)

		linea_array := strings.Split(linea, ",")

		if linea_array[0] == "0"{
			continue
		}

		if linea_array[1] == "U"{
			if(linea_array[2] == username){
				disco.Close()
				cmd.AddConsola("[MIA]@Proyecto2:~$ Ya existe un usuario con ese nombre")
				return
			}
		}

		
	} 

	existe_grupo := false
	cont := 0
	group_id := ""

	for i, linea := range lineas {
		if linea == ""{
			continue
		}
		linea_array := strings.Split(linea, ",")

		if linea_array[0] == "0"{
			continue
		}

		if linea_array[1] == "G"{
			// fmt.Println("Grupo: " + linea_array[2])
			if(linea_array[2] == grupo){
				existe_grupo = true
				cont = i + 1
				group_id = linea_array[0]
				break
			}
		}
		
	} 
	if(existe_grupo != true){
		cmd.AddConsola("[MIA]@Proyecto2:~$ No existe el grupo")
		disco.Close()
		return
	}

	// fmt.Println(group_id)
	nuevo_usuario := group_id + ",U," + username + "," + password

	lineas = append(lineas, "")
	copy(lineas[cont+1:], lineas[cont:])
	lineas[cont] = nuevo_usuario
	
	nuevo_txt := strings.Join(lineas[:], "\n")
	fmt.Println(nuevo_txt)


	content_arr := []byte(nuevo_txt)
	chunks := chunkSlice(content_arr)
	
	pos := 0
	usados := 0

	b_inicio := bytes_to_int(super_block.S_first_blo[:])

	var empty [16]byte

	for i := 0; i < len(chunks); i += 1 {

		if(i == 16){
			break
		}

		bloque_pos := bytes_to_int(inodo_users.I_block[i][:])
		
		if(bloque_pos != -1){
			pos = block_start + bloque_pos*block_size
		} else {
			pos = block_start + b_inicio*block_size
			
			inodo_users.I_block[i] = empty

			copy(inodo_users.I_block[i][:], strconv.Itoa(b_inicio))
			b_inicio +=  1
			usados += 1
		}

		block := Archivo{}
		copy(block.B_content[:], chunks[i])
		block_bytes := cmd.struct_to_bytes(block)
		disco.WriteAt(block_bytes, int64(pos))		
	}

	blo_usados := bytes_to_int(super_block.S_blocks_count[:])
	blo_libres := bytes_to_int(super_block.S_free_blocks_count[:])

	blo_libres = blo_libres - usados
	blo_usados = blo_usados + usados

	copy(super_block.S_blocks_count[:], strconv.Itoa(blo_usados))
	copy(super_block.S_free_blocks_count[:], strconv.Itoa(blo_libres))
	copy(super_block.S_first_blo[:], strconv.Itoa(b_inicio))

	super_bytes := cmd.struct_to_bytes(super_block)
	disco.WriteAt(super_bytes, int64(part_start))

	root_bytes := cmd.struct_to_bytes(inodo_users)
	disco.WriteAt(root_bytes, int64(inicio_root + num_inodo_users*inodo_size))

	cmd.AddConsola("[MIA]@Proyecto2:~$ Creado el usuario: " + username)

	disco.Close()
}


func (cmd *Comandos) Mkgrp(group_name string){

	if( cmd.Root == false){
		cmd.AddConsola("[MIA]@Proyecto2:~$ Usuario no es root")
		return
	}

	id := cmd.Part_id

	part, err_ := cmd.GetMount(id)
	if(err_ == 0){
		cmd.AddConsola("[MIA]@Proyecto2:~$ La particion no ha sido montada")
		return
	}

	disco, err := os.OpenFile(part.Path, os.O_RDWR, 0660)
	if err != nil {
		cmd.msg_error(err)
		// return "ERROR", "[MIA]@Proyecto2:~$ No se ha podido abrir el disco"
	}

	part_start := bytes_to_int(part.Part.Part_start[:])

	super_block := cmd.leerSuper(disco, int64(part_start))

	inodo_size := bytes_to_int(super_block.S_inode_size[:])
	block_size := bytes_to_int(super_block.S_block_size[:])
	block_start := bytes_to_int(super_block.S_block_start[:])

	inicio_root := bytes_to_int(super_block.S_inode_start[:])
	
	carpeta := cmd.leerCarpeta(disco, int64(block_start))

	contenido := carpeta.B_content[0]

	num_inodo_users := bytes_to_int(contenido.B_inodo[:])

	inodo_users := cmd.leerInodo(disco, int64(inicio_root + num_inodo_users*inodo_size))

	contenido_str := ""
	for i := 0; i < 16; i += 1{
		
		bloque_pos := bytes_to_int(inodo_users.I_block[i][:])
		bloque := cmd.leerArchivo(disco, int64(block_start + bloque_pos*block_size))

		content := bytes.Trim(bloque.B_content[:], "\x00")

		contenido_str += string(content[:])
	}

	lineas := strings.Split(contenido_str, "\n")

	existe_grupo := false
	group_id := 0

	for _, linea := range lineas {
		if linea == ""{
			continue
		}
		linea_array := strings.Split(linea, ",")

		if linea_array[0] == "0"{
			continue
		}

		if linea_array[1] == "G"{
			// fmt.Println("Grupo: " + linea_array[2])
			if(linea_array[2] == group_name){
				existe_grupo = true
				break
			}

			id := linea_array[0]
			id_int, _ := strconv.Atoi(id)
			if(group_id < id_int){
				group_id = id_int
			}
		}
		
	} 
	if(existe_grupo == true){
		cmd.AddConsola("[MIA]@Proyecto2:~$ Ya existe un grupo con ese nombre")
		disco.Close()
		return
	}

	group_id = group_id + 1
	// fmt.Println(group_id)
	nuevo_grupo := strconv.Itoa(group_id) + ",G," + group_name

	lineas = append(lineas, nuevo_grupo)
	
	nuevo_txt := strings.Join(lineas[:], "\n")
	// fmt.Println(nuevo_txt)


	content_arr := []byte(nuevo_txt)
	chunks := chunkSlice(content_arr)
	
	pos := 0
	usados := 0
	b_inicio := bytes_to_int(super_block.S_first_blo[:])

	var empty [16]byte

	for i := 0; i < len(chunks); i += 1 {

		if(i == 16){
			break
		}

		bloque_pos := bytes_to_int(inodo_users.I_block[i][:])
		
		if(bloque_pos != -1){
			pos = block_start + bloque_pos*block_size
		} else {
			pos = block_start + b_inicio*block_size
			inodo_users.I_block[i] = empty
			copy(inodo_users.I_block[i][:], strconv.Itoa(b_inicio))
			b_inicio = b_inicio + 1
			usados += 1
		}

		block := Archivo{}
		copy(block.B_content[:], chunks[i])
		block_bytes := cmd.struct_to_bytes(block)
		disco.WriteAt(block_bytes, int64(pos))		
	}


	blo_usados := bytes_to_int(super_block.S_blocks_count[:])
	blo_libres := bytes_to_int(super_block.S_free_blocks_count[:])

	blo_libres = blo_libres - usados
	blo_usados = blo_usados + usados

	copy(super_block.S_blocks_count[:], strconv.Itoa(blo_usados))
	copy(super_block.S_free_blocks_count[:], strconv.Itoa(blo_libres))

	copy(super_block.S_first_blo[:], strconv.Itoa(b_inicio))
	super_bytes := cmd.struct_to_bytes(super_block)
	disco.WriteAt(super_bytes, int64(part_start))

	root_bytes := cmd.struct_to_bytes(inodo_users)
	disco.WriteAt(root_bytes, int64(inicio_root + num_inodo_users*inodo_size))


	// cmd.WriteFile(disco )

	cmd.AddConsola("[MIA]@Proyecto2:~$ Creado el grupo: " + group_name)
	disco.Close()
}


func (cmd *Comandos) Rmgrp(group_name string){

	if( cmd.Root == false){
		cmd.AddConsola("[MIA]@Proyecto2:~$ Usuario no es root")
		return
	}

	id := cmd.Part_id

	part, err_ := cmd.GetMount(id)
	if(err_ == 0){
		cmd.AddConsola("[MIA]@Proyecto2:~$ La particion no ha sido montada")
		return
	}

	disco, err := os.OpenFile(part.Path, os.O_RDWR, 0660)
	if err != nil {
		cmd.msg_error(err)
		// return "ERROR", "[MIA]@Proyecto2:~$ No se ha podido abrir el disco"
	}

	part_start := bytes_to_int(part.Part.Part_start[:])

	super_block := cmd.leerSuper(disco, int64(part_start))

	inodo_size := bytes_to_int(super_block.S_inode_size[:])
	block_size := bytes_to_int(super_block.S_block_size[:])
	block_start := bytes_to_int(super_block.S_block_start[:])

	inicio_root := bytes_to_int(super_block.S_inode_start[:])
	
	carpeta := cmd.leerCarpeta(disco, int64(block_start))

	contenido := carpeta.B_content[0]

	num_inodo_users := bytes_to_int(contenido.B_inodo[:])

	inodo_users := cmd.leerInodo(disco, int64(inicio_root + num_inodo_users*inodo_size))

	contenido_str := ""
	for i := 0; i < 16; i += 1{
		
		bloque_pos := bytes_to_int(inodo_users.I_block[i][:])
		bloque := cmd.leerArchivo(disco, int64(block_start + bloque_pos*block_size))

		content := bytes.Trim(bloque.B_content[:], "\x00")

		contenido_str += string(content[:])
	}

	lineas := strings.Split(contenido_str, "\n")

	group_id := ""

	for _, linea := range lineas {
		if linea == ""{
			continue
		}
		linea_array := strings.Split(linea, ",")

		if linea_array[0] == "0"{
			continue
		}

		if linea_array[1] == "G"{
			// fmt.Println("Grupo: " + linea_array[2])
			if(linea_array[2] == group_name){
				group_id = linea_array[0]
				break
			}
		}
		
	} 

	if(group_id == ""){
		cmd.AddConsola("[MIA]@Proyecto2:~$ El grupo no existe o ya ha sido eliminado")
		disco.Close()
		return
	}

	var nuevas_lineas []string
	for _, linea := range lineas {
		if linea == ""{
			continue
		}
		linea_array := strings.Split(linea, ",")

		if linea_array[0] == group_id{
			linea_array[0] = "0"
		}

		n_linea := strings.Join(linea_array, ",")
		nuevas_lineas = append(nuevas_lineas, n_linea)
	}
	
	nuevo_txt := strings.Join(nuevas_lineas[:], "\n")

	content_arr := []byte(nuevo_txt)
	chunks := chunkSlice(content_arr)
	
	pos := 0

	b_inicio := bytes_to_int(super_block.S_first_blo[:])

	var empty [16]byte

	for i := 0; i < len(chunks); i += 1 {

		if(i == 16){
			break
		}

		bloque_pos := bytes_to_int(inodo_users.I_block[i][:])
		
		if(bloque_pos != -1){
			pos = block_start + bloque_pos*block_size
		} else {
			pos = block_start + b_inicio*block_size
			inodo_users.I_block[i] = empty
			copy(inodo_users.I_block[i][:], strconv.Itoa(b_inicio))
			b_inicio = b_inicio + 1
		}

		block := Archivo{}
		copy(block.B_content[:], chunks[i])
		block_bytes := cmd.struct_to_bytes(block)
		disco.WriteAt(block_bytes, int64(pos))		
	}

	copy(super_block.S_first_blo[:], strconv.Itoa(b_inicio))
	super_bytes := cmd.struct_to_bytes(super_block)
	disco.WriteAt(super_bytes, int64(part_start))

	root_bytes := cmd.struct_to_bytes(inodo_users)
	disco.WriteAt(root_bytes, int64(inicio_root + num_inodo_users*inodo_size))

	disco.Close()

	cmd.AddConsola("[MIA]@Proyecto2:~$ Eliminado el grupo: " + group_name)
}

func (cmd *Comandos) Rmusr(username string){

	if( cmd.Root == false){
		cmd.AddConsola("[MIA]@Proyecto2:~$ Usuario no es root")
		return
	}

	id := cmd.Part_id

	part, err_ := cmd.GetMount(id)
	if(err_ == 0){
		cmd.AddConsola("[MIA]@Proyecto2:~$ La particion no ha sido montada")
		return
	}

	disco, err := os.OpenFile(part.Path, os.O_RDWR, 0660)
	if err != nil {
		cmd.msg_error(err)
		// return "ERROR", "[MIA]@Proyecto2:~$ No se ha podido abrir el disco"
	}

	part_start := bytes_to_int(part.Part.Part_start[:])

	super_block := cmd.leerSuper(disco, int64(part_start))

	inodo_size := bytes_to_int(super_block.S_inode_size[:])
	block_size := bytes_to_int(super_block.S_block_size[:])
	block_start := bytes_to_int(super_block.S_block_start[:])

	inicio_root := bytes_to_int(super_block.S_inode_start[:])
	
	carpeta := cmd.leerCarpeta(disco, int64(block_start))

	contenido := carpeta.B_content[0]

	num_inodo_users := bytes_to_int(contenido.B_inodo[:])

	inodo_users := cmd.leerInodo(disco, int64(inicio_root + num_inodo_users*inodo_size))

	contenido_str := ""
	for i := 0; i < 16; i += 1{
		
		bloque_pos := bytes_to_int(inodo_users.I_block[i][:])
		bloque := cmd.leerArchivo(disco, int64(block_start + bloque_pos*block_size))

		content := bytes.Trim(bloque.B_content[:], "\x00")

		contenido_str += string(content[:])
	}

	lineas := strings.Split(contenido_str, "\n")
	var nuevas_lineas []string
	eliminado := false
	for _, linea := range lineas {
		if linea == ""{
			continue
		}
		linea_array := strings.Split(linea, ",")

		if linea_array[0] == "0"{
			linea = strings.Join(linea_array[:], ",")
			nuevas_lineas = append(nuevas_lineas, linea)
			continue
		}

		if linea_array[1] == "U"{
			// fmt.Println("Grupo: " + linea_array[2])
			if(linea_array[2] == username){
				linea_array[0] = "0"
				eliminado = true
			}
		}

		linea = strings.Join(linea_array[:], ",")

		nuevas_lineas = append(nuevas_lineas, linea)
	} 

	if(!eliminado){
		cmd.AddConsola("[MIA]@Proyecto2:~$ El usuario no existe o ya ha sido eliminado")
		disco.Close()
		return
	}

	nuevo_txt := strings.Join(nuevas_lineas[:], "\n")

	content_arr := []byte(nuevo_txt)
	chunks := chunkSlice(content_arr)
	
	pos := 0

	b_inicio := bytes_to_int(super_block.S_first_blo[:])

	var empty [16]byte

	for i := 0; i < len(chunks); i += 1 {

		if(i == 16){
			break
		}

		bloque_pos := bytes_to_int(inodo_users.I_block[i][:])
		
		if(bloque_pos != -1){
			pos = block_start + bloque_pos*block_size
		} else {
			pos = block_start + b_inicio*block_size
			inodo_users.I_block[i] = empty
			copy(inodo_users.I_block[i][:], strconv.Itoa(b_inicio))
			b_inicio = b_inicio + 1
		}

		block := Archivo{}
		copy(block.B_content[:], chunks[i])
		block_bytes := cmd.struct_to_bytes(block)
		disco.WriteAt(block_bytes, int64(pos))		
	}

	copy(super_block.S_first_blo[:], strconv.Itoa(b_inicio))
	super_bytes := cmd.struct_to_bytes(super_block)
	disco.WriteAt(super_bytes, int64(part_start))

	root_bytes := cmd.struct_to_bytes(inodo_users)
	disco.WriteAt(root_bytes, int64(inicio_root + num_inodo_users*inodo_size))

	disco.Close()

	cmd.AddConsola("[MIA]@Proyecto2:~$ Eliminado el usuario: " + username)
}


func (cmd *Comandos) Mkfile(file_size int, path_file string){

	id := cmd.Part_id

	no_inodos := 0
	no_bloques := 0

	part, err_ := cmd.GetMount(id)
	if(err_ == 0){
		cmd.AddConsola("[MIA]@Proyecto2:~$ La particion no ha sido montada")
		return
	}

	path_name_arr := slicePath(path_file)

	if(len(path_name_arr) > 3){
		cmd.AddConsola("[MIA]@Proyecto2:~$ No existe la ruta")
		return
		
	}

	disco, err := os.OpenFile(part.Path, os.O_RDWR, 0660)
	if err != nil {
		cmd.msg_error(err)
		// fmt.Println("here")
		// return "ERROR", "[MIA]@Proyecto2:~$ No se ha podido abrir el disco"
	}


	

	path_name := path_name_arr[len(path_name_arr)-1]



	currentTime := time.Now()
	date := currentTime.Format("02-01-2006")

	part_start := bytes_to_int(part.Part.Part_start[:])

	super_block := cmd.leerSuper(disco, int64(part_start))

	inodo_start := bytes_to_int(super_block.S_inode_start[:])
	inodo_size := bytes_to_int(super_block.S_inode_size[:])
	block_size := bytes_to_int(super_block.S_block_size[:])
	block_start := bytes_to_int(super_block.S_block_start[:])

	// inicio_root := bytes_to_int(super_block.S_inode_start[:])
	
	b_inicio := bytes_to_int(super_block.S_first_blo[:])
	i_inicio := bytes_to_int(super_block.S_first_ino[:])

	carpeta_root := cmd.leerCarpeta(disco, int64(block_start))

	// root_pos := i_inicio
	home := Inodo{}

	home_pos := 0

	if(part.Creado_home == false){
		// bloque_pos := bytes_to_int(inodo_users.I_block[i][:])
		// home := Inodo{}
		copy(home.I_uid[:], "1")
		copy(home.I_gid[:], "1")
		copy(home.I_size[:], strconv.Itoa(0))
		copy(home.I_atime[:], date)
		copy(home.I_ctime[:], date)
		copy(home.I_type[:], "0")
		copy(home.I_perm[:], "664")
		
		for i := 0; i < 16; i += 1{
			copy(home.I_block[i][:], strconv.Itoa(-1))
		}

		home_cont := Content{}
		copy(home_cont.B_name[:], "home")
		copy(home_cont.B_inodo[:], strconv.Itoa(i_inicio))

		carpeta_root.B_content[1] = home_cont

		home_bytes := cmd.struct_to_bytes(home)
		// disco.WriteAt(home_bytes, int64(inodo_start + inodo_size*i_inicio))

		_, err = disco.WriteAt(home_bytes, int64(inodo_start + inodo_size*i_inicio))
			if err != nil {
			cmd.msg_error(err)
			// fmt.Println(err)
		}

		home_pos = inodo_start + inodo_size*i_inicio

		i_inicio += 1
		
		carpeta_bytes := cmd.struct_to_bytes(carpeta_root)
		// disco.WriteAt(carpeta_bytes, int64(block_start))

		_, err = disco.WriteAt(carpeta_bytes, int64(block_start))
			if err != nil {
			cmd.msg_error(err)
			// fmt.Println(err)
		}

		no_inodos += 1
		part.Creado_home = true

		cmd.PushMount(id, part)

	} else {

		contenido := carpeta_root.B_content[1]
		num_inodo_home := bytes_to_int(contenido.B_inodo[:])
		home = cmd.leerInodo(disco, int64(inodo_start + num_inodo_home*inodo_size))

		home_pos = inodo_start + num_inodo_home*inodo_size
	}
	
	new_file := Inodo{}
	copy(new_file.I_uid[:], "1")
	copy(new_file.I_gid[:], "1")
	copy(new_file.I_size[:], strconv.Itoa(file_size))
	copy(new_file.I_atime[:], date)
	copy(new_file.I_ctime[:], date)
	copy(new_file.I_type[:], "1")
	copy(new_file.I_perm[:], "664")
	no_inodos += 1

	text := createContent(file_size)

	content_arr := []byte(text)
	chunks := chunkSlice(content_arr)

	for i := 0; i < 16; i += 1{
		copy(new_file.I_block[i][:], strconv.Itoa(-1))
	}

	var empty [16]byte
	for i := 0; i < len(chunks); i += 1 {

		if(i == 16){
			break
		}

		pos := block_start + b_inicio*block_size
		new_file.I_block[i] = empty
		copy(new_file.I_block[i][:], strconv.Itoa(b_inicio))
		b_inicio = b_inicio + 1
		no_bloques += 1

		block := Archivo{}
		copy(block.B_content[:], chunks[i])
		block_bytes := cmd.struct_to_bytes(block)
		// disco.WriteAt(block_bytes, int64(pos))
		
		_, err = disco.WriteAt(block_bytes, int64(pos))
			if err != nil {
			cmd.msg_error(err)
			// fmt.Println(err)
		}
		
	}

	new_bytes := cmd.struct_to_bytes(new_file)
	disco.WriteAt(new_bytes, int64(inodo_start + inodo_size*i_inicio))
	

	

	done := false

	fmt.Println(home_pos)

	for i := 0; i < 16; i += 1{
		
		bloque_pos := bytes_to_int(home.I_block[i][:])

		fmt.Println(bloque_pos)

		if(bloque_pos != -1){
			bloque := cmd.leerCarpeta(disco, int64(block_start + bloque_pos*block_size))

			for j := 0; j < 4; j+= 1{
				contenido := bloque.B_content[j]
				nuevo_inodo_pos := bytes_to_int(contenido.B_inodo[:])

				if(nuevo_inodo_pos == -1){
					new_cont := Content{}
					copy(new_cont.B_name[:], path_name)
					copy(new_cont.B_inodo[:], strconv.Itoa(i_inicio))
					i_inicio += 1
					bloque.B_content[j] = new_cont
					done = true

					bloque_bytes := cmd.struct_to_bytes(bloque)
					disco.WriteAt(bloque_bytes, int64(block_start + bloque_pos*block_size))

					break
				}
			}
		} else {
			empty_cont := Content{}
			copy(empty_cont.B_name[:], "")
			copy(empty_cont.B_inodo[:], strconv.Itoa(-1))

			new_cont := Content{}
			copy(new_cont.B_name[:], path_name)
			copy(new_cont.B_inodo[:], strconv.Itoa(i_inicio))
			i_inicio += 1

			new_dir := Carpeta{}
			var cont_arr [4]Content
			cont_arr[0] = new_cont
			cont_arr[1] = empty_cont
			cont_arr[2] = empty_cont
			cont_arr[3] = empty_cont
			new_dir.B_content = cont_arr

			fmt.Println(b_inicio)

			carpeta_bytes := cmd.struct_to_bytes(new_dir)
			disco.WriteAt(carpeta_bytes, int64(block_start + b_inicio*block_size))

			home.I_block[i] = empty
			copy(home.I_block[i][:], strconv.Itoa(b_inicio))

			b_inicio += 1
			no_bloques += 1

			break
		}	

		if(done){
			break
		}

	}



	home_bytes := cmd.struct_to_bytes(home)
	disco.WriteAt(home_bytes, int64(home_pos))


	blo_usados := bytes_to_int(super_block.S_blocks_count[:])
	blo_libres := bytes_to_int(super_block.S_free_blocks_count[:])

	blo_libres = blo_libres - no_bloques
	blo_usados = blo_usados + no_bloques

	copy(super_block.S_blocks_count[:], strconv.Itoa(blo_usados))
	copy(super_block.S_free_blocks_count[:], strconv.Itoa(blo_libres))



	ino_usados := bytes_to_int(super_block.S_inodes_count[:])
	ino_libres := bytes_to_int(super_block.S_free_inodes_count[:])

	ino_libres = ino_libres - no_inodos
	ino_usados = ino_usados + no_inodos

	copy(super_block.S_inodes_count[:], strconv.Itoa(ino_usados))
	copy(super_block.S_free_inodes_count[:], strconv.Itoa(ino_libres))

	copy(super_block.S_first_ino[:], strconv.Itoa(i_inicio))
	copy(super_block.S_first_blo[:], strconv.Itoa(b_inicio))
	super_bytes := cmd.struct_to_bytes(super_block)
	disco.WriteAt(super_bytes, int64(part_start))

	disco.Close()

	cmd.AddConsola("[MIA]@Proyecto2:~$ Creado el archivo: " + path_name)
}


func (cmd *Comandos)ShowFile(file_name string, id string){

	dot := ""
	dot += "digraph G {\n"

	dot += "node [fontname=\"Helvetica,Arial,sans-serif\"]\n"
	dot += "node [shape=box]\n"


	if(file_name == "/users.txt"){
		texto, _ := cmd.GetUsers(id)
		
		dot += "a[label=\" " + texto + "  \"]\n"
		dot += "}"

		cmd.Graph = dot
		return
	}


	part, err_ := cmd.GetMount(id)
	if(err_ == 0){
		cmd.AddConsola("[MIA]@Proyecto2:~$ La particion no ha sido montada")
	}

	disco, err := os.OpenFile(part.Path, os.O_RDWR, 0660)
	if err != nil {
		cmd.msg_error(err)
		// fmt.Println("here")
		// return "ERROR", "[MIA]@Proyecto2:~$ No se ha podido abrir el disco"
	}

	part_start := bytes_to_int(part.Part.Part_start[:])

	super_block := cmd.leerSuper(disco, int64(part_start))

	inodo_start := bytes_to_int(super_block.S_inode_start[:])
	inodo_size := bytes_to_int(super_block.S_inode_size[:])
	block_size := bytes_to_int(super_block.S_block_size[:])
	block_start := bytes_to_int(super_block.S_block_start[:])

	// inicio_root := bytes_to_int(super_block.S_inode_start[:])
	
	// b_inicio := bytes_to_int(super_block.S_first_blo[:])
	// i_inicio := bytes_to_int(super_block.S_first_ino[:])

	carpeta_root := cmd.leerCarpeta(disco, int64(block_start))

	contenido := carpeta_root.B_content[1]
	num_inodo_home := bytes_to_int(contenido.B_inodo[:])

	home := cmd.leerInodo(disco, int64(inodo_start + num_inodo_home*inodo_size))


	for i := 0; i < 16; i += 1{
		
		bloque_pos := bytes_to_int(home.I_block[i][:])

		if(bloque_pos != -1){
			bloque := cmd.leerCarpeta(disco, int64(block_start + bloque_pos*block_size))

			for j := 0; j < 4; j+= 1{

				contenido := bloque.B_content[j]

				if(strings.Contains(string(contenido.B_name[:]), file_name)){

					file_inodo_pos := bytes_to_int(contenido.B_inodo[:])

					file := cmd.leerInodo(disco, int64(inodo_start + file_inodo_pos*inodo_size))
					
					contenido_str := ""

					for i := 0; i < 16; i += 1{
						bloque_pos := bytes_to_int(file.I_block[i][:])

						if(bloque_pos == -1){
							continue
						}

						bloque := cmd.leerArchivo(disco, int64(block_start + bloque_pos*block_size))

						// fmt.Println(strconv.Itoa(bloque_pos) + "------" + strconv.Itoa(block_size) + " --- " +strconv.Itoa(block_start + bloque_pos*block_size))

						content := bytes.Trim(bloque.B_content[:], "\x00")
						
						contenido_str += string(content[:])
					}

					// cmd.AddConsola(contenido_str)

					dot += "a[label=\" " + contenido_str + "  \"]\n"
					dot += "}"

					cmd.Graph = dot
					return			

				}
			}
		} 
	}
	cmd.AddConsola("[MIA]@Proyecto2:~$ No se ha encontrado el archivo")
	dot += "a[label=\"\"]\n"
	dot += "}"

	cmd.Graph = dot

}


func (cmd *Comandos)ReporteTree(id string){

	part, err_ := cmd.GetMount(id)
	if(err_ == 0){
		cmd.AddConsola("[MIA]@Proyecto2:~$ La particion no ha sido montada")
	}

	disco, err := os.OpenFile(part.Path, os.O_RDWR, 0660)
	if err != nil {
		cmd.msg_error(err)
		// fmt.Println("here")
		// return "ERROR", "[MIA]@Proyecto2:~$ No se ha podido abrir el disco"
	}

	part_start := bytes_to_int(part.Part.Part_start[:])

	super_block := cmd.leerSuper(disco, int64(part_start))

	inodo_start := bytes_to_int(super_block.S_inode_start[:])
	inodo_size := bytes_to_int(super_block.S_inode_size[:])
	block_size := bytes_to_int(super_block.S_block_size[:])
	block_start := bytes_to_int(super_block.S_block_start[:])

	inicio_root := bytes_to_int(super_block.S_inode_start[:])
	
	root := cmd.leerInodo(disco, int64(inicio_root))

	data := DiskData{block_start, block_size, inodo_start, inodo_size}


	
	
	dot := "digraph g {\n"
	dot += "fontname=\"Helvetica,Arial,sans-serif\"\n"
	dot +=	"node [fontname=\"Helvetica,Arial,sans-serif\"]\n"
	dot +=	"node [shape = \"record\"]\n"
	dot +=	"graph [\n"
	dot +=	"rankdir = \"LR\"\n"
	dot += "];\n"

	dot = cmd.graficarInodo(dot, 0, root, data, disco)

	dot += "}\n"

	// cmd.AddConsola(dot)

	cmd.Graph = dot

}


func (cmd *Comandos)graficarInodo(dot string, no_inodo int, inodo Inodo, disk DiskData, disco *os.File) string{

	name := "inodo_" + strconv.Itoa(no_inodo)
	dot_inodo := ""

	dot_inodo += name + "[\n"
	dot_inodo += "style=filled;"
	dot_inodo += "color = \"#000000\"\n"
	dot_inodo += "fillcolor = \"#cfe2f3\"\n"
	dot_inodo += "label=\"inodo " + strconv.Itoa(no_inodo) + "|"
	dot_inodo += "{UID|" + bytes_to_string(inodo.I_uid[:]) + "}|"
	dot_inodo += "{GUID|" + bytes_to_string(inodo.I_gid[:]) + "}|"
	dot_inodo += "{SIZE| " + strconv.Itoa(bytes_to_int(inodo.I_size[:])) +"}|"

	dot_inodo += "{LECTURA|" + bytes_to_string(inodo.I_atime[:]) + "}|"
	dot_inodo += "{CREACION|" + bytes_to_string(inodo.I_atime[:]) + "}|"
	dot_inodo += "{MODIFICACION|" + bytes_to_string(inodo.I_atime[:]) + "}"

	tipo := string(inodo.I_type[:])

	


	for i := 0; i < 16; i += 1{

		bloque_pos := bytes_to_int(inodo.I_block[i][:])


		if(bloque_pos == -1){
			// dot_inodo += "|{AP" + strconv.Itoa(i) + " | -1 }"
			continue
		}

		padre := name + ":f" + strconv.Itoa(i) 
		dot_inodo += "|{AP" + strconv.Itoa(i) + " | <f" + strconv.Itoa(i) + "> " + strconv.Itoa(bloque_pos) + " }"

		if(strings.Contains(string(tipo), "1")){

			bloque := cmd.leerArchivo(disco, int64(disk.Block_start + bloque_pos*disk.Block_size))
			content := bytes.Trim(bloque.B_content[:], "\x00")
			dot = cmd.graficarBloqueArchivo(dot, padre, strconv.Itoa(bloque_pos),  bytes_to_string(content[:]))

		} else {

			bloque := cmd.leerCarpeta(disco, int64(disk.Block_start + bloque_pos*disk.Block_size))
			dot = cmd.graficarBloqueCarpeta(dot, padre, strconv.Itoa(bloque_pos), bloque, disk, disco)
		}
		
	} 
	dot_inodo += "|{TIPO|" + tipo + "}|"
	dot_inodo += "{PERMISOS|664}"

	dot_inodo += "\"\n"
	dot_inodo += "]\n"
	dot += dot_inodo

	return dot
}

func (cmd *Comandos)graficarBloqueArchivo(dot string, padre string, id string, contenido string) string {
	name := "bloque_" + id
	dot += name + "[\n"
	dot += "style=filled;"
	dot += "color = \"#000000\"\n"
	dot += "fillcolor = \"#fff2cc\"\n"
	dot += "label=\"bloque " + id + "|"
	dot +=  contenido + "\"\n"
	dot += "]\n"
	dot += padre + "->" + name + "\n"
	return dot
}


func (cmd *Comandos)graficarBloqueCarpeta(dot string, padre string, id string, carpeta Carpeta, disk DiskData, disco *os.File) string {
	
	name := "bloque_" + id
	dot_bloque := name + "[\n"
	dot_bloque += "style=filled;"
	dot_bloque += "color = \"#000000\"\n"
	dot_bloque += "fillcolor = \"#f4cccc\"\n"
	dot_bloque += "label=\"bloque " + id
	
	
	for i := 0; i < 4; i+= 1{
		contenido := carpeta.B_content[i]
		inodo_pos := bytes_to_int(contenido.B_inodo[:])

		if(inodo_pos == -1){
			dot_bloque += "|{   | -1 }"
			continue
		}

		name_arr := bytes.Trim(contenido.B_name[:], "\x00")
		name_cont := bytes_to_string(name_arr[:])
		
		dot_bloque += "|{ " + name_cont + "  | <f" + strconv.Itoa(i) + "> " + strconv.Itoa(inodo_pos) + " }"

		inodo := cmd.leerInodo(disco, int64(disk.Inode_start + inodo_pos*disk.Inode_size))

		dot = cmd.graficarInodo(dot, inodo_pos, inodo, disk, disco)
		
		dot += name + ":f" + strconv.Itoa(i) + "-> inodo_" + strconv.Itoa(inodo_pos) + "\n"

	}

	dot_bloque += "\"\n"
	dot_bloque += "]\n"
	dot += padre + "->" + name + "\n"
	dot += dot_bloque

	return dot
}


func (cmd *Comandos) ReporteSuper(id string){

	part, err_ := cmd.GetMount(id)
	if(err_ == 0){
		cmd.AddConsola("[MIA]@Proyecto2:~$ La particion no ha sido montada")
	}

	disco, err := os.OpenFile(part.Path, os.O_RDWR, 0660)
	if err != nil {
		cmd.AddConsola("[MIA]@Proyecto2:~$ No se ha podido abrir el disco")
	}

	part_start := bytes_to_int(part.Part.Part_start[:])

	super_block := cmd.leerSuper(disco, int64(part_start))


	dot := ""

	dot += "digraph G {\n"; 
    dot += "fontname=\"Helvetica,Arial,sans-serif\"\n"; 
	dot += "node [fontname=\"Helvetica,Arial,sans-serif\"]\n"; 
    dot += "rankdir=TB;\n"; 
    dot += "node [shape=record];\n"; 
    dot += "a[label = <<table border=\"0\" cellborder=\"1\" cellspacing=\"0\" cellpadding=\"4\">\n"; 
    dot += "            <tr> <td bgcolor=\"#800080\">  <font color=\"white\"> <b>REPORTE SUPER BLOQUE</b> </font> </td> <td bgcolor=\"#800080\"></td> </tr>\n";

	dot += "            <tr> <td> <b>s_filesystem_type</b> </td> <td> <b>" + bytes_to_string(super_block.S_filesystem_type[:]) + "</b> </td> </tr>\n"; 
	
	dot += "            <tr> <td bgcolor=\"#f2e5f2\"> <b>s_inodes_count</b> </td> <td bgcolor=\"#f2e5f2\"> <b>" + strconv.Itoa(bytes_to_int(super_block.S_inodes_count[:])) + "</b> </td> </tr>\n"; 
	dot += "            <tr> <td> <b>s_blocks_count</b> </td> <td> <b>" + strconv.Itoa(bytes_to_int(super_block.S_blocks_count[:])) + "</b> </td> </tr>\n"; 
	
	dot += "            <tr> <td bgcolor=\"#f2e5f2\"> <b>s_free_inodes_count</b> </td> <td bgcolor=\"#f2e5f2\"> <b>" + strconv.Itoa(bytes_to_int(super_block.S_free_inodes_count[:])) + "</b> </td> </tr>\n"; 
	dot += "            <tr> <td> <b>s_free_blocks_count</b> </td> <td> <b>" + strconv.Itoa(bytes_to_int(super_block.S_free_blocks_count[:])) + "</b> </td> </tr>\n"; 


	dot += "            <tr> <td bgcolor=\"#f2e5f2\"> <b>s_mtime</b> </td> <td bgcolor=\"#f2e5f2\"> <b>" + bytes_to_string(super_block.S_mtime[:]) + "</b> </td> </tr>\n"; 

	dot += "            <tr> <td> <b>s_mnt_count</b> </td> <td> <b>" + strconv.Itoa(bytes_to_int(super_block.S_mnt_count[:])) + "</b> </td> </tr>\n"; 

	dot += "            <tr> <td bgcolor=\"#f2e5f2\"> <b>s_magic</b> </td> <td bgcolor=\"#f2e5f2\"> <b>" + bytes_to_string(super_block.S_magic[:]) + "</b> </td> </tr>\n"; 

	dot += "            <tr> <td> <b>s_inode_size</b> </td> <td> <b>" + strconv.Itoa(bytes_to_int(super_block.S_inode_size[:])) + "</b> </td> </tr>\n"; 

	dot += "            <tr> <td bgcolor=\"#f2e5f2\"> <b>s_block_size</b> </td> <td bgcolor=\"#f2e5f2\"> <b>" + strconv.Itoa(bytes_to_int(super_block.S_block_size[:])) + "</b> </td> </tr>\n"; 

	dot += "            <tr> <td> <b>s_firts_ino</b> </td> <td> <b>" + strconv.Itoa(bytes_to_int(super_block.S_first_ino[:])) + "</b> </td> </tr>\n"; 

	dot += "            <tr> <td bgcolor=\"#f2e5f2\"> <b>s_first_blo</b> </td> <td bgcolor=\"#f2e5f2\"> <b>" + strconv.Itoa(bytes_to_int(super_block.S_first_blo[:])) + "</b> </td> </tr>\n"; 

	dot += "            <tr> <td> <b>s_bm_inode_start</b> </td> <td> <b>" + strconv.Itoa(bytes_to_int(super_block.S_bm_inode_start[:])) + "</b> </td> </tr>\n"; 

	dot += "            <tr> <td bgcolor=\"#f2e5f2\"> <b>s_bm_block_start</b> </td> <td bgcolor=\"#f2e5f2\"> <b>" + strconv.Itoa(bytes_to_int(super_block.S_bm_block_start[:])) + "</b> </td> </tr>\n"; 


	dot += "            <tr> <td> <b>s_inode_start</b> </td> <td> <b>" + strconv.Itoa(bytes_to_int(super_block.S_inode_start[:])) + "</b> </td> </tr>\n"; 

	dot += "            <tr> <td bgcolor=\"#f2e5f2\"> <b>s_block_start</b> </td> <td bgcolor=\"#f2e5f2\"> <b>" + strconv.Itoa(bytes_to_int(super_block.S_block_start[:])) + "</b> </td> </tr>\n"; 

	dot += "        </table>>\n"; 
    dot += "]\n"; 
    dot += "}\n"; 

	cmd.Graph = dot

	// dot += "            <tr> <td> <b>s_filesystem_type</b> </td> <td> <b>" + bytes_to_string(super_block.S_filesystem_type[:]) + "</b> </td> </tr>\n"; 
	// dot += "            <tr> <td> <b>s_filesystem_type</b> </td> <td> <b>" + bytes_to_string(super_block.S_filesystem_type[:]) + "</b> </td> </tr>\n"; 
	// dot += "            <tr> <td> <b>s_filesystem_type</b> </td> <td> <b>" + bytes_to_string(super_block.S_filesystem_type[:]) + "</b> </td> </tr>\n"; 
	// dot += "            <tr> <td> <b>s_filesystem_type</b> </td> <td> <b>" + bytes_to_string(super_block.S_filesystem_type[:]) + "</b> </td> </tr>\n"; 

	
}


func createContent(size int) string{
	cadena := [10]string{"0","1","2","3","4","5","6","7","8","9"}
	resultado := ""

	j := 0
	for i := 0; i < size; i+=1 {
		resultado = resultado + cadena[j]

		if(j == 9){
			j = 0
		}else{
			j += 1
		}
	}
	return resultado
}


func (cmd *Comandos)WriteFile(disco *os.File, content string, b_libre int) []int {
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


func (cmd *Comandos)leerInodo(disco *os.File, pos int64) Inodo{
	empty_ebr := Inodo{}
	size_ebr := len(cmd.struct_to_bytes(empty_ebr))
	buff_ebr := make([]byte, size_ebr)
	disco.ReadAt(buff_ebr, pos)
	ebr := cmd.bytes_to_inodo(buff_ebr)
	return ebr
}

func (cmd *Comandos)leerSuper(disco *os.File, pos int64) SuperBlock{
	empty_sb := SuperBlock{}
	size_sb := len(cmd.struct_to_bytes(empty_sb))
	buff_sb := make([]byte, size_sb)
	disco.ReadAt(buff_sb, pos)
	sb := cmd.bytes_to_super(buff_sb)
	return sb
}

func (cmd *Comandos)leerCarpeta(disco *os.File, pos int64) Carpeta{
	empty_sb := Carpeta{}
	size_sb := len(cmd.struct_to_bytes(empty_sb))
	buff_sb := make([]byte, size_sb)
	disco.ReadAt(buff_sb, pos)
	sb := cmd.bytes_to_carpeta(buff_sb)
	return sb
}

func (cmd *Comandos)leerArchivo(disco *os.File, pos int64) Archivo{
	empty_sb := Archivo{}
	size_sb := len(cmd.struct_to_bytes(empty_sb))
	buff_sb := make([]byte, size_sb)
	disco.ReadAt(buff_sb, pos)
	sb := cmd.bytes_to_archivo(buff_sb)
	return sb
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


func (cmd *Comandos) RepDisco(id string) {
	dot := ""
	
	for _, m := range cmd.Mounted_list{
		if(m.Id == id){

			disco, err := os.OpenFile(m.Path, os.O_RDWR, 0660)
			if err != nil {
				cmd.msg_error(err)
				return
			}
			
			nombre := m.Path
			nombre_arr := slicePath(nombre)
			nombre = nombre_arr[len(nombre_arr) - 1]

			dot += "digraph G {\n"
            dot += "fontname=\"Helvetica,Arial,sans-serif\"\n"
            dot += "node [fontname=\"Helvetica,Arial,sans-serif\"]\n"
            dot += "rankdir=TB;\n"; 
            dot += " node [shape=record];\n"
			
			dot += "label=\"" + nombre + "\"\n"
			
			dot += "a[label = \" "
            dot += "MBR"
			
			disk_size := bytes_to_int(m.Master.Mbr_tamano[:])
			
			size_mbr := len(cmd.struct_to_bytes(m.Master))
			
			// mbr_ocupado := (size_mbr / disk_size) * 100
			
			ocupado := 0.0
			part_size := 0
			total_ocupado := size_mbr

			

			if(strings.Contains(string(m.Master.Mbr_partition_1.Part_status[:]), "1")){
				part_size = bytes_to_int(m.Master.Mbr_partition_1.Part_size[:])
				total_ocupado += part_size

				if(strings.Contains(string(m.Master.Mbr_partition_1.Part_type[:]), "E")){
					dot += "| {Extendida | {"
					dot += cmd.repExtendida(disco, m.Master.Mbr_partition_1, disk_size, part_size)
					dot += "}}"
				}else{
					ocupado = (float64(part_size) / float64(disk_size))*100
					dot += " | Primaria \\n" + fmt.Sprintf("%.2f", ocupado) + "%"
				}

			}

			if(strings.Contains(string(m.Master.Mbr_partition_2.Part_status[:]), "1")){
				part_size = bytes_to_int(m.Master.Mbr_partition_2.Part_size[:])
				total_ocupado += part_size

				if(strings.Contains(string(m.Master.Mbr_partition_2.Part_type[:]), "E")){
					dot += "| {Extendida | {"
					dot += cmd.repExtendida(disco, m.Master.Mbr_partition_2, disk_size, part_size)
					dot += "}}"
				}else{
					ocupado = (float64(part_size) / float64(disk_size))*100
					dot += " | Primaria \\n" + fmt.Sprintf("%.2f", ocupado) + "%"
				}

				

			}

			if(strings.Contains(string(m.Master.Mbr_partition_3.Part_status[:]), "1")){
				part_size = bytes_to_int(m.Master.Mbr_partition_3.Part_size[:])
				total_ocupado += part_size

				if(strings.Contains(string(m.Master.Mbr_partition_3.Part_type[:]), "E")){
					dot += "| {Extendida | {"
					dot += cmd.repExtendida(disco, m.Master.Mbr_partition_3, disk_size, part_size)
					dot += "}}"
				}else{
					ocupado = (float64(part_size) / float64(disk_size))*100
					dot += " | Primaria \\n" + fmt.Sprintf("%.2f", ocupado) + "%"
				}

			}

			if(strings.Contains(string(m.Master.Mbr_partition_4.Part_status[:]), "1")){
				part_size = bytes_to_int(m.Master.Mbr_partition_4.Part_size[:])
				total_ocupado += part_size

				if(strings.Contains(string(m.Master.Mbr_partition_4.Part_type[:]), "E")){
					dot += "| {Extendida | {"
					dot += cmd.repExtendida(disco, m.Master.Mbr_partition_4, disk_size, part_size)
					dot += "}}"
				}else{
					ocupado = (float64(part_size) / float64(disk_size))*100
					dot += " | Primaria \\n" + fmt.Sprintf("%.2f", ocupado) + "%"
				}

			}
			
			ocupado = (float64(disk_size - total_ocupado) / float64(disk_size))*100
			dot += "| Libre \\n" + fmt.Sprintf("%.2f", ocupado) + "%"

			dot += "\"]"
			dot += "}\n"
		}

	}
	cmd.Graph = dot
}


func (cmd *Comandos)repExtendida(disco *os.File, extendida Partition, size_total int, size_part int) string{
	inicio_part := bytes_to_int(extendida.Part_start[:])
	ebr := cmd.leerEBR(disco, int64(inicio_part))
	
	dot := ""

	logica_size := 0
	ext_size := 0
	ocupado := 0.0

	if(strings.Contains(string(ebr.Part_status[:]), "1")){
		logica_size = bytes_to_int(ebr.Part_size[:])
		ext_size += logica_size

		ocupado = (float64(logica_size) / float64(size_total))*100
		dot += "Logica \\n" + fmt.Sprintf("%.2f", ocupado) + "%"

		for{
			if(strings.Contains(string(ebr.Part_next[:]), "-1")){
				break
			}
			inicio_ebr := bytes_to_int(ebr.Part_next[:])
			ebr = cmd.leerEBR(disco, int64(inicio_ebr))

			logica_size = bytes_to_int(ebr.Part_size[:])
			ext_size += logica_size

			ocupado = (float64(logica_size) / float64(size_total))*100
			dot += "| Logica \\n" + fmt.Sprintf("%.2f", ocupado) + "%"

		}
	}

	ocupado = (float64(size_part - ext_size) / float64(size_total))*100
	dot += "| Libre \\n" + fmt.Sprintf("%.2f", ocupado) + "%"

	return dot
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

func (cmd *Comandos)bytes_to_super(s []byte) SuperBlock {
	p := SuperBlock{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil && err != io.EOF {
		cmd.msg_error(err)
	}
	return p
}

func (cmd *Comandos)bytes_to_inodo(s []byte) Inodo {
	p := Inodo{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil && err != io.EOF {
		cmd.msg_error(err)
	}
	return p
}

func (cmd *Comandos)bytes_to_carpeta(s []byte) Carpeta {
	p := Carpeta{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil && err != io.EOF {
		cmd.msg_error(err)
	}
	return p
}

func (cmd *Comandos)bytes_to_content(s []byte) Content {
	p := Content{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil && err != io.EOF {
		cmd.msg_error(err)
	}
	return p
}

func (cmd *Comandos)bytes_to_archivo(s []byte) Archivo {
	p := Archivo{}
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

func bytes_to_string(s []byte) string{
	s = bytes.Trim(s, "\x00")
	return string(s[:])
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