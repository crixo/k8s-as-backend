FROM mcr.microsoft.com/dotnet/core/sdk:3.1 AS build
WORKDIR /app

# copy csproj and restore as distinct layers
COPY *.sln .
COPY todo-api/*.csproj ./todo-api/
RUN dotnet restore

# copy everything else and build app
COPY todo-api/. ./todo-api/
WORKDIR /app/todo-api
RUN dotnet publish -c Release -o out


FROM mcr.microsoft.com/dotnet/core/aspnet:3.1 AS runtime
WORKDIR /app
COPY --from=build /app/todo-api/out ./
ENTRYPOINT ["dotnet", "TodoApi.dll"]