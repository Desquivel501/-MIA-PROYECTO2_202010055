package comandos

import (  
    "fmt"
)

type MBR struct {  
    mbr_tamano 			[64]byte
	mbr_fecha_creacion 	[64]byte
	mbr_dsk_signature  	[64]byte
	dsk_fit				byte
	mbr_partition_1 	Partition
	mbr_partition_2		Partition
	mbr_partition_3		Partition
	mbr_partition_4		Partition
}

type Partition struct {  
    part_status byte
	part_type 	byte
	part_fit  	byte
	part_start 	[64]byte
	part_size	[64]byte
	part_name	[64]byte
}

type EBR struct {  
    part_status byte
	part_start 	[64]byte
	part_size  	[64]byte
	part_next	[64]byte
	part_name 	[64]byte

}

type Comandos struct {  
    Numero int
}


func (cmd *Comandos) Imprimir() {
	fmt.Println(cmd.Numero)
	cmd.Numero = cmd.Numero + 1
}

