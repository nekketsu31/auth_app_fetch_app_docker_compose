# auth_app

## Description
This is an Restful API for Node.js and Mysql. 

##### Routing         : Express
##### ORM Database    : Sequelize
##### Authentication  : Passport, JWT

## Installation

#### Donwload Code | Clone the Repo

```
git clone {repo_name}
```

#### Install Node Modules
```
yarn or npm install
```

#### Create .env File
You will find a example.env file in the home directory. Paste the contents of that into a file named .env in the same directory. 
Fill in the variables to fit your application

#### Migration Database
```
npx sequelize-cli db:migrate --config "config/config.json" --env "local" up
```

#### Run
```
yarn start
```
or
```
npm start
```