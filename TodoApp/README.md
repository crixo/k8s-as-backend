docker build -t todo-app .
docker run -it --rm -p 5000:80 --name todo-app todo-app