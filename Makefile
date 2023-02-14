NAME := oai
GP := /usr/local/bin

default :: build

format ::
	gofmt -w *.go

build :: format
	go build

cap :: 
	setcap cap_net_admin,cap_net_raw,cap_dac_read_search,cap_sys_ptrace=+eip $(NAME)

install :: 
	mv $(NAME) $(GP)

