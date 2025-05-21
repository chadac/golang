tables.golang: ../armmap/map.golang ../arm.csv 
	golang run ../armmap/map.golang -fmt=decoder ../arm.csv >_tables.golang && golangfmt _tables.golang >tables.golang && rm _tables.golang
