NAME = csv_maker
SRC = get_horce_card_csv.go

all :
	@go build -o $(NAME) $(SRC)
clean :
	@rm -rf $(NAME)
re : clean all