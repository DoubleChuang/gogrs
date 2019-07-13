#TIME = $(shell date +"%Y%m%d%H%M%S")
TIME = $(DATE)
ifeq ($(TIME), )
	TIME = $(shell date +"%Y%m%d")
endif

TARGET_DIR=/run/shm/.gogrscache
BACKUP_DIR=$(PWD)/STOCK_DATA_BACKUP


.PHONY: backup recovery

all:
	@go build main.go

backup:
	@tar cf $(BACKUP_DIR)/$(TIME).tar $(TARGET_DIR)
recovery:
	@tar xvf $(BACKUP_DIR)/$(TIME).tar -C /
rmtmp:
	@find ./ -name '*~' -exec rm -ir {} \;
