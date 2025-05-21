tables.golang: ../s390xmap/map.golang ../s390x.csv
	golang run ../s390xmap/map.golang -fmt=decoder ../s390x.csv >_tables.golang && golangfmt _tables.golang >tables.golang && rm _tables.golang
