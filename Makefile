PACKAGE_NAME=cmd

OUTPUT_FOLDER=out

BINARY_NAME=${OUTPUT_FOLDER}/userportal

VERSION=1.0

ENTRY_POINT=main.go

run : 
	go run ./${PACKAGE_NAME}/${ENTRY_POINT}

build : 
	go build -o ${BINARY_NAME}_${VERSION} ./${PACKAGE_NAME}/${ENTRY_POINT}

clean :
	rm -f ${OUTPUT_FOLDER}/*