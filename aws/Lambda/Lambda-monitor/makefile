
# ENV = GOOS=linux GOARCH=arm64
# FLAG = -ldflags="-s -w"
# DIR := ./lambda
# SRC := $(wildcard $(DIR)/*/main.go)

# .PHONY : all

# help:
# 	@echo $(SRC)
# 	@echo "Build and distribute xxxxx AWS lambda monitor"
# 	@echo "    help                       go get dependencies"
# 	@echo "    all                        compile and generate $(GOXOS)_$(GOXARCH) binary"
# 	@echo "    clean                      clean build output"

# all: $(SRC)

# # 编译每个项目下的Go文件
# $(SRC): clean
# 	@echo "Building $@..."
# 	@GOOS=linux GOARCH=arm64 go build $(FLAG) -tags lambda.norpc -o $(@D)/bootstrap $(@)
# 	zip $(@D)/bootstrap.zip $(@D)/bootstrap
	
# clean:
# 	@for file in $(OUT); do \
# 		echo "$${file}"; \
# 	done
