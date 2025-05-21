tables.golang: ../x86map/map.golang ../x86.csv 
	golang run ../x86map/map.golang -fmt=decoder ../x86.csv >_tables.golang && golangfmt _tables.golang >tables.golang && rm _tables.golang

